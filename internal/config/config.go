package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresURL string `envconfig:"POSTGRES_URL"`
	Host        string `envconfig:"CART_HOST"`
	Port        string `envconfig:"CART_PORT"`
}

func NewConfig() (*Config, error) {
	var c Config

	err := envconfig.Process("cart", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
