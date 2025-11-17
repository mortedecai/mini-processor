package processor

import (
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	ProjectID      string `env:"PUBSUB_PROJECT_ID" validate:"required"`
	TopicID        string `env:"PUBSUB_TOPIC_ID" validate:"required"`
	SubscriptionID string `env:"PUBSUB_SUBSCRIPTION_ID" validate:"required"`
}

func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

// ConfigFromEnv returns a configuration object which has been pre-loaded from the environment.
func ConfigFromEnv() *Config {
	cfg := &Config{}
	env.Parse(cfg)
	return cfg
}
