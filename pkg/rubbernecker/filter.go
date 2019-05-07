package rubbernecker

import "fmt"

type Filter struct {
	StickerName string
	FilterText string
	Exclude bool
}

func (f *Filter) IsApplied(applied []string) bool {
	for _, sticker := range applied {
		if sticker == f.StickerName {
			return true
		}
	}

	return false
}

func (f *Filter) Href() string {
	if f.StickerName == "" {
		return "?"
	}

	key := "include-sticker"
	if f.Exclude {
		key = "exclude-sticker"
	}

	return fmt.Sprintf("?%s=%s", key, f.StickerName)
}
