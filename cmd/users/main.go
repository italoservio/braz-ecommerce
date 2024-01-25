package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type HelloWorld struct {
	Foo   string `json:"foo"`
	Fizz  string `json:"fizz"`
	Hello string `json:"hello"`
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(HelloWorld{
			Foo:   "Bar",
			Fizz:  "Buzz",
			Hello: "World",
		})
	})

	log.Fatal(app.Listen(":3000"))
}
