package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// QuestionRoutes registers routes related to question management
func RegisterQuestionRoutes(app *fiber.App, s *common.Server) {
	questionRepo := repository.NewQuestionRpository(s.DB)
	questionService := service.NewQuestionService(questionRepo)

	app.Get("/hello", handlers.HelloHandlerQuestion(questionService))
	app.Post("/create-question")
	app.Get("/get-question")
	app.Put("/update-question")
	app.Delete("/delete-question")
	app.Get("/get-all-questions")
}
