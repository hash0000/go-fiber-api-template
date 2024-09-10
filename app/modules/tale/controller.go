package tale

import (
	"encoding/json"
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/regex"
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/common/schemas"
	"go-fiber-api-template/app/modules/tale/schema"
	"go-fiber-api-template/app/modules/user"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func getGenerationStats(ctx *fiber.Ctx) error {
	params := ctx.Locals("query_string").(schema.GenerationPeriodStatsSchema)

	connection, err := database.GetConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	result, sqlStatus, err := SelectGenerationPeriodStats(connection, params.DateFrom, params.DateTo)
	if err != nil || sqlStatus != constants.Success {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"count": result})
}

func getGenerationUserStats(ctx *fiber.Ctx) error {
	params := ctx.Locals("query_string").(schema.GenerationPeriodStatsSchema)

	connection, err := database.GetConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	result, sqlStatus, err := SelectGenerationPeriodUser(connection, params.DateFrom, params.DateTo)
	if err != nil || sqlStatus != constants.Success {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"count": result})
}

func downloadFile(ctx *fiber.Ctx) error {
	params := ctx.Locals("params").(schema.GetFileSchema)
	query := ctx.Locals("query_string").(schema.GenerationIdSchema)

	var filePath string

	if query.GenerationId != nil {
		uuid := regex.Uuid.FindString(*query.GenerationId)

		filePath = filepath.Join(constants.StoreDir, uuid, params.FileName)
	} else {
		uuid := regex.Uuid.FindString(params.FileName)

		filePath = filepath.Join(constants.StoreDir, uuid, params.FileName)
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"status": http.StatusNotFound})
	}

	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", params.FileName))

	return ctx.SendFile(filePath)
}

func orderTale(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(schema.OrderTaleSchema)

	status, err := checkIfUserExistsAndHasToken(body.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": http.StatusInternalServerError})
	}
	if status == 1 {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"status": http.StatusNotFound})
	} else if status == 2 {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"status": http.StatusForbidden})
	}

	go func() {
		res := orderPdfTaleService(body, constants.TrialTale.No)
		if res.Status != http.StatusOK {
			if res.Status == http.StatusInternalServerError {
				_, err = user.TokenNumberUpdate(nil, nil, body.ID, 1, constants.SqlMathOperationAdd)
				if err != nil {
					slog.Error(fmt.Sprintf("CRITICAL. Can not return token after error. UserID: %d", body.ID), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
				}
			}

			slog.Error("Error while generating tale. Get status 500", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))

			requestBody, err := json.Marshal(map[string]string{
				"id":  strconv.Itoa(int(body.ID)),
				"url": constants.ErrorGeneratingTaleGeneralRU,
			})
			if err != nil {
				slog.Error("Error marshalling JSON", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			}

			client := resty.New()
			_, err = client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(requestBody).
				Post(constants.Url.UploadErrorTg)

			if err != nil {
				slog.Error(fmt.Sprintf("Error while try send request to: %s", constants.Url.UploadErrorTg), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			}
		}
	}()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"status": http.StatusOK})
}

func fabulaWebhookImgToTxt(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(schema.WebhookPayloadTextToImageSchema)

	var customData schema.WebhookPayloadTextToImageCustomDataSchema
	err := json.Unmarshal([]byte(body.CustomData), &customData)
	if err != nil {
		slog.Error("Error while unmarshaling Fabula CustomData", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return ctx.SendStatus(http.StatusBadRequest)
	}

	sessionDir := filepath.Join(constants.StoreDir, customData.SessionID)
	isExists, err := helpers.IsExistsByPathHelper(sessionDir)
	if isExists == false {
		return ctx.SendStatus(http.StatusOK)
	} else if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	client := resty.New()
	resp, err := client.R().Get(body.Images[0])
	if err != nil || resp.StatusCode() != http.StatusOK {
		slog.Error(fmt.Sprintf("Error while sending request to: %s", body.Images[0]), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return ctx.SendStatus(http.StatusBadRequest)
	}

	fabulaResponse := resp.Body()
	imagePath := filepath.Join(constants.StoreDir, customData.SessionID, fmt.Sprintf("image_%s.jpg", customData.Index))

	if err := os.WriteFile(imagePath, fabulaResponse, 0644); err != nil {
		slog.Error("Error writing file from Fabula", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return ctx.SendStatus(http.StatusOK)
	}

	return ctx.SendStatus(http.StatusOK)
}

func orderTrialTale(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(schema.OrderTaleSchema)

	status, err := checkIfUserAllowTrial(body.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": http.StatusInternalServerError})
	}
	if status == 1 {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"status": http.StatusNotFound})
	} else if status == 2 {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"status": http.StatusForbidden})
	}

	go func() {
		res := orderPdfTaleService(body, constants.TrialTale.Start)
		if res.Status != http.StatusOK {
			if res.Status == http.StatusInternalServerError {
				_, err = user.UpdateUseTrial(nil, nil, body.ID, false)
				if err != nil {
					slog.Error(fmt.Sprintf("CRITICAL. Can not return token after error. UserID: %d", body.ID), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
				}
			}

			slog.Error("Error while generating tale. Get status 500", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))

			requestBody, err := json.Marshal(map[string]string{
				"id":  strconv.Itoa(int(body.ID)),
				"url": constants.ErrorGeneratingTaleGeneralRU,
			})
			if err != nil {
				slog.Error("Error marshalling JSON", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			}

			client := resty.New()
			_, err = client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(requestBody).
				Post(constants.Url.UploadErrorTg)

			if err != nil {
				slog.Error(fmt.Sprintf("Error while try send request to: %s", constants.Url.UploadErrorTg), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			}
		}
	}()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"status": http.StatusOK})
}

func finishTrialTale(ctx *fiber.Ctx) error {
	body := ctx.Locals("body").(schema.FinishTrialTaleSchema)

	status, err := checkIfUserExistsAndHasToken(body.UserId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": http.StatusInternalServerError})
	}
	if status == 1 {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"status": http.StatusNotFound})
	} else if status == 2 {
		return ctx.Status(http.StatusForbidden).JSON(fiber.Map{"status": http.StatusForbidden})
	}

	taleEntity, statusTale, err := SelectTaleWithUserByGenerationId(body.UserId, body.TaleGenerationId)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	if statusTale == constants.NotFound {
		slog.Debug("Tale not found", slog.Any("session_id", body.TaleGenerationId))
		return &responses.MainResponse{Status: http.StatusNotFound}
	}

	go func() {
		res := orderFinishPdfTaleService(body, taleEntity, constants.TrialTale.Finish)
		if res.Status != http.StatusOK {
			if res.Status == http.StatusNotFound {
				return
			}
			if res.Status == http.StatusInternalServerError {
				_, err = user.TokenNumberUpdate(nil, nil, body.UserId, 1, constants.SqlMathOperationAdd)
				if err != nil {
					slog.Error(fmt.Sprintf("CRITICAL. Can not return token after error. UserID: %d", body.UserId), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
				}
			}

			slog.Error("Error while generating tale. Get status 500", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))

			requestBody, err := json.Marshal(map[string]string{
				"id":  strconv.Itoa(int(body.UserId)),
				"url": constants.ErrorGeneratingTaleGeneralRU,
			})
			if err != nil {
				slog.Error("Error marshalling JSON", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			}

			client := resty.New()
			_, err = client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(requestBody).
				Post(constants.Url.UploadErrorTg)

			if err != nil {
				slog.Error(fmt.Sprintf("Error while try send request to: %s", constants.Url.UploadErrorTg), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			}
		}
	}()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"status": http.StatusOK})
}

func getListWithTales(c *fiber.Ctx) error {
	params := c.Locals("params").(schemas.PaginationSchema)
	queryString := c.Locals("query_string").(schema.GetTalesListSchema)

	serviceResult := getListWithTaleS(queryString, params.Page)

	return c.Status(serviceResult.Status).JSON(serviceResult)
}

func getOne(ctx *fiber.Ctx) error {
	params := ctx.Locals("params").(schema.GetOneTaleSchema)

	connection, err := database.GetConnection()
	if err != nil {
		slog.Error("Error getting DB connection", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	result, sqlStatus, err := selectOne(connection, params.ID)
	if sqlStatus == constants.NotFound {
		slog.Error("Error getting DB result", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusNotFound}
	}

	if err != nil || sqlStatus != constants.Success {
		slog.Error("Error getting DB result", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"result": result})
}
