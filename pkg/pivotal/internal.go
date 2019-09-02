package pivotal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	pt "github.com/salsita/go-pivotaltracker/v5/pivotal"
)

var dateBlockerRegex = regexp.MustCompile(`(on or after|after|before|until|on)\s+(\d+\/\d+(?:\/\d+)?|\d+-\d+-\d+|\d+\/\d+\/\d+)`)

var validDateLayouts = []string{
	"2/1/2006",
	"02/01/2006",
	"2006/01/02",
	"2006-01-02",
}

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
	StoryType   string       `json:"story_type"`
}

type blocker struct {
	ID          int        `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	Description string     `json:"description,omitempty"`
	Resolved    bool       `json:"resolved,omitempty"`
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
		state = pt.StoryStateUnstarted
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
			pt.StoryStateUnstarted,
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
	case pt.StoryStateUnstarted:
		return rubbernecker.StatusScheduled.String()
	default:
		return "unknown"
	}
}

func convertBlockersToStickers(blockers []blocker, availableStickers rubbernecker.Stickers) rubbernecker.Stickers {
	var stickers rubbernecker.Stickers
	for _, blocker := range blockers {
		if !blocker.Resolved {
			scheduledDate, err := getScheduledDate(blocker)

			if err != nil {
				if !stickers.Has("scheduled") {
					if sticker, ok := availableStickers.Get("scheduled"); ok {
						sticker.Title = blocker.Description
						sticker.Content = "?/?"
						stickers = append(stickers, sticker)
					}
				}
			} else if scheduledDate != nil {
				if scheduledDate.After(time.Now()) && !stickers.Has("scheduled") {
					if sticker, ok := availableStickers.Get("scheduled"); ok {
						sticker.Title = blocker.Description
						sticker.Content = scheduledDate.Format("2/1")
						stickers = append(stickers, sticker)
					}
				}
			} else {
				if !stickers.Has("blocked") {
					if sticker, ok := availableStickers.Get("blocked"); ok {
						sticker.Title = blocker.Description
						stickers = append(stickers, sticker)
					}
				}
			}
		}
	}

	return stickers
}

func getScheduledDate(blocker blocker) (*time.Time, error) {
	matches := dateBlockerRegex.FindStringSubmatch(blocker.Description)
	if len(matches) > 0 {
		if blocker.CreatedAt == nil {
			panic("a blocker should always have a created_at field")
		}
		date, err := parseScheduledDate(matches[1], matches[2], *blocker.CreatedAt)
		if err != nil {
			return nil, err
		}

		return &date, nil
	}

	return nil, nil
}

func parseScheduledDate(preposition string, date string, minDate time.Time) (time.Time, error) {
	var valid, isShortDate bool
	var result time.Time

	// short date format, like dd/mm
	if regexp.MustCompile(`^\d+\/\d+$`).MatchString(date) {
		isShortDate = true
		date = date + "/" + strconv.Itoa(minDate.Year())
	}

	for _, l := range validDateLayouts {
		if t, err := time.Parse(l, date); err == nil {
			result = t
			valid = true
			break
		}
	}

	if !valid {
		return time.Time{}, fmt.Errorf("unrecognised date format: %s", date)
	}

	// We need to check if the short date format is meant to be in the next year
	// We will consider any dates that are 6 months before the minimum date as a date that should be in the next year
	if isShortDate && result.Before(minDate.AddDate(0, -6, 0)) {
		result = result.AddDate(1, 0, 0)
	}

	if preposition == "after" {
		// Add +1 day
		result = result.AddDate(0, 0, 1)
	}

	return result, nil
}
