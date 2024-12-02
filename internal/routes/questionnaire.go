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
	// repo domain_repository.IQuestionnaireRepository
	qService := service.NewQuestionnaireService(qRepo)

	r.Post("/", handlers.QuestionnaireCreate(qService))
	r.Post("/:qid", handlers.QuestionnaireUpdate(qService))
	r.Get("/:qid", handlers.QuestionnaireGet(qService))
	r.Delete("/:qid", handlers.QuestionnaireDelete(qService))

}
