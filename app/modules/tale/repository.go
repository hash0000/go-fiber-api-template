package tale

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/types"
	"go-fiber-api-template/app/common/types/entities"
	"go-fiber-api-template/app/common/types/jet/neurotales/public/model"
	"log/slog"
)

func InsertBook(ctx context.Context, tx *sql.Tx, values model.Tale) (constants.SqlQueryStatusType, error) {
	query := `
		INSERT INTO tale (user_id, name, file_name, tale_generation_id, is_payed, child_data, background_characters, preferences, moral, open_ai_answer, fabula_img_to_text_json)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	args := []interface{}{values.UserID, values.Name, values.FileName, values.TaleGenerationID, values.IsPayed, values.ChildData, values.BackgroundCharacters, values.Preferences, values.Moral, values.OpenAiAnswer, values.FabulaImgToTextJSON}

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		connection, err := database.GetConnection()
		if err != nil {
			return constants.Unknown, err
		}

		_, err = connection.Exec(query, args...)
	}

	if err != nil {
		status, dbError := helpers.CheckSqlError(err)
		if dbError != nil {
			return status, dbError
		}

		return status, nil
	}

	return constants.Success, nil
}

func SelectGenerationPeriodStats(connection *sql.DB, startDate string, endDate string) (uint32, constants.SqlQueryStatusType, error) {
	args := []interface{}{}
	query := `SELECT COUNT(1) FROM tale `
	if startDate != "" && endDate != "" {
		query += `WHERE created_at > $1 AND created_at <= $2`
		args = append(args, startDate, endDate)
	} else if startDate != "" && endDate == "" {
		query += `WHERE created_at > $1`
		args = append(args, startDate)
	} else if startDate == "" && endDate != "" {
		query += `WHERE created_at <= $1`
		args = append(args, endDate)
	} else {
		query += `WHERE created_at > NOW()::date`
	}

	var count uint32
	err := connection.QueryRow(query, args...).Scan(&count)
	if err != nil {
		status, dbError := helpers.CheckSqlError(err)
		if dbError != nil {
			return 0, status, dbError
		}

		return 0, constants.Unknown, err
	}

	return count, constants.Success, nil
}

func SelectGenerationPeriodUser(connection *sql.DB, startDate string, endDate string) (uint32, constants.SqlQueryStatusType, error) {
	args := []interface{}{}
	query := `SELECT COUNT(DISTINCT user_id) FROM tale `
	if startDate != "" && endDate != "" {
		query += `WHERE created_at > $1 AND created_at <= $2`
		args = append(args, startDate, endDate)
	} else if startDate != "" && endDate == "" {
		query += `WHERE created_at > $1`
		args = append(args, startDate)
	} else if startDate == "" && endDate != "" {
		query += `WHERE created_at <= $1`
		args = append(args, endDate)
	} else {
		query += `WHERE created_at > NOW()::date`
	}

	var count uint32
	err := connection.QueryRow(query, args...).Scan(&count)
	if err != nil {
		status, dbError := helpers.CheckSqlError(err)
		if dbError != nil {
			return 0, status, dbError
		}

		return 0, constants.Unknown, err
	}

	return count, constants.Success, nil
}

type TaleEntityWithUserId struct {
	TaleID       int64  `json:"tale_id,omitempty" sql:"primary_key"`
	GenerationID string `json:"generation_id,omitempty"`
}

func SelectTaleWithUserByGenerationId(userId int64, generationId string) (TaleEntityWithUserId, constants.SqlQueryStatusType, error) {
	var taleID TaleEntityWithUserId

	query := `
	SELECT id as tale_id, tale_generation_id as generation_id
	FROM tale
	WHERE user_id = $1 AND is_payed = false AND tale_generation_id = $2
	LIMIT 1;
	`

	args := []interface{}{userId, generationId}

	connection, err := database.GetConnection()
	if err != nil {
		return TaleEntityWithUserId{}, constants.Unknown, err
	}

	err = connection.QueryRow(query, args...).Scan(&taleID.TaleID, &taleID.GenerationID)
	if err != nil {
		status, dbError := helpers.CheckSqlError(err)
		if dbError != nil {
			return TaleEntityWithUserId{}, status, dbError
		}

		return TaleEntityWithUserId{}, status, nil
	}

	return taleID, constants.Success, nil
}

func updateIsPayedAndFileName(taleID int64, fileName string) (constants.SqlQueryStatusType, error) {
	query := `
		UPDATE tale
		SET is_payed = true, file_name = $1
		WHERE id = $2
	`

	args := []interface{}{fileName, taleID}

	connection, err := database.GetConnection()
	if err != nil {
		return constants.Unknown, err
	}

	_, err = connection.Exec(query, args...)
	if err != nil {
		status, dbError := helpers.CheckSqlError(err)
		if dbError != nil {
			return status, dbError
		}

		return status, nil
	}

	return constants.Success, nil
}

func selectList(connection *sql.DB, dateFrom, dateTo string, limit, page int16, orderByPrepared string) (entities.TalesListWithPaginationType, constants.SqlQueryStatusType, error) {
	var result entities.TalesListWithPaginationType

	rows, err := connection.Query(fmt.Sprintf(selectTalesList, orderByPrepared), dateFrom, dateTo, limit, page)
	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		slog.Error("Error getting SQL connection", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return result, status, err
	}
	defer rows.Close()

	talesList := make([]entities.TaleType, 0)

	for rows.Next() {
		var taleID int64
		var taleName, taleFileName, taleGenerationID, taleChildData, taleBackgroundCharacter, talePreferences, taleMoral sql.NullString
		var taleIsPayed sql.NullBool
		var taleOpenAiAnswer sql.NullString
		var taleFabulaImgToTextJson sql.NullString
		var taleUserID int64
		var taleCreatedAt sql.NullTime

		if err := rows.Scan(
			&taleID,
			&taleName,
			&taleFileName,
			&taleIsPayed,
			&taleGenerationID,
			&taleChildData,
			&taleBackgroundCharacter,
			&talePreferences,
			&taleMoral,
			&taleOpenAiAnswer,
			&taleFabulaImgToTextJson,
			&taleUserID,
			&taleCreatedAt,
		); err != nil {
			slog.Error("Error scanning", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return result, constants.Error, err
		}

		tale := entities.TaleType{
			ID:                  &taleID,
			Name:                helpers.StringPtr(taleName),
			FileName:            helpers.StringPtr(taleFileName),
			IsPayed:             helpers.BoolPtr(taleIsPayed),
			TaleGenerationID:    helpers.StringPtr(taleGenerationID),
			ChildData:           helpers.StringPtr(taleChildData),
			BackgroundCharacter: helpers.StringPtr(taleBackgroundCharacter),
			Preferences:         helpers.StringPtr(talePreferences),
			Moral:               helpers.StringPtr(taleMoral),
			OpenAiAnswer:        helpers.JsonRawMessagePtr(taleOpenAiAnswer),
			FabulaImgToTextJson: helpers.JsonRawMessagePtr(taleFabulaImgToTextJson),
			UserID:              &taleUserID,
			CreatedAt:           helpers.TimePtr(taleCreatedAt),
		}
		talesList = append(talesList, tale)
	}

	queryForCount := `
		SELECT COUNT(*) AS count
		FROM tale
		WHERE tale.created_at BETWEEN $1 AND $2
	`

	var countResult types.SqlCountType
	err = connection.QueryRow(queryForCount, dateFrom, dateTo).Scan(&countResult.Count)
	status, dbError = helpers.CheckSqlError(err)
	if dbError != nil {
		slog.Error("Error getting SQL connection", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return result, status, err
	}

	result = entities.TalesListWithPaginationType{
		Data:  talesList,
		Count: countResult.Count,
	}

	return result, constants.Success, nil
}

func selectOne(connection *sql.DB, ID int64) (entities.TaleType, constants.SqlQueryStatusType, error) {
	var entity entities.TaleType
	var sqlResult []byte

	err := connection.QueryRow(selectOneQ, ID).Scan(&sqlResult)
	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		return entity, status, dbError
	}

	if status != constants.Success {
		return entity, status, nil
	}

	err = json.Unmarshal(sqlResult, &entity)
	if err != nil {
		slog.Error("Error unmarshalling", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return entity, constants.Unknown, nil
	}

	return entity, constants.Success, nil
}
