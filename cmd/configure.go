package cmd

import (
	"github.com/mdreem/json2line/configuration"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "configure",
	Run: func(c *cobra.Command, args []string) {
		configuration.Configuration(c, args)
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolP("show", "s", false, "show current configuration")
	configCmd.AddCommand(
		configuration.FormatterConfigCmd,
		configuration.ReplacementConfigCmd,
		configuration.InitConfigurationFileCmd,
		configuration.BufferSizeConfigCmd,
	)
}
