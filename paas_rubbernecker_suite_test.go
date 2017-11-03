package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPaaSRubbernecker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PaaSRubbernecker Suite")
}
