package db

import (
	"fmt"

	"github.com/QBC8-Team1/magic-survey/config"
	applog "github.com/QBC8-Team1/magic-survey/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config, logger *applog.AppLogger) (*gorm.DB, error) {
	gormLogger := applog.NewGormLogger(logger)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Dbname,
		cfg.Postgres.Port,
		cfg.Postgres.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = deleteAllTablesAndTypes(db)
	if err != nil {
		return nil, fmt.Errorf("failed to delete tables and types: %w", err)
	}

	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	err = seed(db)
	if err != nil {
		return nil, fmt.Errorf("failed to seed: %w", err)
	}

	return db, nil
}
