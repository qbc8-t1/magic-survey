package server

import (
	"github.com/QBC8-Team1/magic-survey/domain/repository"
	"github.com/QBC8-Team1/magic-survey/handler"
	"github.com/QBC8-Team1/magic-survey/service"
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
