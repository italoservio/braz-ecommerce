package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/validation"
	"github.com/italoservio/braz_ecommerce/services/users/app"
)

type CreateUserValidate struct {
	FirstName string `json:"first_name" validate:"required,min=5,max=20"`
	LastName  string `json:"last_name" validate:"required,min=5,max=20"`
	Email     string `json:"email" validate:"required,min=5,max=20"`
	Type      string `json:"type" validate:"required,min=5,max=20"`
	Password  string `json:"password" validate:"required,min=5,max=20"`
}

type UserControllerImpl struct {
	getUserByIdImpl    app.GetUserByIdInterface
	deleteUserByIdImpl app.DeleteUserByIdInterface
	createUserImpl     app.CreateUserInterface
}

func NewUserControllerImpl(
	getUserByIdImpl app.GetUserByIdInterface,
	deleteUserByIdImpl app.DeleteUserByIdInterface,
	createUserImpl app.CreateUserInterface,
) *UserControllerImpl {
	return &UserControllerImpl{
		getUserByIdImpl:    getUserByIdImpl,
		deleteUserByIdImpl: deleteUserByIdImpl,
		createUserImpl:     createUserImpl,
	}
}

func (uc *UserControllerImpl) CreateUser(c *fiber.Ctx) error {
	body := &app.CreateUserInput{}

	if err := c.BodyParser(&body); err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}
	if err := validation.ValidateRequest(c, body); err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	output, err := uc.createUserImpl.Do(&app.CreateUserInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Type:      body.Type,
		Password:  body.Password,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(output)
}

func (uc *UserControllerImpl) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := uc.getUserByIdImpl.Do(id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (uc *UserControllerImpl) DeleteUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	err := uc.deleteUserByIdImpl.Do(id)

	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}
