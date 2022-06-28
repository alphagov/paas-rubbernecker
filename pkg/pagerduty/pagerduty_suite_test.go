package pagerduty_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPagerDuty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rubbernecker PagerDuty Suite")
}
