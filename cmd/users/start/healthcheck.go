package start

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/italoservio/braz_ecommerce/packages/database"
)

type Healthcheck struct {
	Server   string `json:"server"`
	Database string `json:"database"`
}

func HealthCheckEndpoint(db *database.Database) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		healthCheck := Healthcheck{
			Server:   "healthy",
			Database: "healthy",
		}

		err := db.Client().Ping(ctx, nil)

		if err != nil {
			healthCheck.Database = "offline"
		}

		return c.JSON(healthCheck)
	}
}
