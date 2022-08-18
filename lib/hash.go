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

// Calculate the hash of directory (file names + content)
func DirHash(path string, excludedPaths []string) string {
	if !filepath.IsAbs(path) {
		baseDir, err := os.Getwd()
		if err != nil {
			log.Error(err)
		}
		log.Debug(baseDir)
		path = filepath.Join(baseDir, path)
	}
	var allFiles []string = walkDir(path)
	var exlcudedFilesMatch []string = filesToExclude(path, excludedPaths)
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
		fileHash := fileSha256(filesToHash[i])
		relpath, _ := filepath.Rel(path, filesToHash[i])
		hashCombo := fmt.Sprintf("%s %s", relpath, fileHash)
		fileHashes = append(fileHashes, hashCombo)
		log.Debug("hashing: ", hashCombo)
	}
	return mergeAllHashes(fileHashes)
}

func mergeAllHashes(hashes []string) string {
	sort.Strings(hashes)
	var hash string = strings.Join(hashes, "\n")
	return sha256String(strings.NewReader(hash))
}

func fileSha256(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()
	return sha256String(f)
}

func sha256String(f io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error(err)
	}
	return hex.EncodeToString(h.Sum(nil)[:])
}
