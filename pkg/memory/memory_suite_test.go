package memory_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rubbernecker Memory Engine Suite")
}
