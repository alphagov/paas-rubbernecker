package pivotal_test

import (
	"fmt"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/alphagov/paas-rubbernecker/pkg/pivotal"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Pivotal Stories", func() {
	Context("Tracker not setup", func() {
		It("should create a New() tracker", func() {
			pt, err := pivotal.New(123, "test")

			Expect(err).NotTo(HaveOccurred())
			Expect(pt).NotTo(BeNil())
		})
	})

	Context("Tracker setup", func() {
		var (
			pt rubbernecker.ProjectManagementService

			apiURL   = `https://www.pivotaltracker.com/services/v5/projects/123/stories?fields=owner_ids,blockers,transitions,current_state,labels,name,url,created_at,story_type,estimate&filter=state:started`
			response = `[{"blockers": [{"name":1234}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[{"name":"test"}]}]`
		)

		BeforeEach(func() {
			var err error

			pt, err = pivotal.New(123, "test")
			httpmock.Activate()

			Expect(err).NotTo(HaveOccurred())

			pt.AcceptStickers(rubbernecker.Stickers{
				rubbernecker.Sticker{
					Name: "test",
				},
				rubbernecker.Sticker{
					Name: "blocked",
				},
				rubbernecker.Sticker{
					Name: "scheduled",
				},
				rubbernecker.Sticker{
					Name: "knowledge-share",
				},
				rubbernecker.Sticker{
					Name: "zero-points",
				},
			})
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail to FetchCards() stories from an API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(404, ``))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})

			Expect(err).To(HaveOccurred())
		})

		It("should FetchCards() stories from an API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail to FlattenStories() due to faulty API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, `[]`))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})

			Expect(err).NotTo(HaveOccurred())

			cards, err := pt.FlattenStories()

			Expect(err).To(HaveOccurred())
			Expect(cards).To(BeNil())
		})

		It("should FlattenStories() correctly", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})

			Expect(err).NotTo(HaveOccurred())

			cards, err := pt.FlattenStories()

			Expect(err).NotTo(HaveOccurred())
			Expect(len(cards)).To(Equal(1))
			Expect(cards[0].Title).To(Equal("Test Rubbernecker"))
			Expect(len(cards[0].Stickers)).To(Equal(2))
			_, ok := cards[0].Stickers.Get("blocked")
			Expect(ok).To(BeTrue())
		})

		DescribeTable("a scheduled date in a blocker should be parsed correctly",
			func(blockerDescription string, expectedDate time.Time) {
				response = fmt.Sprintf(
					`[{"blockers": [{"created_at":"2199-09-01T12:34:56Z", "description":"%s"}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[{"name":"test"}]}]`,
					blockerDescription,
				)
				httpmock.RegisterResponder("GET", apiURL, httpmock.NewStringResponder(200, response))

				err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})
				Expect(err).NotTo(HaveOccurred())

				cards, err := pt.FlattenStories()
				Expect(err).NotTo(HaveOccurred())
				Expect(cards).To(HaveLen(1))

				stickers := cards[0].Stickers

				Expect(stickers).To(
					ContainElement(MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("scheduled"),
						"Title":   Equal(blockerDescription),
						"Content": Equal(expectedDate.Format("2/1")),
					})),
				)

				Expect(stickers).NotTo(
					ContainElement(MatchFields(IgnoreExtras, Fields{
						"Name": Equal("blocked"),
					})),
				)
			},
			// All date formats
			Entry("until 2/9", "until 2/9", getTimeFromStr("2/9/2199")),
			Entry("until 02/09", "until 02/09", getTimeFromStr("2/9/2199")),
			Entry("until 02/09/2199", "until 02/09/2199", getTimeFromStr("2/9/2199")),
			Entry("until 2199/09/02", "until 2199/09/02", getTimeFromStr("2/9/2199")),
			Entry("until 2199-09-02", "until 2199-09-02", getTimeFromStr("2/9/2199")),

			// Next year
			Entry("until 3/1", "until 3/1", getTimeFromStr("3/1/2020")),

			// All prepositions
			Entry("before 2/9", "before 2/9", getTimeFromStr("2/9/2199")),
			Entry("on 2/9", "on 2/9", getTimeFromStr("2/9/2199")),
			Entry("on or after 2/9", "on or after 2/9", getTimeFromStr("2/9/2199")),

			// After should add an extra day
			Entry("after 2/9", "after 2/9", getTimeFromStr("3/9/2199")),
		)

		DescribeTable("a knowledge-share label in a blocker should be parsed correctly",
			func(blockerDescription string) {
				response = fmt.Sprintf(
					`[{"blockers": [{"created_at":"2199-09-01T12:34:56Z", "description":"%s"}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[{"name":"test"}]}]`,
					blockerDescription,
				)
				httpmock.RegisterResponder("GET", apiURL, httpmock.NewStringResponder(200, response))

				err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})
				Expect(err).NotTo(HaveOccurred())

				cards, err := pt.FlattenStories()
				Expect(err).NotTo(HaveOccurred())
				Expect(cards).To(HaveLen(1))

				stickers := cards[0].Stickers

				Expect(stickers).To(
					ContainElement(MatchFields(IgnoreExtras, Fields{
						"Name":    Equal("knowledge-share"),
						"Title":   Equal("knowledge-share"),
						"Content": Equal(""),
					})),
				)

				Expect(stickers).NotTo(
					ContainElement(MatchFields(IgnoreExtras, Fields{
						"Name": Equal("blocked"),
					})),
				)
			},

			Entry("knowledge-share", "knowledge-share"),
			Entry("knowledge_share", "knowledge_share"),
			Entry("knowledge share", "knowledge share"),
			Entry("knowledgeshare", "knowledgeshare"),
		)

		It("a scheduled sticker should not be added if in the past", func() {
			response = fmt.Sprintf(
				`[{"blockers": [{"created_at":"1970-02-01T12:34:56Z", "description":"xx %s xx"}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[{"name":"test"}]}]`,
				"until 1970-02-01",
			)
			httpmock.RegisterResponder("GET", apiURL, httpmock.NewStringResponder(200, response))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})
			Expect(err).NotTo(HaveOccurred())

			cards, err := pt.FlattenStories()
			Expect(err).NotTo(HaveOccurred())

			_, ok := cards[0].Stickers.Get("scheduled")
			Expect(ok).To(BeFalse())
		})

		It("a scheduled sticker should not be added if resolved", func() {
			response = fmt.Sprintf(
				`[{"blockers": [{"created_at":"2199-02-01T12:34:56Z", "description":"xx %s xx", "resolved": true}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[{"name":"test"}]}]`,
				"until "+time.Now().AddDate(0, 0, 1).Format("2/1/2006"),
			)
			httpmock.RegisterResponder("GET", apiURL, httpmock.NewStringResponder(200, response))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})
			Expect(err).NotTo(HaveOccurred())

			cards, err := pt.FlattenStories()
			Expect(err).NotTo(HaveOccurred())

			_, ok := cards[0].Stickers.Get("blocked")
			Expect(ok).To(BeFalse())

			_, ok = cards[0].Stickers.Get("scheduled")
			Expect(ok).To(BeFalse())
		})

		It("a scheduled blocker should not hide an other blocker", func() {
			response = fmt.Sprintf(
				`[{"blockers": [{"created_at":"2199-02-01T12:34:56Z", "description":"xx %s xx"}, {"created_at":"2199-02-01T12:34:56Z", "description":"other blocker"}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[{"name":"test"}]}]`,
				"until "+time.Now().AddDate(0, 0, 1).Format("2/1/2006"),
			)
			httpmock.RegisterResponder("GET", apiURL, httpmock.NewStringResponder(200, response))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})
			Expect(err).NotTo(HaveOccurred())

			cards, err := pt.FlattenStories()
			Expect(err).NotTo(HaveOccurred())

			_, ok := cards[0].Stickers.Get("blocked")
			Expect(ok).To(BeTrue())

			_, ok = cards[0].Stickers.Get("scheduled")
			Expect(ok).To(BeTrue())
		})

		It("a zero point story should add a sticker", func() {
			response = `[{"estimate": 0, "blockers": [],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[]}]`
			httpmock.RegisterResponder("GET", apiURL, httpmock.NewStringResponder(200, response))

			err := pt.FetchCards(rubbernecker.StatusDoing, map[string]string{})
			Expect(err).NotTo(HaveOccurred())

			cards, err := pt.FlattenStories()
			Expect(err).NotTo(HaveOccurred())

			_, ok := cards[0].Stickers.Get("zero-points")
			Expect(ok).To(BeTrue())
		})
	})

})

func getTimeFromStr(str string) time.Time {
	t, _ := time.Parse("2/1/2006", str)
	return t
}
