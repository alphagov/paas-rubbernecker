package pivotal

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/salsita/go-pivotaltracker/v5/pivotal"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Pivotal internal functionality", func() {
	It("should calculateWorkingDays() correctly", func() {
		days := calculateWorkingDays(time.Date(2017, 10, 30, 12, 0, 0, 0, time.Local), time.Date(2017, 11, 1, 12, 0, 0, 0, time.Local))

		Expect(days).To(Equal(3))
	})

	It("should calculateWorkingDays() over weekend correctly", func() {
		days := calculateWorkingDays(time.Date(2017, 10, 27, 12, 0, 0, 0, time.Local), time.Date(2017, 11, 1, 12, 0, 0, 0, time.Local))

		Expect(days).To(Equal(4))
	})

	It("should fail to calculateInState() due to lack of transitions", func() {
		t := []transition{}

		Expect(calculateInState(t, "started")).To(Equal(0))
	})

	It("should calculateInState() correctly", func() {
		t := []transition{
			transition{
				State:    "started",
				Occurred: time.Now().Add(-7 * 24 * time.Hour),
			},
			transition{
				State:    "finished",
				Occurred: time.Now(),
			},
		}

		Expect(calculateInState(t, "started")).To(Equal(calculateWorkingDays(t[0].Occurred, t[1].Occurred)))
	})

	It("should calculateInState() correctly if it has been restarted", func() {
		t := []transition{
			transition{
				State:    "started", // Started 7 days ago.
				Occurred: time.Now().Add(-7 * 24 * time.Hour),
			},
			transition{
				State:    "finished", // Finished 5 days ago.
				Occurred: time.Now().Add(-5 * 24 * time.Hour),
			},
			transition{
				State:    "rejected", // Someone rejected it shortly after.
				Occurred: time.Now().Add(-4.5 * 24 * time.Hour),
			},
			transition{
				State:    "started", // Restarted 4 days ago.
				Occurred: time.Now().Add(-4 * 24 * time.Hour),
			},
		}

		Expect(calculateInState(t, "started")).To(Equal(calculateWorkingDays(t[3].Occurred, time.Now())))
	})

	It("should composeState() correctly", func() {
		all := composeState(rubbernecker.StatusAll)
		todo := composeState(rubbernecker.StatusScheduled)
		play := composeState(rubbernecker.StatusDoing)
		revi := composeState(rubbernecker.StatusReviewal)
		appr := composeState(rubbernecker.StatusApproval)
		reje := composeState(rubbernecker.StatusRejected)
		done := composeState(rubbernecker.StatusDone)

		Expect(all).To(Equal("unstarted,planned,started,finished,delivered,rejected"))
		Expect(todo).To(Equal(fmt.Sprintf("%s,%s", pivotal.StoryStateUnstarted, pivotal.StoryStatePlanned)))
		Expect(play).To(Equal(pivotal.StoryStateStarted))
		Expect(revi).To(Equal(pivotal.StoryStateFinished))
		Expect(appr).To(Equal(pivotal.StoryStateDelivered))
		Expect(reje).To(Equal(pivotal.StoryStateRejected))
		Expect(done).To(Equal(pivotal.StoryStateAccepted))
	})

	It("should convertState() correctly", func() {
		Expect(convertState(pivotal.StoryStatePlanned)).To(Equal("next"))
		Expect(convertState(pivotal.StoryStateStarted)).To(Equal("doing"))
		Expect(convertState(pivotal.StoryStateFinished)).To(Equal("reviewing"))
		Expect(convertState(pivotal.StoryStateDelivered)).To(Equal("approving"))
		Expect(convertState(pivotal.StoryStateRejected)).To(Equal("rejected"))
		Expect(convertState(pivotal.StoryStateAccepted)).To(Equal("done"))
		Expect(convertState("testing")).To(Equal("unknown"))
	})
})
