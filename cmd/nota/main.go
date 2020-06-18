package main

import (
	"github.com/gofiber/fiber"
	"log"

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

	endpoints.InitEndpoints(app)

	log.Panic(app.Listen(serveAtAddress))

}
