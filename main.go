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
	rootCmd.PersistentFlags().StringP("adhoc", "a", "", "ad hoc format string.")

	if err := rootCmd.Execute(); err != nil {
		printInformationf("could not execute command: %v", err)
		os.Exit(1)
	}

	initConfig(rootCmd)

	adHocFormatString := getString(rootCmd, "adhoc")

	var formattingTemplate *template.Template
	if adHocFormatString == "" {
		formattingTemplate = loadTemplate(rootCmd)
	} else {
		formattingTemplate = createTemplate("adhoc", adHocFormatString)
	}

	replacements := loadReplacements()

	if isInputFromPipe() {
		err := cmd.ProcessInput(os.Stdin, os.Stdout, formattingTemplate, replacements)
		if err != nil {
			printInformationf("could not parse line: %v\n", err)
			os.Exit(1)
		}
	}
}

func getString(rootCmd *cobra.Command, option string) string {
	optionString, err := rootCmd.Flags().GetString(option)
	if err != nil {
		printInformationf("could not fetch %s option: %v\n", option, err)
		os.Exit(1)
	}
	return optionString
}

func loadTemplate(rootCmd *cobra.Command) *template.Template {
	formatter, err := rootCmd.Flags().GetString("formatter")
	formattingTemplate := getFormatter(err, formatter)
	return formattingTemplate
}

func loadReplacements() map[string]string {
	return viper.GetStringMapString("replacements")
}

func getFormatter(err error, formatter string) *template.Template {
	if err != nil {
		printInformationf("could not fetch formatter option: %v\n", err)
		os.Exit(1)
	}
	templates := viper.GetStringMapString("templates")
	formatString := templates[formatter]
	if formatString != "" {
		return createTemplate(formatter, formatString)
	}
	printInformationf("no formatter with the name '%s' defined\n", formatter)
	return nil
}

func createTemplate(formatter string, formatString string) *template.Template {
	parse, err := template.New(formatter).Parse(formatString)
	if err != nil {
		printInformationf("could not parse template: %\n", err)
		os.Exit(1)
	}
	return parse
}

func isInputFromPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		printInformationf("could not read stat\n")
		os.Exit(1)
	}
	return stat.Mode()&os.ModeCharDevice == 0
}

func initConfig(rootCmd *cobra.Command) {
	filePath := getString(rootCmd, "config")
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
