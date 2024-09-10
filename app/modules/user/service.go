package user

import (
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/common/types"
	"go-fiber-api-template/app/common/types/entities"
	"go-fiber-api-template/app/modules/user/schema"
	"net/http"
)

func create(body schema.CreateUserSchema) *responses.MainResponse {
	newUser := entities.UserUpdateableRowType{ID: body.UserID}

	status, err := insert(newUser)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	if status == constants.NotUnique {
		return &responses.MainResponse{Status: http.StatusConflict}
	}

	if body.InviteCode != "" {
		userFromCode, status, err := SelectOneByInviteCode(types.Uuid(body.InviteCode))
		if err != nil {
			return &responses.MainResponse{Status: http.StatusInternalServerError}
		}
		if status == constants.Success {
			TokenNumberUpdate(nil, nil, *userFromCode.ID, 1, constants.SqlMathOperationAdd)
		}
	}

	return &responses.MainResponse{Status: http.StatusCreated}
}

func updateTokenNumber(schema schema.UpdateTokenNumberSchema) *responses.MainResponse {
	TokenNumberUpdate(nil, nil, schema.ID, schema.TokenNumber, constants.SqlMathOperationAdd)

	return &responses.MainResponse{Status: http.StatusOK}
}

func getOne(schema schema.IdSchema) *responses.MainResponse {
	result, status, err := SelectOne(nil, nil, schema.ID)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}
	if status == constants.NotFound {
		return &responses.MainResponse{Status: http.StatusNotFound}
	}

	return &responses.MainResponse{Status: http.StatusOK, Data: result}
}

func getOneWithBooks(schema schema.IdSchema) *responses.MainResponse {
	connection, err := database.GetConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	result, status, err := selectOneWithTales(connection, schema.ID)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	if status == constants.NotFound {
		return &responses.MainResponse{Status: http.StatusNotFound}
	}
	if status == constants.Unknown {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	return &responses.MainResponse{Status: http.StatusOK, Data: result}
}

func getListWithTaleS(schema schema.GetListWithTaleSchema, page int16) *responses.MainResponse {
	connection, err := database.GetConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	var orderByPrepared string
	if schema.OrderBy == "created_at" {
		orderByPrepared = "u.created_at "
	}
	if schema.SortBy == "DESC" {
		orderByPrepared = orderByPrepared + "DESC"
	} else if schema.SortBy == "ASC" {
		orderByPrepared = orderByPrepared + "ASC"

	}

	repResult, dbStatus, err := selectListWithTales(connection, schema.DateFrom, schema.DateTo, schema.Limit, (page-1)*10, orderByPrepared)
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

func updateUserDataS(schema schema.UpdateUserSchema) *responses.MainResponse {
	connection, err := database.GetSqlConnection()
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	rStatus, err := updateUserDataR(connection, nil, nil, schema)

	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}
	if rStatus == constants.NotFound {
		return &responses.MainResponse{Status: http.StatusNotFound}
	} else if rStatus == constants.NoDataToUpdate {
		return &responses.MainResponse{Status: http.StatusBadRequest}
	}

	return &responses.MainResponse{Status: http.StatusOK}
}
