package cmd

import (
	"github.com/spf13/cobra"
)

// aofCmd represents the aofchecker command
var aofCmd = &cobra.Command{
	Use:   "aof",
	Short: "Tools specific for AOF files",
	Long:  `Tools specific for AOF files`,
}

func init() {
	toolsCmd.AddCommand(aofCmd)
}
