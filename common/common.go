package common

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func GetString(rootCmd *cobra.Command, option string) string {
	optionString, err := rootCmd.Flags().GetString(option)
	if err != nil {
		PrintInformationf("could not fetch %s option: %v\n", option, err)
		os.Exit(1)
	}
	return optionString
}

func GetStringArray(rootCmd *cobra.Command, option string) []string {
	optionStringArray, err := rootCmd.Flags().GetStringArray(option)
	if err != nil {
		PrintInformationf("could not fetch %s option: %v\n", option, err)
		os.Exit(1)
	}
	return optionStringArray
}

func GetBoolean(rootCmd *cobra.Command, option string) bool {
	optionString, err := rootCmd.Flags().GetBool(option)
	if err != nil {
		PrintInformationf("could not fetch %s option: %v\n", option, err)
		os.Exit(1)
	}
	return optionString
}

func PrintInformationf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		panic(fmt.Errorf("could not print to stderr: %v", err))
	}
}
