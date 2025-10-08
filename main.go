/*
Print the hash of a folder. You may ignore some files using flags.

Usage:

	dirhash sha256 [flags]

Flags:

	-x, --ignore strings           ignored glob paths
	-h, --help                     help for hash
*/
package main

import (
	"os"

	"github.com/Think-iT-Labs/dirhash/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warn("Invalid LOG_LEVEL, defaulting to info. Valid levels:", log.AllLevels)
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
}

func main() {
	cmd.Execute()
}
