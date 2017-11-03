package rubbernecker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Response", func() {
	var (
		resp rubbernecker.Response
	)

	BeforeEach(func() {
		resp = rubbernecker.Response{}
	})

	It("should convert status to String() correctly", func() {
		Expect(rubbernecker.StatusAll.String()).To(Equal("unknown"))
		Expect(rubbernecker.StatusScheduled.String()).To(Equal("next"))
		Expect(rubbernecker.StatusDoing.String()).To(Equal("doing"))
		Expect(rubbernecker.StatusReviewal.String()).To(Equal("reviewing"))
		Expect(rubbernecker.StatusApproval.String()).To(Equal("approving"))
		Expect(rubbernecker.StatusRejected.String()).To(Equal("rejected"))
		Expect(rubbernecker.StatusDone.String()).To(Equal("done"))
	})

	It("should Filter() stories by status", func() {
		cards := rubbernecker.Cards{
			&rubbernecker.Card{Title: "Test1", Status: "doing"},
			&rubbernecker.Card{Title: "Test2", Status: "reviewing"},
			&rubbernecker.Card{Title: "Test3", Status: "reviewing"},
		}

		doing := cards.Filter(rubbernecker.StatusDoing.String())
		reviewing := cards.Filter(rubbernecker.StatusReviewal.String())

		Expect(len(reviewing)).To(Equal(2))
		Expect(len(doing)).To(Equal(1))
	})
})
