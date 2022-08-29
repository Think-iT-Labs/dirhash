package cmd

import (
	"fmt"

	"github.com/Think-iT-Labs/dirhash/lib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ignoredPaths []string

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringSliceVarP(&ignoredPaths, "ignore", "x", nil, "ignored glob paths")
}

var hashCmd = &cobra.Command{
	Use:   "sha256",
	Short: "Print the sha256 hash of a directory.",
	Long:  `Print the sha256 hash of a directory. You may ignore some files using ignore flag.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory = args[0]
		log.Debug("directory: ", directory)
		log.Debug("ignore: ", ignoredPaths)
		fmt.Println(lib.DirHash(directory, ignoredPaths))
	},
	Example: "dirhash sha256 --ignore=node_modules/** ~/my_node_project",
}
