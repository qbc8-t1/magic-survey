package server

import (
	"errors"
	"github.com/gofiber/fiber/v2"

	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/internal/common"
	"github.com/QBC8-Team1/magic-survey/pkg/db"
	"github.com/QBC8-Team1/magic-survey/pkg/logger"
)

func NewServer(cfg *config.Config) (*common.Server, error) {
	appLogger := logger.NewAppLogger(cfg)

	appLogger.InitLogger(cfg.Logger.Path)
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, "")

	app := fiber.New()

	db, err := db.InitDB(cfg, appLogger)
	if err != nil {
		appLogger.Fatal(errors.Join(errors.New("can't init the db connection"), err))
	}

	s := &common.Server{
		App:    app,
		Logger: appLogger,
		Cfg:    cfg,
		DB:     db,
	}

	registerRoutes(s, cfg.Server.Secret)
	return s, nil
}
