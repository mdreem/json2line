package configuration

import (
	"github.com/mdreem/json2line/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var BufferSizeConfigCmd = &cobra.Command{
	Use: "buffer_size",
	Run: func(c *cobra.Command, args []string) {
		BufferSizeConfiguration(c, args)
	},
}

func init() {
	BufferSizeConfigCmd.PersistentFlags().StringP("set", "S", "", "set buffer size to <VALUE>.")
	BufferSizeConfigCmd.PersistentFlags().BoolP("delete", "d", false, "delete buffer size configuration.")
}

func BufferSizeConfiguration(command *cobra.Command, _ []string) {
	configurationSection := viper.GetStringMapString("configuration")

	bufferSize := common.GetString(command, "set")
	isDeleteBufferSize := common.GetBoolean(command, "delete")

	if bufferSize != "" && isDeleteBufferSize {
		common.PrintInformationf("set and delete are mutually exclusive.\n")
		os.Exit(1)
	}

	if bufferSize != "" {
		configurationSection["buffer_size"] = bufferSize
	}
	if isDeleteBufferSize {
		delete(configurationSection, "buffer_size")
	}

	viper.Set("configuration", configurationSection)

	common.PrintInformationf("writing configuration to: %s\n", viper.ConfigFileUsed())

	err := viper.WriteConfig()
	if err != nil {
		common.PrintInformationf("could not write config: %v\n", err)
		os.Exit(1)
	}
}
