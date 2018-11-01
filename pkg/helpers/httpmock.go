package helpers

import (
	"net/http"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

// Returns a responder that iterates across multiple responders in sequence
func NewCycleResponder(responderCycle ...httpmock.Responder) httpmock.Responder {
	i := 0
	return func(req *http.Request) (*http.Response, error) {
		r := responderCycle[i]
		i = (i + 1) % len(responderCycle)
		return r(req)
	}
}
