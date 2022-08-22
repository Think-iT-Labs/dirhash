package cmd

import (
	"fmt"

	"github.com/Think-iT-Labs/go-dirhash/lib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ignoredPaths []string

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringSliceVarP(&ignoredPaths, "ignore", "x", nil, "ignored glob paths")
}

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Print the hash of a folder.",
	Long:  `Print the hash of a folder. You may ignore some files using flags.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory = args[0]
		log.Debug("directory: ", directory)
		log.Debug("ignore: ", ignoredPaths)
		fmt.Println(lib.DirHash(directory, ignoredPaths))
	},
}
