package main

import (
	"github.com/alphagov/paas-rubbernecker/pkg/rubbernecker"
)

func filterCards(
	cards rubbernecker.Cards,
	includeStickers []string, excludeStickers []string,
) rubbernecker.Cards {
	filteredCards := make(rubbernecker.Cards, 0)

	for _, card := range cards {
		shouldAdd := true

		if len(includeStickers) > 0 {
			shouldAdd = false
			for _, sticker := range card.Stickers {
				for _, includedStickerName := range includeStickers {
					if sticker.Name == includedStickerName {
						shouldAdd = true
					}
				}
			}
		} else if len(excludeStickers) > 0 {
			shouldAdd = true
			for _, sticker := range card.Stickers {
				for _, excludedStickerName := range excludeStickers {
					if sticker.Name == excludedStickerName {
						shouldAdd = false
					}
				}
			}
		}

		if shouldAdd {
			filteredCards = append(filteredCards, card)
		}
	}

	return filteredCards
}
