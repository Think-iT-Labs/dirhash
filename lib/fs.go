package lib

import (
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar"
	log "github.com/sirupsen/logrus"
)

// walkDir walks a directory returning a list of all of its files, following symbolic links
func walkDir(pathToWalk string) ([]string, error) {
	allFiles := []string{}
	visited := make(map[string]bool) // Track visited paths to prevent infinite loops
	
	err := walkDirWithSymlinks(pathToWalk, &allFiles, visited)
	return allFiles, err
}

// walkDirWithSymlinks recursively walks a directory, following symbolic links
func walkDirWithSymlinks(root string, allFiles *[]string, visited map[string]bool) error {
	// Get absolute path to handle relative symlinks correctly
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	
	// Check if we've already visited this path (prevents infinite loops)
	if visited[absRoot] {
		log.Debug("Skipping already visited path: ", absRoot)
		return nil
	}
	visited[absRoot] = true
	
	// Use filepath.WalkDir which doesn't follow symlinks by default
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// If it's a symbolic link, handle it specially
		if d.Type()&os.ModeSymlink != 0 {
			return handleSymlink(path, allFiles, visited)
		}
		
		// Regular file - add to list
		if d.Type().IsRegular() {
			*allFiles = append(*allFiles, path)
		}
		
		return nil
	})
}

// handleSymlink processes a symbolic link, following it if it points to a directory
func handleSymlink(symlinkPath string, allFiles *[]string, visited map[string]bool) error {
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		log.Debug("Failed to read symlink ", symlinkPath, ": ", err)
		return nil // Skip broken symlinks
	}

	var resolvedTarget string
	if filepath.IsAbs(target) {
		// Absolute symlink
		resolvedTarget = target
	} else {
		// Relative symlink - resolve relative to the symlink's directory
		symlinkDir := filepath.Dir(symlinkPath)
		resolvedTarget = filepath.Join(symlinkDir, target)
	}

	resolvedTarget = filepath.Clean(resolvedTarget)
	log.Debug("Following symlink ", symlinkPath, " -> ", resolvedTarget)
	targetInfo, err := os.Stat(resolvedTarget)

	if err != nil {
		log.Debug("Symlink target doesn't exist or is inaccessible: ", resolvedTarget, ": ", err)
		return nil // Skip broken symlinks
	}
	
	if targetInfo.IsDir() {
		// Symlink points to a directory - recursively walk it
		return walkDirWithSymlinks(resolvedTarget, allFiles, visited)
	} else {
		// Symlink points to a file - add it to the list
		*allFiles = append(*allFiles, symlinkPath)
	}
	
	return nil
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
