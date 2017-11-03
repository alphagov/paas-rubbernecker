package redis_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRedisEngine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rubbernecker Redis Engine Suite")
}
