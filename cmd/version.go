package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of go-coco",
	Long:  `All software has versions. This is go-coco's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("go-coco Generator v0.0.2")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
