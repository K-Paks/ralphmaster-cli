package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
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

		comments, err := getComments(taskDoneIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		taskPattern := regexp.MustCompile(fmt.Sprintf(`^\[UNDONE\]\n%d,(.*)$`, taskDoneTask))

		var targetComment *Comment
		for i, c := range comments {
			if taskPattern.MatchString(c.Body) {
				targetComment = &comments[i]
				break
			}
		}

		if targetComment == nil {
			fmt.Fprintf(os.Stderr, "Error: Task #%d not found or already completed\n", taskDoneTask)
			os.Exit(1)
		}

		newBody := strings.Replace(targetComment.Body, "[UNDONE]", "[DONE]", 1)
		newBody += fmt.Sprintf("\n----------------\nWORK DONE:\n%s\nCOMMIT: %s", taskDoneWorkDone, taskDoneCommit)

		if err := editComment(targetComment.ID, newBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task comment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d marked as done\n", taskDoneTask)
	},
}

func editComment(commentID int, body string) error {
	endpoint := fmt.Sprintf("repos/{owner}/{repo}/issues/comments/%d", commentID)
	if repo != "" {
		endpoint = fmt.Sprintf("repos/%s/issues/comments/%d", repo, commentID)
	}

	ghArgs := []string{"api", endpoint, "-X", "PATCH", "-f", fmt.Sprintf("body=%s", body)}
	cmd := exec.Command("gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	taskDoneCmd.Flags().IntVar(&taskDoneIssue, "issue", 0, "Issue number (required)")
	taskDoneCmd.Flags().IntVar(&taskDoneTask, "task", 0, "Task number to mark done (required)")
	taskDoneCmd.Flags().StringVar(&taskDoneCommit, "commit", "", "Commit hash/reference (required)")
	taskDoneCmd.Flags().StringVar(&taskDoneWorkDone, "work-done", "", "Summary of completed work (required)")
	taskCmd.AddCommand(taskDoneCmd)
}
