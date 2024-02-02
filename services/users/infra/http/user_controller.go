package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/services/users/app"
)

type UserControllerImpl struct {
	getUserByIdImpl app.GetUserByIdInterface
}

func NewUserControllerImpl(
	getUserByIdImpl app.GetUserByIdInterface,
) *UserControllerImpl {
	return &UserControllerImpl{
		getUserByIdImpl: getUserByIdImpl,
	}
}

func (uc *UserControllerImpl) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	response := uc.getUserByIdImpl.Do(id)

	return c.JSON(map[string]string{"id": response})
}
