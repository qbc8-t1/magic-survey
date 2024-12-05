package server

import (
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"time"
)

func registerRoutes(app *fiber.App, s *common.Server) {
	limiterCfg := limiter.Config{
		Max:               10,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	}
	app.Use(limiter.New(limiterCfg), compress.New(), middleware.WithLogger(s))
	app.Get("/health", monitor.New())

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
