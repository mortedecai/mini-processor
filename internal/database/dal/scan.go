package dal

import (
	"github.com/censys/scan-takehome/internal/database/models"
)

// Scan represents the actions which can be taken on the Scan database
type Scan interface {
	Upsert(entry *models.ScanEntry) error

	Close()
}
