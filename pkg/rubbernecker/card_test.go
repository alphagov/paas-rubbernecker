package rubbernecker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

var _ = Describe("Response", func() {
	It("should convert status to String() correctly", func() {
		Expect(rubbernecker.StatusAll.String()).To(Equal("unknown"))
		Expect(rubbernecker.StatusScheduled.String()).To(Equal("next"))
		Expect(rubbernecker.StatusDoing.String()).To(Equal("doing"))
		Expect(rubbernecker.StatusReviewal.String()).To(Equal("reviewing"))
		Expect(rubbernecker.StatusApproval.String()).To(Equal("approving"))
		Expect(rubbernecker.StatusRejected.String()).To(Equal("rejected"))
		Expect(rubbernecker.StatusDone.String()).To(Equal("done"))
	})

	It("should Filter() stories by status", func() {
		cards := rubbernecker.Cards{
			&rubbernecker.Card{Title: "Test1", Status: "doing"},
			&rubbernecker.Card{Title: "Test2", Status: "reviewing"},
			&rubbernecker.Card{Title: "Test3", Status: "reviewing"},
		}

		doing := cards.Filter(rubbernecker.StatusDoing.String())
		reviewing := cards.Filter(rubbernecker.StatusReviewal.String())

		Expect(len(reviewing)).To(Equal(2))
		Expect(len(doing)).To(Equal(1))
	})
})

var _ = Describe("Card Filtering", func() {
	Context("Card Filtering By Sticker Names", func() {
		It("should not do anything if there are no filters", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{},
				&rubbernecker.Card{},
				&rubbernecker.Card{},
			)

			includeStickers := []string{}
			excludeStickers := []string{}

			filteredCards := cards.FilterByStickerNames(
				includeStickers,
				excludeStickers,
			)

			Expect(filteredCards).To(HaveLen(3))
		})

		It("should include cards", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{
					Title: "a non-tech card",
					Stickers: []rubbernecker.Sticker{
						rubbernecker.Sticker{Name: "non-tech"},
					},
				},
				&rubbernecker.Card{
					Title: "a tech card",
					Stickers: []rubbernecker.Sticker{
						rubbernecker.Sticker{Name: "tech"},
					},
				},
			)

			includeStickers := []string{"non-tech"}
			excludeStickers := []string{}

			filteredCards := cards.FilterByStickerNames(
				includeStickers,
				excludeStickers,
			)

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("a non-tech card"))
		})

		It("should exclude cards", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{
					Title: "a non-tech card",
					Stickers: []rubbernecker.Sticker{
						rubbernecker.Sticker{Name: "non-tech"},
					},
				},
				&rubbernecker.Card{
					Title: "a bug",
					Stickers: []rubbernecker.Sticker{
						rubbernecker.Sticker{Name: "bug"},
					},
				},
			)

			includeStickers := []string{}
			excludeStickers := []string{"bug"}

			filteredCards := cards.FilterByStickerNames(
				includeStickers,
				excludeStickers,
			)

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("a non-tech card"))
		})

		It("should prioritise include over exclude", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{
					Title: "a non-tech card",
					Stickers: []rubbernecker.Sticker{
						rubbernecker.Sticker{Name: "non-tech"},
					},
				},
				&rubbernecker.Card{
					Title: "a tech card",
					Stickers: []rubbernecker.Sticker{
						rubbernecker.Sticker{Name: "tech"},
					},
				},
			)

			includeStickers := []string{"non-tech"}
			excludeStickers := []string{"tech"}

			filteredCards := cards.FilterByStickerNames(
				includeStickers,
				excludeStickers,
			)

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("a non-tech card"))
		})
	})
	Context("Card Filtering By Text Filters", func() {
		It("should not do anything if there are no filters", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{},
				&rubbernecker.Card{},
				&rubbernecker.Card{},
			)

			filteredCards := cards.FilterByTextFilters([]string{})

			Expect(filteredCards).To(HaveLen(3))
		})

		It("should implement person filters", func() {
			members := make(rubbernecker.Members, 0)
			members[1] = &rubbernecker.Member{Name: "Rubber Necker"}

			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{Title: "a-card", Assignees: members},
				&rubbernecker.Card{},
				&rubbernecker.Card{},
			)

			filteredCards := cards.FilterByTextFilters([]string{
				"person:rubber",
			})

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("a-card"))
		})

		It("should implement multiple person filters", func() {
			members1 := make(rubbernecker.Members, 0)
			members1[1] = &rubbernecker.Member{Name: "Rubber Necker"}

			members2 := make(rubbernecker.Members, 0)
			members2[1] = &rubbernecker.Member{Name: "Necker"}

			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{Title: "a-card", Assignees: members1},
				&rubbernecker.Card{Title: "b-card", Assignees: members2},
				&rubbernecker.Card{},
			)

			filteredCards := cards.FilterByTextFilters([]string{
				"person:necker", "person:rubber",
			})

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("a-card"))
		})

		It("should implement title filters", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{Title: "a-card"},
				&rubbernecker.Card{Title: "b-card"},
			)

			filteredCards := cards.FilterByTextFilters([]string{"title:b"})

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("b-card"))
		})

		It("should implement multiple title filters", func() {
			cards := make(rubbernecker.Cards, 0)
			cards = append(
				cards,
				&rubbernecker.Card{Title: "a-card"},
				&rubbernecker.Card{Title: "b-card"},
			)

			filteredCards := cards.FilterByTextFilters([]string{
				"title:card", "title:b",
			})

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("b-card"))
		})
	})
})
