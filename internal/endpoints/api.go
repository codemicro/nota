package endpoints

import (
	"github.com/codemicro/nota/internal/authentication"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/helpers"
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
)

var (
	allowableMimeTypes = []string{"image/jpeg", "image/png"}
)

// Login/register functions
func apiLogIn(c *fiber.Ctx) {

	// TODO: User verification and claim update

	token := jwt.New(jwt.SigningMethodRS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(authentication.PrivKey)
	if err != nil {
		c.Next(err)
		return
	}

	c.JSON(fiber.Map{"token": t})

}

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
		helpers.BadRequestResponse(c, "ID is not a valid number")
		return
	}

	var session models.Session

	if conn.Find(&session, id).RecordNotFound() {
		helpers.NotFoundResponse(c, "Session ID not found")
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

	// Not using c.BodyParser because this needs everything required
	subject := c.FormValue("subject")
	title := c.FormValue("title")
	time := c.FormValue("timestamp")

	if subject == "" || title == "" {
		helpers.BadRequestResponse(c, "'subject' or 'title' field is missing/blank")
		return
	}

	var submitTime int64
	if time == "" {
		submitTime = 0
	} else {
		var err error
		submitTime, err = strconv.ParseInt(time, 10, 64)
		if err != nil {
			helpers.BadRequestResponse(c, "'timestamp' is not a number")
			return
		}
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

func apiDeleteSession(c *fiber.Ctx) {
	id, hasFailed, httpCode, formedResponseModel := helpers.CheckAndConvertId(c.Params("id"), "session", &models.Session{})

	if hasFailed {
		c.Status(httpCode).JSON(formedResponseModel)
		return
	}

	conn := database.Conn

	var session models.Session
	conn.Find(&session, id)

	conn.Unscoped().Delete(&session)

	c.JSON(models.GenericResponse{
		Status:  "ok",
		Message: "Session deleted successfully",
	})

}

func apiUpdateSession(c *fiber.Ctx) {
	id, hasFailed, httpCode, formedResponseModel := helpers.CheckAndConvertId(c.Params("id"), "session", &models.Session{})

	if hasFailed {
		c.Status(httpCode).JSON(formedResponseModel)
		return
	}

	conn := database.Conn

	var session models.Session
	conn.Find(&session, id)

	var parseTime int64
	if c.FormValue("time") != "" {
		var err error
		parseTime, err = strconv.ParseInt(helpers.SelectiveUpdate(session.Title, c.FormValue("time")), 10, 64)
		if err != nil {
			helpers.BadRequestResponse(c, "time must be a valid number")
			return
		}
		session.Time = parseTime
	}

	session.Title = helpers.SelectiveUpdate(session.Title, c.FormValue("title"))
	session.Subject = helpers.SelectiveUpdate(session.Subject, c.FormValue("subject"))

	conn.Save(&session)
	resp := struct {
		models.GenericResponse
		Values models.Session `json:"session"`
	}{
		GenericResponse: models.GenericResponse{
			Status:  "ok",
			Message: "updated successfully",
		},
		Values: session,
	}

	c.JSON(resp)
}

// File functions
func apiGetAllFiles(c *fiber.Ctx) {
	conn := database.Conn
	var files []models.File
	conn.Find(&files)

	c.JSON(files)
}

func apiGetFile(c *fiber.Ctx) {
	id, hasFailed, httpCode, formedResponseModel := helpers.CheckAndConvertId(c.Params("id"), "file", &models.File{})

	if hasFailed {
		c.Status(httpCode).JSON(formedResponseModel)
		return
	}

	conn := database.Conn
	var file models.File
	conn.Find(&file, id)

	c.JSON(file)
}

func apiAddFile(c *fiber.Ctx) {
	id, hasFailed, httpCode, formedResponseModel := helpers.CheckAndConvertId(c.Params("id"), "session", &models.Session{})

	if hasFailed {
		c.Status(httpCode).JSON(formedResponseModel)
		return
	}

	// Check name exists in request params
	filename := c.FormValue("name")
	if filename == "" {
		helpers.BadRequestResponse(c, "'name' is not present, or is blank")
		return
	}

	// Get multipart image
	imageFile, err := c.FormFile("image")
	if err != nil {
		helpers.BadRequestResponse(c, "No 'image' included")
		return
	}

	// Check file is an image and not some other type of file
	fileHandler, _ := imageFile.Open()
	fileCont, err := ioutil.ReadAll(fileHandler)

	mimeType := http.DetectContentType(fileCont)

	if !helpers.IsStringInSlice(mimeType, allowableMimeTypes) {
		helpers.BadRequestResponse(c, "Bad image format")
		return
	}

	conn := database.Conn

	// Determine image path using random hex data
	var filePath string
	for {
		fileExt, err := helpers.GetFileExtension(filename)
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
		c.Next(err) // This error is now Fiber's problem (pass to predefined error handler)
		return
	}

	responseObject := models.File{
		Name:          filename,
		Path:          filePath,
		ParentSession: id,
	}

	conn.Create(&responseObject) // Insert to database

	c.JSON(responseObject)
}

func apiRotateImage(c *fiber.Ctx) {
	fileId, hasFailed, httpCode, formedResponseModel := helpers.CheckAndConvertId(c.Params("id"), "file", &models.File{})

	if hasFailed {
		c.Status(httpCode).JSON(formedResponseModel)
		return
	}

	conn := database.Conn

	var file models.File
	conn.Find(&file, fileId)

	err := helpers.RotateImage(file.Path)
	if err != nil {
		c.Next(err)
		return
	}

	c.JSON(models.GenericResponse{
		Status:  "ok",
		Message: "Success",
	})

}

func apiDeleteImage(c *fiber.Ctx) {
	fileId, hasFailed, httpCode, formedResponseModel := helpers.CheckAndConvertId(c.Params("id"), "file", &models.File{})
	if hasFailed {
		c.Status(httpCode).JSON(formedResponseModel)
		return
	}

	conn := database.Conn

	var file models.File
	conn.Find(&file, fileId) // Populate a file model based on ID

	err := helpers.DeleteFile(file)
	if err != nil {
		c.Next(err)
		return
	}

	c.JSON(models.GenericResponse{
		Status:  "ok",
		Message: "File deleted successfully",
	})
}

func apiRemoveOrphanedFiles(c *fiber.Ctx) {
	conn := database.Conn
	var files []models.File
	countDeleted := 0
	conn.Find(&files)

	for _, file := range files {
		if conn.Find(&models.Session{}, file.ParentSession).RecordNotFound() {
			err := helpers.DeleteFile(file)
			if err != nil {
				c.Next(err)
				return
			}
			countDeleted++
		}
	}

	c.JSON(models.GenericResponse{
		Status:  "success",
		Message: strconv.Itoa(countDeleted) + " orphaned file(s) deleted",
	})
}
