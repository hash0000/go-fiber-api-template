package tale

import (
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/common/types"
	"go-fiber-api-template/app/common/types/jet/neurotales/public/model"
	"go-fiber-api-template/app/modules/tale/schema"
	"go-fiber-api-template/app/modules/user"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func orderPdfTaleService(schema schema.OrderTaleSchema, trialType int8) *responses.MainResponse {
	var err error

	if trialType == constants.TrialTale.No || trialType == constants.TrialTale.Finish {
		_, err = user.TokenNumberUpdate(nil, nil, schema.ID, 1, constants.SqlMathOperationSubtract)
		if err != nil {
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		}
	} else {
		_, err = user.UpdateUseTrial(nil, nil, schema.ID, true)
		if err != nil {
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		}
	}

	ctx := context.Background()
	connection, err := database.GetConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	tx, err := connection.Begin()
	if err != nil {
		slog.Error("Error starting transaction", slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	var chatGptPrompt string

	// TODO: костыль с опциональной ссылкой
	var jsonFileWithImgPromptFilePath string
	if schema.Url != "" {
		jsonFileWithImgPromptFilePath, err = getFileFromUrl(schema.Url)
		if err != nil {
			tx.Rollback()
			slog.Error("Error running func", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		}
	}

	var jsonPromptFromFabula types.ImgToTextJsonType
	if jsonFileWithImgPromptFilePath != "" {
		jsonPromptFromFabula, err = convertJsonToImgToTextJsonType(jsonFileWithImgPromptFilePath)
		if err != nil {
			tx.Rollback()
			slog.Error("Error running func", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		} else if jsonPromptFromFabula.Description == "" {
			tx.Rollback()
			slog.Error("Error reading image prompt from Fabula. Empty string", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		}
	}

	chatGptPrompt = fmt.Sprintf(constants.GptPrompt, schema.ChildData, jsonPromptFromFabula.Description, schema.BackgroundCharacters, schema.Preferences, schema.Moral)

	//jsonFileWithImgPrompt := filepath.Join(constants.TempDir, "/test_tale.json")
	//tale, err := loadTaleFromFile(jsonFileWithImgPrompt)

	tale, err := getTaleFromChatGpt(chatGptPrompt)
	if err != nil {
		tx.Rollback()
		slog.Error("Error getting tale", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	sessionID := uuid.New().String()
	slog.Debug("Session created", slog.String("sessionID", sessionID))

	storeDir := filepath.Join(constants.StoreDir, sessionID)
	err = os.MkdirAll(storeDir, os.ModePerm)
	if err != nil {
		tx.Rollback()
		slog.Error("Error creating store directory", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	if trialType == constants.TrialTale.Start {
		taleJsonFileName := fmt.Sprintf("tale_%s.json", sessionID)
		err = saveToFile(tale, storeDir, taleJsonFileName)
		if err != nil {
			tx.Rollback()
			slog.Error("Error running func", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		}
	}

	talePath, err := generatePdf(tale, storeDir, sessionID, trialType)
	if err != nil {
		tx.Rollback()

		err = os.RemoveAll(storeDir)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to delete directory %s", storeDir), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		}

		slog.Error("Error running file", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))

		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	var isPayed bool
	if trialType == constants.TrialTale.No {
		isPayed = true
	} else {
		isPayed = false
	}

	openAiAnswerJson, err := json.Marshal(tale)
	if err != nil {
		tx.Rollback()
		slog.Error("Error serializing struct", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	fabulaImgToTextJSON, err := json.Marshal(jsonPromptFromFabula)
	if err != nil {
		tx.Rollback()
		slog.Error("Error serializing struct", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}
	fabulaImgToTextJSONString := string(fabulaImgToTextJSON)

	_, err = InsertBook(ctx, tx, model.Tale{UserID: schema.ID, Name: tale.Title, FileName: talePath, TaleGenerationID: sessionID, IsPayed: isPayed, ChildData: schema.ChildData, BackgroundCharacters: schema.BackgroundCharacters, Preferences: schema.Preferences, Moral: schema.Moral, OpenAiAnswer: string(openAiAnswerJson), FabulaImgToTextJSON: &fabulaImgToTextJSONString})
	if err != nil {
		tx.Rollback()
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	tx.Commit()

	fileDownloadUrl := constants.Url.ApiUrl + talePath

	slog.Debug("Sending tale to Telegram bot")
	sendTaleToTelegramBot(schema.ID, fileDownloadUrl)
	slog.Debug("Successfully sent tale to Telegram bot")

	if os.Getenv("RUN_MODE") == "prod" {
		_ = sendTaleToTelegramAnalytic(fileDownloadUrl)
	}

	return &responses.MainResponse{Status: http.StatusOK}
}

func orderFinishPdfTaleService(schema schema.FinishTrialTaleSchema, taleEntity TaleEntityWithUserId, trialTale int8) *responses.MainResponse {
	tale, err := loadTaleFromFile(taleEntity.GenerationID)
	if err != nil {
		slog.Error("Error getting tale", slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	_, err = user.TokenNumberUpdate(nil, nil, schema.UserId, 1, constants.SqlMathOperationSubtract)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	slog.Debug("Session created", slog.String("sessionID", taleEntity.GenerationID))

	storeDir := filepath.Join(constants.StoreDir, taleEntity.GenerationID)

	fileName, err := generatePdf(tale, storeDir, taleEntity.GenerationID, trialTale)
	if err != nil {
		slog.Error("Error generating tale", slog.Any("info", err))
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	_, err = updateIsPayedAndFileName(taleEntity.TaleID, fileName)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	err = os.Remove(filepath.Join(constants.StoreDir, schema.TaleGenerationId, fmt.Sprintf("trial_tale_%s.pdf", schema.TaleGenerationId)))
	if err != nil {
		slog.Error("Failed to delete file", slog.String("file_name", fileName), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
	}
	err = os.Remove(filepath.Join(constants.StoreDir, schema.TaleGenerationId, fmt.Sprintf("tale_%s.json", schema.TaleGenerationId)))
	if err != nil {
		slog.Error("Failed to delete file", slog.String("file_name", fileName), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
	}

	fileDownloadUrl := constants.Url.ApiUrl + fileName

	slog.Debug("Sending tale to Telegram bot")
	sendTaleToTelegramBot(schema.UserId, fileDownloadUrl)
	slog.Debug("Successfully sent tale to Telegram bot")

	if os.Getenv("RUN_MODE") == "prod" {
		_ = sendTaleToTelegramAnalytic(fileDownloadUrl)
	}

	return &responses.MainResponse{Status: http.StatusOK}
}

func getListWithTaleS(schema schema.GetTalesListSchema, page int16) *responses.MainResponse {
	connection, err := database.GetConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	var orderByPrepared string
	if schema.OrderBy == "created_at" {
		orderByPrepared = "t.created_at "
	}
	if schema.SortBy == "DESC" {
		orderByPrepared = orderByPrepared + "DESC"
	} else if schema.SortBy == "ASC" {
		orderByPrepared = orderByPrepared + "ASC"

	}

	repResult, dbStatus, err := selectList(connection, schema.DateFrom, schema.DateTo, schema.Limit, (page-1)*10, orderByPrepared)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}
	if dbStatus == constants.NotFound {
		return &responses.MainResponse{Status: http.StatusNotFound}
	}

	totalPages := int(0)
	if schema.Limit > 0 {
		totalPages = helpers.CalculateTotalPagesPaginationHelper(repResult.Count, schema.Limit)
	}

	result := map[string]interface{}{
		"count": repResult.Count,
		"data":  repResult.Data,
		"pages": totalPages,
	}

	return &responses.MainResponse{Status: http.StatusOK, Data: result}
}
