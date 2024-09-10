package configs

import "github.com/gofiber/fiber/v2/middleware/cors"

var CorsConfig = cors.New(cors.Config{
	AllowOrigins:  "*",
	AllowMethods:  "GET, POST, PUT, PATCH, DELETE",
	AllowHeaders:  "Origin, Content-Type, Accept, api-token",
	ExposeHeaders: "api-token",
})
