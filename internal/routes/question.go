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

	app.Get("/hello", handlers.HelloQuestionHandler(questionService))
	app.Post("/questions", handlers.CreateQuestionHandler(questionService))
	app.Get("/questions/:id", handlers.GetQuestionHandler(questionService))
	app.Get("/questions", handlers.GetQuestionsHandler(questionService))
	app.Put("/questions/:id", handlers.UpdateQuestionHandler(questionService))
	app.Delete("/questions/:id", handlers.DeleteQuestionHandler(questionService))
}
