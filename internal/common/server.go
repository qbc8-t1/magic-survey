package common

import (
	"context"
	"fmt"
	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	App    *fiber.App
	Logger *logger.AppLogger
	Cfg    *config.Config
	DB     *gorm.DB
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%s", s.Cfg.Server.Host, s.Cfg.Server.Port)
	return s.App.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.App.ShutdownWithContext(ctx); err != nil {
		return err
	}
	appDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("error accessing database connection: %w", err)
	}

	if err := appDB.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}

	return nil
}
