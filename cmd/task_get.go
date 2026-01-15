package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	taskGetIssue int
	taskGetTask  int
)

var taskGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a task by ID from an issue",
	Run: func(cmd *cobra.Command, args []string) {
		if taskGetIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if taskGetTask == 0 {
			fmt.Fprintln(os.Stderr, "Error: --task is required")
			os.Exit(1)
		}

		comments, err := getComments(taskGetIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		idPattern := regexp.MustCompile(`(?m)^id: (\d+)`)

		for _, c := range comments {
			if !isTaskComment(c.Body) {
				continue
			}

			idMatch := idPattern.FindStringSubmatch(c.Body)
			if len(idMatch) < 2 {
				continue
			}

			taskID, err := strconv.Atoi(idMatch[1])
			if err != nil || taskID != taskGetTask {
				continue
			}

			// Found the task, print its details
			printTaskDetails(c.Body)
			return
		}

		fmt.Fprintf(os.Stderr, "Task #%d not found on issue #%d\n", taskGetTask, taskGetIssue)
		os.Exit(1)
	},
}

func printTaskDetails(body string) {
	// Determine status
	status := "UNDONE"
	if strings.HasPrefix(body, "[DONE]") {
		status = "DONE"
	}

	patterns := map[string]*regexp.Regexp{
		"id":         regexp.MustCompile(`(?m)^id: (.+)`),
		"model":      regexp.MustCompile(`(?m)^model: (.+)`),
		"goal":       regexp.MustCompile(`(?m)^goal: (.+)`),
		"comments":   regexp.MustCompile(`(?m)^comments: (.+)`),
		"references": regexp.MustCompile(`(?m)^references: (.+)`),
		"unit-tests": regexp.MustCompile(`(?m)^unit-tests: (.+)`),
		"e2e-tests":  regexp.MustCompile(`(?m)^e2e-tests: (.+)`),
		"work_done":  regexp.MustCompile(`(?m)^work_done: (.+)`),
		"commit":     regexp.MustCompile(`(?m)^commit: (.+)`),
	}

	fmt.Printf("status: %s\n", status)

	// Print fields in a logical order
	fields := []string{"id", "model", "goal", "comments", "references", "unit-tests", "e2e-tests", "work_done", "commit"}
	for _, field := range fields {
		pattern := patterns[field]
		match := pattern.FindStringSubmatch(body)
		if len(match) >= 2 {
			fmt.Printf("%s: %s\n", field, strings.TrimSpace(match[1]))
		}
	}
}

func init() {
	taskGetCmd.Flags().IntVar(&taskGetIssue, "issue", 0, "Issue number (required)")
	taskGetCmd.Flags().IntVar(&taskGetTask, "task", 0, "Task number (required)")
	taskCmd.AddCommand(taskGetCmd)
}
