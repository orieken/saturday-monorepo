package analyzers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/logging"
)

func TestPatternValidator_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create directories
	os.MkdirAll(filepath.Join(tmpDir, "lib", "pages"), 0755)

	// Case 1: Valid Page
	validPage := `
import { BasePage } from './BasePage';
export class ValidPage extends BasePage {}
`
	os.WriteFile(filepath.Join(tmpDir, "lib", "pages", "valid-page.ts"), []byte(validPage), 0644)

	// Case 2: Invalid Naming
	invalidNaming := `
import { BasePage } from './BasePage';
export class WrongName extends BasePage {}
`
	os.WriteFile(filepath.Join(tmpDir, "lib", "pages", "invalid-naming-page.ts"), []byte(invalidNaming), 0644)

	// Case 3: Invalid Inheritance
	invalidInheritance := `
export class LoginPage {}
`
	os.WriteFile(filepath.Join(tmpDir, "lib", "pages", "login-page.ts"), []byte(invalidInheritance), 0644)

	// Run Validator
	logger := logging.NewLogger(os.Stderr)
	validator := NewPatternValidator(logger)

	result, err := validator.Validate(tmpDir)
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected validation to fail")
	}

	if len(result.Issues) == 0 {
		t.Fatal("Expected issues to be found")
	}

	// Check for specific issues
	namingFound := false
	inheritanceFound := false

	for _, issue := range result.Issues {
		if issue.Rule == "NamingConvention" && issue.File == filepath.Join(tmpDir, "lib", "pages", "invalid-naming-page.ts") {
			namingFound = true
		}
		if issue.Rule == "Inheritance" && issue.File == filepath.Join(tmpDir, "lib", "pages", "login-page.ts") {
			inheritanceFound = true
		}
	}

	if !namingFound {
		t.Error("Expected naming convention issue not found")
	}
	if !inheritanceFound {
		t.Error("Expected inheritance issue not found")
	}
}
