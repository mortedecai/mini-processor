package processor_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/censys/scan-takehome/internal/processor"
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

			cfg, err := processor.ConfigFromEnv()
			if expConfig == nil {
				Expect(err).To(HaveOccurred())
				for _, v := range expErrRegex {
					Expect(err.Error()).To(MatchRegexp(v))
				}
				Expect(cfg).To(BeNil())
				return
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).To(Equal(expConfig))
		},
		Entry(
			"All config items valid",
			envMap{"PUBSUB_PROJECT_ID": str_ptr("foo"), "PUBSUB_SUBSCRIPTION_ID": str_ptr("bar")},
			&processor.Config{ProjectID: "foo", SubscriptionID: "bar"},
		),
		Entry(
			"Missing Project ID",
			envMap{"PUBSUB_SUBSCRIPTION_ID": str_ptr("foo")},
			nil,
			`.*'Config\.ProjectID'.* for 'ProjectID' failed on the 'required' tag`,
		),
		Entry(
			"Missing subscription ID",
			envMap{"PUBSUB_PROJECT_ID": str_ptr("foo")},
			nil,
			`.*Config\.SubscriptionID.* for 'SubscriptionID' failed on the 'required' tag`,
		),
		Entry(
			"Missing all required fields",
			envMap{},
			nil,
			`.*Config\.ProjectID.* for 'ProjectID' failed on the 'required' tag`,
			`.*Config\.SubscriptionID.* for 'SubscriptionID' failed on the 'required' tag`,
		),
	)
})
