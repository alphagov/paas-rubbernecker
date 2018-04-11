package rubbernecker

// Support struct will contain any useful information, relevant to our users.
type Support struct {
	Type   string `json:"type,omitempty"`
	Member string `json:"member,omitempty"`
}

// SupportRota will contain a unique list prefixed with a type of support.
type SupportRota map[string]*Support

// Get returns with the given key or an empty value
func (s SupportRota) Get(key string) *Support {
	if support, ok := s[key]; ok {
		return support
	}
	return &Support{
		Type:   key,
		Member: "-",
	}
}

// SupportService interface will establish a standard for any extension handling
// support data.
type SupportService interface {
	FetchSupport() error
	FlattenSupport() (SupportRota, error)
}
