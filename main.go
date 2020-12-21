package main

import (
	"fmt"
	"github.com/mdreem/json2line/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	var rootCmd = &cobra.Command{Use: "json2line"}

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file that should be used")
	rootCmd.PersistentFlags().StringP("formatter", "f", "", "formatter that should be used")

	if err := rootCmd.Execute(); err != nil {
		printInformationf("could not execute command: %v", err)
		os.Exit(1)
	}

	filePath, err := rootCmd.Flags().GetString("config")
	if err != nil {
		printInformationf("could not fetch file option: %v\n", err)
		os.Exit(1)
	}
	loadConfig(filePath)

	formatter, err := rootCmd.Flags().GetString("formatter")
	formattingTemplate := getFormatter(err, formatter)

	if isInputFromPipe() {
		err := cmd.ProcessInput(os.Stdin, os.Stdout, formattingTemplate)
		if err != nil {
			printInformationf("could not parse line: %v\n", err)
			os.Exit(1)
		}
	}
}

func getFormatter(err error, formatter string) *template.Template {
	if err != nil {
		printInformationf("could not fetch formatter option: %v\n", err)
		os.Exit(1)
	}
	templates := viper.GetStringMapString("templates")
	formatString := templates[formatter]
	if formatString != "" {
		parse, err := template.New(formatter).Parse(formatString)
		if err != nil {
			printInformationf("could not parse template: %\n", err)
			os.Exit(1)
		}
		return parse
	}
	printInformationf("no formatter with the name '%s' defined\n", formatter)
	return nil
}

func isInputFromPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		printInformationf("could not read stat\n")
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
			printInformationf("could not find base configuration directory\n")
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

func printInformationf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		panic(fmt.Errorf("could not print to stderr: %v", err))
	}
}
