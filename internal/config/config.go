package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

const envFile = "configs/.env"

var config Config

func init() {
	var err error
	if config, err = newConfig(); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	TemplatesDir    string
	PublicDir       string
	SvelteIndexPath string
	SveltePublicDir string
	Port            string
	PostgresURL     string
	JWTSecret       string
	IsProduction    bool
}

func newConfig() (Config, error) {
	res := Config{}
	err := godotenv.Load(envFile)
	if err != nil {
		return res, err
	}
	err = envconfig.Process("webserver", &res)
	return res, err
}

func GetInstance() Config {
	return config
}
