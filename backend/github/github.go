package github

import (
	"ralphmaster/backend"
)

const backendName = "github"

// GitHubBackend implements the Backend interface using the gh CLI.
type GitHubBackend struct {
	repo string // "owner/repo" format, empty to use git remote detection
}

// New creates a new GitHub backend instance.
func New(cfg backend.Config) (backend.Backend, error) {
	return &GitHubBackend{
		repo: cfg.Repository,
	}, nil
}

// Name returns the backend identifier.
func (g *GitHubBackend) Name() string {
	return backendName
}

// repoArgs returns the --repo flag arguments if a repository is configured.
func (g *GitHubBackend) repoArgs() []string {
	if g.repo != "" {
		return []string{"--repo", g.repo}
	}
	return nil
}

// apiRepoPath returns the API path with either {owner}/{repo} placeholder
// or the explicit repo path.
func (g *GitHubBackend) apiRepoPath() string {
	if g.repo != "" {
		return "repos/" + g.repo
	}
	return "repos/{owner}/{repo}"
}

func init() {
	backend.Register(backendName, New)
}
