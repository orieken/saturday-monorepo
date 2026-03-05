package analyzers

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/orieken/saturday-mcp/internal/models"
)

// ValidationRule interface for all validation rules
type ValidationRule interface {
	Name() string
	Check(fileContent string, filePath string) []models.ValidationIssue
}

// NamingConventionRule checks if class names match their file names and types
type NamingConventionRule struct{}

func (r *NamingConventionRule) Name() string {
	return "NamingConvention"
}

func (r *NamingConventionRule) Check(content string, path string) []models.ValidationIssue {
	var issues []models.ValidationIssue
	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	// Check Page classes
	if strings.Contains(path, "/pages/") {
		// Expect PascalCase class name matching kebab-case filename
		// e.g., "login-page.ts" -> "class LoginPage"
		className := toPascalCase(strings.TrimSuffix(fileName, "-page")) + "Page"
		classRegex := regexp.MustCompile(`class\s+` + className)
		
		if !classRegex.MatchString(content) {
			issues = append(issues, models.ValidationIssue{
				Severity: "warning",
				Message:  "Page class name does not match filename convention. Suggestion: Rename class to " + className,
				File:     path,
				Rule:     r.Name(),
			})
		}
	}

	// Check Flow classes
	if strings.Contains(path, "/flows/") {
		className := toPascalCase(strings.TrimSuffix(fileName, "-flow")) + "Flow"
		classRegex := regexp.MustCompile(`class\s+` + className)
		
		if !classRegex.MatchString(content) {
			issues = append(issues, models.ValidationIssue{
				Severity: "warning",
				Message:  "Flow class name does not match filename convention. Suggestion: Rename class to " + className,
				File:     path,
				Rule:     r.Name(),
			})
		}
	}

	return issues
}

// InheritanceRule checks if classes extend the correct base classes
type InheritanceRule struct{}

func (r *InheritanceRule) Name() string {
	return "Inheritance"
}

func (r *InheritanceRule) Check(content string, path string) []models.ValidationIssue {
	var issues []models.ValidationIssue

	if strings.Contains(path, "/pages/") {
		if !strings.Contains(content, "extends BasePage") {
			issues = append(issues, models.ValidationIssue{
				Severity: "error",
				Message:  "Page class must extend BasePage. Suggestion: Update class to extend BasePage",
				File:     path,
				Rule:     r.Name(),
			})
		}
	}

	if strings.Contains(path, "/flows/") {
		if !strings.Contains(content, "extends BaseFlow") {
			issues = append(issues, models.ValidationIssue{
				Severity: "error",
				Message:  "Flow class must extend BaseFlow. Suggestion: Update class to extend BaseFlow",
				File:     path,
				Rule:     r.Name(),
			})
		}
	}

	return issues
}

// Helper to convert string to PascalCase (simple version for validator)
func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
