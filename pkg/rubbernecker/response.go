package rubbernecker

import (
	"encoding/json"
	"net/http"
	"text/template"
)

// Response will be a standard outcome returned when hitting rubbernecker app.
type Response struct {
	Card            *Card       `json:"card,omitempty"`
	Cards           Cards       `json:"cards,omitempty"`
	Config          *Config     `json:"config,omitempty"`
	Error           string      `json:"error,omitempty"`
	Message         string      `json:"message,omitempty"`
	SupportRota     SupportRota `json:"support,omitempty"`
	TeamMembers     Members     `json:"team_members,omitempty"`
	FreeTeamMembers Members     `json:"free_team_members,omitempty"`
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

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(code)

	t := template.New("Rubbernecker.Template")
	t, err = template.ParseFiles(templateFile...)
	if err != nil {
		return err
	}

	return t.Execute(w, r)
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
