package endpoints

import (
	"github.com/gofiber/fiber"
	"log"
)

func InitEndpoints(app *fiber.App) {
	// Visual endpoints
	app.Get("/", visualIndex)

	// API endpoints
	app.Post("/api/sessions/create/", apiCreateSession)
	app.Get("/api/sessions/", apiGetAllSessions)
	app.Get("/api/sessions/:id/", apiGetSession)
	app.Delete("/api/sessions/:id/", apiDeleteSession)
	app.Patch("/api/sessions/:id/", apiUpdateSession)
	app.Post("/api/sessions/:id/files/add/", apiAddFile)

	app.Get("/api/files/", apiGetAllFiles)
	app.Delete("/api/files/deleteOrphaned", apiRemoveOrphanedFiles)
	app.Get("/api/files/:id/", apiGetFile)
	app.Delete("/api/files/:id/", apiDeleteImage)
	app.Get("/api/files/:id/rotate/", apiRotateImage)

	// Files
	app.Static("/img", "img")

	log.Println("Endpoints all setup")
}
