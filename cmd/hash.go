package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var excludePaths []string

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringSliceVarP(&excludePaths, "excluded-paths", "x", nil, "excluded files")
}

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Print the hash of a folder.",
	Long:  `Print the hash of a folder. You may ignore some files using flags.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory = args[0]
		fmt.Println("directory = " + directory)
		fmt.Printf("excluded-paths = %+v \n", excludePaths)
	},
}
