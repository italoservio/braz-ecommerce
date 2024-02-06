package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/services/users/app"
)

type UserControllerImpl struct {
	getUserByIdImpl    app.GetUserByIdInterface
	deleteUserByIdImpl app.DeleteUserByIdInterface
}

func NewUserControllerImpl(
	getUserByIdImpl app.GetUserByIdInterface,
	deleteUserByIdImpl app.DeleteUserByIdInterface,
) *UserControllerImpl {
	return &UserControllerImpl{
		getUserByIdImpl:    getUserByIdImpl,
		deleteUserByIdImpl: deleteUserByIdImpl,
	}
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
