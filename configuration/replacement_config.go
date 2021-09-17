package configuration

import (
	"github.com/mdreem/json2line/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var ReplacementConfigCmd = &cobra.Command{
	Use: "replacement",
	Run: func(c *cobra.Command, args []string) {
		replacementConfiguration(c, args)
	},
}

func init() {
	ReplacementConfigCmd.PersistentFlags().StringP("add-key", "k", "", "key to add.")
	ReplacementConfigCmd.PersistentFlags().StringP("add-value", "v", "", "value to add.")

	ReplacementConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete replacement configuration. <NAME>")
}

func replacementConfiguration(c *cobra.Command, _ []string) {
	replacementSection := viper.GetStringMapString("replacements")

	reconfigureSection(c, &replacementSection)

	viper.Set("replacements", replacementSection)

	common.PrintInformationf("writing configuration to: %s\n", viper.ConfigFileUsed())

	err := viper.WriteConfig()
	if err != nil {
		common.PrintInformationf("could not write config: %v\n", err)
		os.Exit(1)
	}
}
