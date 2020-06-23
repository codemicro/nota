package main

import (
	"fmt"
	"github.com/codemicro/nota/internal/authentication"
	"github.com/codemicro/nota/internal/config"
	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/endpoints"
	"github.com/codemicro/nota/internal/logging"
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"log"
	"os"
)

var (
	serveAtAddress = 8000
	version        = "0.0.0a"
)

func main() {
	app := fiber.New()

	// Make directories and ignore any errors (assuming they already exist)
	_ = os.Mkdir("img/", os.ModeDir)

	// Load settings and setup logging
	err := config.LoadSettings()
	if err != nil {
		fmt.Println("Error when loading settings. Check the error logs for more detail.")
		logging.ErrorLogger.Fatalln(err)
	}
	logging.InitLogging()

	// Middleware
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{ // Setup access logging
		Next:       nil,
		Format:     "${ip} - ${time} - ${method} \"${url}\" ${protocol} ${status} \"${ua}\"\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		Output:     logging.AccessLogFile,
	}))

	app.Use(middleware.Recover()) // Recover from any panic within the rest of the app

	app.Use(func(c *fiber.Ctx) {
		c.Append("X-Powered-By", "Nota v"+version)
		c.Next()
	})

	app.Settings.ErrorHandler = logging.ErrorHandler

	// Setup database and register endpoints
	database.InitDatabase()
	err = authentication.LoadKeys()
	if err != nil {
		fmt.Println("Error when loading RSA keys. Check the error logs for more detail.")
		logging.ErrorLogger.Fatalln(err)
	}
	endpoints.InitEndpoints(app)

	// 404 handler has to go down here otherwise it overrides everything
	app.Use(func(c *fiber.Ctx) {
		c.Status(fiber.StatusNotFound).JSON(models.GenericResponse{
			Status:  "error",
			Message: "404 - resource not found",
		})
	})

	// Let's run
	log.Panic(app.Listen(serveAtAddress))

}
