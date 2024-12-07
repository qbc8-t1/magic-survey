package routes

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/gofiber/fiber/v2"
)

func RegisterRbacRoutes(router fiber.Router, s *common.Server) {
	rbacRepo := repository.NewRbacRepository(s.DB)
	questionnaireRepo := repository.NewQuestionnaireRepository(s.DB)
	superAdminRepo := repository.NewSuperadminRepository(s.DB)

	rbacService := service.NewRbacService(rbacRepo)
	questionnaireService := service.NewQuestionnaireService(questionnaireRepo)
	superAdminService := service.NewSuperadminService(superAdminRepo)

	// logged in user id here
	router.Get("/users-with-visible-answers", handlers.GetUsersWithVisibleAnswers(rbacService, questionnaireService))
	router.Post("/:userid/make-superadmin", handlers.MakeSuperadmin(*superAdminService))
	router.Get("/can-do", handlers.CanDo(rbacService))
	router.Get("/permissions", handlers.GetAllPermissions(rbacService))
	router.Post("/give-permissions", handlers.GivePermissions(rbacService))
	router.Post("/revoke-permission", handlers.RevokePermission(rbacService))
	router.Get("/can-do", handlers.CanDo(rbacService))
	router.Get("/info", handlers.GetUser(rbacService))
	router.Get("/roles-with-permissions", handlers.GetUserRolesWithPermissions(rbacService))

	router.Post("/make-fake-user", handlers.MakeFakeUser(rbacService))
	router.Post("/make-fake-questionnaire", handlers.MakeFakeQuestionnaire(rbacService))
}
