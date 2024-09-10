package tale

import (
	"go-fiber-api-template/app/common/middlewares"
	"go-fiber-api-template/app/common/schemas"
	"go-fiber-api-template/app/modules/tale/schema"

	"github.com/gofiber/fiber/v2"
)

func Router(router fiber.Router) {
	router.Post("/order", middlewares.CheckAPIToken, middlewares.BodyValidationMiddleware[schema.OrderTaleSchema], orderTale)
	router.Post("/order/trial/start", middlewares.CheckAPIToken, middlewares.BodyValidationMiddleware[schema.OrderTaleSchema], orderTrialTale)
	router.Post("/order/trial/finish", middlewares.CheckAPIToken, middlewares.BodyValidationMiddleware[schema.FinishTrialTaleSchema], finishTrialTale)
	router.Get("/download/:file_name", middlewares.ParamsValidationMiddleware[schema.GetFileSchema], middlewares.QueryStringValidationMiddleware[schema.GenerationIdSchema], downloadFile)
	router.Post("/webhook/fabula", middlewares.BodyValidationMiddleware[schema.WebhookPayloadTextToImageSchema], fabulaWebhookImgToTxt)
	router.Get("/count-by-date", middlewares.QueryStringValidationMiddleware[schema.GenerationPeriodStatsSchema], getGenerationStats)
	router.Get("/user/count-by-date", middlewares.QueryStringValidationMiddleware[schema.GenerationPeriodStatsSchema], getGenerationUserStats)
	router.Get("/get-list/:page<int16>", middlewares.CheckAPIToken, middlewares.ParamsValidationMiddleware[schemas.PaginationSchema], middlewares.QueryStringValidationMiddleware[schema.GetTalesListSchema], getListWithTales)
	router.Get("/get-one/:id<int64>", middlewares.CheckAPIToken, middlewares.ParamsValidationMiddleware[schema.GetOneTaleSchema], getOne)
}
