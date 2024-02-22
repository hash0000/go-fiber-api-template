package url

import (
	"go-fiber-api-template/app/common/middleware"
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/modules/url/schema"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Get("/", middleware.QueryStringValidationMiddleware[schema.RedirectSchema], func(c *fiber.Ctx) error {
		var request = c.Locals("query").(schema.RedirectSchema)

		var result = selectOneR(request)

		var urlString string

		if request.Gid != "" && request.Range != "" {
			urlString = result.Url + "/edit#gid=" + request.Gid + "&range=" + request.Range
		} else if request.Gid != "" {
			urlString = result.Url + "/edit#gid=" + request.Gid
		} else {
			urlString = result.Url
		}

		return &responses.MainResponse{Status: http.StatusCreated, Data: urlString}
	})

	// routeGroup := app.Group("/report")
	// routeGroup.Post("/", middleware.BearerAuthMiddleware, middleware.BodyValidationMiddleware[dto.UpdateUrlDto], updateUrlHandler)
	// routeGroup.Delete("/", middleware.BearerAuthMiddleware, deleteManyHandler)

	// routeGroup.Post("/create", middleware.BodyValidationMiddleware[schema.InsertUserSchema], func(ctx *fiber.Ctx) error {
	// 	var request = insert(ctx.Locals("body").(schema.InsertUserSchema))

	// 	return ctx.Status(request.Status).JSON(request)
	// })

	// routeGroup.Get("/read-one", middleware.QueryStringValidationMiddleware[schema.SelectOneUserSchema], func(ctx *fiber.Ctx) error {
	// 	var request = selectOne(ctx.Locals("query").(schema.SelectOneUserSchema))

	// 	return ctx.Status(request.Status).JSON(request)
	// })
}
