package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	newTitle       string
	newDescription string
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new GitHub issue",
	Run: func(cmd *cobra.Command, args []string) {
		if newTitle == "" {
			fmt.Fprintln(os.Stderr, "Error: --title is required")
			os.Exit(1)
		}
		if newDescription == "" {
			fmt.Fprintln(os.Stderr, "Error: --description is required")
			os.Exit(1)
		}

		ctx := context.Background()
		if err := GetBackend().CreateIssue(ctx, "[UNDONE] "+newTitle, newDescription); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating issue: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	newCmd.Flags().StringVar(&newTitle, "title", "", "Issue title (required)")
	newCmd.Flags().StringVar(&newDescription, "description", "", "Issue description in markdown (required)")
	rootCmd.AddCommand(newCmd)
}
