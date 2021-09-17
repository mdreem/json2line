package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func initConfiguration() {
	filePath := getString(RootCmd, "config")
	loadConfig(filePath)
}

func loadConfig(filePath string) {
	if filePath != "" {
		directory, file := filepath.Split(filePath)
		initializeConfiguration(directory, file)
	} else {
		dir, err := os.UserConfigDir()
		if err != nil {
			printInformationf("could not find base configuration directory\n")
			os.Exit(1)
		}
		initializeConfiguration(filepath.Join(dir, "json2line"), "json2line.toml")
	}
}

func initializeConfiguration(configDir string, configFile string) {
	printInformationf("Loading file '%s' in directory '%s'\n", configFile, configDir)

	viper.SetConfigName(configFile)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configDir)

	err := viper.ReadInConfig()
	switch t := err.(type) {
	case viper.ConfigFileNotFoundError:
		printInformationf("No config file found, using defaults\n")
	case nil:
	default:
		panic(fmt.Errorf("fatal error: (%v) %s", t, err))
	}
}