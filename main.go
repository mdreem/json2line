package main

import (
	"fmt"
	"github.com/mdreem/json2line/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	cmd.Execute()

	initConfig(cmd.RootCmd)

	adHocFormatString := getString(cmd.RootCmd, "adhoc")

	var formattingTemplate *template.Template
	if adHocFormatString == "" {
		formattingTemplate = loadTemplate(cmd.RootCmd)
	} else {
		formattingTemplate = createTemplate("adhoc", adHocFormatString)
	}

	replacements := loadReplacements(cmd.RootCmd)

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

func loadReplacements(rootCmd *cobra.Command) map[string]string {
	configuredReplacements := viper.GetStringMapString("replacements")
	adhocReplacements, err := rootCmd.Flags().GetStringArray("replacement")
	if err != nil {
		printInformationf("could not fetch replacement options: %v\n", err)
		os.Exit(1)
	}

	var replacements = make(map[string]string)
	for k, v := range configuredReplacements {
		replacements[k] = v
	}

	for _, r := range adhocReplacements {
		replacement := strings.Fields(r)
		if len(replacement) != 2 {
			printInformationf("replacement is not of correct format '<FROM> <TO>' where both values are separated via whitespace")
		}

		replacements[replacement[0]] = replacement[1]
	}
	return replacements
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
