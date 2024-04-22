package cmd

import (
	"fmt"
	"github.com/EscanBE/escan-request-redirector/constants"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command, it prints the current version of the binary
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show binary version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constants.APP_NAME)

		fmt.Printf("%-11s %s\n", "Version:", constants.VERSION)
		fmt.Printf("%-11s %s\n", "Commit:", constants.COMMIT_HASH)
		fmt.Printf("%-11s %s\n", "Build date:", constants.BUILD_DATE)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
