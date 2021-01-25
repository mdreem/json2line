package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{Use: "configure", Run: func(c *cobra.Command, args []string) {
	configuration(c, args)
}, TraverseChildren: true}

func configuration(c *cobra.Command, args []string) {
	fmt.Printf("args: %v", args)
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().BoolP("show", "s", false, "show current configuration")
	configCmd.AddCommand(formatterConfigCmd, replacementConfigCmd)
}
