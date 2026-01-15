package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var repo string

var rootCmd = &cobra.Command{
	Use:   "ralphmaster",
	Short: "A CLI tool for LLM-GitHub interactions",
	Long:  `ralphmaster is a CLI tool that helps LLMs use correct interfaces when communicating with GitHub.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repo, "repo", "", "GitHub repository (owner/name). Auto-detected from git remote if not specified.")
}
