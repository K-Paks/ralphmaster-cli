package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	metadataAddIssue    int
	metadataAddBranch   string
	metadataAddOverview string
	metadataAddForce    bool
)

var metadataAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or replace metadata comment on an issue",
	Run: func(cmd *cobra.Command, args []string) {
		if metadataAddIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if metadataAddBranch == "" {
			fmt.Fprintln(os.Stderr, "Error: --branch is required")
			os.Exit(1)
		}
		if metadataAddOverview == "" {
			fmt.Fprintln(os.Stderr, "Error: --overview is required")
			os.Exit(1)
		}

		comments, err := getComments(metadataAddIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		metadataBody := fmt.Sprintf("[METADATA]\nbranch: %s\noverview: %s", metadataAddBranch, metadataAddOverview)

		if len(comments) == 0 {
			if err := addComment(metadataAddIssue, metadataBody); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding metadata comment: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Metadata comment added successfully")
			return
		}

		if len(comments) == 1 && strings.HasPrefix(comments[0].Body, "[METADATA]") {
			if !metadataAddForce {
				fmt.Println("Metadata comment already exists. Use --force to replace.")
				return
			}
			if err := deleteComment(comments[0].ID); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting existing metadata comment: %v\n", err)
				os.Exit(1)
			}
			if err := addComment(metadataAddIssue, metadataBody); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding metadata comment: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Metadata comment replaced successfully")
			return
		}

		fmt.Fprintln(os.Stderr, "Error: Issue has existing comments that are not metadata. Cannot add metadata.")
		os.Exit(1)
	},
}

func init() {
	metadataAddCmd.Flags().IntVar(&metadataAddIssue, "issue", 0, "Issue number (required)")
	metadataAddCmd.Flags().StringVar(&metadataAddBranch, "branch", "", "Branch name for the work (required)")
	metadataAddCmd.Flags().StringVar(&metadataAddOverview, "overview", "", "Single-line task overview (required)")
	metadataAddCmd.Flags().BoolVar(&metadataAddForce, "force", false, "Replace existing metadata comment")
	metadataCmd.AddCommand(metadataAddCmd)
}
