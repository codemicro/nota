package main


import (
	"github.com/codemicro/nota/internal/database"
	"github.com/gofiber/fiber"
	"log"
	"os"

	"github.com/codemicro/nota/internal/endpoints"
)


var (
	serveAtAddress = 8000
	version = "0.0.0a"
)


func main() {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) {
		c.Append("X-Powered-By", "Nota v" + version)
		c.Next()
	})

	database.InitDatabase()
	endpoints.InitEndpoints(app)

	// Make directories and ignore any errors (assuming they already exist)
	_ = os.Mkdir("img/", os.ModeDir)

	log.Panic(app.Listen(serveAtAddress))

}
