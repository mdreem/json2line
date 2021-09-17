package cmd

import (
	"fmt"
	"github.com/mdreem/json2line/common"
	"github.com/mdreem/json2line/configuration/load"
	"github.com/mdreem/json2line/processor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var GitCommit string
var Version string

var RootCmd = &cobra.Command{
	Use: "json2line",
	Run: runCommand,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		common.PrintInformationf("could not execute command: %v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		load.LoadConfig(common.GetString(RootCmd, "config"))
	})

	RootCmd.PersistentFlags().StringP("config", "c", "", "config file that should be used. <FILE>")
	RootCmd.PersistentFlags().StringP("formatter", "f", "", "formatter that should be used. <NAME>")
	RootCmd.PersistentFlags().StringP("adhoc", "o", "", "ad hoc format string.")
	RootCmd.PersistentFlags().StringP("buffer-size", "b", "", "set buffer size for lines.")
	RootCmd.PersistentFlags().StringArrayP("replacement", "r", []string{}, "ad hoc replacements.")
	RootCmd.PersistentFlags().BoolP("version", "V", false, "print version information.")
}

func runCommand(command *cobra.Command, _ []string) {
	if handleVersionFlag(command) {
		return
	}

	adHocFormatString := common.GetString(command, "adhoc")

	var formattingTemplate *template.Template
	if adHocFormatString == "" {
		formattingTemplate = loadTemplate(command)
	} else {
		formattingTemplate = createTemplate("adhoc", adHocFormatString)
	}

	replacements := loadReplacements(command)

	bufferSize := getBufferSize(command)

	if isInputFromPipe() {
		err := processor.ProcessInput(os.Stdin, os.Stdout, formattingTemplate, replacements, bufferSize)
		if err != nil {
			common.PrintInformationf("could not parse line: %v\n", err)
			os.Exit(1)
		}
	}
}

func getBufferSize(command *cobra.Command) int {
	bufferSizeString := common.GetString(command, "buffer-size")
	if bufferSizeString == "" {
		return processor.InitialBufferSize
	}
	bufferSize, err := strconv.Atoi(bufferSizeString)
	if err != nil {
		common.PrintInformationf("could not parse buffer size: %v\n", err)
		os.Exit(1)
	}
	return bufferSize
}

func handleVersionFlag(c *cobra.Command) bool {
	printVersion := common.GetBoolean(c, "version")
	if printVersion {
		fmt.Printf("\nVersion: %s\n", Version)
		fmt.Printf("Commit:  %s\n", GitCommit)
		return true
	}
	return false
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
		common.PrintInformationf("could not fetch replacement options: %v\n", err)
		os.Exit(1)
	}

	var replacements = make(map[string]string)
	for k, v := range configuredReplacements {
		replacements[k] = v
	}

	for _, r := range adhocReplacements {
		replacement := strings.Fields(r)
		if len(replacement) != 2 {
			common.PrintInformationf("replacement is not of correct format '<FROM> <TO>' where both values are separated via whitespace")
		}

		replacements[replacement[0]] = replacement[1]
	}
	return replacements
}

func getFormatter(err error, formatter string) *template.Template {
	if err != nil {
		common.PrintInformationf("could not fetch formatter option: %v\n", err)
		os.Exit(1)
	}
	templates := viper.GetStringMapString("templates")
	formatString := templates[formatter]
	if formatString != "" {
		return createTemplate(formatter, formatString)
	}
	common.PrintInformationf("no formatter with the name '%s' defined\n", formatter)
	return nil
}

func createTemplate(formatter string, formatString string) *template.Template {
	parse, err := template.New(formatter).Parse(formatString)
	if err != nil {
		common.PrintInformationf("could not parse template: %\n", err)
		os.Exit(1)
	}
	return parse
}

func isInputFromPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		common.PrintInformationf("could not read stat\n")
		os.Exit(1)
	}
	return stat.Mode()&os.ModeCharDevice == 0
}
