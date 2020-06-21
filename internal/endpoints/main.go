package endpoints

import (
	"github.com/gofiber/fiber"
	"log"
)

func InitEndpoints(app *fiber.App) {
	// Visual endpoints
	app.Get("/", visualIndex)

	// API endpoints
	app.Post("/api/session/create/", apiCreateSession)
	app.Get("/api/session/", apiGetAllSessions)
	app.Get("/api/session/:id/", apiGetSession)
	app.Post("/api/session/:id/files/add/", apiAddFile)
	app.Get("/api/files/:id/rotate/", apiRotateImage)

	// Files
	log.Println("Endpoints all setup")
}
