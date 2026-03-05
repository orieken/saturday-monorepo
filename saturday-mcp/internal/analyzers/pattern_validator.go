package analyzers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// PatternValidator validates code against Saturday framework patterns
type PatternValidator struct {
	logger *logging.Logger
	rules  []ValidationRule
}

// NewPatternValidator creates a new validator with default rules
func NewPatternValidator(logger *logging.Logger) *PatternValidator {
	return &PatternValidator{
		logger: logger,
		rules: []ValidationRule{
			&NamingConventionRule{},
			&InheritanceRule{},
		},
	}
}

// Validate project structure and code patterns
func (v *PatternValidator) Validate(rootPath string) (*models.ValidationResponse, error) {
	v.logger.Info("Starting pattern validation", "path", rootPath)
	
	var allIssues []models.ValidationIssue

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check TS files
		if filepath.Ext(path) != ".ts" {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			v.logger.Error("Failed to read file", "path", path, "error", err)
			return nil // continue walking
		}

		// Apply rules
		fileContent := string(content)
		for _, rule := range v.rules {
			issues := rule.Check(fileContent, path)
			if len(issues) > 0 {
				allIssues = append(allIssues, issues...)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	response := &models.ValidationResponse{
		Valid:  len(allIssues) == 0,
		Issues: allIssues,
	}

	v.logger.Info("Validation complete", "issues", len(allIssues))
	return response, nil
}
