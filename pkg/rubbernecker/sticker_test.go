package rubbernecker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Response", func() {
	var (
		sticker rubbernecker.Sticker
	)

	BeforeEach(func() {
		sticker = rubbernecker.Sticker{
			Name:    "test",
			Title:   "Test",
			Image:   "/test.png",
			Content: "Only a test. Not to worry!",
			Aliases: []string{"tset", "trial"},
		}
	})

	It("should check if the sticker Matches() name", func() {
		Expect(sticker.Matches("test")).NotTo(BeNil())
		Expect(sticker.Matches("testing")).To(BeNil())
	})

	It("should check if the sticker Matches() one of the aliases", func() {
		Expect(sticker.Matches("tset")).NotTo(BeNil())
		Expect(sticker.Matches("trial")).NotTo(BeNil())
		Expect(sticker.Matches("testing")).To(BeNil())
	})

	Context("when using Regexs", func() {
		BeforeEach(func() {
			sticker = rubbernecker.Sticker{
				Name:    "test",
				Regex:   "test: ([a-zA-Z])$",
				Title:   "Test $1",
				Image:   "/test_$1.png",
				Content: "Only a test. Not to worry! $1",
				Aliases: []string{"tset", "trial"},
				Class:   "class-$1",
			}
		})

		It("should check if the sticker Matches() regex", func() {
			Expect(sticker.Matches("test: a")).NotTo(BeNil())
			Expect(sticker.Matches("test: foo")).To(BeNil())
		})

		It("should expand submatches of Regex in all the UI fields", func() {
			s := sticker.Matches("test: a")
			Expect(s).NotTo(BeNil())
			Expect(s.Title).To(Equal("Test a"))
			Expect(s.Image).To(Equal("/test_a.png"))
			Expect(s.Content).To(Equal("Only a test. Not to worry! a"))
			Expect(s.Class).To(Equal("class-a"))
		})
	})

	It("should establish if the list Has() specific sticker", func() {
		ss := rubbernecker.Stickers{&sticker}

		Expect(ss.Get("trial")).NotTo(BeNil())
		Expect(ss.Get("testing")).To(BeNil())
	})
})
