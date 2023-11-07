package cmd

import "github.com/spf13/cobra"

var multipleLanguageCmd = &cobra.Command{
	Use:   "multi",
	Short: "Multi means to collect all languages to print their code information",
	Run: func(cmd *cobra.Command, args []string) {
		staticAll()
	},
}

func init() {
}

func staticAll() {}
