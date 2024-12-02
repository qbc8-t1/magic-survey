package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// QuestionRoutes registers routes related to question management
func RegisterQuestionRoutes(api fiber.Router, s *common.Server) {
	questionRepo := repository.NewQuestionRpository(s.DB)
	questionService := service.NewQuestionService(questionRepo)

	api.Get("/hello", handlers.HelloQuestionHandler(questionService))
	api.Post("", handlers.CreateQuestionHandler(questionService))
	api.Get("/:id", handlers.GetQuestionHandler(questionService))
	api.Get("", handlers.GetQuestionsHandler(questionService)) //TODO: get all questions of a qustionnaire id
	api.Put("/:id", handlers.UpdateQuestionHandler(questionService))
	api.Delete("/:id", handlers.DeleteQuestionHandler(questionService))
}
