package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"go-fiber-api-template/app/common/helpers"
	"log"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func SetPostgresConnection(connectString string) error {
	var err error
	db, err = sql.Open("postgres", connectString)
	if err != nil {
		slog.Error("Error connecting to Postgres", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		slog.Error("Error pinging database", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		return err
	}

	slog.Info("Successfully connected to the Postgres")

	return nil
}

func GetConnection() (*sql.DB, error) {
	var err error
	if db == nil {
		err = fmt.Errorf("Database connection is not initialized")
		slog.Error("Database connection is not initialized", slog.String("location", helpers.GetFileLine(constants.DeepCallerConstant.Modules)))
		return nil, err
	}
	return db, err
}

func Close() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatalf("error closing database connection: %v", err)
		}
	}
}
