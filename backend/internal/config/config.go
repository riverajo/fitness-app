package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port      string `env:"PORT" envDefault:"8080"`
	AppEnv    string `env:"APP_ENV" envDefault:"development"`
	MongoURI  string `env:"MONGO_URI,required"`
	JWTSecret string `env:"JWT_SECRET,required"`
	FaroURL   string `env:"FARO_URL" envDefault:"http://alloy:12347/collect"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}
	return cfg, nil
}
