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
	taskUpdateIssue      int
	taskUpdateTask       int
	taskUpdateGoal       string
	taskUpdateComments   string
	taskUpdateReferences string
	taskUpdateModel      string
	taskUpdateUnitTests  string
	taskUpdateE2ETests   string
)

var taskUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing task's content",
	Run: func(cmd *cobra.Command, args []string) {
		if taskUpdateIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if taskUpdateTask == 0 {
			fmt.Fprintln(os.Stderr, "Error: --task is required")
			os.Exit(1)
		}

		// Check that at least one field is provided to update
		if taskUpdateGoal == "" && taskUpdateComments == "" && taskUpdateReferences == "" &&
			taskUpdateModel == "" && taskUpdateUnitTests == "" && taskUpdateE2ETests == "" {
			fmt.Fprintln(os.Stderr, "Error: at least one field to update is required (--goal, --comments, --references, --model, --unit-tests, --e2e-tests)")
			os.Exit(1)
		}

		// Validate model if provided
		if taskUpdateModel != "" && !validModels[taskUpdateModel] {
			fmt.Fprintln(os.Stderr, "Error: --model must be one of: opus, sonnet, haiku")
			os.Exit(1)
		}

		ctx := context.Background()
		b := GetBackend()

		comments, err := b.GetComments(ctx, taskUpdateIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		idPattern := regexp.MustCompile(fmt.Sprintf(`(?m)^id: %d$`, taskUpdateTask))

		var targetComment *backend.Comment
		for i, c := range comments {
			if strings.HasPrefix(c.Body, "[UNDONE]") && idPattern.MatchString(c.Body) {
				targetComment = &comments[i]
				break
			}
		}

		if targetComment == nil {
			fmt.Fprintf(os.Stderr, "Error: Task #%d not found or already completed (only UNDONE tasks can be updated)\n", taskUpdateTask)
			os.Exit(1)
		}

		newBody := targetComment.Body

		// Update each field if provided
		if taskUpdateGoal != "" {
			newBody = updateField(newBody, "goal", taskUpdateGoal)
		}
		if taskUpdateComments != "" {
			newBody = updateField(newBody, "comments", taskUpdateComments)
		}
		if taskUpdateReferences != "" {
			newBody = updateField(newBody, "references", taskUpdateReferences)
		}
		if taskUpdateModel != "" {
			newBody = updateField(newBody, "model", taskUpdateModel)
		}
		if taskUpdateUnitTests != "" {
			newBody = updateField(newBody, "unit-tests", taskUpdateUnitTests)
		}
		if taskUpdateE2ETests != "" {
			newBody = updateField(newBody, "e2e-tests", taskUpdateE2ETests)
		}

		if err := b.EditComment(ctx, targetComment.ID, newBody); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating task comment: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d updated successfully\n", taskUpdateTask)
	},
}

func updateField(body, field, value string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^%s: .+$`, regexp.QuoteMeta(field)))
	if pattern.MatchString(body) {
		return pattern.ReplaceAllString(body, fmt.Sprintf("%s: %s", field, value))
	}
	// Field doesn't exist, append it
	return body + fmt.Sprintf("\n%s: %s", field, value)
}

func init() {
	taskUpdateCmd.Flags().IntVar(&taskUpdateIssue, "issue", 0, "Issue number (required)")
	taskUpdateCmd.Flags().IntVar(&taskUpdateTask, "task", 0, "Task ID to update (required)")
	taskUpdateCmd.Flags().StringVar(&taskUpdateGoal, "goal", "", "Update task goal")
	taskUpdateCmd.Flags().StringVar(&taskUpdateComments, "comments", "", "Update comments")
	taskUpdateCmd.Flags().StringVar(&taskUpdateReferences, "references", "", "Update file references")
	taskUpdateCmd.Flags().StringVar(&taskUpdateModel, "model", "", "Update model (opus|sonnet|haiku)")
	taskUpdateCmd.Flags().StringVar(&taskUpdateUnitTests, "unit-tests", "", "Update unit test instructions")
	taskUpdateCmd.Flags().StringVar(&taskUpdateE2ETests, "e2e-tests", "", "Update e2e test instructions")
	taskCmd.AddCommand(taskUpdateCmd)
}
