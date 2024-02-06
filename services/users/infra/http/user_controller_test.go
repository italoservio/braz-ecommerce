package http_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/http"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestUserController_GetUserById(t *testing.T) {
	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGetUserByIdImpl := mocks.NewMockGetUserByIdInterface(ctrl)
		mockDeleteUserByIdImpl := mocks.NewMockDeleteUserByIdInterface(ctrl)
		userController := http.NewUserControllerImpl(mockGetUserByIdImpl, mockDeleteUserByIdImpl)

		id := primitive.NewObjectID().Hex()
		mockGetUserByIdImpl.
			EXPECT().
			Do(id).
			Times(1).
			Return(nil, errors.New(exception.CodeNotFound))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Get("/api/v1/users/:id", userController.GetUserById)
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := primitive.NewObjectID().Hex()
		mockStruct := &app.NewGetUserByIdOutput{
			User:        &domain.User{},
			UserControl: &domain.UserControl{Id: id},
		}

		mockGetUserByIdImpl := mocks.NewMockGetUserByIdInterface(ctrl)
		mockDeleteUserByIdImpl := mocks.NewMockDeleteUserByIdInterface(ctrl)
		userController := http.NewUserControllerImpl(mockGetUserByIdImpl, mockDeleteUserByIdImpl)

		mockGetUserByIdImpl.EXPECT().
			Do(id).
			Times(1).
			Return(mockStruct, nil)

		fbr := fiber.New()
		fbr.Get("/api/v1/users/:id", userController.GetUserById)
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

		var httpResponse app.NewGetUserByIdOutput
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, id, httpResponse.Id, "should return expected response")
	})
}

func TestUserController_DeleteUserById(t *testing.T) {
	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGetUserByIdImpl := mocks.NewMockGetUserByIdInterface(ctrl)
		mockDeleteUserByIdImpl := mocks.NewMockDeleteUserByIdInterface(ctrl)
		userController := http.NewUserControllerImpl(mockGetUserByIdImpl, mockDeleteUserByIdImpl)

		id := primitive.NewObjectID().Hex()
		mockDeleteUserByIdImpl.
			EXPECT().
			Do(id).
			Times(1).
			Return(errors.New(exception.CodeDatabaseFailed))

		fbr := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})
		fbr.Delete("/api/v1/users/:id", userController.DeleteUserById)
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		id := primitive.NewObjectID().Hex()

		mockGetUserByIdImpl := mocks.NewMockGetUserByIdInterface(ctrl)
		mockDeleteUserByIdImpl := mocks.NewMockDeleteUserByIdInterface(ctrl)
		userController := http.NewUserControllerImpl(mockGetUserByIdImpl, mockDeleteUserByIdImpl)

		mockDeleteUserByIdImpl.EXPECT().
			Do(id).
			Times(1).
			Return(nil)

		fbr := fiber.New()
		fbr.Delete("/api/v1/users/:id", userController.DeleteUserById)
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
