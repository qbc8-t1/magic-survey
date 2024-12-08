package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterAnswerRoutes registers routes related to answer management

func RegisterAnswerRoutes(router fiber.Router, s *common.Server) {
	answerRepo := repository.NewAnswerRepository(s.DB)
	userRepo := repository.NewUserRepository(s.DB)
	submissionRepo := repository.NewSubmissionRepository(s.DB)
	questionRepo := repository.NewQuestionRepository(s.DB)
	optionRepo := repository.NewOptionRepository(s.DB)

	answerService := service.NewAnswerService(answerRepo, userRepo, submissionRepo, questionRepo, optionRepo)

	withAuthMiddleware := middleware.WithAuthMiddleware(s.DB, s.Cfg.Secret)

	router.Get("/:id", withAuthMiddleware, handlers.GetAnswerHandler(answerService))
	router.Post("", withAuthMiddleware, handlers.CreateAnswerHandler(answerService))
	router.Put("/:id", withAuthMiddleware, handlers.UpdateAnswerHandler(answerService))
	router.Delete("/:id", withAuthMiddleware, handlers.DeleteAnswerHandler(answerService))
}
