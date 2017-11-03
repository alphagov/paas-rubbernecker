package pivotal_test

import (
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-rubbernecker/pkg/pivotal"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Pivotal Members", func() {
	Context("Tracker setup", func() {
		var (
			pt rubbernecker.MemberService

			apiURL   = `https://www.pivotaltracker.com/services/v5/projects/123/memberships`
			response = `[{"role":"viewer","person":{"id":654321,"name":"non-tester"}},{"role":"owner","person":{"id":123456,"name":"tester"}}]`
		)

		BeforeEach(func() {
			var err error

			pt, err = pivotal.New(123, "test")
			httpmock.Activate()

			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail to FetchMembers() stories from an API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(404, ``))

			err := pt.FetchMembers()

			Expect(err).To(HaveOccurred())
		})

		It("should FetchMembers() stories from an API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err := pt.FetchMembers()

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail to FlattenMembers() due to faulty API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, `[]`))

			err := pt.FetchMembers()

			Expect(err).NotTo(HaveOccurred())

			members, err := pt.FlattenMembers()

			Expect(err).To(HaveOccurred())
			Expect(members).To(BeNil())
		})

		It("should FlattenMembers() correctly", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err := pt.FetchMembers()

			Expect(err).NotTo(HaveOccurred())

			members, err := pt.FlattenMembers()
			m := *members

			Expect(err).NotTo(HaveOccurred())
			Expect(len(m)).To(Equal(1))
			Expect(m[123456].Name).To(Equal("tester"))
		})
	})
})
