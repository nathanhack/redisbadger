package cmd

import (
	"github.com/spf13/cobra"
)

// toolsCmd represents the tools command
var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "RedisBadger commandline tools",
	Long:  `A set of tools used with the RedisBadger files.`,
}

func init() {
	rootCmd.AddCommand(toolsCmd)
}
