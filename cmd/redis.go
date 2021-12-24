package cmd

import (
	"github.com/spf13/cobra"
)

// redisCmd represents the redis command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Redis tools",
	Long:  `Redis tools.`,
}

func init() {
	toolsCmd.AddCommand(redisCmd)
}
