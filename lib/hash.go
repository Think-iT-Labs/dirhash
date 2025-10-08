package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

/*
Dirhash walks into a provided directory and calculates its SHA256 checksum,
based on the checksum and name of th individual files within the directory.
A list of exlcudedPaths as glob patterns can be provided to make Dirhash ignore their matches
*/
func DirHash(path string, ignoredPaths []string) (string, error) {
	// Keep existing behavior for backward compatibility
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
		log.Debug("Using relative path ", baseDir)
		path = filepath.Join(baseDir, path)
	}
	var allFiles, err = walkDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to walk directory %q: %w", path, err)
	}
	exlcudedFilesMatch, err := filesToIgnore(path, ignoredPaths)
	if err != nil {
		return "", fmt.Errorf("failed to process ignore patterns: %w", err)
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		for i := 0; i < len(exlcudedFilesMatch); i++ {
			log.Debug("excluding: ", exlcudedFilesMatch[i])
		}
	}
	var filesToHash = []string{}
	for i := 0; i < len(allFiles); i++ {
		if !slices.Contains(exlcudedFilesMatch, allFiles[i]) {
			filesToHash = append(filesToHash, allFiles[i])
		}
	}
	_, combined, err := hashFilesConcurrently(path, filesToHash)
	if err != nil {
		return "", fmt.Errorf("failed to hash files: %w", err)
	}
	return mergeAllHashes(combined), nil
}

// FileHash represents a per-file digest entry used for JSON output.
type FileHash struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}

// DirHashDetails computes the directory hash and also returns per-file hashes.
// It mirrors DirHash behavior but provides structured details for consumers.
func DirHashDetails(path string, ignoredPaths []string) (string, []FileHash, error) {
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			return "", nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		log.Debug("Using relative path ", baseDir)
		path = filepath.Join(baseDir, path)
	}

	allFiles, err := walkDir(path)
	if err != nil {
		return "", nil, fmt.Errorf("failed to walk directory %q: %w", path, err)
	}
	exlcudedFilesMatch, err := filesToIgnore(path, ignoredPaths)
	if err != nil {
		return "", nil, fmt.Errorf("failed to process ignore patterns: %w", err)
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		for i := 0; i < len(exlcudedFilesMatch); i++ {
			log.Debug("excluding: ", exlcudedFilesMatch[i])
		}
	}

	var filesToHash = []string{}
	for i := 0; i < len(allFiles); i++ {
		if !slices.Contains(exlcudedFilesMatch, allFiles[i]) {
			filesToHash = append(filesToHash, allFiles[i])
		}
	}

	pairs, combined, err := hashFilesConcurrently(path, filesToHash)
	if err != nil {
		return "", nil, fmt.Errorf("failed to hash files: %w", err)
	}
	overall := mergeAllHashes(combined)
	return overall, pairs, nil
}

// mergeAllHashes returns hash of joint slice elements as lines
func mergeAllHashes(hashes []string) string {
	sort.Strings(hashes)
	hash := strings.Join(hashes, "\n")
	return stringSha256(strings.NewReader(hash))
}

// fileSha256 returns the hash of a file
func fileSha256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	return stringSha256(f), nil
}

// stringSha256 returns the SHA256 checksum of a given io.Reader
func stringSha256(f io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error(err)
	}
	return hex.EncodeToString(h.Sum(nil)[:])
}

// hashFilesConcurrently hashes files with a worker pool, returning per-file pairs and combined "<relpath> <hash>" entries.
func hashFilesConcurrently(base string, files []string) ([]FileHash, []string, error) {
	type job struct{ index int }
	type result struct {
		index int
		rel   string
		hash  string
		err   error
	}

	if len(files) == 0 {
		return nil, nil, nil
	}

	workers := runtime.GOMAXPROCS(0)
	if workers < 2 {
		workers = 2
	}
	if workers > 8 {
		workers = 8
	}

	jobs := make(chan job)
	results := make(chan result, len(files))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for j := range jobs {
			f := files[j.index]
			h, err := fileSha256(f)
			rel, _ := filepath.Rel(base, f)
			results <- result{index: j.index, rel: rel, hash: h, err: err}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}
	for i := 0; i < len(files); i++ {
		jobs <- job{index: i}
	}
	close(jobs)
	wg.Wait()
	close(results)

	pairs := make([]FileHash, 0, len(files))
	combined := make([]string, 0, len(files))
	var firstErr error
	for r := range results {
		if r.err != nil && firstErr == nil {
			firstErr = r.err
		}
		pairs = append(pairs, FileHash{Path: r.rel, Hash: r.hash})
		combined = append(combined, fmt.Sprintf("%s %s", r.rel, r.hash))
		log.Debug("hashing: ", r.rel, " ", r.hash)
	}
	return pairs, combined, firstErr
}
