package database

import (
	"database/sql"
	"go-fiber-api-template/app/common/helpers"
)

var GetConnection *sql.DB

func SetCon(connectString string) {
	db, err := sql.Open("postgres", connectString)
	helpers.PanicOnError(err)

	GetConnection = db
}
