package user

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"go-fiber-api-template/app/modules/user/schema"
	"net/http"
	"report-url-redirection/app/common/configs"
	"report-url-redirection/app/common/responses"
	"report-url-redirection/app/common/types"
	"report-url-redirection/app/modules/database/models"
	"report-url-redirection/app/modules/user/dto"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/argon2"
)

func signIn(schema schema.InsertUserSchema) *responses.MainResponse {

	return &responses.MainResponse{Status: http.StatusCreated, Data: fiber.Map{"data": schema}}
}
