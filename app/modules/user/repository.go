package user

import (
	"go-fiber-api-template/app/common/database"
	"go-fiber-api-template/app/common/database/jet/goApi/public/model"
	. "go-fiber-api-template/app/common/database/jet/goApi/public/table"
	"go-fiber-api-template/app/common/helpers"
	"go-fiber-api-template/app/modules/user/schema"
	"time"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
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

func selectOneR(schema schema.SelectOneUserSchema) any {
	var dest struct {
		ID        uuid.UUID `sql:"primary_key" json:"id"`
		Name      string    `json:"name"`
		Phone     string    `json:"phone"`
		CreatedAt time.Time `json:"createdAt"`
	}

	stmt := SELECT(
		User.ID.AS("id"), User.Name.AS("name"), User.Phone.AS("phone"), User.CreatedAt.AS("createdAt"),
	).FROM(
		User,
	).WHERE(
		User.ID.EQ(UUID(schema.Id)),
	).ORDER_BY(
		User.CreatedAt.DESC(),
	).LIMIT(1)

	var err = stmt.Query(database.GetConnection, &dest)

	helpers.PanicOnError(err)

	return dest
}
