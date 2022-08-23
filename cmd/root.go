package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-dirhash",
	Short: "A package for creating a hash of a folder.",
	Long:  `A package for creating a hash of a folder.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
