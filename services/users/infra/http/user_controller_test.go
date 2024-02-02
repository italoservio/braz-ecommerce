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

func TestUserController_NewUserControllerImpl(t *testing.T) {
	t.Run("should mount http exception when receiving an error from app", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGetUserByIdImpl := mocks.NewMockGetUserByIdInterface(ctrl)
		userController := http.NewUserControllerImpl(mockGetUserByIdImpl)

		id := primitive.NewObjectID().Hex()
		mockGetUserByIdImpl.
			EXPECT().
			Do(id).
			Times(1).
			Return(nil, errors.New(exception.CodeNotFound))

		fbr := fiber.New()
		fbr.Get("/api/v1/users/:id", userController.GetUserById)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", id), nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Fatal(err)
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
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
		userController := http.NewUserControllerImpl(mockGetUserByIdImpl)

		mockGetUserByIdImpl.EXPECT().
			Do(id).
			Times(1).
			Return(mockStruct, nil)

		fbr := fiber.New()
		fbr.Get("/api/v1/users/:id", userController.GetUserById)
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", id), nil)

		response, err := fbr.Test(req, -1)
		if err != nil {
			t.Fatal(err)
		}

		bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var httpResponse app.NewGetUserByIdOutput
		json.Unmarshal(bytes, &httpResponse)

		assert.Equal(t, id, httpResponse.Id, "should return expected response")
	})
}
