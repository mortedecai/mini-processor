package psql_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPsql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Psql Suite")
}
