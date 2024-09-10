package user

import (
	"go-fiber-api-template/app/common/middlewares"
	"go-fiber-api-template/app/common/schemas"
	"go-fiber-api-template/app/modules/user/schema"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Post("", middlewares.CheckAPIToken, middlewares.BodyValidationMiddleware[schema.CreateUserSchema], Create)
	router.Patch("/update-token-number", middlewares.CheckAPIToken, middlewares.BodyValidationMiddleware[schema.UpdateTokenNumberSchema], UpdateTokenNumber)
	router.Patch("", middlewares.CheckAPIToken, middlewares.BodyValidationMiddleware[schema.UpdateUserSchema], updateUserData)
	router.Get("/select-one/:id<int64>", middlewares.CheckAPIToken, middlewares.ParamsValidationMiddleware[schema.IdSchema], GetOne)
	router.Get("/select-one-with-tales/:id<int64>", middlewares.CheckAPIToken, middlewares.ParamsValidationMiddleware[schema.IdSchema], GetOneWithBooks)
	router.Get("/get-list/tales/:page<int16>", middlewares.CheckAPIToken, middlewares.ParamsValidationMiddleware[schemas.PaginationSchema], middlewares.QueryStringValidationMiddleware[schema.GetListWithTaleSchema], getListWithTales)
}
