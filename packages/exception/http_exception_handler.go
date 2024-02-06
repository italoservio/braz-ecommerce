package exception

import (
	"github.com/gofiber/fiber/v2"
)

func HttpExceptionHandler(c *fiber.Ctx, err error) error {
	httpException := Http(err.Error())
	return c.Status(httpException.StatusCode).JSON(httpException)
}
