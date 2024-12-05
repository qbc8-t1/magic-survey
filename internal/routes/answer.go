package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterAnswerRoutes registers routes related to answer management
func RegisterAnswerRoutes(router fiber.Router, s *common.Server) {
	answerRepo := repository.NewAnswerRepository(s.DB)
	rbacRepo := repository.NewRbacRepository(s.DB)
	questionnaireRepo := repository.NewQuestionnaireRepository(s.DB)
	questionRepo := repository.NewQuestionRpository(s.DB)

	answerService := service.NewAnswerService(answerRepo)
	rbacService := service.NewRbacService(rbacRepo)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepo)
	questionService := service.NewQuestionService(questionRepo)

	router.Get("/see-another-user-answer", handlers.GetAnswer(*answerService, *rbacService, *questionnaireService, *questionService))
}
