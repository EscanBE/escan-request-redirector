package cmd

import (
	"github.com/spf13/cobra"
)

// startCmd represents the start command, it launches the main business logic of this app
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start request redirect service",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
