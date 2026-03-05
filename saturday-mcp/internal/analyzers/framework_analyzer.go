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

// Regex patterns for detecting components
var (
	pageClassRegex = regexp.MustCompile(`class\s+(\w+?)(?:Page)\s+(?:extends\s+(\w+)\s+)?{`)
	flowClassRegex = regexp.MustCompile(`class\s+(\w+?)(?:Flow)\s+(?:extends\s+(\w+)\s+)?{`)
	stepDefRegex   = regexp.MustCompile(`(Given|When|Then)\(['"](.+?)['"]`)
)

// FrameworkAnalyzer analyzes a Saturday framework project
type FrameworkAnalyzer struct {
	logger *logging.Logger
}

// NewFrameworkAnalyzer creates a new analyzer
func NewFrameworkAnalyzer(logger *logging.Logger) *FrameworkAnalyzer {
	return &FrameworkAnalyzer{
		logger: logger,
	}
}

// Analyze scans the given project path and returns analysis results
func (a *FrameworkAnalyzer) Analyze(rootPath string) (*models.AnalysisResult, error) {
	result := &models.AnalysisResult{
		Structure: models.ProjectStructure{
			RootPath: rootPath,
		},
	}

	a.logger.Info("Starting project analysis", "path", rootPath)

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip node_modules and hidden directories
			if info.Name() == "node_modules" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process TS/JS files
		ext := filepath.Ext(path)
		if ext != ".ts" && ext != ".js" {
			return nil
		}
		
		result.Stats.TotalFiles++

		// Analyze file content
		if err := a.analyzeFile(path, result); err != nil {
			a.logger.Error("Failed to analyze file", "path", path, "error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	a.logger.Info("Analysis complete", 
		"pages", result.Stats.PageCount,
		"flows", result.Stats.FlowCount,
		"steps", result.Stats.StepDefCount,
	)

	return result, nil
}

func (a *FrameworkAnalyzer) analyzeFile(path string, result *models.AnalysisResult) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	// Temporary storage for file info
	var currentSteps []models.Step
	isPage := false
	isFlow := false
	isStepFile := false

	// Check path for structural clues
	relPath := path // ideally relative to root
	if strings.Contains(strings.ToLower(path), "/pages/") {
		result.Structure.PagesPath = filepath.Dir(path)
	} else if strings.Contains(strings.ToLower(path), "/flows/") {
		result.Structure.FlowsPath = filepath.Dir(path)
	} else if strings.Contains(strings.ToLower(path), "/steps/") {
		result.Structure.StepsPath = filepath.Dir(path)
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		// Check for Page Class
		if matches := pageClassRegex.FindStringSubmatch(line); len(matches) > 1 {
			isPage = true
			result.Stats.PageCount++
			result.Inventory.Pages = append(result.Inventory.Pages, models.PageInfo{
				Name: matches[1] + "Page",
				FilePath: relPath,
			})
		}

		// Check for Flow Class
		if matches := flowClassRegex.FindStringSubmatch(line); len(matches) > 1 {
			isFlow = true
			result.Stats.FlowCount++
			result.Inventory.Flows = append(result.Inventory.Flows, models.FlowInfo{
				Name: matches[1] + "Flow",
				FilePath: relPath,
			})
		}

		// Check for Step Definitions
		if matches := stepDefRegex.FindStringSubmatch(line); len(matches) > 2 {
			isStepFile = true
			currentSteps = append(currentSteps, models.Step{
				Type:    matches[1],
				Pattern: matches[2],
				Line:    lineNum,
			})
		}
	}

	if len(currentSteps) > 0 {
		result.Stats.StepDefCount += len(currentSteps)
		result.Inventory.StepDefs = append(result.Inventory.StepDefs, models.StepDefInfo{
			FilePath: relPath,
			Steps:    currentSteps,
		})
	}

    // Heuristics for structure if not explicitly found by path
	if isPage && result.Structure.PagesPath == "" {
        result.Structure.PagesPath = filepath.Dir(path)
    }
    if isFlow && result.Structure.FlowsPath == "" {
        result.Structure.FlowsPath = filepath.Dir(path)
    }
    if isStepFile && result.Structure.StepsPath == "" {
        result.Structure.StepsPath = filepath.Dir(path)
    }

	return scanner.Err()
}
