package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

		ghArgs := []string{"issue", "create", "--title", newTitle, "--body", newDescription}
		if repo != "" {
			ghArgs = append(ghArgs, "--repo", repo)
		}

		ghCmd := exec.Command("gh", ghArgs...)
		ghCmd.Stdout = os.Stdout
		ghCmd.Stderr = os.Stderr

		if err := ghCmd.Run(); err != nil {
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
