package noop

import (
	"github.com/censys/scan-takehome/internal/database"
	"github.com/censys/scan-takehome/internal/database/config"
	"github.com/censys/scan-takehome/internal/database/dal"
	"github.com/censys/scan-takehome/internal/database/models"
	"go.uber.org/zap"
)

const (
	DB_NOOP = "noop"
)

func init() {
	database.RegisterDB(DB_NOOP, New)
}

// DBNoop is a no-operation database implementation that satisfies the dal.Scan interface.
type DBNoop struct{}

func (db *DBNoop) Close()                           {}
func (db *DBNoop) Upsert(_ *models.ScanEntry) error { return nil }

// New creates a new instance of the noop database.
func New(_ *config.Config) (dal.Scan, error) {
	zap.S().Warn("using noop database, no data will be stored")
	return &DBNoop{}, nil
}
