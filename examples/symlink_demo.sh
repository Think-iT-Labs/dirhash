#!/bin/bash

# Symbolic Link Demonstration for dirhash
# This script demonstrates how dirhash handles various types of symbolic links
# Run this on Unix-like systems (Linux, macOS) where symlinks are supported

set -e

echo "=== Dirhash Symbolic Link Demo ==="
echo

# Create a temporary directory for our demo
DEMO_DIR=$(mktemp -d)
echo "Demo directory: $DEMO_DIR"
cd "$DEMO_DIR"

# Build dirhash if not already built
if ! command -v dirhash &> /dev/null; then
    echo "Building dirhash..."
    go build -o dirhash github.com/Think-iT-Labs/dirhash
    export PATH="$PWD:$PATH"
fi

echo
echo "=== Setting up test structure ==="

# Create some regular files and directories
mkdir -p regular_dir/subdir
echo "content1" > regular_dir/file1.txt
echo "content2" > regular_dir/subdir/file2.txt
echo "standalone" > standalone.txt

# Create target files/dirs for symlinks
mkdir -p target_dir
echo "target_content" > target_dir/target_file.txt
echo "external_content" > external_file.txt

echo "Created regular files and directories"

echo
echo "=== Creating symbolic links ==="

# 1. Relative symlink to file
ln -s external_file.txt relative_file_link.txt
echo "Created relative file symlink: relative_file_link.txt -> external_file.txt"

# 2. Absolute symlink to file
ln -s "$PWD/external_file.txt" absolute_file_link.txt
echo "Created absolute file symlink: absolute_file_link.txt -> $PWD/external_file.txt"

# 3. Relative symlink to directory
ln -s target_dir relative_dir_link
echo "Created relative directory symlink: relative_dir_link -> target_dir"

# 4. Absolute symlink to directory
ln -s "$PWD/target_dir" absolute_dir_link
echo "Created absolute directory symlink: absolute_dir_link -> $PWD/target_dir"

# 5. Broken symlink
ln -s nonexistent_file.txt broken_link.txt
echo "Created broken symlink: broken_link.txt -> nonexistent_file.txt"

echo
echo "=== Directory structure ==="
find . -type f -o -type l | sort

echo
echo "=== Testing dirhash with symlinks ==="

echo
echo "1. Basic hash (includes all symlinks):"
dirhash sha256 .

echo
echo "2. Hash with debug logging (shows symlink following):"
LOG_LEVEL=debug dirhash sha256 . 2>&1 | grep -E "(Following symlink|Skipping|Failed to read symlink)" || echo "No symlink debug messages (all processed successfully)"

echo
echo "3. Hash ignoring symlinks:"
dirhash sha256 -x "*_link*" .

echo
echo "4. JSON output showing all files including symlinks:"
dirhash sha256 -o json . | jq '.files[] | select(.path | contains("link")) | {path, hash}'

echo
echo "=== Comparison: Directory without symlinks ==="

# Create comparison directory without symlinks
mkdir -p comparison
cp -r regular_dir comparison/
cp standalone.txt comparison/
cp external_file.txt comparison/
cp -r target_dir comparison/

echo "Hash of directory without symlinks:"
dirhash sha256 comparison

echo
echo "Hash of directory with symlinks (ignoring all symlinks):"
dirhash sha256 -x "*link*" .

echo
echo "=== Cleanup ==="
cd /
rm -rf "$DEMO_DIR"
echo "Demo completed successfully!"
