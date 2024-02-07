package user

import (
	"go-fiber-api-template/app/common/middleware"
	"go-fiber-api-template/app/modules/user/schema"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	routeGroup := router.Group("/user")

	routeGroup.Post("/create", middleware.BodyValidationMiddleware[schema.InsertUserSchema], func(ctx *fiber.Ctx) error {
		var request = insert(ctx.Locals("body").(schema.InsertUserSchema))

		return ctx.Status(request.Status).JSON(request)
	})
}
