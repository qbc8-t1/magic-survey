package server

import (
	"fmt"
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
		fmt.Println(cfg)
		appLogger.Panic("can't init the db connection")
	}

	s := &common.Server{
		App:    app,
		Logger: appLogger,
		Cfg:    cfg,
		DB:     db,
	}

	registerRoutes(app, s)
	return s, nil
}
