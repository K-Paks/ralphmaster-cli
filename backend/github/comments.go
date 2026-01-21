package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"ralphmaster/backend"
)

// GetComments retrieves all comments for an issue.
func (g *GitHubBackend) GetComments(ctx context.Context, issueNumber int) ([]backend.Comment, error) {
	endpoint := fmt.Sprintf("%s/issues/%d/comments", g.apiRepoPath(), issueNumber)
	ghArgs := []string{"api", endpoint}

	out, err := exec.CommandContext(ctx, "gh", ghArgs...).Output()
	if err != nil {
		return nil, err
	}

	var comments []backend.Comment
	if err := json.Unmarshal(out, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

// AddComment adds a new comment to an issue.
func (g *GitHubBackend) AddComment(ctx context.Context, issueNumber int, body string) error {
	ghArgs := []string{"issue", "comment", fmt.Sprintf("%d", issueNumber), "--body", body}
	ghArgs = append(ghArgs, g.repoArgs()...)

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// EditComment updates an existing comment by its ID.
func (g *GitHubBackend) EditComment(ctx context.Context, commentID int, body string) error {
	endpoint := fmt.Sprintf("%s/issues/comments/%d", g.apiRepoPath(), commentID)
	ghArgs := []string{"api", endpoint, "-X", "PATCH", "-f", fmt.Sprintf("body=%s", body)}

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DeleteComment deletes a comment by its ID.
func (g *GitHubBackend) DeleteComment(ctx context.Context, commentID int) error {
	endpoint := fmt.Sprintf("%s/issues/comments/%d", g.apiRepoPath(), commentID)
	ghArgs := []string{"api", endpoint, "-X", "DELETE"}

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
