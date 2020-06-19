package endpoints

import "github.com/gofiber/fiber"

func visualIndex(c *fiber.Ctx) {
	c.Type("html", "utf-8")
	c.Send("<h1>Hello!</h1>")
}
