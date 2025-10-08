package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Think-iT-Labs/dirhash/lib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ignoredPaths []string
var outputFormat string
var ignoreFile string

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.Flags().StringSliceVarP(&ignoredPaths, "ignore", "x", nil, "ignored glob paths")
	hashCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "output format: text|json")
	hashCmd.Flags().StringVar(&ignoreFile, "ignore-file", "", "path to ignore file (overrides .dirhashignore if set)")
}

var hashCmd = &cobra.Command{
	Use:   "sha256",
	Short: "Print the sha256 hash of a directory",
	Long:  "Compute the SHA-256 hash of a directory. You can ignore files or folders using one or more -x/--ignore glob patterns (supports **).",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var directory = args[0]
		log.Debug("directory: ", directory)
		// Load ignore patterns from file if provided, else auto-detect .dirhashignore in target dir
		patterns := append([]string{}, ignoredPaths...)
		if ignoreFile != "" {
			filePatterns, err := readIgnoreFile(ignoreFile)
			if err != nil {
				log.Fatal(err)
			}
			patterns = append(patterns, filePatterns...)
		} else {
			// Resolve .dirhashignore relative to the provided directory
			base := directory
			if !filepath.IsAbs(base) {
				// Use CWD join to mimic lib.DirHash behavior
				cwd, _ := os.Getwd()
				base = filepath.Join(cwd, base)
			}
			defaultIgnore := filepath.Join(base, ".dirhashignore")
			if st, err := os.Stat(defaultIgnore); err == nil && !st.IsDir() {
				filePatterns, err := readIgnoreFile(defaultIgnore)
				if err != nil {
					log.Fatal(err)
				}
				patterns = append(patterns, filePatterns...)
			}
		}
		log.Debug("ignore: ", patterns)
		switch outputFormat {
		case "text":
			hash, err := lib.DirHash(directory, patterns)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(hash)
		case "json":
			overall, details, err := lib.DirHashDetails(directory, patterns)
			if err != nil {
				log.Fatal(err)
			}
			payload := struct {
				Hash   string         `json:"hash"`
				Files  []lib.FileHash `json:"files"`
				Algo   string         `json:"algo"`
				Format string         `json:"format"`
			}{Hash: overall, Files: details, Algo: "sha256", Format: "json"}
			b, err := json.MarshalIndent(payload, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(b))
		default:
			log.Fatalf("unsupported output format: %s", outputFormat)
		}
	},
	Example: "dirhash sha256 -x node_modules/** -x '**/*.log' <directory>",
}

func readIgnoreFile(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	scanner := bufio.NewScanner(f)
	patterns := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return patterns, nil
}
