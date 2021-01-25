package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use: "json2line",
	Run: func(c *cobra.Command, args []string) {},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		printInformationf("could not execute command: %v", err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("config", "c", "", "config file that should be used. <FILE>")
	RootCmd.PersistentFlags().StringP("formatter", "f", "", "formatter that should be used. <NAME>")
	RootCmd.PersistentFlags().StringP("adhoc", "o", "", "ad hoc format string.")
	RootCmd.PersistentFlags().StringArrayP("replacement", "r", []string{}, "ad hoc replacements.")
}
