package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var doneIssue int

type Issue struct {
	Body string `json:"body"`
}

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark an issue as done and close it",
	Run: func(cmd *cobra.Command, args []string) {
		if doneIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}

		issue, err := getIssue(doneIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		newBody := issue.Body
		if strings.HasPrefix(newBody, "[IN PROGRESS]\n") {
			newBody = strings.TrimPrefix(newBody, "[IN PROGRESS]\n")
		}
		newBody = "[DONE]\n" + newBody

		if err := updateIssueBody(doneIssue, newBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating issue: %v\n", err)
			os.Exit(1)
		}

		if err := closeIssue(doneIssue); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Issue #%d marked as done and closed\n", doneIssue)
	},
}

func getIssue(issueNumber int) (*Issue, error) {
	ghArgs := []string{"api", fmt.Sprintf("repos/{owner}/{repo}/issues/%d", issueNumber)}
	if repo != "" {
		ghArgs = []string{"api", fmt.Sprintf("repos/%s/issues/%d", repo, issueNumber)}
	}

	out, err := exec.Command("gh", ghArgs...).Output()
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(out, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

func updateIssueBody(issueNumber int, body string) error {
	ghArgs := []string{"issue", "edit", fmt.Sprintf("%d", issueNumber), "--body", body}
	if repo != "" {
		ghArgs = append(ghArgs, "--repo", repo)
	}

	cmd := exec.Command("gh", ghArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func closeIssue(issueNumber int) error {
	ghArgs := []string{"issue", "close", fmt.Sprintf("%d", issueNumber)}
	if repo != "" {
		ghArgs = append(ghArgs, "--repo", repo)
	}

	cmd := exec.Command("gh", ghArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	doneCmd.Flags().IntVar(&doneIssue, "issue", 0, "Issue number (required)")
	rootCmd.AddCommand(doneCmd)
}
