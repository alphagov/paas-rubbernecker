package rubbernecker

import "regexp"

// Sticker is a rubbernecker definition of labels.
type Sticker struct {
	Name     string
	Label    bool
	Regex    string
	Title    string
	Image    string
	Content  string
	Aliases  []string
	Class    string
	Priority int
}

// Matches will check if the sticker matches the query provided by the extension.
func (s Sticker) Matches(query string) (Sticker, bool) {
	if s.Regex != "" {
		reg := regexp.MustCompile(s.Regex)
		if reg.MatchString(query) {
			sticker := s
			sticker.Title = reg.ReplaceAllString(query, sticker.Title)
			sticker.Image = reg.ReplaceAllString(query, sticker.Image)
			sticker.Content = reg.ReplaceAllString(query, sticker.Content)
			sticker.Class = reg.ReplaceAllString(query, sticker.Class)
			return sticker, true
		}
	}

	if s.Name == query {
		return s, true
	}

	for _, alias := range s.Aliases {
		if alias == query {
			return s, true
		}
	}

	return Sticker{}, false
}

// Stickers is a simple slice of stickers
type Stickers []Sticker

// Has will run a quick check against the slice of stickers to see if one has
// been already added
func (ss Stickers) Has(query string) bool {
	_, ok := ss.Get(query)
	return ok
}

// Get will run a quick check against the slice of stickers to see if one has
// been already added and will return it.
func (ss Stickers) Get(query string) (Sticker, bool) {
	for _, s := range ss {
		if sticker, ok := s.Matches(query); ok {
			return sticker, true
		}
	}

	return Sticker{}, false
}

// Contains returns true if the list has a sticker with the given name
func (ss Stickers) Contains(name string) bool {
	for _, s := range ss {
		if s.Name == name {
			return true
		}
	}
	return false
}

// Len is the number of elements in the collection
// Needed for implementing sort.Interface
func (ss Stickers) Len() int {
	return len(ss)
}

// Swap swaps the elements with indexes i and j
// Needed for implementing sort.Interface
func (ss Stickers) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

// Less reports whether the element with index i should sort before the element with index j
// Needed for implementing sort.Interface
func (ss Stickers) Less(i, j int) bool {
	return ss[i].Priority > ss[j].Priority
}
