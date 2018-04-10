package rubbernecker

// Sticker is a rubbernecker definition of labels.
type Sticker struct {
	Name    string
	Title   string
	Image   string
	Content string
	Aliases []string
}

// Stickers is a simple slice of stickers
type Stickers []*Sticker

// Matches will check if the sticker mathes the query provided by the extension.
func (s *Sticker) Matches(query string) *Sticker {
	if s.Name == query {
		return s
	}

	for _, alias := range s.Aliases {
		if alias == query {
			return s
		}
	}

	return nil
}

// Get will run a quick check against the slice of stickers to see if one has
// been already added and will return it.
func (ss *Stickers) Get(query string) *Sticker {
	for _, s := range *ss {
		if sticker := s.Matches(query); sticker != nil {
			return sticker
		}
	}

	return nil
}
