package github

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	"github.com/google/go-github/v42/github"
)

// FetchPullRequests will make a call to the GitHub API obtaining the
// repositories and pull request the provided token has access to.
func (gh *Schedule) FetchPullRequests() error {
	var pullRequests []*github.PullRequest
	for _, repository := range gh.repositories {
		prs, _, err := gh.Client.PullRequests.List(context.Background(), *repository.Owner.Login, *repository.Name, nil)
		if err != nil {
			return err
		}

		pullRequests = append(pullRequests, prs...)
	}

	gh.pullRequests = pullRequests

	return nil
}

// FlattenPullRequests will make a call to the GitHub API obtaining the
// repositories and pull request the provided token has access to.
func (gh *Schedule) FlattenPullRequests() (rubbernecker.PullRequests, error) {
	var prs []rubbernecker.PullRequest
	for _, pr := range gh.pullRequests {
		draft := *pr.Draft
		if strings.Contains(strings.ToUpper(*pr.Title), "WIP") {
			draft = true
		}

		prs = append(prs, rubbernecker.PullRequest{
			Author:      *pr.User.Login,
			Draft:       draft,
			Number:      *pr.Number,
			OpenForDays: int(math.Floor((time.Now()).Sub(*pr.CreatedAt).Hours() / 24)),
			Title:       *pr.Title,
			URL:         *pr.HTMLURL,
			Repository: rubbernecker.Repository{
				Name:         *pr.GetBase().Repo.Name,
				Organisation: *pr.GetBase().Repo.Owner.Login,
			},
		})
	}

	return prs, nil
}
