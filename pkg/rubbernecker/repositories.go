package rubbernecker

// Repository will be a rubbernecker entity composed of the extension.
type Repository struct {
	Organisation string `json:"organisation"`
	Name         string `json:"name"`
}

// Repositories will be a rubbernecker representative of all repositories.
type Repositories []Repository
