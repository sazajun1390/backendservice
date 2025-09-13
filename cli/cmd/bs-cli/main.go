package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bs-cli",
	Short: "CLI for bs-cli",
}

func Execute() error {
	return rootCmd.Execute()
}
