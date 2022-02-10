package github

import (
	"context"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

// Schedule will hold some internal and external information, such as client and
// contents of the call.
type Schedule struct {
	Client       *github.Client
	pullRequests []*github.PullRequest
	repositories []*github.Repository
}

// New will create an instance of a Schedule and prefill it with the GitHub
// client.
func New(token string) *Schedule {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Schedule{
		Client: github.NewClient(tc),
	}
}
