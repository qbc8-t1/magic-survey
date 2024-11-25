package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/internal/server"
)

func main() {
	configPath := flag.String("c", "config.yml", "Path to the configuration file")
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("load config error: %w", err))
	}

	dbContext := context.Background()

	s, err := server.NewServer(dbContext, conf)
	if err != nil {
		log.Fatal(fmt.Errorf("could not start server: %w", err))
	}

	err = s.Run()
	if err != nil {
		log.Fatalln(fmt.Errorf("error while running server: %w", err))
	}
}
