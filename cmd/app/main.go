package main

import (
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
		log.Panic(fmt.Errorf("load config error: %w", err))
	}

	s, err := server.NewServer(conf)
	if err != nil {
		log.Panic(fmt.Errorf("could not start server: %w", err))
	}

	err = s.Run()
	if err != nil {
		log.Panic(fmt.Errorf("error while running server: %w", err))
	}
}
