package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/italoservio/braz_ecommerce/cmd/users/start"
	"github.com/italoservio/braz_ecommerce/packages/database"
)

func main() {
	app := fiber.New()

	env := start.NewEnv()

	db, err := database.NewDatabase(env.DB_URI, env.DB_NAME)
	if err != nil {
		log.Fatal(err)
	}

	userController := start.InjectionsContainer(db)

	app.Get("/health", start.HealthCheckEndpoint(db))

	api := app.Group("/api")
	api.Use(logger.New(loggerConfig()))

	usersV1 := api.Group("/v1/users")
	usersV1.Get("/:id", userController.GetUserById)

	go func() { log.Fatal(app.Listen(env.PORT)) }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	gracefulShutdown(app, db)
}

func gracefulShutdown(app *fiber.App, db *database.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
	db.Client().Disconnect(ctx)
}

func loggerConfig() logger.Config {
	return logger.Config{
		Format:        "${time} INFO ${method} ${path} ${status} ${latency}\n",
		TimeFormat:    "2006/01/02 15:04:05",
		TimeZone:      "America/Sao_Paulo",
		DisableColors: true,
	}
}
