# Dirhash

[![CI](https://img.shields.io/github/actions/workflow/status/Think-iT-Labs/dirhash/ci.yml?branch=main&label=CI)](https://github.com/Think-iT-Labs/dirhash/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/tag/Think-iT-Labs/dirhash?label=release)](https://github.com/Think-iT-Labs/dirhash/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Think-iT-Labs/dirhash)](https://go.dev/)

Calculating the checksum of a directory made easy.

## Dirhash CLI
Compute a deterministic checksum of a directory by hashing file contents and their relative paths, with support for glob-based exclusions and symbolic link navigation.

### Usage
```sh
dirhash sha256 [--ignore pattern]... <directory>
```

### Examples
```sh
# Basic: hash current directory
dirhash sha256 .

# Ignore common folders and logs
dirhash sha256 -x node_modules/** -x ".git/**" -x "**/*.log" .

# Show which files are processed (debug logging)
LOG_LEVEL=debug dirhash sha256 .

# JSON output with per-file hashes (consumable by CI tools)
dirhash sha256 -o json -x "**/*.log" .
```

### Ignore patterns
- Supports doublestar `**` for recursive matching via `github.com/bmatcuk/doublestar`.
- Patterns are resolved relative to the provided directory argument.
- Examples:
  - `-x .git/**` to exclude the entire Git directory
  - `-x vendor/**` to exclude vendored dependencies
  - `-x "**/*.tmp"` to exclude all temporary files
  - `-x "**/*.log"` to exclude all log files

Tips:
- Quote patterns containing `*` or `**` to avoid shell expansion.
- Combine multiple `-x` flags to refine your selection.

Ignore file:
- Provide patterns via a file with `--ignore-file <path>`.
- If `--ignore-file` is not set, a `.dirhashignore` file in the target directory will be used automatically when present.

### Symbolic Link Support

Dirhash automatically follows symbolic links during directory traversal:

- **File symlinks**: Symbolic links pointing to files are included in the hash using the symlink path
- **Directory symlinks**: Symbolic links pointing to directories are recursively traversed
- **Relative symlinks**: Resolved relative to the symlink's directory
- **Absolute symlinks**: Followed directly to their target
- **Circular symlinks**: Detected and handled gracefully to prevent infinite loops
- **Broken symlinks**: Ignored silently (logged at debug level)

Examples:
```sh
# Hash directory with symlinks (shows debug info about symlink following)
LOG_LEVEL=debug dirhash sha256 /path/with/symlinks

# Ignore symlinks with specific patterns
dirhash sha256 -x "**/*.link" -x "symlink_*" .
```

**Note**: On Windows, creating symbolic links requires administrator privileges, so symlink tests are skipped in the test suite.

### Build locally
 1. First download dependencies
 ```sh
 go mod download
 ```
 2. Build the CLI
 ```sh
 make build-cli
 ```

### Exit codes
- 0: success
- non-zero: fatal error occurred (see stderr for details)

### Notes for Windows users
- Use quotes around patterns with `*` to avoid shell expansion: `-x "**/*.log"`.
- Paths are handled via Go's `filepath`; separators are normalized automatically.
