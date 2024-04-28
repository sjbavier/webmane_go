package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed database with data",
	Long:  "command for seeding database",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Database seeded")
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
