package pivotal

import (
	"strings"
	"time"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	pt "github.com/salsita/go-pivotaltracker/v5/pivotal"
)

type story struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	State       string       `json:"current_state,omitempty"`
	OwnerIds    []int        `json:"owner_ids,omitempty"`
	Labels      []*pt.Label  `json:"labels,omitempty"`
	URL         string       `json:"url,omitempty"`
	Blockers    []blocker    `json:"blockers,omitempty"`
	Transitions []transition `json:"transitions,omitempty"`
	CreatedAt   *time.Time   `json:"created_at,omitempty"`
}

type blocker struct {
	ID       int  `json:"id,omitempty"`
	Resolved bool `json:"resolved,omitempty"`
}

type transition struct {
	State    string    `json:"state,omitempty"`
	Occurred time.Time `json:"occurred_at,omitempty"`
}

type member struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Username string `json:"username"`
}

type membership struct {
	Person member `json:"person"`
	Role   string `json:"role"`
}

func calculateInState(transitions []transition, state string) int {
	var m transition

	if len(transitions) == 0 {
		return 0
	}

	for _, e := range transitions {
		if e.State != state {
			continue
		}

		if e.Occurred.Unix() > m.Occurred.Unix() {
			m = e
		}
	}

	return calculateWorkingDays(m.Occurred, time.Now())
}

func calculateWorkingDays(since, until time.Time) int {
	days := 0

	for {
		if since.After(until) {
			break
		}

		if since.Weekday() != 5 && since.Weekday() != 6 {
			days++
		}

		since = since.Add(24 * time.Hour)
	}

	return days
}

func composeState(status rubbernecker.Status) string {
	var state string

	switch status {
	case rubbernecker.StatusScheduled:
		state = pt.StoryStatePlanned
	case rubbernecker.StatusDoing:
		state = pt.StoryStateStarted
	case rubbernecker.StatusReviewal:
		state = pt.StoryStateFinished
	case rubbernecker.StatusApproval:
		state = pt.StoryStateDelivered
	case rubbernecker.StatusRejected:
		state = pt.StoryStateRejected
	case rubbernecker.StatusDone:
		state = pt.StoryStateAccepted
	default:
		state = strings.Join([]string{
			pt.StoryStateStarted,
			pt.StoryStateFinished,
			pt.StoryStateDelivered,
			pt.StoryStateRejected,
		}, ",")
	}

	return state
}

func convertState(state string) string {
	switch state {
	case pt.StoryStateStarted:
		return rubbernecker.StatusDoing.String()
	case pt.StoryStateFinished:
		return rubbernecker.StatusReviewal.String()
	case pt.StoryStateDelivered:
		return rubbernecker.StatusApproval.String()
	case pt.StoryStateAccepted:
		return rubbernecker.StatusDone.String()
	case pt.StoryStateRejected:
		return rubbernecker.StatusRejected.String()
	case pt.StoryStatePlanned:
		return rubbernecker.StatusScheduled.String()
	default:
		return "unknown"
	}
}
