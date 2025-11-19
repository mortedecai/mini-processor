package noop_test

import (
	"github.com/censys/scan-takehome/internal/database/models"
	"github.com/censys/scan-takehome/internal/database/noop"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	"github.com/censys/scan-takehome/internal/database"
	_ "github.com/censys/scan-takehome/internal/database/noop"
)

var _ = Describe("Noop", func() {
	var (
		envMap := EnvMap{
			"DATABASE_TYPE":     StringPointer("noop"),
			"DATABASE_HOST":     StringPointer("localhost"),
			"DATABASE_USER":     StringPointer("user"),
			"DATABASE_PASSWORD": StringPointer("password"),
			"DATABASE_PORT":     StringPointer("5432"),
			"DATABASE_NAME":     StringPointer("dbname"),
		}
		restoreMap EnvMap
	)
	BeforeEach(func() {
		restoreMap = envMap.SetupEnv()
	})
	AfterEach(func() {
		restoreMap.SetupEnv()
	})
	It("should register itself in the database package", func() {
		db, err := database.New()
		Expect(err).ToNot(HaveOccurred())
		Expect(db).ToNot(BeNil())
		Expect(db).To(BeAssignableToTypeOf(&noop.DBNoop{}))
		db.Close()
	})
	It("should return nil on upsert regardless of content", func() {
		db, err := database.New()
		Expect(err).ToNot(HaveOccurred())
		Expect(db).ToNot(BeNil())
		Expect(db).To(BeAssignableToTypeOf(&noop.DBNoop{}))
		err = db.Upsert(nil)
		Expect(err).ToNot(HaveOccurred())
		err = db.Upsert(&models.ScanEntry{})
		Expect(err).ToNot(HaveOccurred())
		db.Close()

	})
})
