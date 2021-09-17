package configuration

import (
	"github.com/mdreem/json2line/common"
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

func formatterConfiguration(c *cobra.Command, _ []string) {
	replacementSection := viper.GetStringMapString("templates")

	reconfigureSection(c, &replacementSection)

	viper.Set("templates", replacementSection)

	common.PrintInformationf("writing configuration to: %s\n", viper.ConfigFileUsed())

	err := viper.WriteConfig()
	if err != nil {
		common.PrintInformationf("could not write config: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	formatterConfigCmd.PersistentFlags().StringP("add-key", "k", "", "key to add.")
	formatterConfigCmd.PersistentFlags().StringP("add-value", "v", "", "value to add.")

	formatterConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete formatter configuration. <NAME>")
}
