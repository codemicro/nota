package endpoints

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/helpers"
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
)

var (
	allowableMimeTypes = []string{"image/jpeg", "image/png"}
)

// Session functions
func apiGetAllSessions(c *fiber.Ctx) {
	conn := database.Conn
	var sessions []models.Session
	conn.Find(&sessions)

	c.JSON(sessions)
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

func apiCreateSession(c *fiber.Ctx) {

	subject := c.FormValue("subject")
	title := c.FormValue("title")
	time := c.FormValue("timestamp")

	if subject == "" || title == "" {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "'subject' or 'title' field is missing/blank",
		})
		return
	}

	var submitTime int32
	if time == "" {
		submitTime = 0
	} else {
		st, err := strconv.ParseInt(time, 10, 32)
		if err != nil {
			c.Status(400).JSON(models.GenericResponse{
				Status:  "error",
				Message: "'timestamp' is not a number",
			})
			return
		}
		submitTime = int32(st)
	}

	conn := database.Conn
	session := models.Session{
		Time:    submitTime,
		Subject: subject,
		Title:   title,
	}

	conn.Create(&session)

	c.JSON(session)
}

// File functions
func apiAddFile(c *fiber.Ctx) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)

	// Check the ID is a valid number
	if err != nil {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "ID is not a valid number",
		})
		return
	}

	// Check name exists in request params
	filename := c.FormValue("name")
	if filename == "" {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "'name' is not present, or is blank",
		})
		return
	}

	conn := database.Conn

	// check there's a session with the specified ID in the db
	var count int
	conn.Model(&models.Session{}).Where("id = ?", int(id)).Count(&count)

	if count == 0 {
		c.Status(404).JSON(models.GenericResponse{
			Status:  "error",
			Message: "Session ID not found",
		})
		return
	}

	// Get multipart image
	imageFile, err := c.FormFile("image")
	if err != nil {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "No 'image' included",
		})
		return
	}

	// Check file is an image and not some other type of file
	fileHandler, _ := imageFile.Open()
	fileCont, err := ioutil.ReadAll(fileHandler)

	mimeType := http.DetectContentType(fileCont)

	if !helpers.IsStringInSlice(mimeType, allowableMimeTypes) {
		c.Status(400).JSON(models.GenericResponse{
			Status:  "error",
			Message: "Bad image format",
		})
		return
	}

	// Determine image path
	var filePath string
	for {
		var fileExt string
		fileExt, err = helpers.GetFileExtension(filename)
		if err != nil {
			fileExt, _ = helpers.MimeTypeToFileExt(mimeType)
		}
		randomFilename, _ := helpers.RandomHex(20)
		filePath = randomFilename + "." + fileExt

		var num int
		conn.Model(&models.File{}).Where(&models.File{Path: filePath}).Count(&num)

		if num == 0 {
			break
		}
	}

	filePath = "img/" + filePath

	err = helpers.SaveBytesToDisk(filePath, fileCont)

	if err != nil {
		c.Status(500).JSON(models.GenericResponse{
			Status:  "error",
			Message: "Internal server error - check back later",
		})
		return
	}

	responseObject := models.File{
		Name:          filename,
		Path:          filePath,
		ParentSession: id,
	}

	// Insert

	conn.Create(&responseObject)

	c.JSON(responseObject)
}
