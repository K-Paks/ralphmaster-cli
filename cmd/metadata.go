package cmd

import (
	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Manage metadata comments on GitHub issues",
	Long:  `Metadata commands for adding and reading metadata on GitHub issues.`,
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
