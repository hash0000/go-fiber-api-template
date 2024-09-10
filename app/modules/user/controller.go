package user

import (
	"go-fiber-api-template/app/common/schemas"
	"go-fiber-api-template/app/modules/user/schema"

	"github.com/gofiber/fiber/v2"
)

func Create(c *fiber.Ctx) error {
	var res = create(c.Locals("body").(schema.CreateUserSchema))

	return c.Status(res.Status).JSON(res)
}

func UpdateTokenNumber(ctx *fiber.Ctx) error {
	var request = updateTokenNumber(ctx.Locals("body").(schema.UpdateTokenNumberSchema))

	return ctx.Status(request.Status).JSON(request)
}

func GetOne(c *fiber.Ctx) error {
	var request = getOne(c.Locals("params").(schema.IdSchema))

	return c.Status(request.Status).JSON(request)
}

func GetOneWithBooks(ctx *fiber.Ctx) error {
	var request = getOneWithBooks(ctx.Locals("params").(schema.IdSchema))

	return ctx.Status(request.Status).JSON(request)
}

func getListWithTales(c *fiber.Ctx) error {
	params := c.Locals("params").(schemas.PaginationSchema)
	queryString := c.Locals("query_string").(schema.GetListWithTaleSchema)

	serviceResult := getListWithTaleS(queryString, params.Page)

	return c.Status(serviceResult.Status).JSON(serviceResult)
}

func updateUserData(c *fiber.Ctx) error {
	body := c.Locals("body").(schema.UpdateUserSchema)

	serviceResult := updateUserDataS(body)

	return c.Status(serviceResult.Status).JSON(serviceResult)
}
