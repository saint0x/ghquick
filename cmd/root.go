package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghquick",
	Short: "ghquick - Lightning fast GitHub operations with AI-powered automation",
	Long: `ghquick is a CLI tool that automates GitHub operations with AI assistance.
It optimizes for speed and developer experience, making git operations instant.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.ghquick.yaml)")
}
