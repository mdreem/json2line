package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var replacementConfigCmd = &cobra.Command{
	Use: "replacement",
	Run: func(c *cobra.Command, args []string) {
		replacementConfiguration(c, args)
	},
}

func replacementConfiguration(c *cobra.Command, args []string) {
	printInformationf("writing configuration to: %s\n", viper.ConfigFileUsed())

	err := viper.WriteConfig()
	if err != nil {
		printInformationf("could not write config: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	replacementConfigCmd.PersistentFlags().StringArrayP("add", "a", []string{}, "add replacement configuration. <NAME> <PATTERN>")
	replacementConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete replacement configuration. <NAME>")
}
