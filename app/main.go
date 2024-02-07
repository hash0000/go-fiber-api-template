package main

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/modules/user"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			fmt.Printf(err.Error())

			return ctx.Status(code).JSON(&responses.MainResponse{Status: http.StatusInternalServerError})
		},
	})
	app.Use(recover.New())
	database.SetCon(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", 5432, "postgres", "root", "goApi"))
	defer database.GetConnection.Close()

	user.Router(app)

	port := "4444"

	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
