package server

import (
	"github.com/QBC8-Team1/magic-survey/internal/domain/repository"
	"github.com/QBC8-Team1/magic-survey/internal/handler"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, server *Server) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

}

func QuestionRoutes(app *fiber.App, server *Server) {
	questionRepo := repository.NewQuestionRpository(server.db)
	questionService := service.NewQuestionService(questionRepo)

	app.Get("/hello", handler.HelloHandlerQuestion(questionService))
}

func QuestionnaireRoutes(app *fiber.App, server *Server) {
	questionnaireRepo := repository.NewQuestionnaireRpository(server.db)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepo)

	app.Get("/hello", handler.HelloHandlerQuestionnaire(questionnaireService))
}
