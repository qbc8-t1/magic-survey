package server

import "github.com/gofiber/fiber/v2"

func registerRoutes(app *fiber.App, server *Server) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

}
