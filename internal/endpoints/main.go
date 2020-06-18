package endpoints

import (
	"github.com/gofiber/fiber"
	"log"
)


func InitEndpoints(app *fiber.App) {
	// Visual endpoints
	app.Get("/", visualIndex)

	// API endpoints
	app.Post("/api/session/create", apiCreateSession)

	log.Println("Setup endpoints")
}