package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use: "configure",
	Run: func(c *cobra.Command, args []string) {
		configuration(c, args)
	},
}

func configuration(c *cobra.Command, args []string) {
	show := getBoolean(c, "show")

	if show {
		fmt.Printf("Templates:\n")
		for key, value := range viper.GetStringMapString("templates") {
			fmt.Printf("'%v' -> '%v'\n", key, value)
		}
		fmt.Printf("\nReplacements:\n")
		for key, value := range viper.GetStringMapString("replacements") {
			fmt.Printf("'%v' -> '%v'\n", key, value)
		}
	}
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolP("show", "s", false, "show current configuration")
	configCmd.AddCommand(formatterConfigCmd, replacementConfigCmd)
}
