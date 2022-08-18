package main

import (
	"os"

	"github.com/Think-iT-Labs/go-dirhash/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
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
	cmd.Execute()
}
