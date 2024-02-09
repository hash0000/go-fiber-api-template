package user

import (
	"go-fiber-api-template/app/common/responses"
	"go-fiber-api-template/app/modules/user/schema"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func insert(schema schema.InsertUserSchema) *responses.MainResponse {

	result := insertR(schema)

	return &responses.MainResponse{Status: http.StatusCreated, Data: fiber.Map{"data": result}}
}

func selectOne(schema schema.SelectOneUserSchema) *responses.MainResponse {

	result := selectOneR(schema)

	return &responses.MainResponse{Status: http.StatusCreated, Data: fiber.Map{"data": result}}
}
