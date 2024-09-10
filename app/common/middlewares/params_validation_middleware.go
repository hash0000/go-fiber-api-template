package middlewares

import (
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/common/types"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateP = validator.New()

func errorParserP(dto any) []types.ValidationErrorType {
	var errors []types.ValidationErrorType
	if err := validateP.Struct(dto); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errorItem types.ValidationErrorType
			errorItem.Property = err.Field()
			errorItem.Message = err.ActualTag()
			errors = append(errors, errorItem)
		}
	}
	return errors
}

func ParamsValidationMiddleware[T any](ctx *fiber.Ctx) error {
	var dto T
	var err error

	err = ctx.ParamsParser(&dto)
	if err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity)
		return ctx.Status(http.StatusUnprocessableEntity).JSON(responses.MainResponse{Status: http.StatusUnprocessableEntity})
	}

	var errors = errorParserP(dto)
	if len(errors) != 0 {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(responses.MainResponse{Status: http.StatusUnprocessableEntity, ErrorInfo: errors})
	}

	ctx.Locals("params", dto)
	return ctx.Next()
}
