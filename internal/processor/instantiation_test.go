package processor_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/censys/scan-takehome/internal/processor"
)

const (
	VAR_PUBSUB_HOST = "PUBSUB_EMULATOR_HOST"
)

var _ = Describe("Instantiation (integration test)", func() {
	var (
		envVars     envMap
		restoreVars envMap
	)
	AfterEach(func() {
		restoreVars.SetupEnv()
	})
	Context("invalid configuration", func() {
		BeforeEach(func() {
			// With nil for the host, an error for subscription not existing should be thrown.
			envVars = envMap{VAR_PROJECT_ID: nil, VAR_SUBSCRIPTION_ID: str_ptr("scan-sub"), VAR_TOPIC_ID: str_ptr("scan-topic"), VAR_PUBSUB_HOST: str_ptr("localhost")}
			restoreVars = envVars.SetupEnv()
		})
		It("should return an error with no valid configuration", func() {
			proc, err := processor.New(processor.ConfigFromEnv())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp(`.*'Config\.ProjectID'.* for 'ProjectID' failed on the 'required' tag`))
			Expect(proc).To(BeNil())
		})
	})
	Context("invalid emulator configuration", func() {
		BeforeEach(func() {
			// With nil for the host, an error for subscription not existing should be thrown.
			envVars = envMap{VAR_PROJECT_ID: str_ptr("test-project"), VAR_SUBSCRIPTION_ID: str_ptr("scan-sub"), VAR_TOPIC_ID: str_ptr("scan-topic"), VAR_PUBSUB_HOST: str_ptr("localhost")}
			restoreVars = envVars.SetupEnv()
		})
		It("should return an error with no valid emulator location", func() {
			proc, err := processor.New(processor.ConfigFromEnv())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp(`subscription does not exist`))
			Expect(proc).To(BeNil())
		})
	})
	Context("valid environment configuration with pubsub emulator running", func() {
		BeforeEach(func() {
			envVars = envMap{VAR_PROJECT_ID: str_ptr("test-project"), VAR_SUBSCRIPTION_ID: str_ptr("scan-sub"), VAR_TOPIC_ID: str_ptr("scan-topic"), VAR_PUBSUB_HOST: str_ptr("localhost:8085")}
			restoreVars = envVars.SetupEnv()
		})
		It("should return a valid processor", func() {
			proc, err := processor.New(processor.ConfigFromEnv())
			Expect(err).ToNot(HaveOccurred())
			Expect(proc).ToNot(BeNil())
		})
	})
})
