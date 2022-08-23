package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

/*
Dirhash walks into a provided directory and calculates its SHA256 checksum,
based on the checksum and name of th individual files within the directory.
A list of exlcudedPaths as glob patterns can be provided to make Dirhash ignore their matches
*/
func DirHash(path string, ignoredPaths []string) string {
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("Using relative path ", baseDir)
		path = filepath.Join(baseDir, path)
	}
	var allFiles, err = walkDir(path)
	if err != nil {
		log.Fatal(err)
	}
	exlcudedFilesMatch, err := filesToIgnore(path, ignoredPaths)
	if err != nil {
		log.Fatal(err)
	}
	if log.GetLevel() > log.DebugLevel {
		for i := 0; i < len(exlcudedFilesMatch); i++ {
			log.Debug("excluding: ", exlcudedFilesMatch[i])
		}
	}
	var filesToHash []string = []string{}
	for i := 0; i < len(allFiles); i++ {
		if !slices.Contains(exlcudedFilesMatch, allFiles[i]) {
			filesToHash = append(filesToHash, allFiles[i])
		}
	}
	var fileHashes []string = []string{}
	for i := 0; i < len(filesToHash); i++ {
		fileHash, err := fileSha256(filesToHash[i])
		if err != nil {
			log.Fatal(fmt.Sprintf("Error hashing file %s : %s", filesToHash[i], err))
		}
		relpath, _ := filepath.Rel(path, filesToHash[i])
		hashCombo := fmt.Sprintf("%s %s", relpath, fileHash)
		fileHashes = append(fileHashes, hashCombo)
		log.Debug("hashing: ", hashCombo)
	}
	return mergeAllHashes(fileHashes)
}

// mergeAllHashes returns hash of joint slice elements as lines
func mergeAllHashes(hashes []string) string {
	sort.Strings(hashes)
	var hash string = strings.Join(hashes, "\n")
	return stringSha256(strings.NewReader(hash))
}

// fileSha256 returns the hash of a file
func fileSha256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		defer f.Close()
	}
	return stringSha256(f), err
}

// stringSha256 returns the SHA256 checksum of a given io.Reader
func stringSha256(f io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error(err)
	}
	return hex.EncodeToString(h.Sum(nil)[:])
}
