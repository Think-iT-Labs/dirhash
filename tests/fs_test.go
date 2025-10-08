package tests

import (
	"os"
	"path/filepath"
	"testing"

	lib "github.com/Think-iT-Labs/dirhash/lib"
)

func mustWrite(t *testing.T, p, s string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

// ? Verifies that ignoring a directory via patterns produces the same hash as hashing
// ? a mirror directory that simply does not contain that directory
func TestIgnoreDirectoryBehavior(t *testing.T) {
	src := t.TempDir()
	mustWrite(t, filepath.Join(src, "keep.txt"), "ok")
	mustWrite(t, filepath.Join(src, "dir", "a.txt"), "a")
	mustWrite(t, filepath.Join(src, "dir", "b.txt"), "b")

	mirror := t.TempDir()
	mustWrite(t, filepath.Join(mirror, "keep.txt"), "ok")

	ignoredHash, err := lib.DirHash(src, []string{"dir/**"})
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	mirrorHash, err := lib.DirHash(mirror, nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}

	if ignoredHash != mirrorHash {
		t.Fatalf("ignoring dir/** should match mirror without dir, got %q vs %q", ignoredHash, mirrorHash)
	}
}

// ? Verifies doublestar glob patterns (e.g., **/*.log) affect the final hash as expected
func TestGlobStarMatchesAffectHash(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "keep.txt"), "ok")
	mustWrite(t, filepath.Join(dir, "logs", "app.log"), "x")

	mirror := t.TempDir()
	mustWrite(t, filepath.Join(mirror, "keep.txt"), "ok")

	ignored, err := lib.DirHash(dir, []string{"**/*.log"})
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	expected, err := lib.DirHash(mirror, nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	if ignored != expected {
		t.Fatalf("ignoring **/*.log should match mirror without log file, got %q vs %q", ignored, expected)
	}
}
