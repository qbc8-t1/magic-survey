package server

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, s *common.Server) {
	app.Get("/health", handlers.HealthCheck)

	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	questions := api.Group("/questions")
	answers := api.Group("/answers")

	routes.RegisterUserRoutes(auth, s)
	routes.RegisterQuestionRoutes(questions, s)
	routes.RegisterAnswerRoutes(answers, s)
}
