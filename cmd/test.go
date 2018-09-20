package cmd

import (
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run go and solidity tests",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runStub("test"); err != nil {
			Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
