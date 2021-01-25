package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var replacementConfigCmd = &cobra.Command{
	Use: "replacement", Run: func(c *cobra.Command, args []string) {
		replacementConfiguration(c, args)
	},
}

func replacementConfiguration(c *cobra.Command, args []string) {
	fmt.Printf("args: %v", args)
}

func init() {
	replacementConfigCmd.PersistentFlags().StringArrayP("add", "a", []string{}, "add replacement configuration. <NAME> <PATTERN>")
	replacementConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete replacement configuration. <NAME>")
}
