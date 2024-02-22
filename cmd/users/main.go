package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fbrlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/italoservio/braz_ecommerce/cmd/users/start"
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
)

func main() {
	app := fiber.New(fiber.Config{ErrorHandler: exception.HttpExceptionHandler})

	env := start.NewEnv()

	db, err := database.NewDatabase(env.DB_URI, env.DB_NAME)
	if err != nil {
		log.Fatal(err)
	}

	userController := start.InjectionsContainer(db)

	app.Get("/health", start.HealthCheckEndpoint(db))

	api := app.Group("/api")
	api.Use(fbrlogger.New(loggerConfig()))
	api.Use(requestid.New(requestid.Config{
		Header:     string(logger.CorrelationId),
		ContextKey: string(logger.CorrelationId),
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	usersV1 := api.Group("/v1/users")
	usersV1.Post("/", userController.CreateUser)
	usersV1.Get("/", userController.GetUserPaginated)
	usersV1.Get("/:id", userController.GetUserById)
	usersV1.Delete("/:id", userController.DeleteUserById)

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

func loggerConfig() fbrlogger.Config {
	return fbrlogger.Config{
		Format:        "${time} INFO ${locals:X-Correlation-ID} ${method} ${path} ${status} ${latency}\n",
		TimeFormat:    "2006/01/02 15:04:05",
		DisableColors: true,
	}
}
