package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// ImprovementAnalyzer analyzes code for potential improvements
type ImprovementAnalyzer struct {
	logger *logging.Logger
}

// NewImprovementAnalyzer creates a new improvement analyzer
func NewImprovementAnalyzer(logger *logging.Logger) *ImprovementAnalyzer {
	return &ImprovementAnalyzer{
		logger: logger,
	}
}

// Analyze suggests improvements for the project
func (a *ImprovementAnalyzer) Analyze(rootPath string) (*models.ImprovementResponse, error) {
	a.logger.Info("Starting improvement analysis", "path", rootPath)
	
	var allSuggestions []models.ValidationIssue

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

		// Only check TS/JS files
		ext := filepath.Ext(path)
		if ext != ".ts" && ext != ".js" {
			return nil
		}

		// Read file content and check line by line
		suggestions, err := a.analyzeFile(path)
		if err != nil {
			a.logger.Error("Failed to analyze file", "path", path, "error", err)
			return nil
		}

		allSuggestions = append(allSuggestions, suggestions...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	response := &models.ImprovementResponse{
		Suggestions: allSuggestions,
		Summary: map[string]interface{}{
			"totalSuggestions": len(allSuggestions),
		},
	}

	return response, nil
}

func (a *ImprovementAnalyzer) analyzeFile(path string) ([]models.ValidationIssue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var suggestions []models.ValidationIssue
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Compile regex for patterns
	hardWaitRegex := regexp.MustCompile(`waitForTimeout\(\d+\)`)
	consoleLogRegex := regexp.MustCompile(`console\.log\(`)
	anyTypeRegex := regexp.MustCompile(`: any`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check: Hard waits
		if hardWaitRegex.MatchString(line) {
			suggestions = append(suggestions, models.ValidationIssue{
				Severity: "warning",
				Message:  "Avoid using waitForTimeout. Use waitForSelector or assertions instead.",
				File:     path,
				Line:     lineNum,
				Rule:     "no-hard-waits",
			})
		}

		// Check: Console logs
		if consoleLogRegex.MatchString(line) {
			suggestions = append(suggestions, models.ValidationIssue{
				Severity: "info",
				Message:  "Remove console.log statements before merging.",
				File:     path,
				Line:     lineNum,
				Rule:     "no-console-log",
			})
		}

		// Check: Explicit 'any' type
		if anyTypeRegex.MatchString(line) {
			suggestions = append(suggestions, models.ValidationIssue{
				Severity: "warning",
				Message:  "Avoid using 'any' type. Define strict interfaces.",
				File:     path,
				Line:     lineNum,
				Rule:     "no-any",
			})
		}
	}

	return suggestions, scanner.Err()
}
