package main

// Run this test inside the main package to validate the main method currently

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {
	const (
		VAR_PROJECT_ID      = "PUBSUB_PROJECT_ID"
		VAR_SUBSCRIPTION_ID = "PUBSUB_SUBSCRIPTION_ID"
		VAR_TOPIC_ID        = "TOPIC_ID"
	)
	var (
		origProject      *string = nil
		origTopic        *string = nil
		origSubscription *string = nil
	)
	BeforeEach(func() {
		if v, found := os.LookupEnv(VAR_PROJECT_ID); found {
			origProject = &v
			os.Unsetenv(VAR_PROJECT_ID)
		}
		if v, found := os.LookupEnv(VAR_TOPIC_ID); found {
			origTopic = &v
			os.Unsetenv(VAR_TOPIC_ID)
		}
		if v, found := os.LookupEnv(VAR_SUBSCRIPTION_ID); found {
			origSubscription = &v
			os.Unsetenv(VAR_SUBSCRIPTION_ID)
		}
	})
	AfterEach(func() {
		if origProject != nil {
			os.Setenv(VAR_PROJECT_ID, *origProject)
		}
		if origTopic != nil {
			os.Setenv(VAR_TOPIC_ID, *origTopic)
		}
		if origSubscription != nil {
			os.Setenv(VAR_SUBSCRIPTION_ID, *origSubscription)
		}
	})
	It("should panic due to invalid configuration", func() {
		Expect(main).To(Panic())
	})
})
