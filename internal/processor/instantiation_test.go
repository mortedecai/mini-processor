package processor_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/censys/scan-takehome/internal/processor"
)

var _ = Describe("Instantiation", func() {
	It("should return a not yet implemented error before implementation", func() {
		proc, err := processor.New()
		Expect(proc).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(errors.New("not yet implemented")))
	})
})
