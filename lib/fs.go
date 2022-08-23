package lib

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	log "github.com/sirupsen/logrus"
)

// walkDir walks a directory returning a list of all of its files
func walkDir(pathToWalk string) ([]string, error) {
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
	return allFiles, err
}

// filesToIgnore returns the files to ignore from a parent path and a the list of glob patterns
func filesToIgnore(path string, ingoredPaths []string) ([]string, error) {
	ignoreMatches := pathsToIgnore(path, ingoredPaths)
	files := []string{}
	for i := 0; i < len(ignoreMatches); i++ {
		fileStat, err := os.Stat(ignoreMatches[i])
		if err != nil {
			continue
		}
		if fileStat.IsDir() {
			ignoredDirFiles, err := walkDir(ignoreMatches[i])
			if err != nil {
				return nil, err
			}
			files = append(files, ignoredDirFiles...)
		} else {
			files = append(files, ignoreMatches[i])
		}
	}
	return files, nil
}

// pathsToIgnore returns the paths to ignore from a parent path and a the list of glob patterns
func pathsToIgnore(path string, ingoredPaths []string) []string {
	allMatches := []string{}
	for i := 0; i < len(ingoredPaths); i++ {
		pattern := ingoredPaths[i]
		pattern = filepath.Join(path, pattern)
		matches, err := doublestar.Glob(pattern)
		if err != nil {
			log.Error("Unable to find Glob matches", err)
		}
		allMatches = append(allMatches, matches...)
	}
	return allMatches
}
