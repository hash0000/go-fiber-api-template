package middlewares

import (
	"go-fiber-api-template/app/common/responses"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func RequestValidationMiddleware(c *fiber.Ctx) error {
	if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
		if c.Get("Content-Type") != "application/json" {
			return &responses.MainResponse{Status: http.StatusNotFound}
		}
	}

	return c.Next()
}
