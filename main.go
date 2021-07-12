package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alphagov/paas-rubbernecker/pkg/pagerduty"
	"github.com/alphagov/paas-rubbernecker/pkg/pivotal"
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	"github.com/carlescere/scheduler"
	"github.com/gorilla/mux"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	etag      time.Time
	cards     rubbernecker.Cards
	doneCards rubbernecker.Cards
	members   rubbernecker.Members
	support   = rubbernecker.SupportRota{
		"in-hours":           &rubbernecker.Support{},
		"in-hours-comms":     &rubbernecker.Support{},
		"out-of-hours":       &rubbernecker.Support{},
		"out-of-hours-comms": &rubbernecker.Support{},
		"escalations":        &rubbernecker.Support{},
		"autom8":             &rubbernecker.Support{},
	}

	verbose = kingpin.Flag("verbose", "Will enable the DEBUG logging level.").Default("false").Short('v').OverrideDefaultFromEnvar("DEBUG").Bool()
	port    = kingpin.Flag("port", "Port the application should listen for the traffic on.").Default("8080").Short('p').OverrideDefaultFromEnvar("PORT").Int64()

	pivotalProjectID   = kingpin.Flag("pivotal-project", "Pivotal Tracker project ID rubbernecker will be using.").OverrideDefaultFromEnvar("PIVOTAL_TRACKER_PROJECT_ID").Int64()
	pivotalAPIToken    = kingpin.Flag("pivotal-token", "Pivotal Tracker API token rubbernecker will use to communicate with Pivotal API.").OverrideDefaultFromEnvar("PIVOTAL_TRACKER_API_TOKEN").String()
	pagerdutyAuthToken = kingpin.Flag("pagerduty-token", "PagerDuty auth token rubbernecker will use to communicate with PagerDuty API.").OverrideDefaultFromEnvar("PAGERDUTY_AUTHTOKEN").String()
)

func setupLogger() {
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("Setting up logger.")

	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}

	log.SetFormatter(formatter)
}

func combineCards(collections ...rubbernecker.Cards) rubbernecker.Cards {
	var all = rubbernecker.Cards{}

	for _, collection := range collections {
		for _, card := range collection {
			all = append(all, card)
		}
	}

	return all
}

func fetchStories(pt *pivotal.Tracker) error {
	if members == nil {
		return fmt.Errorf("rubbernecker: could not find any members")
	}

	err := pt.FetchCards(rubbernecker.StatusAll, map[string]string{})
	if err != nil {
		return err
	}

	c, err := pt.FlattenStories()
	if err != nil {
		log.Debug(err)
	}

	year, month, day := time.Now().Date()
	past := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -5).UnixNano() / int64(time.Millisecond)
	err = pt.FetchCards(rubbernecker.StatusDone, map[string]string{
		"accepted_after": fmt.Sprintf("%d", past),
	})
	if err != nil {
		return err
	}

	d, err := pt.FlattenStories()
	if err != nil {
		log.Debug(err)
	}
	d.Reverse()

	for s, story := range c {
		if story.Assignees == nil {
			continue
		}

		for i, a := range story.Assignees {
			if member, ok := members[a.ID]; ok {
				(c[s].Assignees)[i] = member
			}
		}
	}

	if !reflect.DeepEqual(cards, c) {
		cards = c
		etag = time.Now()
	}

	if !reflect.DeepEqual(doneCards, d) {
		doneCards = d
		etag = time.Now()
	}

	log.Debug("Stories have been fetched.")

	return nil
}

func fetchSupport(pd *pagerduty.Schedule) error {
	if pd.Client == nil {
		return fmt.Errorf("PAGERDUTY_AUTHTOKEN is not set, support rota will not be fetched")
	}

	err := pd.FetchSupport()
	if err != nil {
		return err
	}

	s, err := pd.FlattenSupport()
	if err != nil {
		return err
	}

	s = formatSupportNames(s)

	if !reflect.DeepEqual(support, s) {
		support = s
		etag = time.Now()
	}

	log.Debug("Support Rota have been fetched.")

	return nil
}

func fetchUsers(pt *pivotal.Tracker) error {
	err := pt.FetchMembers()
	if err != nil {
		return err
	}

	m, err := pt.FlattenMembers()
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(members, m) {
		members = m
		etag = time.Now()
	}

	log.Debug("Team Members have been fetched.")

	return nil
}

func formatSupportNames(s rubbernecker.SupportRota) rubbernecker.SupportRota {
	return rubbernecker.SupportRota{
		"in-hours":           s.Get("PaaS team rota - in hours"),
		"in-hours-comms":     s.Get("PaaS team rota - comms lead (in Hours)"),
		"out-of-hours":       s.Get("PaaS team rota - out of hours"),
		"out-of-hours-comms": s.Get("PaaS team rota - comms lead (OOH)"),
		"autom8":             s.Get("Techops Autom8"),
		"escalations":        s.Get("GaaP SCS Escalation"),
	}
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := rubbernecker.Response{Message: "OK"}
	resp.JSON(200, w)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	resp := rubbernecker.Response{}
	et := strconv.FormatInt(etag.Unix(), 10)

	if r.Header.Get("If-None-Match") == et {
		resp.
			JSON(http.StatusNotModified, w)

		return
	}

	filterQueries := r.URL.Query()["filter"]

	filteredCards := cards.FilterBy(filterQueries)
	filteredDoneCards := doneCards.FilterBy(filterQueries)

	squads := []rubbernecker.Squad{
		{DisplayText: "PaaS Squad", QueryText: "filter=squad:paas"},
		{DisplayText: "Disco Squad", QueryText: "filter=squad:disco"},
	}

	resp.
		WithConfig(&rubbernecker.Config{
			ReviewalLimit: 4,
			ApprovalLimit: 5,
		}).
		WithCards(combineCards(filteredCards, filteredDoneCards), false).
		WithSampleCard(&rubbernecker.Card{}).
		WithTeamMembers(members).
		WithFreeTeamMembers().
		WithFilters(rubbernecker.DefaultFilterSet()).
		WithAppliedFilterQueries(filterQueries).
		WithTextFilters(filterQueries).
		WithSquads(squads).
		WithSupport(support)

	if strings.Contains(r.Header.Get("Accept"), "json") {
		w.Header().Set("ETag", et)

		err = resp.JSON(http.StatusOK, w)
	} else {
		err = resp.Template(http.StatusOK, w,
			"./build/views/sticker.html",
			"./build/views/card.html",
			"./build/views/index.html",
		)
	}

	if err != nil {
		log.Error(err)
	}
}

func main() {
	kingpin.Parse()
	setupLogger()

	var pd = &pagerduty.Schedule{
		Client: nil,
	}
	if pagerdutyAuthToken != nil && *pagerdutyAuthToken != "" {
		pd = pagerduty.New(*pagerdutyAuthToken)
	}

	pt, err := pivotal.New(*pivotalProjectID, *pivotalAPIToken)
	if err != nil {
		log.Fatal(err)
	}

	stickers, err := ioutil.ReadFile("stickers.yml")
	if err != nil {
		log.Fatal(err)
	}

	var approvedStickers rubbernecker.Stickers
	err = yaml.Unmarshal(stickers, &approvedStickers)
	if err != nil {
		log.Fatal(err)
	}

	pt.AcceptStickers(approvedStickers)

	// We have to fetch the users synchronously first as the fetchStories call depends on it
	if err := fetchUsers(pt); err != nil {
		log.Error(err)
	}

	scheduler.Every(1).Hours().NotImmediately().Run(func() {
		if err := fetchUsers(pt); err != nil {
			log.Error(err)
		}
	})

	scheduler.Every(5).Minutes().Run(func() {
		if err := fetchSupport(pd); err != nil {
			log.Error(err)
		}
	})

	scheduler.Every(20).Seconds().Run(func() {
		if err := fetchStories(pt); err != nil {
			log.Error(err)
		}
	})

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/state", indexHandler)
	r.HandleFunc("/health-check", healthcheckHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}
