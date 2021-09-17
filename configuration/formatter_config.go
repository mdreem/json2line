package configuration

import (
	"github.com/mdreem/json2line/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var FormatterConfigCmd = &cobra.Command{
	Use: "formatter",
	Run: func(c *cobra.Command, args []string) {
		FormatterConfiguration(c, args)
	},
}

func init() {
	FormatterConfigCmd.PersistentFlags().StringP("add-key", "k", "", "key to add.")
	FormatterConfigCmd.PersistentFlags().StringP("add-value", "v", "", "value to add.")

	FormatterConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete formatter configuration. <NAME>")
}

func FormatterConfiguration(c *cobra.Command, _ []string) {
	replacementSection := viper.GetStringMapString("templates")

	reconfigureSectionViaCommand(c, &replacementSection)

	viper.Set("templates", replacementSection)

	common.PrintInformationf("writing configuration to: %s\n", viper.ConfigFileUsed())

	err := viper.WriteConfig()
	if err != nil {
		common.PrintInformationf("could not write config: %v\n", err)
		os.Exit(1)
	}
}
