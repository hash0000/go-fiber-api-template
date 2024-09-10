package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/helpers"
	"log/slog"
)

var dbStatement *sql.DB

func SetSqlConnection(connectString string) error {
	var err error
	dbStatement, err = sql.Open("postgres", connectString)
	if err != nil {
		slog.Error("Error connecting to Postgres", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	err = dbStatement.Ping()
	if err != nil {
		slog.Error("Error pinging database", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	slog.Info("Successfully connected to the Postgres")

	return nil
}

func GetSqlConnection() (*sql.DB, error) {
	var err error
	if dbStatement == nil {
		err = fmt.Errorf("Database connection is not initialized")
		slog.Error("Database connection is not initialized", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
		return nil, err
	}
	return dbStatement, err
}

func CloseSqlConnection() {
	if db != nil {
		err := db.Close()
		if err != nil {
			slog.Error("Error closing database connection", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))

		}
	}
}
