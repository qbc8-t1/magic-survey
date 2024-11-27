package server

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, s *Server) {
	app.Get("/health", handlers.HealthCheck)

	api := app.Group("/api")
	auth := api.Group("/v1/auth")

	auth.Post("signup", handlers.UserCreate(s.db))
}
