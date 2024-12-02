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

	api := app.Group("/api")
	auth := api.Group("/v1/auth")

	routes.RegisterUserRoutes(auth, s)

	rbac := api.Group("/rbac")
	routes.RegisterRbacRoutes(rbac, s)
}
