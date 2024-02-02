package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"report-url-redirection/app/common/responses"
	"report-url-redirection/app/common/types"
)

var validate = validator.New()

func errorParser(dto any) []types.ValidationErrorType {
	var errors []types.ValidationErrorType
	if err := validate.Struct(dto); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errorItem types.ValidationErrorType
			errorItem.Property = err.Field()
			errorItem.Code = err.ActualTag()
			errors = append(errors, errorItem)
		}
	}
	return errors
}

func BodyValidationMiddleware[T any](ctx *fiber.Ctx) error {
	var dto T
	var err error

	err = ctx.BodyParser(&dto)
	if err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity)
		return ctx.JSON(responses.MainResponse{Status: http.StatusUnprocessableEntity})
	}

	var errors = errorParser(dto)
	if len(errors) != 0 {
		return ctx.JSON(responses.MainResponse{Status: http.StatusUnprocessableEntity, ValidationError: errors})
	}

	ctx.Locals("body", dto)
	return ctx.Next()
}
