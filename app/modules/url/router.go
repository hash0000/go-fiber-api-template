package url

import (
	"go-fiber-api-template/app/common/middleware"
	"go-fiber-api-template/app/modules/url/schema"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Get("/", middleware.QueryStringValidationMiddleware[schema.RedirectSchema], func(c *fiber.Ctx) error {
		var request = c.Locals("query").(schema.RedirectSchema)

		result := insertR(request)

		return &responses.MainResponse{Status: http.StatusCreated, Data: result}
	})

	routeGroup := app.Group("/report")
	routeGroup.Post("/", middleware.BearerAuthMiddleware, middleware.BodyValidationMiddleware[dto.UpdateUrlDto], updateUrlHandler)
	routeGroup.Delete("/", middleware.BearerAuthMiddleware, deleteManyHandler)

	routeGroup.Post("/create", middleware.BodyValidationMiddleware[schema.InsertUserSchema], func(ctx *fiber.Ctx) error {
		var request = insert(ctx.Locals("body").(schema.InsertUserSchema))

		return ctx.Status(request.Status).JSON(request)
	})

	routeGroup.Get("/read-one", middleware.QueryStringValidationMiddleware[schema.SelectOneUserSchema], func(ctx *fiber.Ctx) error {
		var request = selectOne(ctx.Locals("query").(schema.SelectOneUserSchema))

		return ctx.Status(request.Status).JSON(request)
	})
}
