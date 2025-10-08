package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dirhash",
	Short: "Compute a stable checksum for a directory",
	Long:  "Dirhash computes a deterministic checksum of a directory by hashing file contents and paths, with support for glob-based exclusions.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// SetVersion allows injecting version at build time via ldflags.
// Example: go build -ldflags "-X 'github.com/Think-iT-Labs/dirhash/cmd.version=v1.0.0'"
var version = "dev"

func init() {
	rootCmd.Version = version
}
