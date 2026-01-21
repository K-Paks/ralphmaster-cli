package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var startIssue int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Mark an issue as in progress",
	Run: func(cmd *cobra.Command, args []string) {
		if startIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}

		ctx := context.Background()
		b := GetBackend()

		issue, err := b.GetIssue(ctx, startIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		newBody := "[IN PROGRESS]\n" + issue.Body

		if err := b.UpdateIssueBody(ctx, startIssue, newBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Issue #%d marked as in progress\n", startIssue)
	},
}

func init() {
	startCmd.Flags().IntVar(&startIssue, "issue", 0, "Issue number (required)")
	rootCmd.AddCommand(startCmd)
}
