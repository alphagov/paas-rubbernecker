package github

import (
	"context"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

// FetchRepositories will make a call to the GitHub API obtaining the
// repositories and pull request the provided token has access to.
func (gh *Schedule) FetchRepositories() error {
	repos, _, err := gh.Client.Repositories.List(context.Background(), "", nil)
	if err != nil {
		return err
	}

	gh.repositories = repos

	return nil
}

// FlattenRepositories will make a call to the GitHub API obtaining the
// repositories and pull request the provided token has access to.
func (gh *Schedule) FlattenRepositories() (rubbernecker.Repositories, error) {
	var repos []rubbernecker.Repository
	for _, repo := range gh.repositories {
		repos = append(repos, rubbernecker.Repository{
			Name:         *repo.Name,
			Organisation: *repo.Owner.Login,
		})
	}

	return repos, nil
}
