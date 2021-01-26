package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var formatterConfigCmd = &cobra.Command{
	Use: "formatter",
	Run: func(c *cobra.Command, args []string) {
		formatterConfiguration(c, args)
	},
}

func formatterConfiguration(c *cobra.Command, args []string) {
	printInformationf("writing configuration to: %s\n", viper.ConfigFileUsed())

	err := viper.WriteConfig()
	if err != nil {
		printInformationf("could not write config: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	formatterConfigCmd.PersistentFlags().StringArrayP("add", "a", []string{}, "add formatter configuration. <NAME> <PATTERN>")
	formatterConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete formatter configuration. <NAME>")
}
