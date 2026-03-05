package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_GeneratePage(t *testing.T) {
	// Build CLI if not already built
	buildCmd := exec.Command("go", "build", "-o", "../../bin/saturday-test", "../../cmd/cli")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("../../bin/saturday-test")

	// Run generate page command
	cmd := exec.Command("../../bin/saturday-test", "generate", "page", "TestPage", 
		"--path", "/test",
		"--elements", "username:#user:input")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	
	// Verify output contains expected code
	if !strings.Contains(outputStr, "export class TestPagePage") {
		t.Error("Expected class declaration in output")
	}
	if !strings.Contains(outputStr, "extends BasePage") {
		t.Error("Expected BasePage extension in output")
	}
}

func TestCLI_GeneratePageWithWrite(t *testing.T) {
	// Build CLI if not already built
	buildCmd := exec.Command("go", "build", "-o", "../../bin/saturday-test", "../../cmd/cli")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("../../bin/saturday-test")

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "cli-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Run generate page command with write
	cmd := exec.Command("../../bin/saturday-test", "generate", "page", "TestPage", 
		"--path", "/test",
		"--elements", "username:#user:input",
		"--write",
		"-o", tmpDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, "lib", "pages", "test-page-page.ts")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}

	// Read and verify content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "export class TestPagePage") {
		t.Error("Generated file missing class declaration")
	}
}

func TestCLI_Help(t *testing.T) {
	// Build CLI if not already built
	buildCmd := exec.Command("go", "build", "-o", "../../bin/saturday-test", "../../cmd/cli")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}
	defer os.Remove("../../bin/saturday-test")

	// Run help command
	cmd := exec.Command("../../bin/saturday-test", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	outputStr := string(output)
	
	// Verify help output
	expectedCommands := []string{"generate", "analyze", "validate", "migrate", "docs"}
	for _, cmdName := range expectedCommands {
		if !strings.Contains(outputStr, cmdName) {
			t.Errorf("Help output missing command: %s", cmdName)
		}
	}
}
