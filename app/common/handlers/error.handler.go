package handlers

import (
	"go-fiber-api-template/app/common/responses"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var customCode int

	if mainResp, ok := err.(*responses.MainResponse); ok {
		customCode = mainResp.Status
	} else {
		customCode = http.StatusNotFound

		if strings.Contains(err.Error(), "Cannot") == false {
			log.Printf("Error: %v", err)
		}
	}

	return ctx.Status(customCode).JSON(&responses.MainResponse{
		Status: customCode,
	})
}
