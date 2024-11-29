package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

func QuestionRoutes(app *fiber.App, s *common.Server) {
	questionRepo := repository.NewQuestionRpository(s.DB)
	questionService := service.NewQuestionService(questionRepo)

	app.Get("/hello", handlers.HelloHandlerQuestion(questionService))
}
