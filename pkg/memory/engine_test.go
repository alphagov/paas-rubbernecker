package memory_test

import (
	"github.com/alphagov/paas-rubbernecker/pkg/memory"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memory Engine", func() {
	Context("Engine is not setup", func() {
		It("should successfully SetupEngine()", func() {
			me := memory.SetupEngine()

			Expect(me).NotTo(BeNil())
		})
	})

	var (
		me       rubbernecker.PersistanceEngine = memory.SetupEngine()
		arn                                     = "rubbernecker.pkg.memory.engine.test"
		arnValue                                = 123
	)

	It("should Put() value successfully", func() {
		err := me.Put(arn, arnValue)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should Get() value successfully", func() {
		value, err := me.Get(arn)

		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(arnValue))
	})

	It("should fail to Get() value", func() {
		value, err := me.Get("arn")

		Expect(err).To(HaveOccurred())
		Expect(value).To(BeNil())
	})
})
