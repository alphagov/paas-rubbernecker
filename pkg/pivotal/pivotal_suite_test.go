package pivotal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPivotal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rubbernecker PivotalTracker Extension Suite")
}
