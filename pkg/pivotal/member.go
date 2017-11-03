package pivotal

import (
	"fmt"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

// FetchMembers will contact the PivotalTracker API to get the list of members.
func (t *Tracker) FetchMembers() error {
	path := fmt.Sprintf("projects/%d/memberships", t.projectID)

	req, err := t.client.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	t.members = []*membership{}
	_, err = t.client.Do(req, &t.members)
	if err != nil {
		return err
	}

	return nil
}

// FlattenMembers will convert the PivotalTracker memberships into rubbernecker
// users.
func (t *Tracker) FlattenMembers() (*rubbernecker.Members, error) {
	if len(t.members) == 0 {
		return nil, fmt.Errorf("pivotal extension: no members to be flattened")
	}

	members := rubbernecker.Members{}

	for _, m := range t.members {
		if m.Role == "viewer" {
			continue
		}

		members[m.Person.ID] = &rubbernecker.Member{
			ID:    m.Person.ID,
			Name:  m.Person.Name,
			Email: m.Person.Email,
		}
	}

	return &members, nil
}
