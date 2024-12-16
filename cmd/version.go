package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gkwa/myland/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of myland",
	Long:  `All software has versions. This is myland's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
