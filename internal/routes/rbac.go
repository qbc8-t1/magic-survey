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
	rbacService := service.NewRbacService(rbacRepo)

	// logged in user id here

	router.Get("/users-with-visible-answers", handlers.GetUsersWithVisibleAnswers(*rbacService))
	router.Post("/:userid/make-superadmin", handlers.MakeSuperadmin(*rbacService))
	router.Get("/can-do", handlers.CanDo(*rbacService))
	router.Get("/permissions", handlers.GetAllPermissions(*rbacService))
	router.Post("/:userid/give-permissions", handlers.GivePermissions(*rbacService))
	router.Post("/make-fake-user", handlers.MakeFakeUser(*rbacService))
	router.Post("/:userid/make-fake-questionnaire", handlers.MakeFakeQuestionnaire(*rbacService))
	router.Get("/:userid/roles-with-permissions", handlers.GetUserRolesWithPermissions(*rbacService))
	router.Post("/:userid/revoke-permission", handlers.RevokePermission(*rbacService))
	router.Get("/:userid", handlers.GetUser(*rbacService))
}
