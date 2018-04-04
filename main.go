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
	cards     *rubbernecker.Cards
	doneCards *rubbernecker.Cards
	members   *rubbernecker.Members
	support   *rubbernecker.SupportRota

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

func combineCards(collections ...*rubbernecker.Cards) *rubbernecker.Cards {
	var all = rubbernecker.Cards{}

	for _, collection := range collections {
		for _, card := range *collection {
			all = append(all, card)
		}
	}

	return &all
}

func fetchStories(pt *pivotal.Tracker) error {
	if members == nil {
		return fmt.Errorf("rubbernecker: could not find any members")
	}

	err := pt.FetchCards(rubbernecker.StatusAll)
	if err != nil {
		return err
	}

	c, err := pt.FlattenStories()
	if err != nil {
		return err
	}

	err = pt.FetchCards(rubbernecker.StatusDone)
	if err != nil {
		return err
	}

	d, err := pt.FlattenStories()
	if err != nil {
		return err
	}

	for s, story := range *c {
		if story.Assignees == nil {
			continue
		}

		for i, a := range *story.Assignees {
			(*(*c)[s].Assignees)[i] = (*members)[a.ID]
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

	return nil
}

func fetchSupport(pd *pagerduty.Schedule) error {
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
		support = &s
		etag = time.Now()
	}

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

	return nil
}

func formatSupportNames(s rubbernecker.SupportRota) rubbernecker.SupportRota {
	return rubbernecker.SupportRota{
		"in-hours":     s["PaaS team rota - in hours"],
		"out-of-hours": s["PaaS team rota - out of hours"],
		"escalations":  s["PaaS team Escalations - out of hours"],
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

	resp.
		WithConfig(&rubbernecker.Config{
			ReviewalLimit: 4,
			ApprovalLimit: 5,
		}).
		WithCards(combineCards(cards, doneCards), false).
		WithTeamMembers(members).
		WithFreeTeamMembers().
		WithSupport(support)

	if strings.Contains(r.Header.Get("Accept"), "json") {
		w.Header().Set("ETag", et)

		err = resp.JSON(http.StatusOK, w)
	} else {
		err = resp.Template(http.StatusOK, w,
			"./build/views/index.html",
			"./build/views/card.html",
			"./build/views/sticker.html",
		)
	}

	if err != nil {
		log.Error(err)
	}
}

func main() {
	kingpin.Parse()
	setupLogger()

	pd := pagerduty.New(*pagerdutyAuthToken)
	pt, err := pivotal.New(*pivotalProjectID, *pivotalAPIToken)
	if err != nil {
		log.Fatal(err)
	}

	stickers, err := ioutil.ReadFile("stickers.yml")
	if err != nil {
		log.Fatal(err)
	}

	var approvedStickers *rubbernecker.Stickers
	err = yaml.Unmarshal(stickers, &approvedStickers)
	if err != nil {
		log.Fatal(err)
	}

	pt.AcceptStickers(approvedStickers)

	scheduler.Every(1).Hours().Run(func() {
		err := fetchUsers(pt)
		if err != nil {
			log.Error(err)
		}

		log.Debug("Team Members have been fetched.")
	})

	scheduler.Every(6).Hours().Run(func() {
		err := fetchSupport(pd)
		if err != nil {
			log.Error(err)
		}

		log.Debug("Support Rota have been fetched.")
	})

	// This is only a procaution as the stories rely on the members to be fetched
	// first. Applying NotImmediately() method to the scheduler will make sure,
	// there isn't a race condition between the two.
	log.Info("Will fetch stories in 20 seconds.")
	scheduler.Every(20).Seconds().NotImmediately().Run(func() {
		err := fetchStories(pt)
		if err != nil {
			log.Error(err)
		}

		log.Debug("Stories have been fetched.")
	})

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/state", indexHandler)
	r.HandleFunc("/health-check", healthcheckHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

	http.ListenAndServe(fmt.Sprintf(":%d", *port), r)
}
