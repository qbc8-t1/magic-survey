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

	api := s.App.Group("/api/v1")

	middleware.RegisterRbacMiddlewares(api, s.DB)

	auth := api.Group("/auth")
	questions := api.Group("/questions")
	answers := api.Group("/answers")
	options := api.Group("/options")
	rbac := api.Group("/rbac")
	core := api.Group("/core")

	superadmin := api.Group("/superadmin")

	questionnaire := api.Group("/questionnaires/:questionnaire_id")
	// questionnareQuestions := questionnaire.Group("/questions")
	// questionnaireAnswers := questionnaire.Group("/answers")
	// questionnaireOptions := questionnaire.Group("/options")

	routes.RegisterVisibleAnswersRoutes(questionnaire, s)
	routes.RegisterQuestionRoutes(questions, s)
	routes.RegisterAnswerRoutes(answers, s)
	routes.RegisterOptionRoutes(options, s)
	routes.RegisterRbacRoutes(rbac, s)
	routes.RegisterCoreRoutes(core, s)
	routes.RegisterSuperadminRoutes(superadmin, s)
	routes.RegisterUserRoutes(auth, s)
}
