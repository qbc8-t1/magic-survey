package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/pkg/logger"
)

type Server struct {
	app    *fiber.App
	logger *logger.AppLogger
	cfg    *config.Config
	db     any // TODO: fix it when we add a db
}

func NewServer(cfg *config.Config) (*Server, error) {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	app := fiber.New()

	s := &Server{
		app:    app,
		logger: appLogger,
		cfg:    cfg,
		db:     nil,
	}

	registerRoutes(app, s)
	return s, nil
}

func (s *Server) Run() error {
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)

	return s.app.Listen(addr)
}
