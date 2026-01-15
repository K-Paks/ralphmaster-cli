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
	taskAddUnitTests  string
	taskAddE2ETests   string
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
		if taskAddComments == "" {
			fmt.Fprintln(os.Stderr, "Error: --comments is required")
			os.Exit(1)
		}
		if taskAddReferences == "" {
			fmt.Fprintln(os.Stderr, "Error: --references is required")
			os.Exit(1)
		}

		comments, err := getComments(taskAddIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		nextNumber := findNextTaskNumber(comments)

		taskBody := fmt.Sprintf("[UNDONE]\nid: %d\nmodel: %s\ngoal: %s\ncomments: %s\nreferences: %s",
			nextNumber, taskAddModel, taskAddGoal, taskAddComments, taskAddReferences)

		if taskAddUnitTests != "" {
			taskBody += fmt.Sprintf("\nunit-tests: %s", taskAddUnitTests)
		}
		if taskAddE2ETests != "" {
			taskBody += fmt.Sprintf("\ne2e-tests: %s", taskAddE2ETests)
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
	taskAddCmd.Flags().StringVar(&taskAddComments, "comments", "", "Additional comments (required)")
	taskAddCmd.Flags().StringVar(&taskAddReferences, "references", "", "File references, e.g., path/to/file.ts#10-17 (required)")
	taskAddCmd.Flags().StringVar(&taskAddUnitTests, "unit-tests", "", "Unit test instructions (optional)")
	taskAddCmd.Flags().StringVar(&taskAddE2ETests, "e2e-tests", "", "E2E test instructions (optional)")
	taskCmd.AddCommand(taskAddCmd)
}
