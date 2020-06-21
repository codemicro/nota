package main

import (
	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"log"
	"os"

	"github.com/codemicro/nota/internal/endpoints"
)

var (
	serveAtAddress = 8000
	version        = "0.0.0a"
)

func main() {
	app := fiber.New()

	app.Use(middleware.Recover()) // Recover from any panic within the rest of the app

	app.Use(func(c *fiber.Ctx) {
		c.Append("X-Powered-By", "Nota v"+version)
		c.Next()
	})

	app.Use(func(c *fiber.Ctx) {
		c.Status(fiber.StatusNotFound).JSON(models.GenericResponse{
			Status:  "error",
			Message: "404 - resource not found",
		})
	})

	app.Settings.ErrorHandler = func(ctx *fiber.Ctx, err error) {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Return HTTP response
		ctx.Status(code).JSON(models.GenericResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	database.InitDatabase()
	endpoints.InitEndpoints(app)

	// Make directories and ignore any errors (assuming they already exist)
	_ = os.Mkdir("img/", os.ModeDir)

	log.Panic(app.Listen(serveAtAddress))

}
