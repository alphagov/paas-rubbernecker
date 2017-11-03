package rubbernecker

// Status is treated as an enum for the story status codes.
type Status int

const (
	// StatusAll should bring up all the cards from a
	// ProjectManagementService.Fetch call.
	StatusAll Status = iota
	// StatusScheduled should only filter the stories that are not in the
	// StatusStarted state, but prioritised into backlock.
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
	Assignees *Members `json:"assignees"`
	Elapsed   int      `json:"in_play"`
	Status    string   `json:"status"`
	Stickers  Stickers `json:"stickers"`
	Title     string   `json:"title"`
	URL       string   `json:"url"`
}

// Cards will be a rubbernecker representative of all cards.
type Cards []*Card

// ProjectManagementService is an interface that should force each extension to
// flatten their story into rubbernecker format.
type ProjectManagementService interface {
	AcceptStickers(*Stickers)
	FetchCards(Status) error
	FlattenStories() (*Cards, error)
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
func (c *Cards) Filter(s string) Cards {
	tmp := Cards{}

	if c == nil {
		return tmp
	}

	for _, card := range *c {
		if card.Status == s {
			tmp = append(tmp, card)
		}
	}

	return tmp
}
