package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const envFile = ".env"

var config Config

type Config struct {
	Port          int
	IsProduction  bool
	PostgresURL   string
	JWTSecret     string
	TesterApiIp   string
	TesterApiPort int
}

func LoadConfig() error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}
	return envconfig.Process("webserver", &config)
}

func GetInstance() Config {
	return config
}
