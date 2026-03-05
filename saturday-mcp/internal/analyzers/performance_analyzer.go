package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// PerformanceAnalyzer analyzes code for performance bottlenecks
type PerformanceAnalyzer struct {
	logger *logging.Logger
}

// NewPerformanceAnalyzer creates a new performance analyzer
func NewPerformanceAnalyzer(logger *logging.Logger) *PerformanceAnalyzer {
	return &PerformanceAnalyzer{
		logger: logger,
	}
}

// Analyze scans the project for performance issues concurrently
func (a *PerformanceAnalyzer) Analyze(rootPath string) (*models.ImprovementResponse, error) {
	a.logger.Info("Starting performance analysis", "path", rootPath)
	start := time.Now()

	var allIssues []models.ValidationIssue
	var issuesMutex sync.Mutex
	var wg sync.WaitGroup

	// Channel to limit concurrency
	sem := make(chan struct{}, 10) // Limit to 10 concurrent file reads

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

		ext := filepath.Ext(path)
		if ext != ".ts" && ext != ".js" {
			return nil
		}

		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			sem <- struct{}{} // Acquire token
			defer func() { <-sem }() // Release token (ensure it runs)

			issues, err := a.analyzeFile(filePath)
			if err != nil {
				a.logger.Error("Failed to analyze file", "path", filePath, "error", err)
				return
			}

			if len(issues) > 0 {
				issuesMutex.Lock()
				allIssues = append(allIssues, issues...)
				issuesMutex.Unlock()
			}
		}(path)

		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	duration := time.Since(start)
	
	response := &models.ImprovementResponse{
		Suggestions: allIssues,
		Summary: map[string]interface{}{
			"scanDurationMs": duration.Milliseconds(),
			"issuesFound":    len(allIssues),
		},
	}

	return response, nil
}

func (a *PerformanceAnalyzer) analyzeFile(path string) ([]models.ValidationIssue, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var issues []models.ValidationIssue

	// check file size
	// Rule: Large files (> 500 lines)
	// We'll count lines while scanning

	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	// Regex patterns
	// Naive check for await inside loops (e.g. for ... await ...)
	// This is hard to do line-by-line w/o AST, but we can check for `await` and `for (` proximity or indentation context
	// Let's stick to simpler regex for MVP:
	// 1. Long selectors in quotes
	// 2. High timeouts
	
	// 1. Long selectors in quotes (starts with alphanumeric, . or #, > 100 chars)
	longSelectorRegex := regexp.MustCompile(`["']([a-zA-Z0-9#\.][^"']{100,})["']`)
	highTimeoutRegex := regexp.MustCompile(`timeout:\s*([0-9]{5,})`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Rule: Long Selectors
		if matches := longSelectorRegex.FindStringSubmatch(line); len(matches) > 1 {
			issues = append(issues, models.ValidationIssue{
				Severity: "warning",
				Message:  "Selector is too long and complex (>100 chars). Performance may degrade.",
				File:     path,
				Line:     lineNum,
				Rule:     "simplify-selector",
			})
		}

		// Rule: Excessive Timeout (>9999ms i.e. 5 digits)
		if matches := highTimeoutRegex.FindStringSubmatch(line); len(matches) > 1 {
			issues = append(issues, models.ValidationIssue{
				Severity: "warning",
				Message:  "Hardcoded timeout > 10s detected. Use default timeouts configured in playwright.config.",
				File:     path,
				Line:     lineNum,
				Rule:     "avoid-high-timeouts",
			})
		}
	}

	// Rule: Large File
	if lineNum > 500 {
		issues = append(issues, models.ValidationIssue{
			Severity: "info",
			Message:  "File exceeds 500 lines. Consider splitting into multiple files.",
			File:     path,
			Line:     0,
			Rule:     "large-file",
		})
	}

	// Performance calculation for large file size (bytes)
	if info.Size() > 50*1024 { // 50KB
		issues = append(issues, models.ValidationIssue{
			Severity: "info",
			Message:  "File size exceeds 50KB.",
			File:     path,
			Line:     0,
			Rule:     "large-file-size",
		})
	}

	return issues, scanner.Err()
}
