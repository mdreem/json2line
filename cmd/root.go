package cmd

import (
	"fmt"
	"github.com/mdreem/json2line/common"
	"github.com/mdreem/json2line/configuration/load"
	"github.com/mdreem/json2line/processor"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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

func isInputFromPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		common.PrintInformationf("could not read stat\n")
		os.Exit(1)
	}
	return stat.Mode()&os.ModeCharDevice == 0
}
