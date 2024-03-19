package http

import (
	"errors"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"github.com/italoservio/braz_ecommerce/packages/validation"
	"github.com/italoservio/braz_ecommerce/services/users/app"
)

type UserControllerImpl struct {
	logger               logger.LoggerInterface
	getUserByIdImpl      app.GetUserByIdInterface
	deleteUserByIdImpl   app.DeleteUserByIdInterface
	createUserImpl       app.CreateUserInterface
	getUserPaginatedImpl app.GetUserPaginatedInterface
	updateUserByIdImpl   app.UpdateUserByIdInterface
}

func NewUserControllerImpl(
	logger logger.LoggerInterface,
	getUserByIdImpl app.GetUserByIdInterface,
	deleteUserByIdImpl app.DeleteUserByIdInterface,
	createUserImpl app.CreateUserInterface,
	getUserPaginatedImpl app.GetUserPaginatedInterface,
	updateUserByIdImpl app.UpdateUserByIdInterface,
) *UserControllerImpl {
	return &UserControllerImpl{
		logger:               logger,
		getUserByIdImpl:      getUserByIdImpl,
		deleteUserByIdImpl:   deleteUserByIdImpl,
		createUserImpl:       createUserImpl,
		getUserPaginatedImpl: getUserPaginatedImpl,
		updateUserByIdImpl:   updateUserByIdImpl,
	}
}

type CreateUserPayload struct {
	FirstName string `json:"first_name" validate:"required,min=5,max=20"`
	LastName  string `json:"last_name" validate:"required,min=5,max=20"`
	Email     string `json:"email" validate:"required,min=5,max=20"`
	Type      string `json:"type" validate:"required,min=5,max=20"`
	Password  string `json:"password" validate:"required,min=5,max=20"`
}

func (uc *UserControllerImpl) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()
	body := &app.CreateUserInput{}

	if err := c.BodyParser(&body); err != nil {
		uc.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	if err := validation.ValidateRequest(c, body); err != nil {
		uc.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	output, err := uc.createUserImpl.Do(ctx, &app.CreateUserInput{
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

func (uc *UserControllerImpl) UpdateUserById(c *fiber.Ctx) error {
	ctx := c.Context()
	body := &app.UpdateUserByIdInput{}

	if err := c.BodyParser(&body); err != nil {
		uc.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	id := c.Params("id")

	output, err := uc.updateUserByIdImpl.Do(ctx, id, &app.UpdateUserByIdInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Type:      body.Type,
		Password:  body.Password,
	})

	if err != nil {
		return err
	}

	return c.JSON(output)
}

type GetUserByIdPayload struct {
	Deleted bool `query:"deleted"`
}

func (uc *UserControllerImpl) GetUserById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	queryParams := GetUserByIdPayload{}

	err := c.QueryParser(&queryParams)

	if err != nil {
		uc.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	user, err := uc.getUserByIdImpl.Do(ctx, &app.GetUserByIdInput{
		Id:      id,
		Deleted: queryParams.Deleted,
	})
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (uc *UserControllerImpl) DeleteUserById(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	err := uc.deleteUserByIdImpl.Do(ctx, id)

	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusNoContent)
}

type GetUserPaginatedPayload struct {
	Page    int      `query:"page" validate:"required,number,gt=0"`
	PerPage int      `query:"per_page" validate:"required,number,gt=0,lte=100"`
	Emails  []string `query:"email" validate:"omitempty,dive,email"`
	Ids     []string `query:"id" validate:"omitempty,dive,mongodb"`
	Deleted bool     `query:"deleted"`
}

func (uc *UserControllerImpl) GetUserPaginated(c *fiber.Ctx) error {
	ctx := c.Context()
	queryParams := GetUserPaginatedPayload{}

	err := c.QueryParser(&queryParams)
	if err != nil {
		uc.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	if err := validation.ValidateRequest(c, queryParams); err != nil {
		uc.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeValidationFailed)
	}

	output, err := uc.getUserPaginatedImpl.Do(ctx, &app.GetUserPaginatedInput{
		Page:    queryParams.Page,
		PerPage: queryParams.PerPage,
		Emails:  queryParams.Emails,
		Ids:     queryParams.Ids,
		Deleted: queryParams.Deleted,
	})
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(output)
}
