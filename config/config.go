package config

import (
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = ""

type Configuration struct {
	HTTPServer
	Database
}

type HTTPServer struct {
	Port int `envconfig:"PORT" default:"3000"`
}

type Database struct {
	DatabaseURL    string `envconfig:"DATABASE_URL" required:"true"`
	DatabaseName   string `envconfig:"DATABASE_NAME" default:"EWallet"`
	CollectionName string `envconfig:"COLLECTION_NAME" default:"wallets"`
}

func Load() (Configuration, error) {
	var cfg Configuration
	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
