package server

import (
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, s *common.Server) {
	app.Get("/health", middleware.WithAuthMiddleware(s.DB, s.Cfg.Server.Secret), handlers.HealthCheck)

	api := app.Group("/api")

	auth := api.Group("/v1/auth")

	middleware.RegisterRbacMiddlewares(api, s.DB)

	routes.RegisterUserRoutes(auth, s)

	rbac := api.Group("/rbac")
	routes.RegisterRbacRoutes(rbac, s)

	superadminGroup := api.Group("/superadmin")
	routes.RegisterSuperadminRoutes(superadminGroup, s)

	questionnaireGroup := api.Group("/questionnaires/:questionnaire_id")
	routes.RegisterAnswerRoutes(questionnaireGroup, s)
}
