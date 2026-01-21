package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"ralphmaster/backend"
)

var (
	taskDoneIssue    int
	taskDoneTask     int
	taskDoneCommit   string
	taskDoneWorkDone string
)

var taskDoneCmd = &cobra.Command{
	Use:   "done",
	Short: "Mark a task as completed",
	Run: func(cmd *cobra.Command, args []string) {
		if taskDoneIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if taskDoneTask == 0 {
			fmt.Fprintln(os.Stderr, "Error: --task is required")
			os.Exit(1)
		}
		if taskDoneCommit == "" {
			fmt.Fprintln(os.Stderr, "Error: --commit is required")
			os.Exit(1)
		}
		if taskDoneWorkDone == "" {
			fmt.Fprintln(os.Stderr, "Error: --work-done is required")
			os.Exit(1)
		}

		ctx := context.Background()
		b := GetBackend()

		comments, err := b.GetComments(ctx, taskDoneIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		idPattern := regexp.MustCompile(fmt.Sprintf(`(?m)^id: %d$`, taskDoneTask))

		var targetComment *backend.Comment
		for i, c := range comments {
			if strings.HasPrefix(c.Body, "[UNDONE]") && idPattern.MatchString(c.Body) {
				targetComment = &comments[i]
				break
			}
		}

		if targetComment == nil {
			fmt.Fprintf(os.Stderr, "Error: Task #%d not found or already completed\n", taskDoneTask)
			os.Exit(1)
		}

		newBody := strings.Replace(targetComment.Body, "[UNDONE]", "[DONE]", 1)
		newBody += fmt.Sprintf("\n\nwork_done: %s\ncommit: %s", taskDoneWorkDone, taskDoneCommit)

		if err := b.EditComment(ctx, targetComment.ID, newBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task comment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d marked as done\n", taskDoneTask)
	},
}

func init() {
	taskDoneCmd.Flags().IntVar(&taskDoneIssue, "issue", 0, "Issue number (required)")
	taskDoneCmd.Flags().IntVar(&taskDoneTask, "task", 0, "Task number to mark done (required)")
	taskDoneCmd.Flags().StringVar(&taskDoneCommit, "commit", "", "Commit hash/reference (required)")
	taskDoneCmd.Flags().StringVar(&taskDoneWorkDone, "work-done", "", "Summary of completed work (required)")
	taskCmd.AddCommand(taskDoneCmd)
}
