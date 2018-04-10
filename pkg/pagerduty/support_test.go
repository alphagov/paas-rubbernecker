package pagerduty_test

import (
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-rubbernecker/pkg/pagerduty"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("PagerDuty", func() {
	Context("Schedule not setup", func() {
		It("should create a New() schedule", func() {
			pd := pagerduty.New("test")

			Expect(pd).NotTo(BeNil())
		})
	})

	Context("Schedule setup", func() {
		var (
			pd rubbernecker.SupportService

			apiURL   = `https://api.pagerduty.com/oncalls`
			response = `{"oncalls":[{"user":{"summary":"tester"},"schedule":{"summary":"test"}},{"user":{"summary":"tester"}}]}`
		)

		BeforeEach(func() {
			pd = pagerduty.New("test")
			httpmock.Activate()

			Expect(pd).NotTo(BeNil())
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail to FetchSupport() from an API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(404, ``))

			err := pd.FetchSupport()

			Expect(err).To(HaveOccurred())
		})

		It("should FetchSupport() from an API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err := pd.FetchSupport()

			Expect(err).NotTo(HaveOccurred())
		})

		It("should FlattenStories() correctly", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err := pd.FetchSupport()

			Expect(err).NotTo(HaveOccurred())

			support, err := pd.FlattenSupport()
			s := support

			Expect(err).NotTo(HaveOccurred())
			Expect(len(s)).To(Equal(1))
			Expect(s["test"].Member).To(Equal("tester"))
		})
	})
})
