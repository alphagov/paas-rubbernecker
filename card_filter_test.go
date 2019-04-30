package main

import (
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Card Filtering", func() {
	Context("Card Filtering", func() {
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

			filteredCards := filterCards(cards, includeStickers, excludeStickers)

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

			filteredCards := filterCards(cards, includeStickers, excludeStickers)

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

			filteredCards := filterCards(cards, includeStickers, excludeStickers)

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

			filteredCards := filterCards(cards, includeStickers, excludeStickers)

			Expect(filteredCards).To(HaveLen(1))
			Expect(filteredCards[0].Title).To(Equal("a non-tech card"))
		})
	})
})
