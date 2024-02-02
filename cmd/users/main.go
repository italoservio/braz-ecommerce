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

	app.Use(logger.New())

	env := start.NewEnv()

	db, err := database.NewDatabase(env.DB_URI, env.DB_NAME)
	if err != nil {
		log.Fatal(err)
	}

	userController := start.InjectionsContainer(db)

	app.Get("/health", start.HealthCheckEndpoint(db))

	api := app.Group("/api")
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
