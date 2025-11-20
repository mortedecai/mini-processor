package processor_test

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	"github.com/censys/scan-takehome/internal/database"
	"github.com/censys/scan-takehome/internal/database/config"
	"github.com/censys/scan-takehome/internal/database/models"
	_ "github.com/censys/scan-takehome/internal/database/noop"
	_ "github.com/censys/scan-takehome/internal/database/psql"
	"github.com/censys/scan-takehome/internal/processor"
	"github.com/censys/scan-takehome/pkg/scanning"
)

const (
	VAR_PUBSUB_HOST = "PUBSUB_EMULATOR_HOST"
)

var _ = Describe("Instantiation (integration test)", func() {
	var (
		envVars     EnvMap
		restoreVars EnvMap
	)
	AfterEach(func() {
		restoreVars.SetupEnv()
	})
	Context("invalid configuration", func() {
		BeforeEach(func() {
			// With nil for the host, an error for subscription not existing should be thrown.
			envVars = EnvMap{VAR_PROJECT_ID: nil, VAR_SUBSCRIPTION_ID: StringPointer("scan-sub"), VAR_TOPIC_ID: StringPointer("scan-topic"), VAR_PUBSUB_HOST: StringPointer("localhost")}
			restoreVars = envVars.SetupEnv()
		})
		It("should return an error with no valid configuration", func() {
			// Should error out before the db is needed so passing nil for this test
			proc, err := processor.New(processor.ConfigFromEnv(), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp(`.*'Config\.ProjectID'.* for 'ProjectID' failed on the 'required' tag`))
			Expect(proc).To(BeNil())
		})
	})
	Context("invalid emulator configuration", func() {
		BeforeEach(func() {
			// With nil for the host, an error for subscription not existing should be thrown.
			envVars = EnvMap{VAR_PROJECT_ID: StringPointer("test-project"), VAR_SUBSCRIPTION_ID: StringPointer("scan-sub"), VAR_TOPIC_ID: StringPointer("scan-topic"), VAR_PUBSUB_HOST: StringPointer("localhost")}
			restoreVars = envVars.SetupEnv()
		})
		It("should return an error with no valid emulator location", func() {
			// Should error out before the db is needed so passing nil for this test
			proc, err := processor.New(processor.ConfigFromEnv(), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp(`subscription does not exist`))
			Expect(proc).To(BeNil())
		})
	})
	Context("valid environment configuration with pubsub emulator running", func() {
		BeforeEach(func() {
			envVars = EnvMap{}
			envVars[VAR_PROJECT_ID] = StringPointer("test-project")
			envVars[VAR_SUBSCRIPTION_ID] = StringPointer("scan-sub")
			envVars[VAR_TOPIC_ID] = StringPointer("scan-topic")
			envVars[VAR_PUBSUB_HOST] = StringPointer("localhost:8085")
			restoreVars = envVars.SetupEnv()
		})
		It("should return a valid processor", func() {
			// Processing will not be tested here, so passing nil for db is acceptable
			proc, err := processor.New(processor.ConfigFromEnv(), nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(proc).ToNot(BeNil())
		})
	})
	Describe("Processor Start / Stop Testing", func() {
		var (
			envMap = EnvMap{
				VAR_PROJECT_ID:      StringPointer("test-project"),
				VAR_SUBSCRIPTION_ID: StringPointer("scan-sub"),
				VAR_TOPIC_ID:        StringPointer("scan-topic"),
				VAR_PUBSUB_HOST:     StringPointer("localhost:8085"),
				"DATABASE_TYPE":     StringPointer("noop"),
				"DATABASE_HOST":     StringPointer("localhost"),
				"DATABASE_USER":     StringPointer("censysTest"),
				"DATABASE_PASSWORD": StringPointer("censysS4mpl3!"),
				"DATABASE_PORT":     StringPointer("5432"),
				"DATABASE_NAME":     StringPointer("censys_data"),
			}
			restoreVars EnvMap
		)
		BeforeEach(func() {
			restoreVars = envMap.SetupEnv()
		})
		AfterEach(func() {
			restoreVars.SetupEnv()
		})
		It("should start and stop the processor without error", func() {
			db, err := database.New()
			Expect(err).ToNot(HaveOccurred())
			proc, err := processor.New(processor.ConfigFromEnv(), db)
			Expect(err).ToNot(HaveOccurred())
			Expect(proc).ToNot(BeNil())
			go proc.Start()
			// Inject a 5 second delay to allow for startup
			time.Sleep(5 * time.Second)
			proc.Stop()
		})
	})
	Describe("Integration Testing For Processing", func() {
		var (
			envVars = EnvMap{
				VAR_PROJECT_ID:      StringPointer("test-project"),
				VAR_SUBSCRIPTION_ID: StringPointer("scan-sub"),
				VAR_TOPIC_ID:        StringPointer("scan-topic"),
				VAR_PUBSUB_HOST:     StringPointer("localhost:8085"),
				"DATABASE_TYPE":     StringPointer("postgres"),
				"DATABASE_HOST":     StringPointer("localhost"),
				"DATABASE_USER":     StringPointer("censysTest"),
				"DATABASE_PASSWORD": StringPointer("censysS4mpl3!"),
				"DATABASE_PORT":     StringPointer("5432"),
				"DATABASE_NAME":     StringPointer("censys_data"),
			}
			restoreMap     EnvMap
			scan           scanning.Scan
			pgxPool        *pgxpool.Pool
			terminatingErr error
			ctx            = context.Background()
			v1data         = scanning.V1Data{ResponseBytesUtf8: []byte("v1 response")}
			v2data         = scanning.V2Data{ResponseStr: "v2 response"}
		)
		BeforeEach(func() {
			restoreMap = envVars.SetupEnv()
			scan = scanning.Scan{
				Ip:        "192.168.1.1",
				Port:      80,
				Service:   "http",
				Timestamp: 1625247600,
			}
			cfg := config.ConfigFromEnv()
			pgxPool, terminatingErr = pgxpool.New(ctx, cfg.ConnectionString())
		})
		AfterEach(func() {
			restoreMap.SetupEnv()
			pgxPool.Exec(ctx, `DELETE FROM scan_data;`)
			pgxPool.Close()
		})
		It("should handle a V1Data scan message", func() {
			Expect(terminatingErr).ToNot(HaveOccurred())
			db, err := database.New()
			Expect(err).ToNot(HaveOccurred())
			proc, err := processor.New(processor.ConfigFromEnv(), db)
			// This test requires a running Pub/Sub emulator with a topic "scan-topic" and subscription "scan-sub"
			scan.DataVersion = scanning.V1
			scan.Data = &v1data
			data, err := json.Marshal(scan)
			Expect(err).ToNot(HaveOccurred())
			msg := &pubsub.Message{
				Data: data,
			}
			entry, err := models.NewScanEntry(scan)
			Expect(err).ToNot(HaveOccurred())

			proc.HandleMessage(context.Background(), msg)
			rows, checkErr := pgxPool.Query(ctx, `SELECT ip, port, service, scan_date, response FROM scan_data WHERE ip=$1 AND port=$2 AND service=$3`, entry.IP, entry.Port, entry.Service)
			Expect(checkErr).ToNot(HaveOccurred())
			defer rows.Close()
			Expect(rows.Next()).To(BeTrue())
			var fetchedEntry models.ScanEntry
			checkErr = rows.Scan(&fetchedEntry.IP, &fetchedEntry.Port, &fetchedEntry.Service, &fetchedEntry.ScanTimestamp, &fetchedEntry.Response)
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(fetchedEntry).To(Equal(*entry))
		})
		It("should handle a V2Data scan message", func() {
			Expect(terminatingErr).ToNot(HaveOccurred())
			db, err := database.New()
			Expect(err).ToNot(HaveOccurred())
			proc, err := processor.New(processor.ConfigFromEnv(), db)
			// This test requires a running Pub/Sub emulator with a topic "scan-topic" and subscription "scan-sub"
			scan.DataVersion = scanning.V2
			scan.Data = &v2data
			data, err := json.Marshal(scan)
			Expect(err).ToNot(HaveOccurred())
			msg := &pubsub.Message{
				Data: data,
			}
			entry, err := models.NewScanEntry(scan)
			Expect(err).ToNot(HaveOccurred())

			proc.HandleMessage(context.Background(), msg)
			rows, checkErr := pgxPool.Query(ctx, `SELECT ip, port, service, scan_date, response FROM scan_data WHERE ip=$1 AND port=$2 AND service=$3`, entry.IP, entry.Port, entry.Service)
			Expect(checkErr).ToNot(HaveOccurred())
			defer rows.Close()
			Expect(rows.Next()).To(BeTrue())
			var fetchedEntry models.ScanEntry
			checkErr = rows.Scan(&fetchedEntry.IP, &fetchedEntry.Port, &fetchedEntry.Service, &fetchedEntry.ScanTimestamp, &fetchedEntry.Response)
			Expect(checkErr).ToNot(HaveOccurred())
			Expect(fetchedEntry).To(Equal(*entry))
			proc.HandleMessage(context.Background(), msg)
		})
	})
})
