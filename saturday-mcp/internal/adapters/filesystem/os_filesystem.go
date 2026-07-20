// Package filesystem contains the os.WriteFile-backed adapter that
// implements domain.FileWriter. It is the only place in the codebase
// that touches os.WriteFile through the domain.FileWriter seam — the
// higher-level internal/filewriter package remains for the generator
// tools that need baseDir + WriteMode semantics.
package filesystem

import (
	"fmt"
	"os"
)

// OSFileSystem is the trivial os.WriteFile wrapper that satisfies
// domain.FileWriter. The struct is stateless — kept as a struct rather
// than a package function so callers can inject a fake for unit tests
// without needing a monkey-patch or a build tag.
type OSFileSystem struct{}

// NewOSFileSystem returns the adapter. There is nothing to configure —
// mode is fixed at 0644 (matching the pre-refactor os.WriteFile call in
// GenerateDocumentationTool.Execute), and permission or ownership
// concerns are left to the OS umask.
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// WriteFile writes content to the absolute path. It wraps any os error
// so the caller gets a stable error prefix rather than the raw syscall
// message — matches the wrapping style filewriter.FileWriter.WriteFile
// uses for the same operation.
func (fs *OSFileSystem) WriteFile(path string, content []byte) error {
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
