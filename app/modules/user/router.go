package user

import (
	"report-url-redirection/app/common/middleware"
	"report-url-redirection/app/modules/user/dto"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	routeGroup := router.Group("/user")

	routeGroup.Post("/sign-in", middleware.BodyValidationMiddleware[dto.SignInDto], func(ctx *fiber.Ctx) error {
		var data = signIn(ctx.Locals("body").(dto.SignInDto))

		return ctx.Status(data.Status).JSON(data)
	})

	routeGroup.Post("/create", middleware.BearerAuthMiddleware, middleware.BodyValidationMiddleware[dto.CreateDto], func(ctx *fiber.Ctx) error {
		var data = signIn(ctx.Locals("body").(dto.SignInDto))

		return ctx.Status(data.Status).JSON(data)
	})
}
