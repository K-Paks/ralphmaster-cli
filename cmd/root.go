package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"ralphmaster/backend"
	_ "ralphmaster/backend/github" // Register GitHub backend
)

var (
	repo        string
	backendType string
	activeBackend backend.Backend
)

var rootCmd = &cobra.Command{
	Use:   "ralphmaster",
	Short: "A CLI tool for LLM-GitHub interactions",
	Long:  `ralphmaster is a CLI tool that helps LLMs use correct interfaces when communicating with GitHub.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg := backend.Config{
			Type:       backendType,
			Repository: repo,
		}
		b, err := backend.New(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize backend: %w", err)
		}
		activeBackend = b
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// GetBackend returns the initialized backend for use by subcommands.
func GetBackend() backend.Backend {
	return activeBackend
}

// stripStatusPrefix removes any status prefix ([UNDONE], [IN PROGRESS], [DONE]) from a title.
func stripStatusPrefix(title string) string {
	prefixes := []string{"[UNDONE] ", "[IN PROGRESS] ", "[DONE] "}
	for _, p := range prefixes {
		if strings.HasPrefix(title, p) {
			return strings.TrimPrefix(title, p)
		}
	}
	return title
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repo, "repo", "", "GitHub repository (owner/name). Auto-detected from git remote if not specified.")
	rootCmd.PersistentFlags().StringVar(&backendType, "backend", "github", "Backend type: github")
}
