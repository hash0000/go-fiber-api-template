package user

import (
	"github.com/google/uuid"
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/database/jet/goApi/public/model"
	. "go-fiber-api-template/app/common/database/jet/goApi/public/table"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/modules/user/schema"
)

func insertR(schema schema.InsertUserSchema) any {
	user := model.User{
		Name:  schema.Name,
		Phone: schema.Phone,
	}

	stmt := User.INSERT(User.Name, User.Phone).MODEL(user).RETURNING(User.ID, User.Name)

	var dest struct {
		Id   uuid.UUID `sql:"primary_key"`
		Name string
	}

	err := stmt.Query(database.GetConnection, &user)

	helpers.PanicOnError(err)

	return dest
}
