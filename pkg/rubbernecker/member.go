package rubbernecker

// Member will be a rubbernecker entity composed of the extension.
type Member struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Members will be a rubbernecker representative of all members.
type Members map[int]*Member

// MemberService interface will establish a standard for any extension handling
// support data.
type MemberService interface {
	FetchMembers() error
	FlattenMembers() (*Members, error)
}
