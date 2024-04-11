package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "go-coco",
		Short: "A generator for wire-based Micro Service",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
