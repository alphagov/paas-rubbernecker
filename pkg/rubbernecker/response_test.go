package rubbernecker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Response", func() {
	var (
		resp rubbernecker.Response
	)

	BeforeEach(func() {
		resp = rubbernecker.Response{}
	})

	It("should setup the response WithCards() collection", func() {
		resp.WithCards(rubbernecker.Cards{
			&rubbernecker.Card{Title: "Test Case"},
		}, false)

		Expect(resp.Cards).NotTo(BeNil())
		Expect(resp.Card).To(BeNil())
	})

	It("should setup the response WithCards() on a single entity", func() {
		resp.WithCards(rubbernecker.Cards{
			&rubbernecker.Card{Title: "Test Case"},
		}, true)

		Expect(resp.Cards).To(BeNil())
		Expect(resp.Card).NotTo(BeNil())
	})

	It("should setup the response WithConfig() collection", func() {
		resp.WithConfig(&rubbernecker.Config{
			ApprovalLimit: 100,
		})

		Expect(resp.Config).NotTo(BeNil())
	})

	It("should setup the response WithError()", func() {
		resp.WithError(fmt.Errorf("test case: unknonw error"))

		Expect(resp.Error).NotTo(BeNil())
	})

	It("should setup the response WithSupport() rota collection", func() {
		resp.WithSupport(rubbernecker.SupportRota{
			"Test Case": &rubbernecker.Support{Type: "Test Case"},
		})

		Expect(resp.SupportRota).NotTo(BeNil())
	})

	It("should setup the response WithTeamMembers()", func() {
		resp.WithTeamMembers(rubbernecker.Members{
			0: &rubbernecker.Member{Name: "Test Case"},
		})

		Expect(resp.TeamMembers).NotTo(BeNil())
	})

	It("should setup the response WithFreeTeamMembers()", func() {
		mem := rubbernecker.Member{ID: 1234, Name: "Tester"}
		mems := rubbernecker.Members{1234: &mem, 4321: &rubbernecker.Member{ID: 4321, Name: "Free"}}
		card := rubbernecker.Card{Title: "Test", Assignees: rubbernecker.Members{1234: &mem}}
		cards := rubbernecker.Cards{&card}

		resp.
			WithCards(cards, false).
			WithTeamMembers(mems).
			WithFreeTeamMembers()

		Expect(resp.FreeTeamMembers).NotTo(BeNil())
		Expect(len(resp.FreeTeamMembers)).To(Equal(1))
	})

	It("should compose a JSON() response", func() {
		req, err := http.NewRequest("GET", "/500", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := resp.WithError(fmt.Errorf("test case: unknown error")).JSON(500, w)

			Expect(err).NotTo(HaveOccurred())
		})
		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(500))
		Expect(rr.Body.String()).To(ContainSubstring(`{"error":"test case: unknown error"}`))
	})

	It("should compose a Template() response", func() {
		req, err := http.NewRequest("GET", "/", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := resp.WithError(fmt.Errorf("test case: unknown error")).Template(500, w, "./test/index.html")

			Expect(err).NotTo(HaveOccurred())
		})
		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(500))
		Expect(rr.Body.String()).To(ContainSubstring(`<!doctype html>`))
	})

	It("should fail to compose a Template() response due to missing file", func() {
		req, err := http.NewRequest("GET", "/", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := resp.WithError(fmt.Errorf("test case: unknown error")).Template(500, w, "./test/failed.html")

			Expect(err).To(HaveOccurred())
		})
		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(500))
		Expect(rr.Body.String()).To(ContainSubstring("Rubbernecker could not parse templates"))
	})

	It("should not error when assignees are valid", func() {

		req, err := http.NewRequest("GET", "/", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp.WithCards(rubbernecker.Cards{
				&rubbernecker.Card{
					Title: "Test Case",
					Assignees: rubbernecker.Members{
						1: &rubbernecker.Member{Name: "fieldexists"},
					},
				},
			}, false).Template(200, w, "./test/cards.html")
		})
		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(200))
		Expect(rr.Body.String()).To(ContainSubstring("fieldexists"))
	})

	It("should error when assignees are invalid", func() {

		req, err := http.NewRequest("GET", "/", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp.WithCards(rubbernecker.Cards{
				&rubbernecker.Card{
					Title: "Test Case",
					Assignees: rubbernecker.Members{
						1: nil,
					},
				},
			}, false).Template(200, w, "./test/cards.html")
		})
		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(500))
		Expect(rr.Body.String()).To(ContainSubstring("Rubbernecker could not render template"))
	})

	It("should escape HTML in templates", func() {

		req, err := http.NewRequest("GET", "/", nil)
		Expect(err).NotTo(HaveOccurred())

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp.WithCards(rubbernecker.Cards{
				&rubbernecker.Card{
					Title: "<script>alert(1)</script>",
				},
			}, false).Template(200, w, "./test/cards.html")
		})
		handler.ServeHTTP(rr, req)

		Expect(rr.Code).To(Equal(200))

		Expect(rr.Body.String()).NotTo(
			ContainSubstring("<script>alert(1)</script>"),
		)

		Expect(rr.Body.String()).To(
			ContainSubstring("&lt;script&gt;alert(1)&lt;/script&gt;"),
		)
	})
})
