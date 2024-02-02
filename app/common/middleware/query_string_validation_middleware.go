package middleware

import (
	"net/http"
	"report-url-redirection/app/common/responses"
	"report-url-redirection/app/common/types"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateQ = validator.New()

func errorParserQ(dto any) []types.ValidationErrorType {
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

func QueryStringValidationMiddleware[T any](ctx *fiber.Ctx) error {
	var dto T
	var err error

	err = ctx.QueryParser(&dto)
	if err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity)
		return ctx.JSON(responses.MainResponse{Status: http.StatusUnprocessableEntity})
	}

	var errors = errorParser(dto)
	if len(errors) != 0 {
		return ctx.JSON(responses.MainResponse{Status: http.StatusUnprocessableEntity, ValidationError: errors})
	}

	ctx.Locals("query", dto)
	return ctx.Next()
}
