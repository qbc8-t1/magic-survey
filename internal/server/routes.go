package server

import (
	"time"

	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/internal/middleware"
	"github.com/QBC8-Team1/magic-survey/internal/routes"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func registerRoutes(s *common.Server, secret string) {
	limiterCfg := limiter.Config{
		Max:               10,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.FixedWindow{},
	}

	// this is just for learning purpose, not useful in real-world project
	CORSCfg := cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}

	s.App.Use(middleware.WithLoggingMiddleware(s.Logger), limiter.New(limiterCfg), cors.New(CORSCfg), compress.New())
	s.App.Get("/health", monitor.New())

	versionApi := s.App.Group("/api/v1")

	auth := versionApi.Group("/auth")
	routes.RegisterUserRoutes(auth, s)

	api := versionApi.Group("/")
	api.Use(middleware.WithAuthMiddleware(s.DB, secret))
	middleware.RegisterRbacMiddlewares(api, s.DB)

	rbac := api.Group("/rbac")
	questions := api.Group("/questions")
	answers := api.Group("/answers")
	options := api.Group("/options")
	questionnaires := api.Group("/questionnaires")
	superadmin := api.Group("/superadmin")
	core := api.Group("/core")

	routes.RegisterVisibleAnswersRoutes(questionnaires, s)
	routes.RegisterQuestionnaireRoutes(questionnaires, s)
	routes.RegisterQuestionRoutes(questions, s)
	routes.RegisterAnswerRoutes(answers, s)
	routes.RegisterOptionRoutes(options, s)
	routes.RegisterRbacRoutes(rbac, s)
	routes.RegisterCoreRoutes(core, s)
	routes.RegisterSuperadminRoutes(superadmin, s)
	routes.RegisterUserRoutes(auth, s)
}
