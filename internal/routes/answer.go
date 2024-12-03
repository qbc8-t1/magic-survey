package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterAnswerRoutes registers routes related to answer management
func RegisterAnswerRoutes(api fiber.Router, s *common.Server) {
	answerRepo := repository.NewAnswerRepository(s.DB)
	userRepo := repository.NewUserRepository(s.DB)
	submissionRepo := repository.NewSubmissionRepository(s.DB)
	questionRepo := repository.NewQuestionRpository(s.DB)
	optionRepo := repository.NewOptionRepository(s.DB)

	answerService := service.NewAnswerService(answerRepo, userRepo, submissionRepo, questionRepo, optionRepo)

	api.Get("/:id", handlers.GetAnswerHandler(answerService))
	api.Post("", handlers.CreateAnswerHandler(answerService))
	api.Put("/:id", handlers.UpdateAnswerHandler(answerService))
	api.Delete("/:id", handlers.DeleteAnswerHandler(answerService))
}
