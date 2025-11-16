package main

// Run this test inside the main package to validate the main method currently

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Processor", func() {
	var (
		origStdOut *os.File
		testLog    *os.File
	)
	BeforeEach(func() {
		origStdOut = os.Stdout
		if logFile, err := os.CreateTemp("", "processor_*.log"); err == nil {
			testLog = logFile
			os.Stdout = testLog
		} else {
			// Unrecoverable error
			Fail(fmt.Sprintf("failed to create temporary log file: %s\n", err.Error()))
		}
	})
	AfterEach(func() {
		os.Stdout = origStdOut
		// Nothing we can do if we can't close the test log file
		_ = testLog.Close()
	})
	It("should print a log message and exit", func() {
		main()
		data, err := os.ReadFile(testLog.Name())
		Expect(err).ToNot(HaveOccurred())
		// add a line return for the fmt.Println usage.
		Expect(string(data)).To(Equal(msg + "\n"))
	})
})
