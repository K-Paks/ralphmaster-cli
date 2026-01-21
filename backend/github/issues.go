package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"ralphmaster/backend"
)

// CreateIssue creates a new issue with the given title and description.
func (g *GitHubBackend) CreateIssue(ctx context.Context, title, description string) error {
	ghArgs := []string{"issue", "create", "--title", title, "--body", description}
	ghArgs = append(ghArgs, g.repoArgs()...)

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetIssue retrieves an issue by its number.
func (g *GitHubBackend) GetIssue(ctx context.Context, issueNumber int) (*backend.Issue, error) {
	endpoint := fmt.Sprintf("%s/issues/%d", g.apiRepoPath(), issueNumber)
	ghArgs := []string{"api", endpoint}

	out, err := exec.CommandContext(ctx, "gh", ghArgs...).Output()
	if err != nil {
		return nil, err
	}

	var issue backend.Issue
	if err := json.Unmarshal(out, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// UpdateIssueBody updates the body of an issue.
func (g *GitHubBackend) UpdateIssueBody(ctx context.Context, issueNumber int, body string) error {
	ghArgs := []string{"issue", "edit", fmt.Sprintf("%d", issueNumber), "--body", body}
	ghArgs = append(ghArgs, g.repoArgs()...)

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// CloseIssue closes an issue by its number.
func (g *GitHubBackend) CloseIssue(ctx context.Context, issueNumber int) error {
	ghArgs := []string{"issue", "close", fmt.Sprintf("%d", issueNumber)}
	ghArgs = append(ghArgs, g.repoArgs()...)

	cmd := exec.CommandContext(ctx, "gh", ghArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
