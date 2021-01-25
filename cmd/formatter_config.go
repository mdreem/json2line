package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var formatterConfigCmd = &cobra.Command{
	Use: "formatter",
	Run: func(c *cobra.Command, args []string) {
		formatterConfiguration(c, args)
	},
}

func formatterConfiguration(c *cobra.Command, args []string) {
	fmt.Printf("args: %v", args)
}

func init() {
	formatterConfigCmd.PersistentFlags().StringArrayP("add", "a", []string{}, "add formatter configuration. <NAME> <PATTERN>")
	formatterConfigCmd.PersistentFlags().StringP("delete", "d", "", "delete formatter configuration. <NAME>")
}
