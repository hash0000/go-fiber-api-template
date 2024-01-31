package user

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")

func signIn(dto dto.SignInDto) *responses.MainResponse {
	var ctxDb, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userModel models.User

	err := userCollection.FindOne(ctxDb, fiber.Map{"username": dto.Username}).Decode(&userModel)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusNotFound}
	}

	match, err := comparePasswordAndHash(dto.Password, userModel.Password)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}
	if match != true {
		return &responses.MainResponse{Status: http.StatusForbidden}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userModel.Username
	claims["userId"] = userModel.Id
	claims["grandAccess"] = userModel.GrandAccess
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(configs.GetEnvVar("JWT_SECRET")))
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	return &responses.MainResponse{Status: http.StatusCreated, Data: fiber.Map{"accessToken": t}}
}

func create(dto dto.CreateDto, jwtClaims *types.JwtClaimType) *responses.MainResponse {
	if jwtClaims.GrandAccess == false {
		return &responses.MainResponse{Status: http.StatusForbidden}
	}

	var ctxDb, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userModel models.User

	err := userCollection.FindOne(ctxDb, fiber.Map{"username": dto.Username}).Decode(&userModel)
	if err == nil {
		return &responses.MainResponse{Status: http.StatusConflict}
	}

	hash, err := generateFromPassword(dto.Password, configs.Argon2Params)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		Username:    dto.Username,
		Password:    hash,
		GrandAccess: false,
	}

	_, err = userCollection.InsertOne(ctxDb, newUser)
	if err != nil {
		return &responses.MainResponse{Status: http.StatusInternalServerError}
	}

	return &responses.MainResponse{Status: http.StatusCreated}
}

func generateFromPassword(password string, p *types.Argon2ParamType) (encodedHash string, err error) {
	var salt = []byte(configs.GetEnvVar("ARGON2_SALT"))
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *types.Argon2ParamType, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, err
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, err
	}

	p = &types.Argon2ParamType{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
