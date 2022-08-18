package lib

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	log "github.com/sirupsen/logrus"
)

func walkDir(pathToWalk string) []string {
	allFiles := []string{}
	err := filepath.Walk(pathToWalk,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				allFiles = append(allFiles, path)
			}
			return nil
		})
	if err != nil {
		log.Error(err)
	}
	return allFiles
}

func filesToExclude(path string, excludedPaths []string) []string {
	excludeMatches := pathsToExclude(path, excludedPaths)
	files := []string{}
	for i := 0; i < len(excludeMatches); i++ {
		fileStat, err := os.Stat(excludeMatches[i])
		if err != nil {
			continue
		}
		if fileStat.IsDir() {
			files = append(files, walkDir(excludeMatches[i])...)
		} else {
			files = append(files, excludeMatches[i])
		}
	}
	return files
}

func pathsToExclude(path string, excludedPaths []string) []string {
	allMatches := []string{}
	for i := 0; i < len(excludedPaths); i++ {
		pattern := excludedPaths[i]
		pattern = filepath.Join(path, pattern)
		matches, err := doublestar.Glob(pattern)
		if err != nil {
			log.Error(err)
		}
		allMatches = append(allMatches, matches...)
	}
	return allMatches
}
