package pagerduty

import (
	"time"

	pd "github.com/PagerDuty/go-pagerduty"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

// Schedule will hold some internal and external information, such as client and
// contents of the call.
type Schedule struct {
	Client  *pd.Client
	content []pd.OnCall
}

// New will create an instance of a Schedule and prefill it with the PagerDuty
// client.
func New(token string) *Schedule {
	return &Schedule{
		Client: pd.NewClient(token),
	}
}

// FetchSupport will make a call to the PagerDuty API obtaining the response and
// storing it in the Schedule for future use.
func (p *Schedule) FetchSupport() error {
	opts := pd.ListOnCallOptions{
		Since: time.Now().String(),
		Until: time.Now().Add(24 * time.Hour).String(),
	}

	res, err := p.Client.ListOnCalls(opts)
	if err != nil {
		return err
	}

	p.content = res.OnCalls

	return nil
}

// FlattenSupport should convert the stored response from PagerDuty and convert
// it into rubbernecker compatible SupportRota.
func (p *Schedule) FlattenSupport() (rubbernecker.SupportRota, error) {
	support := rubbernecker.SupportRota{}

	for _, oncall := range p.content {
		if _, ok := support[oncall.Schedule.Summary]; oncall.Schedule.Summary == "" || ok {
			continue
		}

		support[oncall.Schedule.Summary] = &rubbernecker.Support{
			Type:   oncall.Schedule.Summary,
			Member: oncall.User.Summary,
		}
	}

	return support, nil
}
