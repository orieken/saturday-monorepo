package filewriter

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteMode determines how to handle existing files
type WriteMode string

const (
	WriteModeOverwrite WriteMode = "overwrite" // Overwrite existing files
	WriteModeSkip      WriteMode = "skip"      // Skip if file exists
	WriteModeFail      WriteMode = "fail"      // Fail if file exists
)

// FileWriter handles writing generated code to files
type FileWriter struct {
	baseDir   string
	writeMode WriteMode
	dryRun    bool
}

// NewFileWriter creates a new file writer
func NewFileWriter(baseDir string, writeMode WriteMode, dryRun bool) *FileWriter {
	return &FileWriter{
		baseDir:   baseDir,
		writeMode: writeMode,
		dryRun:    dryRun,
	}
}

// WriteFile writes content to a file
func (w *FileWriter) WriteFile(relativePath string, content string) error {
	// Validate path
	if err := w.validatePath(relativePath); err != nil {
		return err
	}

	// Build full path
	fullPath := filepath.Join(w.baseDir, relativePath)

	// Check if file exists
	exists, err := w.fileExists(fullPath)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	// Handle existing file based on write mode
	if exists {
		switch w.writeMode {
		case WriteModeSkip:
			return fmt.Errorf("file already exists (skip mode): %s", relativePath)
		case WriteModeFail:
			return fmt.Errorf("file already exists (fail mode): %s", relativePath)
		case WriteModeOverwrite:
			// Continue to write
		}
	}

	// Dry run mode - don't actually write
	if w.dryRun {
		return nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// WriteMultipleFiles writes multiple files
func (w *FileWriter) WriteMultipleFiles(files map[string]string) error {
	for path, content := range files {
		if err := w.WriteFile(path, content); err != nil {
			return err
		}
	}
	return nil
}

// fileExists checks if a file exists
func (w *FileWriter) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// validatePath validates that the path is safe
func (w *FileWriter) validatePath(path string) error {
	// Prevent directory traversal
	cleanPath := filepath.Clean(path)
	if filepath.IsAbs(cleanPath) {
		return fmt.Errorf("absolute paths not allowed: %s", path)
	}

	// Check for parent directory references
	if filepath.HasPrefix(cleanPath, "..") {
		return fmt.Errorf("parent directory references not allowed: %s", path)
	}

	return nil
}

// GetFullPath returns the full path for a relative path
func (w *FileWriter) GetFullPath(relativePath string) (string, error) {
	if err := w.validatePath(relativePath); err != nil {
		return "", err
	}
	return filepath.Join(w.baseDir, relativePath), nil
}

// SetWriteMode changes the write mode
func (w *FileWriter) SetWriteMode(mode WriteMode) {
	w.writeMode = mode
}

// SetDryRun enables or disables dry run mode
func (w *FileWriter) SetDryRun(dryRun bool) {
	w.dryRun = dryRun
}

// IsDryRun returns whether dry run mode is enabled
func (w *FileWriter) IsDryRun() bool {
	return w.dryRun
}
