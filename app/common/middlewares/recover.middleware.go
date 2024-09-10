package middlewares

import (
	"go-fiber-api-template/app/common/responses"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func RecoverMiddleware(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic occurred: %v\n", r)
			log.Println(string(debug.Stack()))

			c.Status(http.StatusInternalServerError).JSON(&responses.MainResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
	}()

	return c.Next()
}
