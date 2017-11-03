package rubbernecker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRubbernecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rubbernecker Suite")
}
