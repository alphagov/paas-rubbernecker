package rubbernecker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

// Response will be a standard outcome returned when hitting rubbernecker app.
type Response struct {
	Card                 *Card       `json:"card,omitempty"`
	Cards                Cards       `json:"cards,omitempty"`
	SampleCard           *Card       `json:"sample_card,omitempty"`
	Config               *Config     `json:"config,omitempty"`
	Error                string      `json:"error,omitempty"`
	Message              string      `json:"message,omitempty"`
	SupportRota          SupportRota `json:"support,omitempty"`
	TeamMembers          Members     `json:"team_members,omitempty"`
	FreeTeamMembers      Members     `json:"free_team_members,omitempty"`
	Filters              []Filter    `json:"filers,omitempty"`
	AppliedFilterQueries []string    `json:"applied_filters,omitempty"`
	TextFilters          string      `json:"text_filters,omitempty"`
}

// JSON function will execute the response to our HTTP writer.
func (r *Response) JSON(code int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(r)
}

// Template function will execute the response to our HTTP writer providing it
// with HTML.
func (r *Response) Template(code int, w http.ResponseWriter, templateFile ...string) error {
	var err error

	t, err := template.New(path.Base(templateFile[len(templateFile)-1])).Funcs(
		template.FuncMap{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
			"safeURL": func(s string) template.URL {
				return template.URL(s)
			},
		},
	).ParseFiles(templateFile...)

	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		fmt.Fprintf(w, "Rubbernecker could not parse templates:\n%s", err)
		return err
	}

	b := &bytes.Buffer{}
	err = t.Execute(b, r)

	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		fmt.Fprintf(w, "Rubbernecker could not render template:\n%s\n", err)

		return err
	} else {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(code)
		w.Write(b.Bytes())

		return nil
	}
}

// WithSampleCard will set a sample card used for template generation
func (r *Response) WithSampleCard(card *Card) *Response {
	r.SampleCard = card
	return r
}

// WithCards will set a collection/single card for the current response.
func (r *Response) WithCards(cards Cards, single bool) *Response {
	if single {
		c := cards
		r.Card = c[0]
		return r
	}

	r.Cards = cards
	return r
}

// WithConfig will set a configuration that will be returned in a response.
func (r *Response) WithConfig(config *Config) *Response {
	r.Config = config
	return r
}

// WithError will set an error for the current response.
func (r *Response) WithError(err error) *Response {
	r.Error = err.Error()
	return r
}

// WithSupport will set either rota or a single support data for the current
// response.
func (r *Response) WithSupport(rota SupportRota) *Response {
	r.SupportRota = rota
	return r
}

// WithTeamMembers will set the allocated parameter for the current response.
func (r *Response) WithTeamMembers(members Members) *Response {
	r.TeamMembers = members
	return r
}

// WithFreeTeamMembers should prepare a list of team members that are free to
// pickup new work.
func (r *Response) WithFreeTeamMembers() *Response {
	if r.TeamMembers != nil && r.Cards != nil {
		free := Members{}
		for id, member := range r.TeamMembers {
			free[id] = member
		}

		for _, card := range r.Cards {
			for _, assignee := range card.Assignees {
				if assignee != nil {
					delete(free, assignee.ID)
				}
			}
		}

		r.FreeTeamMembers = free
	}

	return r
}

// WithFilters will set te filters param for the current response
func (r *Response) WithFilters(filters []Filter) *Response {
	r.Filters = filters
	return r
}

// WithAppliedFilterQueries will set the applied filter queries param for the current response.
func (r *Response) WithAppliedFilterQueries(queries []string) *Response {
	r.AppliedFilterQueries = queries
	return r
}

// WithTextFilters will set the text-filters param for the current response.
func (r *Response) WithTextFilters(filters []string) *Response {
	// We do not include stickers in the text filters
	textFilters := make([]string, 0)

	for _, f := range filters {
		if !strings.HasPrefix(f, "sticker:") &&
			!strings.HasPrefix(f, "not-sticker:") {

			textFilters = append(textFilters, f)
		}
	}

	r.AppliedFilterQueries = filters
	return r
}
