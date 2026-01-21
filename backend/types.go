package backend

// Issue represents an issue from any backend.
type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Comment represents a comment on an issue from any backend.
type Comment struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// Config holds configuration for backend initialization.
type Config struct {
	Type       string // Backend type: "github", "linear", "file", etc.
	Repository string // For GitHub: "owner/repo"
}
