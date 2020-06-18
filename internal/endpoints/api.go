package endpoints

import (
	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
)


func apiCreateSession(c *fiber.Ctx) {
	db := database.Conn
	session := models.Session{
		Time:    1592491944,
		Subject: "Maths",
		Files:   nil,
		Title:   "Quadratic equations",
	}
	db.Create(&session)

	if err := c.JSON(session); err != nil {
		c.Status(500).JSON(&models.GenericResponse{
			Status:  "error",
			Message: "An internal server error occurred. Try again later.",
		})
	}
}
