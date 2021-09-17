package configuration

import (
	"fmt"
	"github.com/mdreem/json2line/cmd"
	"github.com/mdreem/json2line/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var configCmd = &cobra.Command{
	Use: "configure",
	Run: func(c *cobra.Command, args []string) {
		configuration(c, args)
	},
}

func configuration(c *cobra.Command, _ []string) {
	show := common.GetBoolean(c, "show")

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
	cmd.RootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolP("show", "s", false, "show current configuration")
	configCmd.AddCommand(formatterConfigCmd, replacementConfigCmd, initConfigurationFileCmd)
}

func reconfigureSection(c *cobra.Command, section *map[string]string) {
	addKey := strings.ToLower(common.GetString(c, "add-key"))
	addValue := common.GetString(c, "add-value")

	if (addKey != "" && addValue == "") || (addKey == "" && addValue != "") {
		common.PrintInformationf("add option needs key and value option")
		os.Exit(1)
	}

	deleteOption := common.GetString(c, "delete")

	if addKey != "" && deleteOption != "" {
		common.PrintInformationf("add and delete are mutually exclusive\n")
		os.Exit(1)
	}

	if deleteOption != "" {
		delete(*section, deleteOption)
	}

	if addKey != "" {
		(*section)[addKey] = addValue
	}
}
