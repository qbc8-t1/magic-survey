package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterAnswerRoutes registers routes related to answer management
func RegisterVisibleAnswersRoutes(router fiber.Router, s *common.Server) {
	answerRepo := repository.NewAnswerRepository(s.DB)
	userRepo := repository.NewUserRepository(s.DB)
	submissionRepo := repository.NewSubmissionRepository(s.DB)
	questionRepo := repository.NewQuestionRepository(s.DB)
	optionRepo := repository.NewOptionRepository(s.DB)
	rbacRepo := repository.NewRbacRepository(s.DB)
	questionnaireRepo := repository.NewQuestionnaireRepository(s.DB)

	answerService := service.NewAnswerService(answerRepo, userRepo, submissionRepo, questionRepo, optionRepo)
	rbacService := service.NewRbacService(rbacRepo)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepo)
	questionService := service.NewQuestionService(questionRepo, questionnaireRepo)

	router.Get("/see-another-user-answer", handlers.GetAnotherUserAnswer(answerService, rbacService, questionService, questionnaireService))
	router.Get("/users-with-visible-answers", handlers.GetUsersWithVisibleAnswers(rbacService, questionnaireService))
}
