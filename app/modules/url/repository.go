package url

import (
	"go-fiber-api-template/app/common/database"
	. "go-fiber-api-template/app/common/database/jet/tsFastifyTemplate/public/table"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/common/types/entities"
	"go-fiber-api-template/app/modules/url/schema"

	. "github.com/go-jet/jet/v2/postgres"
)

func selectOneR(schema schema.RedirectSchema) entities.Url {
	var dest = entities.Url{}

	stmt := SELECT(URL.URL).FROM(URL).ORDER_BY(URL.URL.DESC()).LIMIT(1)

	var err = stmt.Query(database.GetConnection, &dest)

	helpers.PanicOnError(err)

	return dest
}

// func insertR(schema schema.InsertUserSchema) entities.User {
// 	newUser := model.User{
// 		Name:  schema.Name,
// 		Phone: schema.Phone,
// 	}

// 	dest := entities.User{}

// 	stmt := User.INSERT(User.Name, User.Phone).MODEL(newUser).RETURNING(User.AllColumns)

// 	err := stmt.Query(database.GetConnection, &dest)

// 	helpers.PanicOnError(err)

// 	return dest
// }
