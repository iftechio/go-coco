package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade go-coco",
	Long:  `Install latest version of go-coco`,
	Run: func(cmd *cobra.Command, args []string) {
		report(flagSuccess, "upgrading...")
		err := exec.Command("go", "install", "github.com/iftechio/go-coco@latest").Run()
		report(flagSuccess, "done")
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
