package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOSFileSystem_WriteFile_RoundTrip(t *testing.T) {
	fs := NewOSFileSystem()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")
	content := []byte("hello saturday")

	if err := fs.WriteFile(path, content); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q want %q", got, content)
	}
}

func TestOSFileSystem_WriteFile_Overwrites(t *testing.T) {
	fs := NewOSFileSystem()

	dir := t.TempDir()
	path := filepath.Join(dir, "over.txt")

	if err := fs.WriteFile(path, []byte("original")); err != nil {
		t.Fatalf("first write failed: %v", err)
	}
	if err := fs.WriteFile(path, []byte("replaced")); err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(got) != "replaced" {
		t.Errorf("expected overwritten content, got %q", got)
	}
}

func TestOSFileSystem_WriteFile_EmptyContent(t *testing.T) {
	fs := NewOSFileSystem()

	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	if err := fs.WriteFile(path, []byte{}); err != nil {
		t.Fatalf("WriteFile with empty content should succeed, got: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("expected 0-byte file, got %d bytes", info.Size())
	}
}

func TestOSFileSystem_WriteFile_MissingParentDirReturnsError(t *testing.T) {
	// The adapter's contract explicitly says it does NOT create parent
	// dirs — the caller owns mkdir. Verify that missing-parent surfaces
	// as an error (with the documented wrap), never a panic.
	fs := NewOSFileSystem()

	dir := t.TempDir()
	path := filepath.Join(dir, "does_not_exist", "child.txt")

	err := fs.WriteFile(path, []byte("payload"))
	if err == nil {
		t.Fatal("expected error when parent directory is missing")
	}
	if !strings.Contains(err.Error(), "failed to write file") {
		t.Errorf("expected wrapped error prefix, got %q", err.Error())
	}
}

func TestOSFileSystem_WriteFile_PermissionErrorReturnsError(t *testing.T) {
	fs := NewOSFileSystem()

	// A read-only directory should reject the write. Skip on Windows
	// (chmod semantics differ) — that's the standard Go stdlib approach
	// for this shape of test.
	dir := t.TempDir()
	readOnly := filepath.Join(dir, "readonly")
	if err := os.Mkdir(readOnly, 0500); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	// Skip if the current process can bypass permission bits (running
	// as root). CI runs as a normal user; local dev usually does too.
	if os.Geteuid() == 0 {
		t.Skip("running as root — permission bits ignored")
	}

	err := fs.WriteFile(filepath.Join(readOnly, "child.txt"), []byte("x"))
	if err == nil {
		t.Fatal("expected permission error writing under 0500 dir")
	}
	if !strings.Contains(err.Error(), "failed to write file") {
		t.Errorf("expected wrapped error prefix, got %q", err.Error())
	}
}

func TestNewOSFileSystem_ReturnsNonNil(t *testing.T) {
	if NewOSFileSystem() == nil {
		t.Fatal("NewOSFileSystem returned nil")
	}
}
