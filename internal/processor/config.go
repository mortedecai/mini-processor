package processor

import (
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	ProjectID      string `env:"PUBSUB_PROJECT_ID" validate:"required"`
	SubscriptionID string `env:"PUBSUB_SUBSCRIPTION_ID" validate:"required"`
}

func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func ConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
