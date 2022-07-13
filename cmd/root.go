package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-dirhash",
	Short: "A package for creating a hash of a folder.",
	Long:  `A package for creating a hash of a folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Runned")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
