package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	"github.com/censys/scan-takehome/internal/database/config"
)

const (
	VAR_DB_HOST     = "DATABASE_HOST"
	VAR_DB_USER     = "DATABASE_USER"
	VAR_DB_PASSWORD = "DATABASE_PASSWORD"
	VAR_DB_PORT     = "DATABASE_PORT"
	VAR_DB_NAME     = "DATABASE_NAME"
)

var (
	hostPtr = StringPointer("localhost")
	userPtr = StringPointer("user")
	passPtr = StringPointer("password")
	portPtr = StringPointer("5432")
	namePtr = StringPointer("dbname")
)

var _ = DescribeTable("Config validation testing",
	func(envMap EnvMap, expectedErrRegex ...string) {
		restoreMap := envMap.SetupEnv()
		defer restoreMap.SetupEnv()
		err := config.ConfigFromEnv().Validate()
		if len(expectedErrRegex) == 0 {
			Expect(err).ToNot(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
			for _, regex := range expectedErrRegex {
				Expect(err.Error()).To(MatchRegexp(regex))
			}
		}
	},
	Entry("valid config",
		EnvMap{
			VAR_DB_HOST:     hostPtr,
			VAR_DB_USER:     userPtr,
			VAR_DB_PASSWORD: passPtr,
			VAR_DB_PORT:     portPtr,
			VAR_DB_NAME:     namePtr,
		},
	),
	Entry("missing host",
		EnvMap{
			VAR_DB_HOST:     nil,
			VAR_DB_USER:     userPtr,
			VAR_DB_PASSWORD: passPtr,
			VAR_DB_PORT:     portPtr,
			VAR_DB_NAME:     namePtr,
		},
		`.*'Config\.Host'.*Field validation for 'Host' failed on the 'required' tag`,
	),
	Entry("missing user",
		EnvMap{
			VAR_DB_HOST:     hostPtr,
			VAR_DB_USER:     nil,
			VAR_DB_PASSWORD: passPtr,
			VAR_DB_PORT:     portPtr,
			VAR_DB_NAME:     namePtr,
		},
		`.*'Config\.User'.*Field validation for 'User' failed on the 'required' tag`,
	),
	Entry("missing password",
		EnvMap{
			VAR_DB_HOST:     hostPtr,
			VAR_DB_USER:     userPtr,
			VAR_DB_PASSWORD: nil,
			VAR_DB_PORT:     portPtr,
			VAR_DB_NAME:     namePtr,
		},
		`.*'Config\.Pass'.*Field validation for 'Pass' failed on the 'required' tag`,
	),
	Entry("missing port",
		EnvMap{
			VAR_DB_HOST:     hostPtr,
			VAR_DB_USER:     userPtr,
			VAR_DB_PASSWORD: passPtr,
			VAR_DB_PORT:     nil,
			VAR_DB_NAME:     namePtr,
		},
		`.*'Config\.Port'.*Field validation for 'Port' failed on the 'required' tag`,
	),
	Entry("invalid port",
		EnvMap{
			VAR_DB_HOST:     hostPtr,
			VAR_DB_USER:     userPtr,
			VAR_DB_PASSWORD: passPtr,
			VAR_DB_PORT:     StringPointer("70000"),
			VAR_DB_NAME:     namePtr,
		},
		`.*'Config\.Port'.*Field validation for 'Port' failed on the 'port' tag`,
	),
	Entry("missing db name",
		EnvMap{
			VAR_DB_HOST:     hostPtr,
			VAR_DB_USER:     userPtr,
			VAR_DB_PASSWORD: passPtr,
			VAR_DB_PORT:     portPtr,
			VAR_DB_NAME:     nil,
		},
		`.*'Config\.DBName'.*Field validation for 'DBName' failed on the 'required' tag`,
	),
)
