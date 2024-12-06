package server

import (
	"github.com/QBC8-Team1/magic-survey/handlers"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, s *common.Server) {
	app.Get("/health", middleware.WithAuthMiddleware(s.DB, s.Cfg.Server.Secret), handlers.HealthCheck)

	api := app.Group("/api/v1")
	middleware.RegisterRbacMiddlewares(api, s.DB)

	auth := api.Group("/auth")
	rbac := api.Group("/rbac")
	superadmin := api.Group("/superadmin")

	routes.RegisterRbacRoutes(rbac, s)
	routes.RegisterSuperadminRoutes(superadmin, s)

	routes.RegisterUserRoutes(auth, s)

	questionnaire := api.Group("/questionnaires/:questionnaire_id")
	questions := questionnaire.Group("/questions")
	answers := questionnaire.Group("/answers")
	options := questionnaire.Group("/options")

	routes.RegisterVisibleAnswersRoutes(questionnaire, s)
	routes.RegisterQuestionRoutes(questions, s)
	routes.RegisterAnswerRoutes(answers, s)
	routes.RegisterOptionRoutes(options, s)
}
