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

		taskPattern := regexp.MustCompile(`^\[UNDONE\]\n(\d+),(\w+),([^,]*),([^,]*),(.*)$`)
		foundTasks := false

		for _, c := range comments {
			if !strings.HasPrefix(c.Body, "[UNDONE]") {
				continue
			}

			matches := taskPattern.FindStringSubmatch(c.Body)
			if len(matches) >= 4 {
				foundTasks = true
				number := matches[1]
				model := matches[2]
				goal := matches[3]
				fmt.Printf("#%s [%s] %s\n", number, model, goal)
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
