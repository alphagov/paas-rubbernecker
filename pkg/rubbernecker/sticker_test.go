package rubbernecker_test

import (
	. "github.com/onsi/ginkgo/v2"
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
		var ok bool

		_, ok = sticker.Matches("test")
		Expect(ok).To(BeTrue())
		_, ok = sticker.Matches("testing")
		Expect(ok).To(BeFalse())
	})

	It("should check if the sticker Matches() one of the aliases", func() {
		var ok bool

		_, ok = sticker.Matches("tset")
		Expect(ok).To(BeTrue())

		_, ok = sticker.Matches("trial")
		Expect(ok).To(BeTrue())

		_, ok = sticker.Matches("testing")
		Expect(ok).To(BeFalse())
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
			var ok bool

			_, ok = sticker.Matches("test: a")
			Expect(ok).To(BeTrue())

			_, ok = sticker.Matches("test: foo")
			Expect(ok).To(BeFalse())
		})

		It("should expand submatches of Regex in all the UI fields", func() {
			s, ok := sticker.Matches("test: a")
			Expect(ok).To(BeTrue())
			Expect(s.Title).To(Equal("Test a"))
			Expect(s.Image).To(Equal("/test_a.png"))
			Expect(s.Content).To(Equal("Only a test. Not to worry! a"))
			Expect(s.Class).To(Equal("class-a"))
		})
	})

	It("should establish if the list Has() specific sticker", func() {
		ss := rubbernecker.Stickers{sticker}

		s, ok := ss.Get("trial")
		Expect(ok).To(BeTrue())
		Expect(s).To(Equal(sticker))

		_, ok = ss.Get("testing")
		Expect(ok).To(BeFalse())
	})
})
