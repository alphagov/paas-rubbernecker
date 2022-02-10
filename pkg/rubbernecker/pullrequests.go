package rubbernecker

// PullRequest will be a rubbernecker entity composed of the extension.
type PullRequest struct {
	Author      string     `json:"author"`
	Draft       bool       `json:"draft"`
	Number      int        `json:"number"`
	OpenForDays int        `json:"openForDays"`
	Repository  Repository `json:"repository"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
}

// PullRequests will be a rubbernecker representative of all PullRequests.
type PullRequests []PullRequest
