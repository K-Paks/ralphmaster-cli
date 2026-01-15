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
	taskAddIssue      int
	taskAddModel      string
	taskAddGoal       string
	taskAddComments   string
	taskAddReferences string
)

var validModels = map[string]bool{
	"opus":   true,
	"sonnet": true,
	"haiku":  true,
}

var taskAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task to an issue",
	Run: func(cmd *cobra.Command, args []string) {
		if taskAddIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if taskAddModel == "" {
			fmt.Fprintln(os.Stderr, "Error: --model is required")
			os.Exit(1)
		}
		if !validModels[taskAddModel] {
			fmt.Fprintln(os.Stderr, "Error: --model must be one of: opus, sonnet, haiku")
			os.Exit(1)
		}
		if taskAddGoal == "" {
			fmt.Fprintln(os.Stderr, "Error: --goal is required")
			os.Exit(1)
		}

		comments, err := getComments(taskAddIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		nextNumber := findNextTaskNumber(comments)

		taskBody := fmt.Sprintf("[UNDONE]\nid: %d\nmodel: %s\ngoal: %s",
			nextNumber, taskAddModel, taskAddGoal)
		if taskAddComments != "" {
			taskBody += fmt.Sprintf("\ncomments: %s", taskAddComments)
		}
		if taskAddReferences != "" {
			taskBody += fmt.Sprintf("\nreferences: %s", taskAddReferences)
		}

		if err := addComment(taskAddIssue, taskBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding task comment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d added successfully\n", nextNumber)
	},
}

func findNextTaskNumber(comments []Comment) int {
	maxNumber := 0
	idPattern := regexp.MustCompile(`(?m)^id: (\d+)`)

	for _, c := range comments {
		if !isTaskComment(c.Body) {
			continue
		}
		matches := idPattern.FindStringSubmatch(c.Body)
		if len(matches) >= 2 {
			num, err := strconv.Atoi(matches[1])
			if err == nil && num > maxNumber {
				maxNumber = num
			}
		}
	}
	return maxNumber + 1
}

func isTaskComment(body string) bool {
	return strings.HasPrefix(body, "[UNDONE]") || strings.HasPrefix(body, "[DONE]")
}

func init() {
	taskAddCmd.Flags().IntVar(&taskAddIssue, "issue", 0, "Issue number (required)")
	taskAddCmd.Flags().StringVar(&taskAddModel, "model", "", "Model: opus|sonnet|haiku (required)")
	taskAddCmd.Flags().StringVar(&taskAddGoal, "goal", "", "Task goal description (required)")
	taskAddCmd.Flags().StringVar(&taskAddComments, "comments", "", "Additional comments")
	taskAddCmd.Flags().StringVar(&taskAddReferences, "references", "", "File references (e.g., path/to/file.ts#10-17)")
	taskCmd.AddCommand(taskAddCmd)
}
