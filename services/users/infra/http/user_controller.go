package http

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/validation"
	"github.com/italoservio/braz_ecommerce/services/users/app"
)

type UserControllerImpl struct {
	getUserByIdImpl      app.GetUserByIdInterface
	deleteUserByIdImpl   app.DeleteUserByIdInterface
	createUserImpl       app.CreateUserInterface
	getUserPaginatedImpl app.GetUserPaginatedInterface
	updateUserImpl       app.UpdateUserInterface
}

func NewUserControllerImpl(
	getUserByIdImpl app.GetUserByIdInterface,
	deleteUserByIdImpl app.DeleteUserByIdInterface,
	createUserImpl app.CreateUserInterface,
	getUserPaginatedImpl app.GetUserPaginatedInterface,
	updateUserImpl app.UpdateUserInterface,
) *UserControllerImpl {
	return &UserControllerImpl{
		getUserByIdImpl:      getUserByIdImpl,
		deleteUserByIdImpl:   deleteUserByIdImpl,
		createUserImpl:       createUserImpl,
		updateUserImpl:       updateUserImpl,
		getUserPaginatedImpl: getUserPaginatedImpl,
	}
}

type CreateUserPayload struct {
	FirstName string `json:"first_name" validate:"required,min=5,max=20"`
	LastName  string `json:"last_name" validate:"required,min=5,max=20"`
	Email     string `json:"email" validate:"required,min=5,max=20"`
	Type      string `json:"type" validate:"required,min=5,max=20"`
	Password  string `json:"password" validate:"required,min=5,max=20"`
}

type GetUserPaginatedPayload struct {
	Page    int      `query:"page" validate:"required,number,gt=0"`
	PerPage int      `query:"per_page" validate:"required,number,gt=0,lte=100"`
	Emails  []string `query:"email" validate:"omitempty,dive,email"`
	Ids     []string `query:"id" validate:"omitempty,dive,mongodb"`
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

func (uc *UserControllerImpl) UpdateUser(c *fiber.Ctx) error {

	body := &app.UpdateUserInput{}

	if err := c.BodyParser(&body); err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	id := c.Params("id")

	updateUser, err := uc.updateUserImpl.Do(&app.UpdateUserInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Type:      body.Type,
		Password:  body.Password,
		UpdatedAt: time.Now(),
	}, id)

	if err != nil {
		return err
	}

	return c.JSON(updateUser)
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

func (uc *UserControllerImpl) GetUserPaginated(c *fiber.Ctx) error {
	queryParams := GetUserPaginatedPayload{}

	fmt.Println(c.Params("page"))

	err := c.QueryParser(&queryParams)
	if err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	fmt.Println(queryParams)

	if err := validation.ValidateRequest(c, queryParams); err != nil {
		slog.Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	output, err := uc.getUserPaginatedImpl.Do(&app.GetUserPaginatedInput{
		Page:    queryParams.Page,
		PerPage: queryParams.PerPage,
		Emails:  queryParams.Emails,
		Ids:     queryParams.Ids,
	})
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(output)
}
