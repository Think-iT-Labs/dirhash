package tests

import (
	"os"
	"path/filepath"
	"testing"

	lib "github.com/Think-iT-Labs/dirhash/lib"
)

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
}

func TestDirHashDeterministic(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "hello")
	writeFile(t, filepath.Join(dir, "b.txt"), "world")

	h1, err := lib.DirHash(dir, nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	h2, err := lib.DirHash(dir, nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	if h1 == "" || h2 == "" {
		t.Fatalf("hash should not be empty: %q %q", h1, h2)
	}
	if h1 != h2 {
		t.Fatalf("hashes must match for unchanged directory: %q vs %q", h1, h2)
	}
}

func TestDirHashIgnores(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "keep.txt"), "keep")
	writeFile(t, filepath.Join(dir, "debug.log"), "ignore me")
	writeFile(t, filepath.Join(dir, "sub", "temp.tmp"), "ignore me too")

	base, err := lib.DirHash(dir, nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	if base == "" {
		t.Fatal("baseline hash empty")
	}
	ignored, err := lib.DirHash(dir, []string{"**/*.log", "**/*.tmp"})
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	if ignored == "" {
		t.Fatal("ignored hash empty")
	}
	if ignored == base {
		t.Fatal("hash should change when files are ignored")
	}
}

func TestDirHashRelativeVsAbsolute(t *testing.T) {
	work := t.TempDir()
	project := filepath.Join(work, "proj")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeFile(t, filepath.Join(project, "a.txt"), "same")
	writeFile(t, filepath.Join(project, "dir", "b.txt"), "same2")

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldwd) })
	if err := os.Chdir(work); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	rel, err := lib.DirHash("proj", nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	abs, err := lib.DirHash(project, nil)
	if err != nil {
		t.Fatalf("DirHash failed: %v", err)
	}
	if rel != abs {
		t.Fatalf("relative and absolute hashing should match: %q vs %q", rel, abs)
	}
}
