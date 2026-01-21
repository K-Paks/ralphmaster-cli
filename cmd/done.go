package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var doneIssue int

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark an issue as done and close it",
	Run: func(cmd *cobra.Command, args []string) {
		if doneIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}

		ctx := context.Background()
		b := GetBackend()

		issue, err := b.GetIssue(ctx, doneIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		newBody := issue.Body
		if strings.HasPrefix(newBody, "[IN PROGRESS]\n") {
			newBody = strings.TrimPrefix(newBody, "[IN PROGRESS]\n")
		}
		newBody = "[DONE]\n" + newBody

		if err := b.UpdateIssueBody(ctx, doneIssue, newBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating issue: %v\n", err)
			os.Exit(1)
		}

		if err := b.CloseIssue(ctx, doneIssue); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Issue #%d marked as done and closed\n", doneIssue)
	},
}

func init() {
	doneCmd.Flags().IntVar(&doneIssue, "issue", 0, "Issue number (required)")
	rootCmd.AddCommand(doneCmd)
}
