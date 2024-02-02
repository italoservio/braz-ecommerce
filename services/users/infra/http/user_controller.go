package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/exception"
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

	user, err := uc.getUserByIdImpl.Do(id)
	if err != nil {
		httpException := exception.Http(err.Error())
		return c.Status(httpException.StatusCode).JSON(httpException)
	}

	return c.JSON(user)
}
