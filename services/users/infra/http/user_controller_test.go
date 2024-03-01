package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/http"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type TestingDependencies_TestUserController struct {
	ctx                      context.Context
	ctrl                     *gomock.Controller
	mockLoggerImpl           *mocks.MockLoggerInterface
	mockGetUserByIdImpl      *mocks.MockGetUserByIdInterface
	mockDeleteUserByIdImpl   *mocks.MockDeleteUserByIdInterface
	mockCreateUserImpl       *mocks.MockCreateUserInterface
	mockGetUserPaginatedImpl *mocks.MockGetUserPaginatedInterface
	mockUpdateUserByIdImpl   *mocks.MockUpdateUserByIdInterface
	userController           *http.UserControllerImpl
}

func BeforeEach_TestUserController(t *testing.T) *TestingDependencies_TestUserController {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)

	mockLoggerImpl := mocks.NewMockLoggerInterface(ctrl)
	mockGetUserByIdImpl := mocks.NewMockGetUserByIdInterface(ctrl)
	mockDeleteUserByIdImpl := mocks.NewMockDeleteUserByIdInterface(ctrl)
	mockCreateUserImpl := mocks.NewMockCreateUserInterface(ctrl)
	mockGetUserPaginatedImpl := mocks.NewMockGetUserPaginatedInterface(ctrl)
	mockUpdateUserByIdImpl := mocks.NewMockUpdateUserByIdInterface(ctrl)

	mockLoggerImpl.
		EXPECT().
		WithCtx(gomock.Any()).
		AnyTimes().
		Return(&logger.Logger{})

	userController := http.NewUserControllerImpl(
		mockLoggerImpl,
		mockGetUserByIdImpl,
		mockDeleteUserByIdImpl,
		mockCreateUserImpl,
		mockGetUserPaginatedImpl,
		mockUpdateUserByIdImpl,
	)

	return &TestingDependencies_TestUserController{
		ctx:                      ctx,
		ctrl:                     ctrl,
		mockLoggerImpl:           mockLoggerImpl,
		mockGetUserByIdImpl:      mockGetUserByIdImpl,
		mockDeleteUserByIdImpl:   mockDeleteUserByIdImpl,
		mockCreateUserImpl:       mockCreateUserImpl,
		mockGetUserPaginatedImpl: mockGetUserPaginatedImpl,
		userController:           userController,
		mockUpdateUserByIdImpl:   mockUpdateUserByIdImpl,
	}
}

func TestUserController_GetUserById(t *testing.T) {
	deps := BeforeEach_TestUserController(t)
	defer deps.ctrl.Finish()

	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		id := primitive.NewObjectID().Hex()
		deps.mockGetUserByIdImpl.
			EXPECT().
			Do(gomock.Any(), id).
			Times(1).
			Return(nil, errors.New(exception.CodeNotFound))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Get("/api/v1/users/:id", deps.userController.GetUserById)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", id), nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 404, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Entity not found", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should return user struct when received from app", func(t *testing.T) {
		id := primitive.NewObjectID().Hex()
		mockStruct := &app.GetUserByIdOutput{
			UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
				DatabaseIdentifier: &database.DatabaseIdentifier{Id: id},
			},
		}

		deps.mockGetUserByIdImpl.EXPECT().
			Do(gomock.Any(), id).
			Times(1).
			Return(mockStruct, nil)

		fbr := fiber.New()
		fbr.Get("/api/v1/users/:id", deps.userController.GetUserById)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", id), nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse app.GetUserByIdOutput
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, id, httpResponse.Id, "should return expected response")
	})
}

func TestUserController_DeleteUserById(t *testing.T) {
	deps := BeforeEach_TestUserController(t)
	defer deps.ctrl.Finish()

	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		id := primitive.NewObjectID().Hex()
		deps.mockDeleteUserByIdImpl.
			EXPECT().
			Do(gomock.Any(), id).
			Times(1).
			Return(errors.New(exception.CodeDatabaseFailed))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Delete("/api/v1/users/:id", deps.userController.DeleteUserById)
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", id), nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 500, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Failed to communicate with database", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should not return error when error returned from the app is nil", func(t *testing.T) {
		id := primitive.NewObjectID().Hex()

		deps.mockDeleteUserByIdImpl.EXPECT().
			Do(gomock.Any(), id).
			Times(1).
			Return(nil)

		fbr := fiber.New()
		fbr.Delete("/api/v1/users/:id", deps.userController.DeleteUserById)
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", id), nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Equal(t, len(bytes), 0, "should return NO_CONTENT")
	})
}

func TestUserController_CreateUser(t *testing.T) {
	deps := BeforeEach_TestUserController(t)
	defer deps.ctrl.Finish()

	t.Run("should mount the http exception when there is an error in BodyParser", func(t *testing.T) {
		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Post("/api/v1/users/", deps.userController.CreateUser)
		req := httptest.NewRequest("POST", "/api/v1/users/", nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 400, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Invalid input for one or more required attributes", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should mount the http exception when there is an error in ValidationRequest", func(t *testing.T) {
		payload := &app.CreateUserInput{
			FirstName: "username",
			LastName:  "userlastname",
		}

		body, _ := json.Marshal(payload)
		reader := strings.NewReader(string(body))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Post("/api/v1/users/", deps.userController.CreateUser)
		req := httptest.NewRequest("POST", "/api/v1/users/", io.Reader(reader))
		req.Header.Set("Content-Type", "application/json")
		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 400, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Invalid input for one or more required attributes", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		payload := &app.CreateUserInput{
			FirstName: "username",
			LastName:  "userlastname",
			Email:     "foobar@domain.com",
			Type:      "admin",
			Password:  "something",
		}
		body, _ := json.Marshal(payload)
		reader := strings.NewReader(string(body))

		deps.mockCreateUserImpl.
			EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, errors.New(exception.CodeDatabaseFailed))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Post("/api/v1/users/", deps.userController.CreateUser)
		req := httptest.NewRequest("POST", "/api/v1/users/", io.Reader(reader))
		req.Header.Set("Content-Type", "application/json")
		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 500, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Failed to communicate with database", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should return empty error when successfully executed on ValidationRequest and createUser", func(t *testing.T) {
		payload := &app.CreateUserInput{
			FirstName: "username",
			LastName:  "userlastname",
			Email:     "foobar@domain.com",
			Type:      "admin",
			Password:  "something",
		}
		body, _ := json.Marshal(payload)
		reader := strings.NewReader(string(body))

		mockStruct := &app.CreateUserOutput{
			DatabaseIdentifier: &database.DatabaseIdentifier{Id: "123"},
		}

		deps.mockCreateUserImpl.EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockStruct, nil)

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Post("/api/v1/users/", deps.userController.CreateUser)
		req := httptest.NewRequest("POST", "/api/v1/users/", io.Reader(reader))
		req.Header.Set("Content-Type", "application/json")
		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse app.CreateUserOutput
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, "123", httpResponse.Id, "should return expected response")
	})
}

func TestUserController_GetUserPaginated(t *testing.T) {
	deps := BeforeEach_TestUserController(t)
	defer deps.ctrl.Finish()

	const getUserPaginatedEndpoint = "/api/v1/users"

	t.Run("should mount the http exception when there is an error parsing query params", func(t *testing.T) {
		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Get(getUserPaginatedEndpoint, deps.userController.GetUserPaginated)

		req := httptest.NewRequest("GET", "/api/v1/users?page=ABCD", nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {

			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 400, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Invalid input for one or more required attributes", httpResponse.ErrorMessage, "should return expected error message")

	})

	t.Run("should mount the http exception when there is an error validating query params", func(t *testing.T) {
		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Get(getUserPaginatedEndpoint, deps.userController.GetUserPaginated)
		req := httptest.NewRequest("GET", "/api/v1/users", nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 400, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Invalid input for one or more required attributes", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		deps.mockGetUserPaginatedImpl.
			EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Times(1).
			Return(nil, errors.New(exception.CodeDatabaseFailed))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Get(getUserPaginatedEndpoint, deps.userController.GetUserPaginated)

		req := httptest.NewRequest("GET", "/api/v1/users?page=1&per_page=10", nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 500, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Failed to communicate with database", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should return empty error when successfully executed on getUserPaginatedImpl", func(t *testing.T) {
		mockStruct := &database.PaginatedSlice[app.GetUserPaginatedOutput]{
			Items: &[]app.GetUserPaginatedOutput{{
				UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
					DatabaseIdentifier: &database.DatabaseIdentifier{Id: "123"},
				},
			}},
			Page:    1,
			PerPage: 10,
		}

		deps.mockGetUserPaginatedImpl.
			EXPECT().
			Do(gomock.Any(), gomock.Any()).
			Times(1).
			Return(mockStruct, nil)

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Get(getUserPaginatedEndpoint, deps.userController.GetUserPaginated)

		req := httptest.NewRequest("GET", "/api/v1/users?page=1&per_page=10", nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		httpResponse := database.PaginatedSlice[app.GetUserPaginatedOutput]{}
		json.Unmarshal(bytes, &httpResponse)

		items := *httpResponse.Items

		assert.Equal(t, 1, httpResponse.Page, "should return expected response")
		assert.Equal(t, 10, httpResponse.PerPage, "should return expected response")
		assert.Equal(t, "123", items[0].Id, "should return expected response")
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	deps := BeforeEach_TestUserController(t)
	defer deps.ctrl.Finish()

	t.Run("should mount the http exception when there is an error in BodyParser", func(t *testing.T) {
		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Patch("/api/v1/users/", deps.userController.UpdateUserById)
		req := httptest.NewRequest("PATCH", "/api/v1/users/", nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 400, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Invalid input for one or more required attributes", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should must raise the http exception when receiving an error from the method do", func(t *testing.T) {
		payload := &app.UpdateUserByIdInput{
			FirstName: "username",
			LastName:  "userlastname",
			Email:     "foobar@domain.com",
			Type:      "admin",
			Password:  "something",
			UpdatedAt: time.Now(),
		}

		body, _ := json.Marshal(payload)
		reader := strings.NewReader(string(body))

		deps.mockUpdateUserByIdImpl.
			EXPECT().
			Do(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, errors.New(exception.CodeDatabaseFailed))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Patch("/api/v1/users/", deps.userController.UpdateUserById)
		req := httptest.NewRequest("PATCH", "/api/v1/users/", io.Reader(reader))
		req.Header.Set("Content-Type", "application/json")
		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		var httpResponse exception.HTTPException
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, 500, httpResponse.StatusCode, "should return expected status code")
		assert.Equal(t, "Failed to communicate with database", httpResponse.ErrorMessage, "should return expected error message")
	})

	t.Run("should return empty error when successfully executed on ValidationRequest and updateUser", func(t *testing.T) {
		id := primitive.NewObjectID().Hex()

		payload := &app.UpdateUserByIdInput{
			FirstName: "username",
			LastName:  "userlastname",
			Email:     "foobar@domain.com",
			Type:      "admin",
			Password:  "something",
			UpdatedAt: time.Now(),
		}

		mockStruct := &app.UpdateUserByIdOutput{
			UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
				DatabaseIdentifier: &database.DatabaseIdentifier{Id: id},
				User: &domain.User{FirstName: "username",
					LastName: "userlastname",
					Email:    "foobar@domain.com",
					Type:     "admin"},
				DatabaseTimestamp: &database.DatabaseTimestamp{CreatedAt: time.Now()},
			},
		}

		body, _ := json.Marshal(payload)
		reader := strings.NewReader(string(body))

		deps.mockUpdateUserByIdImpl.
			EXPECT().
			Do(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(mockStruct, nil)

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Patch("/api/v1/users/:id", deps.userController.UpdateUserById)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/users/%s", id), io.Reader(reader))
		req.Header.Set("Content-Type", "application/json")
		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}
		var httpResponse app.UpdateUserByIdOutput
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, id, httpResponse.Id, "should return expected response")
	})
}
