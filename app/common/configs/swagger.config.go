package configs

import "github.com/gofiber/contrib/swagger"

var SwaggerConfig = swagger.Config{
	BasePath: "/",
	FilePath: "./external/documentation/swagger.yaml",
	Path:     "docs",
	Title:    "Template docs",
}
