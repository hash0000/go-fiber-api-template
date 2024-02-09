package user

import (
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/database/jet/goApi/public/model"
	. "go-fiber-api-template/app/common/database/jet/goApi/public/table"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/types/entities"
	"go-fiber-api-template/app/modules/user/schema"

	. "github.com/go-jet/jet/v2/postgres"
)

func insertR(schema schema.InsertUserSchema) entities.User {
	newUser := model.User{
		Name:  schema.Name,
		Phone: schema.Phone,
	}

	dest := entities.User{}

	stmt := User.INSERT(User.Name, User.Phone).MODEL(newUser).RETURNING(User.AllColumns)

	err := stmt.Query(database.GetConnection, &dest)

	helpers.PanicOnError(err)

	return dest
}

func selectOneR(schema schema.SelectOneUserSchema) entities.User {
	var dest = entities.User{}

	stmt := SELECT(
		User.ID, User.Name, User.Phone,
	).FROM(
		User,
	).WHERE(
		User.ID.EQ(UUID(schema.Id)),
	).LIMIT(1)

	var err = stmt.Query(database.GetConnection, &dest)

	helpers.PanicOnError(err)

	return dest
}
