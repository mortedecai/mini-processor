package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
)

// Config holds the configuration database for a postgres connection.
// Note: This is a minimal configuration setup due to time constraints.
type Config struct {
	// DBType indicates the type of database (e.g., "postgres").
	// This field is included for extensibility to support multiple database types in the future.
	// It will be used directly as part of the connection string (e.g. postgres:// if DATABASE_TYPE=postgres).
	DBType string `env:"DATABASE_TYPE" validate:"required"`
	Host   string `env:"DATABASE_HOST" validate:"required"`
	User   string `env:"DATABASE_USER" validate:"required"`
	Pass   string `env:"DATABASE_PASSWORD" validate:"required"`
	Port   uint   `env:"DATABASE_PORT" validate:"required,port"`
	DBName string `env:"DATABASE_NAME" validate:"required"`
}

// ConfigFromEnv configures the postgres from the environment.
func ConfigFromEnv() *Config {
	cfg := &Config{}
	env.Parse(cfg)
	return cfg
}

func (c *Config) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func (c *Config) ConnectionString() string {
	// Normally sslmode would be configurable, but for this take home assignment we will disable it.
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable", c.DBType, c.User, c.Pass, c.Host, c.Port, c.DBName)
}
