package endpoints

import (
	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
	"strconv"
)


func apiCreateSession(c *fiber.Ctx) {
	// TODO: Parse values
	conn := database.Conn
	session := models.Session{
		Time:    1592491944,
		Subject: "Maths",
		Title:   "Quadratic equations",
	}

	conn.Create(&session)

	c.JSON(session)
}


func apiGetSession(c *fiber.Ctx) {
	conn := database.Conn

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "ID is not a valid number",
		})
		return
	}

	var session models.Session

	if conn.Find(&session, id).RecordNotFound() {
		c.Status(404).JSON(models.GenericResponse{
			Status:  "error",
			Message: "Session ID not found",
		})
		return
	}

	var files []models.File
	conn.Where(&models.File{ParentSession: id}).Find(&files)

	c.JSON(models.SessionWithFiles{
		Session: session,
		Files:   files,
	})
}


func apiAddFile(c *fiber.Ctx) {
	// TODO: Support file upload

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)

	// Check the ID is a valid number
	if err != nil {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "ID is not a valid number",
		})
		return
	}

	// Check name exists

	filename := c.FormValue("name")
	if filename == ""{
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "'name' is not present, or is blank",
		})
		return
	}

	conn := database.Conn

	// check there's a session with the specified ID

	var sessions []models.Session
	var count int
	conn.Model(&sessions).Where("id = ?", int(id)).Count(&count)

	if count == 0 {
		c.Status(404).JSON(models.GenericResponse{
			Status:  "error",
			Message: "Session ID not found",
		})
		return
	}

	file := models.File{
		Name:  filename,
		Path:  filename,
		ParentSession: id,
	}

	// Insert

	conn.Create(&file)

	c.JSON(file)
}

func apiGetAllSessions(c *fiber.Ctx) {
	conn := database.Conn
	var sessions []models.Session
	conn.Find(&sessions)

	c.JSON(sessions)
}