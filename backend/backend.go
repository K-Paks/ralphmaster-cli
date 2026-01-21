package backend

import "context"

// Backend defines the interface for issue tracking operations.
// Implementations can target different platforms (GitHub, Linear, Jira, etc.)
type Backend interface {
	// Issue operations
	CreateIssue(ctx context.Context, title, description string) error
	GetIssue(ctx context.Context, issueNumber int) (*Issue, error)
	UpdateIssueTitle(ctx context.Context, issueNumber int, title string) error
	UpdateIssueBody(ctx context.Context, issueNumber int, body string) error
	CloseIssue(ctx context.Context, issueNumber int) error

	// Comment operations
	GetComments(ctx context.Context, issueNumber int) ([]Comment, error)
	AddComment(ctx context.Context, issueNumber int, body string) error
	EditComment(ctx context.Context, commentID int, body string) error
	DeleteComment(ctx context.Context, commentID int) error

	// Name returns the backend identifier
	Name() string
}
