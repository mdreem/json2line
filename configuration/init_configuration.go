package configuration

import (
	"github.com/mdreem/json2line/common"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var InitConfigurationFileCmd = &cobra.Command{
	Use: "init",
	Run: func(c *cobra.Command, args []string) {
		initConfigurationFile(c, args)
	},
}

func initConfigurationFile(_ *cobra.Command, _ []string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		common.PrintInformationf("could not find base configuration directory: %v\n", err)
		os.Exit(1)
	}

	configDir := filepath.Join(dir, "json2line")
	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		common.PrintInformationf("could not create base configuration directory (%s) %v\n", configDir, err)
		os.Exit(1)
	}

	configFile := filepath.Join(configDir, "json2line.toml")
	_, err = os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		common.PrintInformationf("could not create configuration file (%s) %v\n", configFile, err)
		os.Exit(1)
	}

	common.PrintInformationf("Created file: %s\n", configFile)
}
