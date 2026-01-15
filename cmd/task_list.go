package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var taskListIssue int

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List undone tasks on an issue",
	Run: func(cmd *cobra.Command, args []string) {
		if taskListIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}

		comments, err := getComments(taskListIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		idPattern := regexp.MustCompile(`(?m)^id: (\d+)`)
		modelPattern := regexp.MustCompile(`(?m)^model: (\w+)`)
		goalPattern := regexp.MustCompile(`(?m)^goal: (.+)`)
		foundTasks := false

		for _, c := range comments {
			if !strings.HasPrefix(c.Body, "[UNDONE]") {
				continue
			}

			idMatch := idPattern.FindStringSubmatch(c.Body)
			modelMatch := modelPattern.FindStringSubmatch(c.Body)
			goalMatch := goalPattern.FindStringSubmatch(c.Body)

			if len(idMatch) >= 2 && len(modelMatch) >= 2 && len(goalMatch) >= 2 {
				foundTasks = true
				fmt.Printf("#%s [%s] %s\n", idMatch[1], modelMatch[1], goalMatch[1])
			}
		}

		if !foundTasks {
			fmt.Println("No undone tasks found")
		}
	},
}

func init() {
	taskListCmd.Flags().IntVar(&taskListIssue, "issue", 0, "Issue number (required)")
	taskCmd.AddCommand(taskListCmd)
}
