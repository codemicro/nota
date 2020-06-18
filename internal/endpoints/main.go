package endpoints


import "github.com/gofiber/fiber"


func InitEndpoints(app *fiber.App) {
	// Visual endpoints
	app.Get("/", visualIndex)

	// API endpoints
}