package helpers

import (
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"log/slog"
	"os"
	"strconv"
)

func GetPostgresConString() (string, error) {
	pgPort, err := strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		slog.Error("Error while converting PG_PORT to int", slog.String("location", GetFileLine(constants.DeepCallerConstant.CommonRoot)), slog.Any("info", err))
		return "", fmt.Errorf("error while converting PG_PORT to int")
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("PG_HOST"), pgPort, os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"), os.Getenv("PG_DB_NAME")), nil
}
