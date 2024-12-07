package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes registers routes related to user management
func RegisterQuestionnaireRoutes(r fiber.Router, s *common.Server) {
	qRepo := repository.NewQuestionnaireRepository(s.DB)
	userRepo := repository.NewUserRepository(s.DB)
	// repo domain_repository.IQuestionnaireRepository
	qService := service.NewQuestionnaireService(qRepo, userRepo)

	router := r.Group("/questionnaires")
	router.Post("/", handlers.QuestionnaireCreate(qService))
	router.Post("/:qid", handlers.QuestionnaireUpdate(qService))
	router.Get("/:qid", handlers.QuestionnaireGet(qService))
	router.Delete("/:qid", handlers.QuestionnaireDelete(qService))

}
