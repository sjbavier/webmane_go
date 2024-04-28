package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "webmane",
	Short: "webmane command interface",
	Long:  "webmane tools for the command line",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("webmane cmd:")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
