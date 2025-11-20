package psql

import (
	"context"

	"github.com/censys/scan-takehome/internal/database"
	"github.com/censys/scan-takehome/internal/database/config"
	"github.com/censys/scan-takehome/internal/database/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/censys/scan-takehome/internal/database/dal"
)

const (
	InsertStmt     = "INSERT INTO scan_data(ip, port, service, scan_date, response) VALUES ($1, $2, $3, $4, $5)"
	OnConflictStmt = "ON CONFLICT (ip, port, service) DO UPDATE SET scan_date = EXCLUDED.scan_date, response = EXCLUDED.response WHERE scan_data.ip = EXCLUDED.ip AND scan_data.port = EXCLUDED.port AND scan_data.service = EXCLUDED.service AND scan_data.scan_date < EXCLUDED.scan_date"
	UpsertStmt     = InsertStmt + " " + OnConflictStmt
)

func init() {
	database.RegisterDB("postgres", New)
}

type psqlDB struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Config) (dal.Scan, error) {
	db := &psqlDB{}
	var err error
	if db.pool, err = pgxpool.New(context.Background(), cfg.ConnectionString()); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *psqlDB) Close() {
	db.pool.Close()
}

func (db *psqlDB) Upsert(entry *models.ScanEntry) error {
	// We really don't need a transaction here, but using one to keep the code extensible for future changes
	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		zap.S().Errorw("failed to begin transaction", "error", err, "entry", entry)
		return err
	}
	// The UpsertStmt uses an ON CONFLICT setup to overwrite existing entries only if the new scan_date is more recent
	_, err = tx.Exec(context.Background(), UpsertStmt, entry.IP, entry.Port, entry.Service, entry.ScanTimestamp, entry.Response)
	if err != nil {
		zap.S().Errorw("failed to upsert scan entry", "error", err, "entry", entry)
		tx.Rollback(context.Background())
		return err
	}
	err = tx.Commit(context.Background())
	return err
}
