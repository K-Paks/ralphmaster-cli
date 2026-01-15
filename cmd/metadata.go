package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

type Comment struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Manage metadata comments on GitHub issues",
	Long:  `Metadata commands for adding and reading metadata on GitHub issues.`,
}

func getComments(issueNumber int) ([]Comment, error) {
	ghArgs := []string{"api", fmt.Sprintf("repos/{owner}/{repo}/issues/%d/comments", issueNumber)}
	if repo != "" {
		ghArgs = []string{"api", fmt.Sprintf("repos/%s/issues/%d/comments", repo, issueNumber)}
	}

	out, err := exec.Command("gh", ghArgs...).Output()
	if err != nil {
		return nil, err
	}

	var comments []Comment
	if err := json.Unmarshal(out, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func addComment(issueNumber int, body string) error {
	ghArgs := []string{"issue", "comment", fmt.Sprintf("%d", issueNumber), "--body", body}
	if repo != "" {
		ghArgs = append(ghArgs, "--repo", repo)
	}

	cmd := exec.Command("gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func deleteComment(commentID int) error {
	endpoint := fmt.Sprintf("repos/{owner}/{repo}/issues/comments/%d", commentID)
	if repo != "" {
		endpoint = fmt.Sprintf("repos/%s/issues/comments/%d", repo, commentID)
	}

	ghArgs := []string{"api", endpoint, "-X", "DELETE"}
	cmd := exec.Command("gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
