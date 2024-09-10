package middlewares

import (
	"go-fiber-api-template/app/common/responses"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func CheckAPIToken(c *fiber.Ctx) error {
	headerToken := c.Get("api-token")
	envToken := os.Getenv("API_TOKEN")

	if headerToken != envToken {
		return &responses.MainResponse{Status: http.StatusForbidden, Message: "Api token not found"}
	}

	return c.Next()
}
