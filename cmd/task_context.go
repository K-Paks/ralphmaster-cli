package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	taskContextIssue int
	taskContextTask  int
)

var taskContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Get context for a task including issue details and previous work",
	Run: func(cmd *cobra.Command, args []string) {
		if taskContextIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if taskContextTask == 0 {
			fmt.Fprintln(os.Stderr, "Error: --task is required")
			os.Exit(1)
		}

		ctx := context.Background()
		b := GetBackend()

		// Fetch issue details
		issue, err := b.GetIssue(ctx, taskContextIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		// Fetch comments
		comments, err := b.GetComments(ctx, taskContextIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		// Determine task status
		idPattern := regexp.MustCompile(`(?m)^id: (\d+)`)
		taskStatus := "[IN PROGRESS]"
		for _, c := range comments {
			if !strings.HasPrefix(c.Body, "[DONE]") && !strings.HasPrefix(c.Body, "[UNDONE]") {
				continue
			}
			idMatch := idPattern.FindStringSubmatch(c.Body)
			if len(idMatch) >= 2 {
				taskID, _ := strconv.Atoi(idMatch[1])
				if taskID == taskContextTask {
					if strings.HasPrefix(c.Body, "[DONE]") {
						taskStatus = "[DONE]"
					}
					break
				}
			}
		}

		// Output issue details with status
		fmt.Printf("title: %s %s\n", taskStatus, issue.Title)
		fmt.Printf("description: %s\n", issue.Body)

		// Find and parse metadata comment
		branchPattern := regexp.MustCompile(`(?m)^branch: (.+)`)
		overviewPattern := regexp.MustCompile(`(?m)^overview: (.+)`)

		for _, c := range comments {
			if strings.HasPrefix(c.Body, "[METADATA]") {
				if match := branchPattern.FindStringSubmatch(c.Body); len(match) >= 2 {
					fmt.Printf("branch: %s\n", strings.TrimSpace(match[1]))
				}
				if match := overviewPattern.FindStringSubmatch(c.Body); len(match) >= 2 {
					fmt.Printf("overview: %s\n", strings.TrimSpace(match[1]))
				}
				break
			}
		}

		// Collect work_done from completed tasks with id < taskContextTask
		workDonePattern := regexp.MustCompile(`(?m)^work_done: (.+)`)

		var previousWork []string
		for _, c := range comments {
			if !strings.HasPrefix(c.Body, "[DONE]") {
				continue
			}

			idMatch := idPattern.FindStringSubmatch(c.Body)
			if len(idMatch) < 2 {
				continue
			}

			taskID, err := strconv.Atoi(idMatch[1])
			if err != nil || taskID >= taskContextTask {
				continue
			}

			if match := workDonePattern.FindStringSubmatch(c.Body); len(match) >= 2 {
				previousWork = append(previousWork, fmt.Sprintf("- task %d: %s", taskID, strings.TrimSpace(match[1])))
			}
		}

		if len(previousWork) > 0 {
			fmt.Println("\nprevious_work:")
			for _, work := range previousWork {
				fmt.Println(work)
			}
		}
	},
}

func init() {
	taskContextCmd.Flags().IntVar(&taskContextIssue, "issue", 0, "Issue number (required)")
	taskContextCmd.Flags().IntVar(&taskContextTask, "task", 0, "Task number (context is for tasks with id < this)")
	taskCmd.AddCommand(taskContextCmd)
}
