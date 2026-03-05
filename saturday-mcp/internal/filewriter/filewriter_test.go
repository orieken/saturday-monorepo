package filewriter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileWriter_WriteFile(t *testing.T) {
	// Create temp directory for tests
	tmpDir := t.TempDir()

	t.Run("Write new file", func(t *testing.T) {
		writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)
		
		err := writer.WriteFile("test.txt", "Hello, World!")
		if err != nil {
			t.Fatalf("Expected successful write, got error: %v", err)
		}

		// Verify file was written
		content, err := os.ReadFile(filepath.Join(tmpDir, "test.txt"))
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(content) != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got %q", string(content))
		}
	})

	t.Run("Write file in subdirectory", func(t *testing.T) {
		writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)
		
		err := writer.WriteFile("subdir/nested/file.txt", "Nested content")
		if err != nil {
			t.Fatalf("Expected successful write, got error: %v", err)
		}

		// Verify file was written
		content, err := os.ReadFile(filepath.Join(tmpDir, "subdir/nested/file.txt"))
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(content) != "Nested content" {
			t.Errorf("Expected 'Nested content', got %q", string(content))
		}
	})

	t.Run("Overwrite mode - overwrites existing file", func(t *testing.T) {
		writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)
		
		// Write initial file
		writer.WriteFile("overwrite.txt", "Original")
		
		// Overwrite
		err := writer.WriteFile("overwrite.txt", "Updated")
		if err != nil {
			t.Fatalf("Expected successful overwrite, got error: %v", err)
		}

		// Verify content was updated
		content, err := os.ReadFile(filepath.Join(tmpDir, "overwrite.txt"))
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(content) != "Updated" {
			t.Errorf("Expected 'Updated', got %q", string(content))
		}
	})

	t.Run("Skip mode - skips existing file", func(t *testing.T) {
		writer := NewFileWriter(tmpDir, WriteModeSkip, false)
		
		// Write initial file
		writer.SetWriteMode(WriteModeOverwrite)
		writer.WriteFile("skip.txt", "Original")
		
		// Try to write again with skip mode
		writer.SetWriteMode(WriteModeSkip)
		err := writer.WriteFile("skip.txt", "Updated")
		if err == nil {
			t.Error("Expected error for existing file in skip mode")
		}

		// Verify content wasn't changed
		content, err := os.ReadFile(filepath.Join(tmpDir, "skip.txt"))
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if string(content) != "Original" {
			t.Errorf("Expected 'Original', got %q", string(content))
		}
	})

	t.Run("Fail mode - fails on existing file", func(t *testing.T) {
		writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)
		
		// Write initial file
		writer.WriteFile("fail.txt", "Original")
		
		// Try to write again with fail mode
		writer.SetWriteMode(WriteModeFail)
		err := writer.WriteFile("fail.txt", "Updated")
		if err == nil {
			t.Error("Expected error for existing file in fail mode")
		}
	})

	t.Run("Dry run mode - doesn't write file", func(t *testing.T) {
		writer := NewFileWriter(tmpDir, WriteModeOverwrite, true)
		
		err := writer.WriteFile("dryrun.txt", "Should not be written")
		if err != nil {
			t.Fatalf("Expected no error in dry run, got: %v", err)
		}

		// Verify file was NOT written
		_, err = os.ReadFile(filepath.Join(tmpDir, "dryrun.txt"))
		if err == nil {
			t.Error("Expected file to not exist in dry run mode")
		}
	})
}

func TestFileWriter_WriteMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)

	files := map[string]string{
		"file1.txt":       "Content 1",
		"dir/file2.txt":   "Content 2",
		"dir/file3.txt":   "Content 3",
	}

	err := writer.WriteMultipleFiles(files)
	if err != nil {
		t.Fatalf("Expected successful write, got error: %v", err)
	}

	// Verify all files were written
	for path, expectedContent := range files {
		content, err := os.ReadFile(filepath.Join(tmpDir, path))
		if err != nil {
			t.Errorf("Failed to read %s: %v", path, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("File %s: expected %q, got %q", path, expectedContent, string(content))
		}
	}
}

func TestFileWriter_ValidatePath(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)

	tests := []struct {
		name      string
		path      string
		shouldErr bool
	}{
		{"Valid relative path", "file.txt", false},
		{"Valid nested path", "dir/subdir/file.txt", false},
		{"Absolute path", "/etc/passwd", true},
		{"Parent directory reference", "../file.txt", true},
		{"Parent in middle", "dir/../../file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writer.validatePath(tt.path)
			if tt.shouldErr && err == nil {
				t.Error("Expected error for invalid path")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected valid path, got error: %v", err)
			}
		})
	}
}

func TestFileWriter_GetFullPath(t *testing.T) {
	tmpDir := t.TempDir()
	writer := NewFileWriter(tmpDir, WriteModeOverwrite, false)

	t.Run("Valid path", func(t *testing.T) {
		fullPath, err := writer.GetFullPath("test.txt")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expected := filepath.Join(tmpDir, "test.txt")
		if fullPath != expected {
			t.Errorf("Expected %q, got %q", expected, fullPath)
		}
	})

	t.Run("Invalid path", func(t *testing.T) {
		_, err := writer.GetFullPath("../test.txt")
		if err == nil {
			t.Error("Expected error for invalid path")
		}
	})
}

func TestFileWriter_SettersGetters(t *testing.T) {
	writer := NewFileWriter(t.TempDir(), WriteModeOverwrite, false)

	t.Run("SetWriteMode", func(t *testing.T) {
		writer.SetWriteMode(WriteModeSkip)
		// Verify by attempting to write
		// (mode is private, so we test behavior)
	})

	t.Run("SetDryRun and IsDryRun", func(t *testing.T) {
		writer.SetDryRun(true)
		if !writer.IsDryRun() {
			t.Error("Expected dry run to be enabled")
		}

		writer.SetDryRun(false)
		if writer.IsDryRun() {
			t.Error("Expected dry run to be disabled")
		}
	})
}
