package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type Healthcheck struct {
	Server   string `json:"server"`
	Database string `json:"database"`
}

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(Healthcheck{
			Server:   "ok",
			Database: "n/a",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
