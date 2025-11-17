package processor_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/censys/scan-takehome/internal/processor"
)

const (
	VAR_PROJECT_ID      = "PUBSUB_PROJECT_ID"
	VAR_SUBSCRIPTION_ID = "PUBSUB_SUBSCRIPTION_ID"
	VAR_TOPIC_ID        = "PUBSUB_TOPIC_ID"
)

type envMap map[string]*string

func str_ptr(s string) *string {
	return &s
}

func (e envMap) SetupEnv() envMap {
	restoreMap := envMap{}
	for k, v := range e {
		if ev, found := os.LookupEnv(k); found {
			restoreMap[k] = str_ptr(ev)
		}
		var err error
		if v != nil {
			err = os.Setenv(k, *v)
		} else {
			err = os.Unsetenv(k)
		}
		Expect(err).ToNot(HaveOccurred())
	}
	return restoreMap
}

var _ = Describe("Config", func() {
	DescribeTable("Configuration from environment validation",
		func(vars envMap, expConfig *processor.Config, expErrRegex ...string) {
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
			envMap{VAR_PROJECT_ID: str_ptr("foo"), VAR_SUBSCRIPTION_ID: str_ptr("bar"), VAR_TOPIC_ID: str_ptr("baz")},
			&processor.Config{ProjectID: "foo", SubscriptionID: "bar", TopicID: "baz"},
		),
		Entry(
			"Missing Project ID",
			envMap{VAR_PROJECT_ID: nil, VAR_SUBSCRIPTION_ID: str_ptr("foo"), VAR_TOPIC_ID: str_ptr("baz")},
			&processor.Config{SubscriptionID: "foo", TopicID: "baz"},
			`.*'Config\.ProjectID'.* for 'ProjectID' failed on the 'required' tag`,
		),
		Entry(
			"Missing subscription ID",
			envMap{VAR_PROJECT_ID: str_ptr("foo"), VAR_SUBSCRIPTION_ID: nil, VAR_TOPIC_ID: str_ptr("baz")},
			&processor.Config{ProjectID: "foo", TopicID: "baz"},
			`.*Config\.SubscriptionID.* for 'SubscriptionID' failed on the 'required' tag`,
		),
		Entry(
			"Missing topic ID",
			envMap{VAR_PROJECT_ID: str_ptr("foo"), VAR_SUBSCRIPTION_ID: str_ptr("bar"), VAR_TOPIC_ID: nil},
			&processor.Config{ProjectID: "foo", SubscriptionID: "bar"},
			`.*Config\.TopicID.* for 'TopicID' failed on the 'required' tag`,
		),
		Entry(
			"Missing all required fields",
			envMap{VAR_PROJECT_ID: nil, VAR_SUBSCRIPTION_ID: nil, VAR_TOPIC_ID: nil},
			&processor.Config{},
			`.*Config\.ProjectID.* for 'ProjectID' failed on the 'required' tag`,
			`.*Config\.SubscriptionID.* for 'SubscriptionID' failed on the 'required' tag`,
			`.*Config\.TopicID.* for 'TopicID' failed on the 'required' tag`,
		),
	)
})
