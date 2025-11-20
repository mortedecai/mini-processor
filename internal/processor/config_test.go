package processor_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/censys/scan-takehome/_test"
	"github.com/censys/scan-takehome/internal/processor"
)

const (
	VAR_PROJECT_ID      = "PUBSUB_PROJECT_ID"
	VAR_SUBSCRIPTION_ID = "PUBSUB_SUBSCRIPTION_ID"
	VAR_TOPIC_ID        = "PUBSUB_TOPIC_ID"
)

var _ = Describe("Config", func() {
	DescribeTable("Configuration from environment validation",
		func(vars EnvMap, expConfig *processor.Config, expErrRegex ...string) {
			restoreMap := vars.SetupEnv()
			defer restoreMap.SetupEnv()

			cfg := processor.ConfigFromEnv()
			err := cfg.Validate()

			Expect(cfg).To(Equal(expConfig))
			if len(expErrRegex) == 0 {
				Expect(err).ToNot(HaveOccurred())
				return
			}
			Expect(err).To(HaveOccurred())
			for _, v := range expErrRegex {
				Expect(err.Error()).To(MatchRegexp(v))
			}
		},
		Entry(
			"All config items valid",
			EnvMap{VAR_PROJECT_ID: StringPointer("foo"), VAR_SUBSCRIPTION_ID: StringPointer("bar"), VAR_TOPIC_ID: StringPointer("baz")},
			&processor.Config{ProjectID: "foo", SubscriptionID: "bar", TopicID: "baz"},
		),
		Entry(
			"Missing Project ID",
			EnvMap{VAR_PROJECT_ID: nil, VAR_SUBSCRIPTION_ID: StringPointer("foo"), VAR_TOPIC_ID: StringPointer("baz")},
			&processor.Config{SubscriptionID: "foo", TopicID: "baz"},
			`.*'Config\.ProjectID'.* for 'ProjectID' failed on the 'required' tag`,
		),
		Entry(
			"Missing subscription ID",
			EnvMap{VAR_PROJECT_ID: StringPointer("foo"), VAR_SUBSCRIPTION_ID: nil, VAR_TOPIC_ID: StringPointer("baz")},
			&processor.Config{ProjectID: "foo", TopicID: "baz"},
			`.*Config\.SubscriptionID.* for 'SubscriptionID' failed on the 'required' tag`,
		),
		Entry(
			"Missing topic ID",
			EnvMap{VAR_PROJECT_ID: StringPointer("foo"), VAR_SUBSCRIPTION_ID: StringPointer("bar"), VAR_TOPIC_ID: nil},
			&processor.Config{ProjectID: "foo", SubscriptionID: "bar"},
			`.*Config\.TopicID.* for 'TopicID' failed on the 'required' tag`,
		),
		Entry(
			"Missing all required fields",
			EnvMap{VAR_PROJECT_ID: nil, VAR_SUBSCRIPTION_ID: nil, VAR_TOPIC_ID: nil},
			&processor.Config{},
			`.*Config\.ProjectID.* for 'ProjectID' failed on the 'required' tag`,
			`.*Config\.SubscriptionID.* for 'SubscriptionID' failed on the 'required' tag`,
			`.*Config\.TopicID.* for 'TopicID' failed on the 'required' tag`,
		),
	)
})
