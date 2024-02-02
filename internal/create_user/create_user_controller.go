package create_user

import "github.com/gofiber/fiber/v2"

type UserController struct{ createUserService *CreateUserService }

func NewUserController(createUserService *CreateUserService) *UserController {
	return &UserController{createUserService}
}

func (cr *UserController) CreateUser(c *fiber.Ctx) error {
	return nil
}
