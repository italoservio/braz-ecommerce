package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/internal/dtos"
	"github.com/italoservio/braz_ecommerce/internal/services"
)

type Healthcheck struct {
	Server   string `json:"server"`
	Database string `json:"database"`
}

func RoutesUsers(app *fiber.App) {
	users := app.Group("/user")

	users.Get("/create", func(c *fiber.Ctx) error {

		payload := dtos.DTOCreateUserReq{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		services.CreateUser(&payload)

		return nil
	})
}
