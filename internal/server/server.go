package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	echoMiddlewares "github.com/labstack/echo/v4/middleware"

	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/pkg/logger"
)

type Server struct {
	router *echo.Echo
	logger *logger.AppLogger
	cfg    *config.Config
	db     any // TODO: fix it when we add a db
}

func NewServer(cfg *config.Config) (*Server, error) {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	e := echo.New()

	s := &Server{
		router: e,
		logger: appLogger,
		cfg:    cfg,
		db:     nil,
	}

	registerRoutes(e, s)
	return s, nil
}

func (s *Server) Run() error {
	s.router.Use(echoMiddlewares.CORSWithConfig(echoMiddlewares.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)

	if s.cfg.Server.Mode == config.Development {
		return s.router.Start(addr)
	}
	return s.router.Start(addr)

}
