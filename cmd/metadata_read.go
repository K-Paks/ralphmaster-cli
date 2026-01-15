package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var metadataReadIssue int

var metadataReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read metadata comment from an issue",
	Run: func(cmd *cobra.Command, args []string) {
		if metadataReadIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}

		comments, err := getComments(metadataReadIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		for _, c := range comments {
			if strings.HasPrefix(c.Body, "[METADATA]") {
				fmt.Println(c.Body)
				return
			}
		}

		fmt.Fprintln(os.Stderr, "No metadata comment found on this issue")
		os.Exit(1)
	},
}

func init() {
	metadataReadCmd.Flags().IntVar(&metadataReadIssue, "issue", 0, "Issue number (required)")
	metadataCmd.AddCommand(metadataReadCmd)
}
