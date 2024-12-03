package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/internal/server"
	"gorm.io/gorm"
)

var (
	testDB *gorm.DB
)

func TestMain(m *testing.M) {
	conf, err := config.LoadConfig("../test_config.yml")
	if err != nil {
		log.Panic(fmt.Errorf("load config error: %w", err))
	}

	s, err := server.NewServer(conf)
	if err != nil {
		log.Panic(fmt.Errorf("could not start server: %w", err))
	}

	testDB = s.DB

	go func() {
		if err := s.Run(); err != nil {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	code := m.Run()

	os.Exit(code)
}
