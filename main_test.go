package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/alphagov/paas-rubbernecker/pkg/helpers"
	"github.com/alphagov/paas-rubbernecker/pkg/pagerduty"
	"github.com/alphagov/paas-rubbernecker/pkg/pivotal"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var _ = Describe("Main", func() {
	Context("provided everything has been setup correctly", func() {
		var (
			err error
			pt  *pivotal.Tracker
			pd  *pagerduty.Schedule

			year, month, day = time.Now().Date()
			past             = time.Date(year, month, day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -5).UnixNano() / int64(time.Millisecond)

			apiURL          = `https://www.pivotaltracker.com/services/v5/projects/123456/stories?fields=owner_ids,blockers,transitions,current_state,labels,name,url,created_at,story_type,estimate&filter=state:unstarted,started,finished,delivered,rejected`
			apiURLAccepted  = fmt.Sprintf(`https://www.pivotaltracker.com/services/v5/projects/123456/stories?fields=owner_ids,blockers,transitions,current_state,labels,name,url,created_at,story_type,estimate&accepted_after=%d`, past)
			apiURLMembers   = `https://www.pivotaltracker.com/services/v5/projects/123456/memberships`
			apiURLSupport   = `https://api.pagerduty.com/oncalls`
			response        = `[{"blockers": [{"name":1234}],"transitions": [],"name": "Test Rubbernecker","current_state": "started","url": "http://localhost/story/show/561","owner_ids":[1234],"labels":[], "story_type": "feature"}]`
			responseMembers = `[{"person":{"id":1234,"name":"Tester"}}]`
		)

		BeforeEach(func() {
			pt, err = pivotal.New(123456, "qwerty123456")

			Expect(err).NotTo(HaveOccurred())
			Expect(pt).NotTo(BeNil())

			pd = pagerduty.New("qwerty123456")

			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail to fetchStories() due to non-responsive API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(500, ``))

			err = fetchStories(pt)

			Expect(err).To(HaveOccurred())
		})

		It("should fail to fetchStories() due to faulty API", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, `[]`))

			err = fetchStories(pt)

			Expect(err).To(HaveOccurred())
		})

		It("should fail to fetchStories() due to lack of members", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))

			err = fetchStories(pt)

			Expect(err).To(HaveOccurred())
		})

		It("should fail to fetchUsers() due to faulty API", func() {
			httpmock.RegisterResponder("GET", apiURLMembers,
				httpmock.NewStringResponder(200, `[]`))

			err = fetchUsers(pt)

			Expect(err).To(HaveOccurred())
		})

		It("should fetchUsers() successfully", func() {
			httpmock.RegisterResponder("GET", apiURLMembers,
				httpmock.NewStringResponder(200, responseMembers))

			err = fetchUsers(pt)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fetchStories() successfully", func() {
			httpmock.RegisterResponder("GET", apiURL,
				httpmock.NewStringResponder(200, response))
			httpmock.RegisterResponder("GET", apiURLAccepted,
				httpmock.NewStringResponder(200, response))

			err = fetchStories(pt)

			Expect(err).NotTo(HaveOccurred())
		})

		It("should fail to fetchSupport() due to faulty API", func() {
			httpmock.RegisterResponder("GET", apiURLSupport,
				httpmock.NewStringResponder(200, `[]`))

			err = fetchSupport(pd)

			Expect(err).To(HaveOccurred())
		})

		It("should fetch all on-call schedules", func() {
			resp := `{"oncalls":[
				{"user":{"summary":"X"},"schedule":{"summary":"PaaS team rota - out of hours"}},
				{"user":{"summary":"Y"},"schedule":{"summary":"PaaS team rota - in hours"}},
				{"user":{"summary":"Z"},"schedule":{"summary":"GaaP SCS Escalation"}}
			]}`
			httpmock.RegisterResponder("GET", apiURLSupport, httpmock.NewStringResponder(200, resp))

			err = fetchSupport(pd)
			Expect(err).NotTo(HaveOccurred())

			Expect(support).To(Equal(rubbernecker.SupportRota(map[string]*rubbernecker.Support{
				"in-hours-comms": {
					Type:   "PaaS team rota - comms lead (in Hours)",
					Member: "-",
				},
				"out-of-hours": {
					Type:   "PaaS team rota - out of hours",
					Member: "X",
				},
				"out-of-hours-comms": {
					Type:   "PaaS team rota - comms lead (OOH)",
					Member: "-",
				},
				"in-hours": {
					Type:   "PaaS team rota - in hours",
					Member: "Y",
				},
				"autom8": {
					Type: "Techops Autom8",
					Member: "-",
				},
				"escalations": {
					Type:   "GaaP SCS Escalation",
					Member: "Z",
				},
			})))
		})

		It("should fetch multiple pages of on-call schedules", func() {
			resp1 := `{
				"oncalls":[
					{"user":{"summary":"X"},"schedule":{"summary":"PaaS team rota - out of hours"}},
					{"user":{"summary":"A"},"schedule":{"summary":"other team foo"}},
					{"user":{"summary":"B"},"schedule":{"summary":"other team bar"}}
				],
				"limit": 3,
				"offset": 0,
				"more": true
			}`
			resp2 := `{
				"oncalls":[
					{"user":{"summary":"C"},"schedule":{"summary":"different team foo"}},
					{"user":{"summary":"Y"},"schedule":{"summary":"PaaS team rota - in hours"}},
					{"user":{"summary":"D"},"schedule":{"summary":"different team bar"}}
				],
				"limit": 3,
				"offset": 3,
				"more": true
			}`
			resp3 := `{
				"oncalls":[
					{"user":{"summary":"Z"},"schedule":{"summary":"GaaP SCS Escalation"}}
				],
				"limit": 1,
				"offset": 6,
				"more": false
			}`
			httpmock.RegisterResponder("GET", apiURLSupport,
				helpers.NewCycleResponder(
					httpmock.NewStringResponder(200, resp1),
					httpmock.NewStringResponder(200, resp2),
					httpmock.NewStringResponder(200, resp3),
				),
			)

			err = fetchSupport(pd)
			Expect(err).NotTo(HaveOccurred())

			Expect(support).To(Equal(rubbernecker.SupportRota(map[string]*rubbernecker.Support{
				"in-hours-comms": {
					Type:   "PaaS team rota - comms lead (in Hours)",
					Member: "-",
				},
				"out-of-hours": {
					Type:   "PaaS team rota - out of hours",
					Member: "X",
				},
				"out-of-hours-comms": {
					Type:   "PaaS team rota - comms lead (OOH)",
					Member: "-",
				},
				"in-hours": {
					Type:   "PaaS team rota - in hours",
					Member: "Y",
				},
				"autom8": {
					Type: "Techops Autom8",
					Member: "-",
				},
				"escalations": {
					Type:   "GaaP SCS Escalation",
					Member: "Z",
				},
			})))
		})

		It("should handle when there isn't anyone on call for a certain schedule", func() {
			resp := `{"oncalls":[
				{"user":{"summary":"X"},"schedule":{"summary":"PaaS team rota - out of hours"}},
				{"user":{"summary":"Z"},"schedule":{"summary":"GaaP SCS Escalation"}}
			]}`
			httpmock.RegisterResponder("GET", apiURLSupport, httpmock.NewStringResponder(200, resp))

			err = fetchSupport(pd)
			Expect(err).NotTo(HaveOccurred())

			Expect(support).To(Equal(rubbernecker.SupportRota(map[string]*rubbernecker.Support{
				"in-hours-comms": {
					Type:   "PaaS team rota - comms lead (in Hours)",
					Member: "-",
				},
				"out-of-hours": {
					Type:   "PaaS team rota - out of hours",
					Member: "X",
				},
				"out-of-hours-comms": {
					Type:   "PaaS team rota - comms lead (OOH)",
					Member: "-",
				},
				"in-hours": {
					Type:   "PaaS team rota - in hours",
					Member: "-",
				},
				"autom8": {
					Type: "Techops Autom8",
					Member: "-",
				},
				"escalations": {
					Type:   "GaaP SCS Escalation",
					Member: "Z",
				},
			})))
		})

		It("should deal healthcheckHandler() correctly", func() {
			req, err := http.NewRequest("GET", "/health-check", nil)
			Expect(err).NotTo(HaveOccurred())

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(healthcheckHandler)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Body.String()).To(ContainSubstring(`{"message":"OK"}`))
		})

		It("should deal indexHandler() correctly expecting Not Modified", func() {
			req, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Add("Accept", "application/json")
			req.Header.Add("If-None-Match", strconv.FormatInt(etag.Unix(), 10))

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(indexHandler)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusNotModified))
			Expect(rr.Body.String()).To(ContainSubstring(`{}`))
		})

		It("should deal indexHandler() correctly expecting JSON", func() {
			req, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Add("Accept", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(indexHandler)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(ContainSubstring("application/json"))
			Expect(rr.Body.String()).To(ContainSubstring(`"title":"Test Rubbernecker"`))
		})

		It("should deal indexHandler() correctly expecting HTML", func() {
			req, err := http.NewRequest("GET", "/", nil)
			Expect(err).NotTo(HaveOccurred())
			req.Header.Add("Accept", "text/html")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(indexHandler)
			handler.ServeHTTP(rr, req)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(ContainSubstring("text/html"))
			Expect(rr.Body.String()).To(ContainSubstring(`<!doctype html>`))
		})
	})

})
