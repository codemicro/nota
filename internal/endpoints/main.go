package endpoints

import (
	"github.com/codemicro/nota/internal/authentication"
	"github.com/codemicro/nota/internal/helpers"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	"log"
)

func InitEndpoints(app *fiber.App) {
	// Visual endpoints
	app.Get("/", visualIndex)

	// Files
	app.Static("/img", "img")

	// API endpoints that do not require authentication
	app.Post("/api/login/", apiLogIn)
	app.Post("/api/register/", apiRegister)

	// Add JWT authentication
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: authentication.PubKey,
		ErrorHandler: func(c *fiber.Ctx, _ error) {
			helpers.UnauthorisedResponse(c, "Invalid or expired token")
		},
		SigningMethod: "RS256",
		ContextKey:    "jwt",
		TokenLookup:   "cookie:token",
	}))

	// All other API endpoints
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

	log.Println("Endpoints all setup")
}
