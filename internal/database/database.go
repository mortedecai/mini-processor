package database

import (
	"fmt"
	"os"

	"github.com/censys/scan-takehome/internal/database/config"
	"github.com/censys/scan-takehome/internal/database/dal"
)

type dbInitializers map[string]func(cfg *config.Config) (dal.Scan, error)

var (
	initializers = dbInitializers{}
)

func RegisterDB(dbType string, initializer func(cfg *config.Config) (dal.Scan, error)) {
	initializers[dbType] = initializer
}

func New() (dal.Scan, error) {
	var dbType string
	var found bool

	if dbType, found = os.LookupEnv("DATABASE_TYPE"); !found {
		return nil, fmt.Errorf("DATABASE_TYPE environment variable not set")
	}

	initializer, found := initializers[dbType]
	if !found {
		return nil, fmt.Errorf("unknown database type: %s", dbType)
	}

	cfg := config.ConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	return initializer(cfg)
}
