package logging

import (
	"github.com/codemicro/nota/internal/models"
	"github.com/gofiber/fiber"
	"log"
	"os"
)

var (
	AccessLogFile *os.File
	errorLogFile  *os.File
	ErrorLogger   *log.Logger
)

func InitLogging() {
	_ = os.Mkdir("logs/", os.ModeDir) // Ensure log dir exists
	// Open files
	AccessLogFile, _ = os.OpenFile("logs/access.log", os.O_APPEND|os.O_CREATE, 0644)
	errorLogFile, _ = os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE, 0644)

	ErrorLogger = log.New(errorLogFile, "error ", log.LstdFlags) // Register ErrorLogger
}

func ErrorHandler(ctx *fiber.Ctx, err error) {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	ErrorLogger.Println(err.Error())

	var message string

	if code == 500 {
		message = "Internal server error - check back later"
	} else {
		message = err.Error()
	}

	// Return HTTP response
	ctx.Status(code).JSON(models.GenericResponse{
		Status:  "error",
		Message: message,
	})
}
