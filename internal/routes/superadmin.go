package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

// RegisterSuperadminRoutes registers routes related to superadmin functionality
func RegisterSuperadminRoutes(auth fiber.Router, s *common.Server) {
	superadminRepo := repository.NewSuperadminRepository(s.DB)
	superadminService := service.NewSuperadminService(superadminRepo)

	auth.Post("/make-superadmin", handlers.MakeSuperadmin(*superadminService))
	auth.Post("/limit-user-questionnaires-count", handlers.LimitUserQuestionnaireCount(*superadminService))
}
