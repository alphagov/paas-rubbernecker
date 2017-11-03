package rubbernecker

// Support struct will contain any useful information, relevant to our users.
type Support struct {
	Type   string `json:"type,omitempty"`
	Member string `json:"member,omitempty"`
}

// SupportRota will contain a unique list prefixed with a type of support.
type SupportRota map[string]*Support

// SupportService interface will establish a standard for any extension handling
// support data.
type SupportService interface {
	FetchSupport() error
	FlattenSupport() (SupportRota, error)
}
