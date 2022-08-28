/*
Print the hash of a folder. You may ignore some files using flags.

Usage:

	dirhash hash [flags]

Flags:

	-x, --ignore strings           ignored glob paths
	-h, --help                     help for hash
*/
package main

import (
	"os"

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
		log.Error("LOG_LEVEL environment variable should be one of:", log.AllLevels)
	}
	log.SetLevel(lvl)
}

func main() {
	Execute()
}
