package pivotal

import (
	"fmt"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	pt "github.com/salsita/go-pivotaltracker/v5/pivotal"
)

// Tracker will be responsible for acting as the story resource returned
// by the API.
type Tracker struct {
	client    *pt.Client
	projectID int64
	stories   []*story
	stickers  *rubbernecker.Stickers
	members   []*membership
}

// New will compose a Tracker struct ready to use by the rubbernecker.
func New(projectID int64, token string) (*Tracker, error) {
	return &Tracker{
		client:    pt.NewClient(token),
		projectID: projectID,
		stickers:  &rubbernecker.Stickers{},
	}, nil
}

// AcceptStickers will make a note of enabled stickers in the application and
// attempt to assign them to each story.
func (t *Tracker) AcceptStickers(stickers *rubbernecker.Stickers) {
	t.stickers = stickers
}

// FetchCards will fetch the stories from PivotalTracker.
func (t *Tracker) FetchCards(status rubbernecker.Status) error {
	fields := "owner_ids,blockers,transitions,current_state,labels,name,url,created_at"
	path := fmt.Sprintf("projects/%d/stories?fields=%s&filter=state:%s", t.projectID, fields, composeState(status))

	req, err := t.client.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	t.stories = []*story{}
	_, err = t.client.Do(req, &t.stories)
	if err != nil {
		return err
	}

	return nil
}

// FlattenStories function will take what we have so far and convert it into the
// rubbernecker standard.
func (t *Tracker) FlattenStories() (*rubbernecker.Cards, error) {
	if len(t.stories) == 0 {
		return nil, fmt.Errorf("pivotal extension: no stories to be flattened")
	}

	stories := rubbernecker.Cards{}

	for _, s := range t.stories {
		stickers := rubbernecker.Stickers{}

		for _, l := range s.Labels {
			if sticker := t.stickers.Get(l.Name); sticker != nil {
				stickers = append(stickers, sticker)
			}
		}

		if len(s.Blockers) > 0 && stickers.Get("blocked") == nil {
			if sticker := t.stickers.Get("blocked"); sticker != nil {
				stickers = append(stickers, sticker)
			}
		}

		assignees := rubbernecker.Members{}

		for _, id := range s.OwnerIds {
			assignees[id] = &rubbernecker.Member{ID: id}
		}

		stories = append(stories, &rubbernecker.Card{
			ID:        s.ID,
			Assignees: &assignees,
			Elapsed:   calculateInState(s.Transitions, s.State),
			Status:    convertState(s.State),
			Stickers:  stickers,
			Title:     s.Name,
			URL:       s.URL,
		})
	}

	return &stories, nil
}
