package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/common/types"
	"go-fiber-api-template/app/common/types/entities"
	"go-fiber-api-template/app/modules/user/schema"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

func insert(schema entities.UserUpdateableRowType) (constants.SqlQueryStatusType, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return constants.Unknown, err
	}

	_, err = User.INSERT(User.ID).MODEL(schema).Exec(connection)

	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		return status, dbError
	}

	return status, nil
}

func TokenNumberUpdate(ctx context.Context, tx *sql.Tx, ID int64, tokenNumber int16, operation string) (constants.SqlQueryStatusType, error) {
	var updateExpr postgres.IntegerExpression

	if operation == constants.SqlMathOperationAdd {
		updateExpr = User.TokenNumber.ADD(postgres.Int16(tokenNumber))
	} else if operation == constants.SqlMathOperationSubtract {
		updateExpr = User.TokenNumber.SUB(postgres.Int16(tokenNumber))
	} else {
		panic(&responses.MainResponse{Status: http.StatusInternalServerError})
	}

	var err error
	connection, err := database.GetConnection()
	if err != nil {
		return constants.Unknown, err
	}

	if tx != nil {
		_, err = User.UPDATE().SET(
			User.TokenNumber.SET(updateExpr),
			User.IsPayedTale.SET(postgres.Bool(true)),
		).WHERE(User.ID.EQ(postgres.Int64(ID))).ExecContext(ctx, tx)
	} else {
		_, err = User.UPDATE().SET(
			User.TokenNumber.SET(updateExpr),
			User.IsPayedTale.SET(postgres.Bool(true)),
		).WHERE(User.ID.EQ(postgres.Int64(ID))).Exec(connection)
	}

	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		return status, dbError
	}

	return status, nil
}

func SelectOne(ctx context.Context, tx *sql.Tx, id int64) (entities.User, constants.SqlQueryStatusType, error) {
	dest := entities.User{}

	var err error
	connection, err := database.GetConnection()
	if err != nil {
		return dest, constants.Unknown, err
	}

	stmt := postgres.SELECT(
		User.AllColumns,
	).FROM(
		User,
	).WHERE(
		User.ID.EQ(postgres.Int64(id)),
	).LIMIT(1)

	if tx != nil {
		err = stmt.QueryContext(ctx, tx, &dest)
	} else {
		err = stmt.Query(connection, &dest)
	}

	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		return dest, status, dbError
	}

	return dest, status, nil
}

func SelectOneByInviteCode(inviteCode types.Uuid) (entities.User, constants.SqlQueryStatusType, error) {
	var dest entities.User

	connection, err := database.GetConnection()
	if err != nil {
		return dest, constants.Unknown, err
	}

	stmt := postgres.SELECT(
		User.AllColumns,
	).FROM(
		User,
	).WHERE(
		User.InviteCode.EQ(postgres.UUID(inviteCode)),
	).LIMIT(1)

	err = stmt.Query(connection, &dest)

	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		return dest, status, dbError
	}

	return dest, status, dbError
}

func UpdateUseTrial(ctx context.Context, tx *sql.Tx, userID int64, trial bool) (constants.SqlQueryStatusType, error) {
	query := `
		UPDATE "user"
		SET use_trial = $1
		WHERE id = $2
	`
	var err error

	connection, err := database.GetConnection()
	if err != nil {
		return constants.Unknown, err
	}

	if tx != nil {
		_, err = tx.ExecContext(ctx, query, trial, userID)
	} else {
		_, err = connection.Exec(query, trial, userID)
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

type book struct {
	ID               *int64     `json:"id"`
	Name             *string    `json:"name"`
	FileName         string     `json:"file_name"`
	TaleGenerationID string     `json:"tale_generation_id"`
	IsPayed          *bool      `json:"is_payed"`
	CreatedAt        *time.Time `json:"created_at"`
}

type UserWithBooks struct {
	ID          *int64     `json:"id,omitempty" sql:"primary_key"`
	TokenNumber *int16     `json:"token_number,omitempty"`
	UseTrial    *bool      `json:"use_trial,omitempty"`
	InviteCode  *uuid.UUID `json:"invite_code,omitempty"`
	CreatedAt   *time.Time `json:"created_at"`
	Books       *[]book    `json:"books"`
}

func selectOneWithTales(connection *sql.DB, id int64) (entities.UserWithTalesType, constants.SqlQueryStatusType, error) {
	var userWithTales entities.UserWithTalesType
	var sqlResult []byte

	err := connection.QueryRow(selectOneWithTalesQuery, id).Scan(&sqlResult)
	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		return userWithTales, status, dbError
	}

	if status != constants.Success {
		return userWithTales, status, nil
	}

	err = json.Unmarshal(sqlResult, &userWithTales)
	if err != nil {
		slog.Error("Error unmarshalling", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return userWithTales, constants.Unknown, nil
	}

	return userWithTales, constants.Success, nil
}

func selectListWithTales(connection *sql.DB, dateFrom, dateTo string, limit, page int16, orderByPrepared string) (entities.UserWithTalesPaginationType, constants.SqlQueryStatusType, error) {
	var usersList entities.UserWithTalesPaginationType

	rows, err := connection.Query(fmt.Sprintf(selectUsersWithTalesQuery, orderByPrepared, orderByPrepared), dateFrom, dateTo, limit, page)
	status, dbError := helpers.CheckSqlError(err)
	if dbError != nil {
		slog.Error("Error getting SQL connection", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return usersList, status, err
	}
	defer rows.Close()

	var currentUserID int64
	var currentUser *entities.UserWithTalesType
	users := make([]entities.UserWithTalesType, 0)

	for rows.Next() {
		var userID int64
		var taleID sql.NullInt64

		var tokenNumber sql.NullInt16
		var useTrial sql.NullBool
		var inviteCode sql.NullString
		var isPayedTale sql.NullBool
		var firstName, lastName, telegramUsername sql.NullString
		var userCreatedAt sql.NullTime

		var taleName, taleFileName, taleGenerationID, taleChildData, taleBackgroundCharacter, talePreferences, taleMoral sql.NullString
		var taleIsPayed sql.NullBool
		var taleOpenAiAnswer sql.NullString
		var taleFabulaImgToTextJson sql.NullString
		var taleCreatedAt sql.NullTime

		if err := rows.Scan(
			&userID,
			&tokenNumber,
			&useTrial,
			&inviteCode,
			&isPayedTale,
			&firstName,
			&lastName,
			&telegramUsername,
			&userCreatedAt,
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
			&taleCreatedAt,
		); err != nil {
			slog.Error("Error scanning", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
			return usersList, constants.Error, err
		}

		if currentUser == nil || currentUserID != userID {
			if currentUser != nil {
				users = append(users, *currentUser)
			}
			currentUserID = userID
			currentUser = &entities.UserWithTalesType{
				ID:               &userID,
				TokenNumber:      helpers.Int16Ptr(tokenNumber),
				UseTrial:         helpers.BoolPtr(useTrial),
				InviteCode:       helpers.UuidPtr(inviteCode),
				IsPayedTale:      helpers.BoolPtr(isPayedTale),
				FirstName:        helpers.StringPtr(firstName),
				LastName:         helpers.StringPtr(lastName),
				TelegramUsername: helpers.StringPtr(telegramUsername),
				CreatedAt:        helpers.TimePtr(userCreatedAt),
				Tales:            []entities.TaleType{},
			}
		}

		if taleID.Valid {
			tale := entities.TaleType{
				ID:                  helpers.Int64Ptr(taleID),
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
				CreatedAt:           helpers.TimePtr(taleCreatedAt),
			}
			currentUser.Tales = append(currentUser.Tales, tale)
		}
	}

	if currentUser != nil {
		users = append(users, *currentUser)
	}

	queryForCount := `
		SELECT COUNT(*) AS count
		FROM "user"
		WHERE "user".created_at BETWEEN $1 AND $2
	`

	var countResult types.SqlCountType
	err = connection.QueryRow(queryForCount, dateFrom, dateTo).Scan(&countResult.Count)
	status, dbError = helpers.CheckSqlError(err)
	if dbError != nil {
		slog.Error("Error getting SQL connection", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return usersList, status, err
	}

	result := entities.UserWithTalesPaginationType{
		Data:  users,
		Count: countResult.Count,
	}

	return result, constants.Success, nil
}

func updateUserDataR(connection *sql.DB, ctx context.Context, tx *sql.Tx, schema schema.UpdateUserSchema) (constants.SqlQueryStatusType, error) {
	q := `UPDATE "user" SET `
	qParts := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)
	argIndex := 1

	if schema.FirstNamePresent == true {
		if schema.FirstName != nil {
			qParts = append(qParts, fmt.Sprintf(`first_name = $%d`, argIndex))
			args = append(args, schema.FirstName)
			argIndex++
		} else {
			qParts = append(qParts, fmt.Sprintf(`first_name = NULL`))
		}
	}
	if schema.LastNamePresent == true {
		if schema.LastName != nil {
			qParts = append(qParts, fmt.Sprintf(`last_name = $%d`, argIndex))
			args = append(args, *schema.LastName)
			argIndex++
		} else {
			qParts = append(qParts, fmt.Sprintf(`last_name = NULL`))
		}
	}
	if schema.TelegramUsernamePresent == true {
		if schema.LastName != nil {
			qParts = append(qParts, fmt.Sprintf(`telegram_username = $%d`, argIndex))
			args = append(args, *schema.TelegramUsername)
			argIndex++
		} else {
			qParts = append(qParts, fmt.Sprintf(`telegram_username = NULL`))
		}
	}

	if len(qParts) == 0 {
		return constants.NoDataToUpdate, nil
	}

	q += strings.Join(qParts, ", ") + fmt.Sprintf(` WHERE id = $%d`, argIndex)
	args = append(args, schema.UserId)

	_, err := connection.Exec(q, args...)
	if err != nil {
		slog.Error("Error while querying SQL", slog.String("query", q), slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return constants.Error, err
	}

	return constants.Success, nil
}
