package main

import (
	"fmt"
	"go-fiber-api-template/app/common/configs"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/handlers"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/middlewares"
	"go-fiber-api-template/app/modules/tale"
	"go-fiber-api-template/app/modules/user"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})
	appPort := os.Getenv("APP_PORT")

	app.Use(middlewares.RecoverMiddleware)
	app.Use(configs.CorsConfig)
	app.Use(middlewares.RequestValidationMiddleware)
	app.Use(swagger.New(configs.SwaggerConfig))

	if appPort == "" {
		err := godotenv.Load()
		if err != nil {
			slog.Error("Error loading env", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.AppRoot)), slog.Any("info", err))
			panic(err)
		}
		appPort = os.Getenv("APP_PORT")
		if appPort == "" {
			slog.Error("APP_PORT is empty", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.AppRoot)))
			panic(fmt.Errorf("APP_PORT is empty"))
		}
	}

	constants.Load()

	if os.Getenv("LOGGER_LEVEL") == "info" {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	} else if os.Getenv("LOGGER_LEVEL") == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug mode enabled")
	}

	postgresConString, err := helpers.GetPostgresConString()
	if err != nil {
		panic(err)
	}
	err = database.SetPostgresConnection(postgresConString)
	if err != nil {
		panic(err)
	}
	defer database.Close()
	err = database.SetSqlConnection(postgresConString)
	if err != nil {
		panic(err)
	}
	defer database.CloseSqlConnection()

	user.Router(app.Group("/user"))
	tale.Router(app.Group("/tale"))

	slog.Info(fmt.Sprintf("Application started at: %s", time.Now().Format(time.RFC1123)))

	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
		slog.Error("Error starting server", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.AppRoot)), slog.Any("info", err))
		panic("Error starting server")
	}
}
