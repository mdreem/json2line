package main

import (
	"fmt"
	"github.com/mdreem/json2line/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func main() {
	var rootCmd = &cobra.Command{Use: "json2line", Short: "json2line converts json to a formatted line of text"}

	rootCmd.PersistentFlags().StringP("file", "f", "", "config file that should be used")
	if err := rootCmd.Execute(); err != nil {
		printInformation("could not execute command: %v", err)
		os.Exit(1)
	}

	filePath, err := rootCmd.Flags().GetString("file")
	if err != nil {
		printInformation("could not fetch file option: %v", err)
		os.Exit(1)
	}
	loadConfig(filePath)

	if isInputFromPipe() {
		err := cmd.ProcessInput(os.Stdin, os.Stdout)
		if err != nil {
			printInformation("could not parse line: %v", err)
			os.Exit(1)
		}
	}
}

func isInputFromPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		printInformation("could not read stat")
		os.Exit(1)
	}
	return stat.Mode()&os.ModeCharDevice == 0
}

func loadConfig(filePath string) {
	if filePath != "" {
		directory, file := filepath.Split(filePath)
		initializeConfiguration(directory, file)
	} else {
		dir, err := os.UserConfigDir()
		if err != nil {
			printInformation("could not find base configuration directory")
		}
		initializeConfiguration(filepath.Join(dir, "json2line"), "json2line.toml")
	}
}

func initializeConfiguration(configDir string, configFile string) {
	printInformation("Loading file '%s' in directory '%s'\n", configFile, configDir)

	viper.SetConfigName(configFile)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configDir)

	err := viper.ReadInConfig()
	switch t := err.(type) {
	case viper.ConfigFileNotFoundError:
		fmt.Println("No config file found, using defaults")
	case nil:
	default:
		panic(fmt.Errorf("fatal error: (%v) %s", t, err))
	}
}

func printInformation(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		panic(fmt.Errorf("could not print to stderr: %v", err))
	}
}
