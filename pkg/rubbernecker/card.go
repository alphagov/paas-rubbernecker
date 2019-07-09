package rubbernecker

import (
	"strings"
)

// Status is treated as an enum for the story status codes.
type Status int

const (
	// StatusAll should bring up all the cards from a
	// ProjectManagementService.Fetch call.
	StatusAll Status = iota
	// StatusScheduled should only filter the stories that are not in the
	// StatusStarted state, but prioritised into backlog.
	StatusScheduled
	// StatusDoing should only filter the stories that are in play.
	StatusDoing
	// StatusReviewal should only filter the stories that are in progress of
	// reviewal.
	StatusReviewal
	// StatusApproval should only filter the stories that are in progress of
	// approval.
	StatusApproval
	// StatusRejected should only filter the stories that have been rejected upon
	// approval.
	StatusRejected
	// StatusDone should only filter the stories that are done.
	StatusDone
)

// Card will be a rubbernecker entity composed of the extension.
type Card struct {
	ID        int      `json:"id"`
	Assignees Members  `json:"assignees"`
	Elapsed   int      `json:"in_play"`
	Status    string   `json:"status"`
	Stickers  Stickers `json:"stickers"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	StoryType string   `json:"story_type"`
}

// Cards will be a rubbernecker representative of all cards.
type Cards []*Card

// Reverse reverses the order of the cards in place
func (c Cards) Reverse() {
	for i := len(c)/2 - 1; i >= 0; i-- {
		opp := len(c) - 1 - i
		c[i], c[opp] = c[opp], c[i]
	}
}

// ProjectManagementService is an interface that should force each extension to
// flatten their story into rubbernecker format.
type ProjectManagementService interface {
	AcceptStickers(Stickers)
	FetchCards(Status, map[string]string) error
	FlattenStories() (Cards, error)
}

func (s Status) String() string {
	switch s {
	case StatusDoing:
		return "doing"
	case StatusReviewal:
		return "reviewing"
	case StatusApproval:
		return "approving"
	case StatusDone:
		return "done"
	case StatusRejected:
		return "rejected"
	case StatusScheduled:
		return "next"
	default:
		return "unknown"
	}
}

// Filter the cards by status.
func (c Cards) Filter(s string) Cards {
	tmp := Cards{}

	if c == nil {
		return tmp
	}

	for _, card := range c {
		if card.Status == s {
			tmp = append(tmp, card)
		}
	}

	return tmp
}

func (c Cards) FilterBy(filters []string) Cards {
	if len(filters) == 0 {
		return c
	}

	filter := strings.ToLower(filters[0])

	filteredCards := make(Cards, 0)

	for _, card := range c {
		shouldAdd := false

		if strings.HasPrefix(filter, "person:") {
			memberName := strings.ToLower(strings.Replace(filter, "person:", "", -1))
			for _, member := range card.Assignees {
				if strings.Contains(strings.ToLower(member.Name), memberName) {
					shouldAdd = true
				}
			}
		} else if strings.HasPrefix(filter, "title:") {
			title := strings.ToLower(strings.Replace(filter, "title:", "", -1))
			if strings.Contains(strings.ToLower(card.Title), title) {
				shouldAdd = true
			}
		} else if strings.HasPrefix(filter, "sticker:") {
			sname := strings.ToLower(strings.Replace(filter, "sticker:", "", -1))

			for _, sticker := range card.Stickers {
				if strings.HasPrefix(sticker.Name, sname) {
					shouldAdd = true
				}
			}
		} else if strings.HasPrefix(filter, "not-sticker:") {
			sname := strings.ToLower(strings.Replace(filter, "not-sticker:", "",- 1))
			shouldAdd = true

			for _, sticker := range card.Stickers {
				if strings.HasPrefix(sticker.Name, sname) {
					shouldAdd = false
				}
			}
		} else {
			shouldAdd = true
		}

		if shouldAdd {
			filteredCards = append(filteredCards, card)
		}
	}

	return filteredCards.FilterBy(filters[1:])
}
