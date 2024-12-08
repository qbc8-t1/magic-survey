package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerMode string

const (
	Development ServerMode = "development"
	Production  ServerMode = "production"
)

type Logger struct {
	Level string
	Path  string
}

type Server struct {
	Host                        string
	Port                        string
	Mode                        ServerMode
	AppVersion                  string
	Secret                      string
	AuthExpMinute               uint `json:"authExpMin"`
	AuthRefreshMinute           uint `json:"authExpRefreshMin"`
	MailPass                    string
	FromMail                    string
	MaxSecondForChangeBirthdate int
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SslMode  string
	Clean    bool
}

type Config struct {
	Server
	Postgres
	Logger
}

func LoadConfig(configPath string) (*Config, error) {
	var config Config
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err := v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
