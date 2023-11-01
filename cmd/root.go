package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codeLines",
	Short: "CodeLines is a tool to count lines of code of many programming languages",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(singleLanguageCmd)
}
