package config

import (
	"os"
)

type Config struct {
	PostgresDSN string
	ExternalUrl string
	Port        string
}

func ReadConfig() *Config {
	return &Config{
		PostgresDSN: os.Getenv("POSTGRESDSN"),
		ExternalUrl: os.Getenv("EXTERNALURL"),
		Port:        os.Getenv("PORT"),
	}
}
