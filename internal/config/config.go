package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config represents envs from the config.env file.
type Config struct {
	PostgresURL string `envconfig:"POSTGRES_URL"` // PostgresURL is database connection string
	Host        string `envconfig:"CART_HOST"`    // Host is an application IP address
	Port        string `envconfig:"CART_PORT"`    // Port is an application port
}

// NewConfig is a constructor for Config struct.
func NewConfig() (*Config, error) {
	var c Config

	err := envconfig.Process("cart", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
