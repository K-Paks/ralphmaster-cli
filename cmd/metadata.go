package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	metadataIssue    int
	metadataBranch   string
	metadataOverview string
	metadataForce    bool
)

type Comment struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Add or replace metadata comment on an issue",
	Run: func(cmd *cobra.Command, args []string) {
		if metadataIssue == 0 {
			fmt.Fprintln(os.Stderr, "Error: --issue is required")
			os.Exit(1)
		}
		if metadataBranch == "" {
			fmt.Fprintln(os.Stderr, "Error: --branch is required")
			os.Exit(1)
		}
		if metadataOverview == "" {
			fmt.Fprintln(os.Stderr, "Error: --overview is required")
			os.Exit(1)
		}

		comments, err := getComments(metadataIssue)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		metadataBody := fmt.Sprintf("[METADATA]\nbranch: %s\noverview: %s", metadataBranch, metadataOverview)

		if len(comments) == 0 {
			if err := addComment(metadataIssue, metadataBody); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding metadata comment: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Metadata comment added successfully")
			return
		}

		if len(comments) == 1 && strings.HasPrefix(comments[0].Body, "[METADATA]") {
			if !metadataForce {
				fmt.Println("Metadata comment already exists. Use --force to replace.")
				return
			}
			if err := deleteComment(comments[0].ID); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting existing metadata comment: %v\n", err)
				os.Exit(1)
			}
			if err := addComment(metadataIssue, metadataBody); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding metadata comment: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Metadata comment replaced successfully")
			return
		}

		fmt.Fprintln(os.Stderr, "Error: Issue has existing comments that are not metadata. Cannot add metadata.")
		os.Exit(1)
	},
}

func getComments(issueNumber int) ([]Comment, error) {
	ghArgs := []string{"api", fmt.Sprintf("repos/{owner}/{repo}/issues/%d/comments", issueNumber)}
	if repo != "" {
		ghArgs = []string{"api", fmt.Sprintf("repos/%s/issues/%d/comments", repo, issueNumber)}
	}

	out, err := exec.Command("gh", ghArgs...).Output()
	if err != nil {
		return nil, err
	}

	var comments []Comment
	if err := json.Unmarshal(out, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func addComment(issueNumber int, body string) error {
	ghArgs := []string{"issue", "comment", fmt.Sprintf("%d", issueNumber), "--body", body}
	if repo != "" {
		ghArgs = append(ghArgs, "--repo", repo)
	}

	cmd := exec.Command("gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func deleteComment(commentID int) error {
	endpoint := fmt.Sprintf("repos/{owner}/{repo}/issues/comments/%d", commentID)
	if repo != "" {
		endpoint = fmt.Sprintf("repos/%s/issues/comments/%d", repo, commentID)
	}

	ghArgs := []string{"api", endpoint, "-X", "DELETE"}
	cmd := exec.Command("gh", ghArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	metadataCmd.Flags().IntVar(&metadataIssue, "issue", 0, "Issue number (required)")
	metadataCmd.Flags().StringVar(&metadataBranch, "branch", "", "Branch name for the work (required)")
	metadataCmd.Flags().StringVar(&metadataOverview, "overview", "", "Single-line task overview (required)")
	metadataCmd.Flags().BoolVar(&metadataForce, "force", false, "Replace existing metadata comment")
	rootCmd.AddCommand(metadataCmd)
}
