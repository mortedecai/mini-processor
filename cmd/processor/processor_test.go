package main

// Run this test inside the main package to validate the main method currently

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	_ "github.com/censys/scan-takehome/internal/database/noop"
)

var _ = Describe("Processor", func() {
	const (
		VAR_PROJECT_ID      = "PUBSUB_PROJECT_ID"
		VAR_SUBSCRIPTION_ID = "PUBSUB_SUBSCRIPTION_ID"
		VAR_TOPIC_ID        = "TOPIC_ID"
		VAR_DB_HOST         = "DATABASE_HOST"
		VAR_DB_USER         = "DATABASE_USER"
		VAR_DB_PASSWORD     = "DATABASE_PASSWORD"
		VAR_DB_PORT         = "DATABASE_PORT"
		VAR_DB_NAME         = "DATABASE_NAME"
		VAR_DB_TYPE         = "DATABASE_TYPE"
	)
	var (
		restoreMap EnvMap
	)
	AfterEach(func() {
		// This AfterEach will be called after each test block below;
		// There is no need to repeat this call in each context.
		restoreMap.SetupEnv()
	})
	Context("invalid database configuration", func() {
		var (
			envMap = EnvMap{
				VAR_DB_HOST:     nil,
				VAR_DB_USER:     nil,
				VAR_DB_PASSWORD: nil,
				VAR_DB_PORT:     nil,
				VAR_DB_NAME:     nil,
				VAR_DB_TYPE:     StringPointer("noop"),
			}
		)
		BeforeEach(func() {
			restoreMap = envMap.SetupEnv()
		})
		It("should panic due to invalid db configuration", func() {
			Expect(main).To(Panic())
		})
	})
	Context("invalid processor configuration", func() {
		var (
			envMap = EnvMap{
				VAR_DB_HOST:         StringPointer("localhost"),
				VAR_DB_USER:         StringPointer("testUser"),
				VAR_DB_PASSWORD:     StringPointer("testPass"),
				VAR_DB_PORT:         StringPointer("5432"),
				VAR_DB_NAME:         StringPointer("scans"),
				VAR_DB_TYPE:         StringPointer("noop"),
				VAR_PROJECT_ID:      nil,
				VAR_SUBSCRIPTION_ID: nil,
				VAR_TOPIC_ID:        nil,
			}
		)
		BeforeEach(func() {
			restoreMap = envMap.SetupEnv()
		})
		It("should panic due to invalid processor configuration", func() {
			Expect(main).To(Panic())
		})
	})
})
