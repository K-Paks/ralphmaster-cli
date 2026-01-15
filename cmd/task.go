package cmd

import (
	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage tasks on GitHub issues",
	Long:  `Task commands for adding, listing, and completing tasks on GitHub issues.`,
}

func init() {
	rootCmd.AddCommand(taskCmd)
}
