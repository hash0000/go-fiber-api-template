package main

import (
	"errors"
	"fmt"
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

	user.Router(app)

	port := 4444

	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
