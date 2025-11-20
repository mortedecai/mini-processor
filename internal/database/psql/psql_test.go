package psql_test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	"github.com/censys/scan-takehome/internal/database"
	"github.com/censys/scan-takehome/internal/database/config"
	"github.com/censys/scan-takehome/internal/database/dal"
	"github.com/censys/scan-takehome/internal/database/models"
	"github.com/censys/scan-takehome/internal/database/psql"
)

var _ = Describe("PSQL Database", func() {
	It("should return a psql error if the config is missing params", func() {
		envMap := EnvMap{
			"DATABASE_TYPE":     StringPointer("postgres"),
			"DATABASE_HOST":     nil,
			"DATABASE_USER":     nil,
			"DATABASE_PASSWORD": nil,
			"DATABASE_PORT":     nil,
			"DATABASE_NAME":     nil,
		}
		restoreMap := envMap.SetupEnv()
		defer restoreMap.SetupEnv()
		cfg := config.ConfigFromEnv()

		db, err := psql.New(cfg)
		Expect(err).To(HaveOccurred())
		Expect(db).To(BeNil())
	})
	It("integration testing postgres registration and initialization", func() {
		envMap := EnvMap{
			"DATABASE_TYPE":     StringPointer("postgres"),
			"DATABASE_HOST":     StringPointer("localhost"),
			"DATABASE_USER":     StringPointer("censysTest"),
			"DATABASE_PASSWORD": StringPointer("censysS4mpl3!"),
			"DATABASE_PORT":     StringPointer("5432"),
			"DATABASE_NAME":     StringPointer("censys_data"),
		}
		restoreMap := envMap.SetupEnv()
		defer restoreMap.SetupEnv()
		db, err := database.New()
		Expect(err).ToNot(HaveOccurred())
		Expect(db).ToNot(BeNil())
		cfg := config.ConfigFromEnv()

		expDb, err := psql.New(cfg)
		Expect(db).To(BeAssignableToTypeOf(expDb))
		Expect(err).ToNot(HaveOccurred())
		db.Close()
		expDb.Close()
	})
	// Use an ordered describe here for integration tests that build on each other
	Describe("Integration Testing Upsert", Ordered, func() {
		var (
			envMap = EnvMap{
				"DATABASE_TYPE":     StringPointer("postgres"),
				"DATABASE_HOST":     StringPointer("localhost"),
				"DATABASE_USER":     StringPointer("censysTest"),
				"DATABASE_PASSWORD": StringPointer("censysS4mpl3!"),
				"DATABASE_PORT":     StringPointer("5432"),
				"DATABASE_NAME":     StringPointer("censys_data"),
			}
			restoreMap EnvMap
			db         dal.Scan
			pgxPool    *pgxpool.Pool
			// terminatingErr existing indicates that no subsequent tests can succeed
			terminatingErr error
			ctx            = context.Background()
		)
		const (
			persistedResponse1 = "HTTP/1.1 200 OK"
			persistedResponse2 = "HTTP/1.1 404 Not Found"
		)
		BeforeAll(func() {
			restoreMap = envMap.SetupEnv()
			cfg := config.ConfigFromEnv()
			pgxPool, terminatingErr = pgxpool.New(ctx, cfg.ConnectionString())
			Expect(terminatingErr).ToNot(HaveOccurred())
			db, terminatingErr = database.New()
			// Clean up any existing test data in case the database was not reset
			_, _ = pgxPool.Exec(ctx, `DELETE FROM scan_data;`)
		})
		AfterAll(func() {
			restoreMap.SetupEnv()
			db.Close()
			pgxPool.Close()
		})
		It("should successfully insert an initial scan entry", func() {
			// The databse connection should succeed.
			Expect(terminatingErr).ToNot(HaveOccurred())
			entry := &models.ScanEntry{
				IP:            "192.168.0.1",
				Port:          80,
				Service:       "http",
				ScanTimestamp: 5,
				Response:      "HTTP/1.1 200 OK",
			}
			terminatingErr = db.Upsert(entry)
			Expect(terminatingErr).ToNot(HaveOccurred())
			rows, checkErr := pgxPool.Query(ctx, `SELECT ip, port, service, scan_date, response FROM scan_data WHERE ip=$1 AND port=$2 AND service=$3`, entry.IP, entry.Port, entry.Service)
			Expect(checkErr).ToNot(HaveOccurred())
			defer rows.Close()
			Expect(rows.Next()).To(BeTrue())
			var fetchedEntry models.ScanEntry
			checkErr = rows.Scan(&fetchedEntry.IP, &fetchedEntry.Port, &fetchedEntry.Service, &fetchedEntry.ScanTimestamp, &fetchedEntry.Response)
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(fetchedEntry).To(Equal(*entry))
		})
		It("should not overwrite an existing entry with an older timestamp", func() {
			if terminatingErr != nil {
				Skip("previous test(s) failed or were skipped due to an early error")
			}
			entry := &models.ScanEntry{
				IP:            "192.168.0.1",
				Port:          80,
				Service:       "http",
				ScanTimestamp: 4,
				Response:      "HTTP/1.1 418 I'm a teapot",
			}
			terminatingErr = db.Upsert(entry)
			Expect(terminatingErr).ToNot(HaveOccurred())
			rows, checkErr := pgxPool.Query(ctx, `SELECT ip, port, service, scan_date, response FROM scan_data WHERE ip=$1 AND port=$2 AND service=$3`, entry.IP, entry.Port, entry.Service)
			Expect(checkErr).ToNot(HaveOccurred())
			defer rows.Close()
			Expect(rows.Next()).To(BeTrue())
			var fetchedEntry models.ScanEntry
			checkErr = rows.Scan(&fetchedEntry.IP, &fetchedEntry.Port, &fetchedEntry.Service, &fetchedEntry.ScanTimestamp, &fetchedEntry.Response)
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(fetchedEntry).ToNot(Equal(*entry))
			// Validate the original entry was returned
			entry.Response = persistedResponse1
			entry.ScanTimestamp = 5
			Expect(fetchedEntry).To(Equal(*entry))
		})
		It("should not overwrite an existing entry with an equal timestamp", func() {
			if terminatingErr != nil {
				Skip("previous test(s) failed or were skipped due to an early error")
			}
			entry := &models.ScanEntry{
				IP:            "192.168.0.1",
				Port:          80,
				Service:       "http",
				ScanTimestamp: 5,
				Response:      "HTTP/1.1 418 I'm a teapot",
			}
			terminatingErr = db.Upsert(entry)
			Expect(terminatingErr).ToNot(HaveOccurred())
			rows, checkErr := pgxPool.Query(ctx, `SELECT ip, port, service, scan_date, response FROM scan_data WHERE ip=$1 AND port=$2 AND service=$3`, entry.IP, entry.Port, entry.Service)
			Expect(checkErr).ToNot(HaveOccurred())
			defer rows.Close()
			Expect(rows.Next()).To(BeTrue())
			var fetchedEntry models.ScanEntry
			checkErr = rows.Scan(&fetchedEntry.IP, &fetchedEntry.Port, &fetchedEntry.Service, &fetchedEntry.ScanTimestamp, &fetchedEntry.Response)
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(fetchedEntry).ToNot(Equal(*entry))
			// Validate the original entry was returned
			entry.Response = persistedResponse1
			entry.ScanTimestamp = 5
			Expect(fetchedEntry).To(Equal(*entry))
		})
		It("should overwrite an existing entry with a newer timestamp", func() {
			if terminatingErr != nil {
				Skip("previous test(s) failed or were skipped due to an early error")
			}
			entry := &models.ScanEntry{
				IP:            "192.168.0.1",
				Port:          80,
				Service:       "http",
				ScanTimestamp: 6,
				Response:      persistedResponse2,
			}
			terminatingErr = db.Upsert(entry)
			Expect(terminatingErr).ToNot(HaveOccurred())
			rows, checkErr := pgxPool.Query(ctx, `SELECT ip, port, service, scan_date, response FROM scan_data WHERE ip=$1 AND port=$2 AND service=$3`, entry.IP, entry.Port, entry.Service)
			Expect(checkErr).ToNot(HaveOccurred())
			defer rows.Close()
			Expect(rows.Next()).To(BeTrue())
			var fetchedEntry models.ScanEntry
			checkErr = rows.Scan(&fetchedEntry.IP, &fetchedEntry.Port, &fetchedEntry.Service, &fetchedEntry.ScanTimestamp, &fetchedEntry.Response)
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(fetchedEntry).To(Equal(*entry))
		})
	})
})
