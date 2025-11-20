package database_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	"github.com/censys/scan-takehome/internal/database"
	"github.com/censys/scan-takehome/internal/database/dal"
	"github.com/censys/scan-takehome/internal/database/noop"
)

var _ = Describe("Database Initialization", func() {
	DescribeTable("Database initialization testing",
		func(envMap EnvMap, expectedDB dal.Scan, expectedErrRegex ...string) {
			restoreMap := envMap.SetupEnv()
			defer restoreMap.SetupEnv()
			db, err := database.New()
			if len(expectedErrRegex) == 0 {
				Expect(err).ToNot(HaveOccurred())
				Expect(db).ToNot(BeNil())
				Expect(db).To(BeAssignableToTypeOf(expectedDB))
			} else {
				Expect(err).To(HaveOccurred())
				for _, regex := range expectedErrRegex {
					Expect(err.Error()).To(MatchRegexp(regex))
				}
			}
		},
		Entry("Database type not set", EnvMap{"DATABASE_TYPE": nil}, nil, `DATABASE_TYPE environment variable not set`),
		Entry("Unknown database type", EnvMap{"DATABASE_TYPE": StringPointer("unknown_db")}, nil, `unknown database type: unknown_db`),
		Entry("Invalid database config", EnvMap{
			"DATABASE_TYPE":     StringPointer("noop"),
			"DATABASE_HOST":     nil,
			"DATABASE_USER":     StringPointer("user"),
			"DATABASE_PASSWORD": StringPointer("password"),
			"DATABASE_PORT":     StringPointer("5432"),
			"DATABASE_NAME":     StringPointer("dbname"),
		}, nil, `invalid database configuration: .*'Config\.Host'.*Field validation for 'Host' failed on the 'required' tag`),
		Entry("Valid noop database config", EnvMap{
			"DATABASE_TYPE":     StringPointer("noop"),
			"DATABASE_HOST":     StringPointer("localhost"),
			"DATABASE_USER":     StringPointer("user"),
			"DATABASE_PASSWORD": StringPointer("password"),
			"DATABASE_PORT":     StringPointer("5432"),
			"DATABASE_NAME":     StringPointer("dbname"),
		}, func() dal.Scan { db, _ := noop.New(nil); return db }()),
	)
})
