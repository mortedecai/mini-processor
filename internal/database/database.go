package database

import (
	"fmt"

	"github.com/censys/scan-takehome/internal/database/config"
	"github.com/censys/scan-takehome/internal/database/dal"
)

type dbInitializers map[string]func(cfg *config.Config) (dal.Scan, error)

var (
	initializers = dbInitializers{}
)

// RegisterDB registers a database initializer for a given database type.
// The database type (DATABASE_TYPE environment variable) is used to look up the appropriate initializer function.
// Additionally it will be used as part of the connection string (e.g. postgres:// if DATABASE_TYPE=postgres).
func RegisterDB(dbType string, initializer func(cfg *config.Config) (dal.Scan, error)) {
	initializers[dbType] = initializer
}

func New() (dal.Scan, error) {
	var found bool

	cfg := config.ConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	initializer, found := initializers[cfg.DBType]
	if !found {
		return nil, fmt.Errorf("unknown database type: %s", cfg.DBType)
	}

	return initializer(cfg)
}
