package rubbernecker_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRubbernecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rubbernecker Suite")
}
