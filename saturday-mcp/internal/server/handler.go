package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/executor"
	"github.com/orieken/saturday-mcp/internal/filewriter"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/observability"
	"github.com/orieken/saturday-mcp/internal/prompts"
	"github.com/orieken/saturday-mcp/internal/resources"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// Handler manages the Saturday MCP server operations
type Handler struct {
	logger    *logging.Logger
	generator        *generators.Generator
	analyzer         *analyzers.FrameworkAnalyzer
	validator        *analyzers.PatternValidator
	resourceProvider *resources.Provider
	promptProvider   *prompts.Provider
	improvementAnalyzer *analyzers.ImprovementAnalyzer
	performanceAnalyzer *analyzers.PerformanceAnalyzer
	graphAnalyzer       *analyzers.GraphAnalyzer
	logAnalyzer         *analyzers.TestLogAnalyzer
	usageAnalyzer       *analyzers.UsageAnalyzer
	testExecutor        *executor.TestExecutor
}

// NewHandler creates a new server handler
func NewHandler(logger *logging.Logger) (*Handler, error) {
	// Initialize template system
	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	// Load all templates
	if err := loader.LoadAll(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	// Initialize validator
	validator := validators.NewValidator()

	// Initialize generators
	siteGen := generators.NewSiteGenerator(processor, validator)
	pageGen := generators.NewPageGenerator(processor, validator)
	flowGen := generators.NewFlowGenerator(processor, validator)
	stepGen := generators.NewStepGenerator(processor, validator)
	elementGen := generators.NewElementGenerator(processor, validator)
	serviceGen := generators.NewServiceGenerator(processor, validator)
	migrationGen := generators.NewMigrationGenerator(processor, validator)
	docGen := generators.NewDocumentationGenerator(processor, validator)
	generator := generators.NewGenerator(siteGen, pageGen, flowGen, stepGen, elementGen, serviceGen, migrationGen, docGen)

	// Initialize analyzer and validator
	analyzer := analyzers.NewFrameworkAnalyzer(logger)
	validatorTool := analyzers.NewPatternValidator(logger)
	improvementAnalyzer := analyzers.NewImprovementAnalyzer(logger)
	performanceAnalyzer := analyzers.NewPerformanceAnalyzer(logger)
	graphAnalyzer := analyzers.NewGraphAnalyzer(logger)
	logAnalyzer := analyzers.NewTestLogAnalyzer(logger)
	usageAnalyzer := analyzers.NewUsageAnalyzer(logger)
	resourceProvider := resources.NewProvider(logger, loader)
	promptProvider := prompts.NewProvider(logger)
	testExecutor := executor.NewTestExecutor(logger)

	return &Handler{
		logger:              logger,
		generator:           generator,
		analyzer:            analyzer,
		validator:           validatorTool,
		resourceProvider:    resourceProvider,
		promptProvider:      promptProvider,
		improvementAnalyzer: improvementAnalyzer,
		performanceAnalyzer: performanceAnalyzer,
		graphAnalyzer:       graphAnalyzer,
		logAnalyzer:         logAnalyzer,
		usageAnalyzer:       usageAnalyzer,
		testExecutor:        testExecutor,
	}, nil
}





// RegisterTools registers all available tools with the MCP server
func (h *Handler) RegisterTools(s *server.MCPServer) error {
	h.logger.Info("Registering tools")

	// Register list_tools endpoint
	s.AddTool(mcp.Tool{
		Name:        "list_tools",
		Description: "List all available Saturday framework generation and analysis tools",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, h.handleListTools)

	// Register generate_site tool
	s.AddTool(mcp.Tool{
		Name:        "generate_site",
		Description: "Generate a Site class with page and flow registration",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"name", "baseUrl", "pages"},
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the site class",
				},
				"baseUrl": map[string]interface{}{
					"type":        "string",
					"description": "Base URL for the site",
				},
				"pages": map[string]interface{}{
					"type":        "array",
					"description": "List of page names to register",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"flows": map[string]interface{}{
					"type":        "array",
					"description": "Optional list of flow names",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the site",
				},
				"writeToFile": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to write the generated code to a file (default: false)",
					"default":     false,
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Base directory for output files (required if writeToFile is true)",
				},
			},
		},
	}, h.handleGenerateSite)

	// Register generate_page tool
	s.AddTool(mcp.Tool{
		Name:        "generate_page",
		Description: "Generate a Page class with element registration",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"name", "path", "elements"},
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the page class",
				},
				"path": map[string]interface{}{
					"type":        "string",
					"description": "URL path for the page",
				},
				"elements": map[string]interface{}{
					"type":        "array",
					"description": "List of elements on the page",
					"items": map[string]interface{}{
						"type": "object",
						"required": []string{"name", "selector"},
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type":        "string",
								"description": "Element name",
							},
							"selector": map[string]interface{}{
								"type":        "string",
								"description": "CSS selector for the element",
							},
							"type": map[string]interface{}{
								"type":        "string",
								"description": "Element type (button, input, link, select, checkbox, radio)",
								"enum":        []string{"button", "input", "link", "select", "checkbox", "radio"},
							},
						},
					},
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the page",
				},
				"writeToFile": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to write the generated code to a file (default: false)",
					"default":     false,
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Base directory for output files (required if writeToFile is true)",
				},
			},
		},
	}, h.handleGeneratePage)

	// Register generate_flow tool
	s.AddTool(mcp.Tool{
		Name:        "generate_flow",
		Description: "Generate a Flow class for multi-step user journeys",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"name", "steps"},
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the flow class",
				},
				"steps": map[string]interface{}{
					"type":        "array",
					"description": "List of step method names in the flow",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description of the flow",
				},
				"writeToFile": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to write the generated code to a file (default: false)",
					"default":     false,
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Base directory for output files (required if writeToFile is true)",
				},
			},
		},
	}, h.handleGenerateFlow)

	// Register generate_steps tool
	s.AddTool(mcp.Tool{
		Name:        "generate_steps",
		Description: "Generate Cucumber step definitions from Gherkin patterns",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"steps"},
			Properties: map[string]interface{}{
				"steps": map[string]interface{}{
					"type":        "array",
					"description": "List of step definitions",
					"items": map[string]interface{}{
						"type":     "object",
						"required": []string{"type", "pattern"},
						"properties": map[string]interface{}{
							"type": map[string]interface{}{
								"type":        "string",
								"description": "Step type (Given, When, Then)",
								"enum":        []string{"Given", "When", "Then"},
							},
							"pattern": map[string]interface{}{
								"type":        "string",
								"description": "Gherkin step pattern",
							},
							"parameters": map[string]interface{}{
								"type":        "string",
								"description": "Optional parameter types",
							},
						},
					},
				},
				"language": map[string]interface{}{
					"type":        "string",
					"description": "Target language (typescript or javascript)",
					"enum":        []string{"typescript", "javascript"},
					"default":     "typescript",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description",
				},
				"writeToFile": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to write the generated code to a file (default: false)",
					"default":     false,
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Base directory for output files (required if writeToFile is true)",
				},
			},
		},
	}, h.handleGenerateSteps)

	// Register generate_element tool
	s.AddTool(mcp.Tool{
		Name:        "generate_element",
		Description: "Generate a custom Element/Component class",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"name", "rootSelector"},
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the element component",
				},
				"rootSelector": map[string]interface{}{
					"type":        "string",
					"description": "Selector for the root element",
				},
				"methods": map[string]interface{}{
					"type":        "array",
					"description": "Optional list of method names",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description",
				},
				"writeToFile": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to write the generated code to a file (default: false)",
					"default":     false,
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Base directory for output files (required if writeToFile is true)",
				},
			},
		},
	}, h.handleGenerateElement)

	// Register generate_service tool
	s.AddTool(mcp.Tool{
		Name:        "generate_service",
		Description: "Generate an API Service class",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"name", "baseUrl", "endpoints"},
			Properties: map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the service",
				},
				"baseUrl": map[string]interface{}{
					"type":        "string",
					"description": "Base URL for the service",
				},
				"endpoints": map[string]interface{}{
					"type":        "array",
					"description": "List of API endpoints",
					"items": map[string]interface{}{
						"type": "object",
						"required": []string{"name", "method", "path"},
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "string",
							},
							"method": map[string]interface{}{
								"type": "string",
								"enum": []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
							},
							"path": map[string]interface{}{
								"type": "string",
							},
						},
					},
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Optional description",
				},
				"writeToFile": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether to write the generated code to a file (default: false)",
					"default":     false,
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Base directory for output files (required if writeToFile is true)",
				},
			},
		},
	}, h.handleGenerateService)

	// Register migrate_code tool
	s.AddTool(mcp.Tool{
		Name:        "migrate_code",
		Description: "Analyze and migrate legacy code to Saturday Framework patterns",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"sourceCode"},
			Properties: map[string]interface{}{
				"sourceCode": map[string]interface{}{
					"type":        "string",
					"description": "Legacy Playwright source code to migrate",
				},
				"type": map[string]interface{}{
					"type":        "string",
					"description": "Type of migration (page or test)",
					"enum":        []string{"page", "test"},
					"default":     "page",
				},
			},
		},
	}, h.handleMigrateCode)

	// Register analyze_performance tool
	s.AddTool(mcp.Tool{
		Name:        "analyze_performance",
		Description: "Analyze code for performance bottlenecks (concurrent)",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
			},
		},
	}, h.handleAnalyzePerformance)

	// Register generate_documentation tool
	s.AddTool(mcp.Tool{
		Name:        "generate_documentation",
		Description: "Generate markdown documentation for the project",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath", "outputPath"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path for the output markdown file",
				},
			},
		},
	}, h.handleGenerateDocumentation)

	// Register analyze_framework tool
	s.AddTool(mcp.Tool{
		Name:        "analyze_framework",
		Description: "Analyze existing framework structure and patterns",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
				"patterns": map[string]interface{}{
					"type":        "array",
					"description": "Optional list of patterns to look for",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}, h.handleAnalyzeFramework)

	// Register validate_patterns tool
	s.AddTool(mcp.Tool{
		Name:        "validate_patterns",
		Description: "Validate code against Saturday framework patterns",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
				"checkTypes": map[string]interface{}{
					"type":        "array",
					"description": "Optional kinds of checks to run (naming, inheritance, structure)",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}, h.handleValidatePatterns)

	// Register suggest_improvements tool
	s.AddTool(mcp.Tool{
		Name:        "suggest_improvements",
		Description: "Suggest code improvements based on Saturday framework best practices",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
			},
		},
	}, h.handleSuggestImprovements)

	// Register analyze_impact tool
	s.AddTool(mcp.Tool{
		Name:        "analyze_impact",
		Description: "Analyze the impact of modifying a specific file (dependency graph)",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath", "targetFile"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
				"targetFile": map[string]interface{}{
					"type":        "string",
					"description": "Relative path of the file to modify",
				},
			},
		},
	}, h.handleAnalyzeImpact)

	// Register run_tests tool
	s.AddTool(mcp.Tool{
		Name:        "run_tests",
		Description: "Execute tests and capture output",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"projectPath"},
			Properties: map[string]interface{}{
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to the project root",
				},
				"command": map[string]interface{}{
					"type":        "string",
					"description": "Test command (default: npx playwright test)",
				},
				"filter": map[string]interface{}{
					"type":        "string",
					"description": "Grep filter for tests",
				},
			},
		},
	}, h.handleRunTests)

	// Register parse_test_failure tool
	s.AddTool(mcp.Tool{
		Name:        "parse_test_failure",
		Description: "Parse standard output from Playwright tests to identify failing files and lines",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"output"},
			Properties: map[string]interface{}{
				"output": map[string]interface{}{
					"type":        "string",
					"description": "Raw output string from run_tests",
				},
			},
		},
	}, h.handleParseTestFailure)

	// Register prioritize_tests tool
	s.AddTool(mcp.Tool{
		Name:        "prioritize_tests",
		Description: "Rank test coverage needs based on production usage metrics",
		InputSchema: mcp.ToolInputSchema{
			Type:     "object",
			Required: []string{"metricsFile"},
			Properties: map[string]interface{}{
				"metricsFile": map[string]interface{}{
					"type":        "string",
					"description": "Absolute path to metrics.json file",
				},
				"projectPath": map[string]interface{}{
					"type":        "string",
					"description": "Optional path to project to match files",
				},
			},
		},
	}, h.handlePrioritizeTests)

	h.logger.Info("Tools registered successfully")
	return nil
}

// ... existing code ...

// handleRunTests executes tests
func (h *Handler) handleRunTests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling run_tests request")

	// Parse arguments
	args := request.Params.Arguments
	var req models.TestExecutionRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	if req.ProjectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	result, err := h.testExecutor.Run(req)
	if err != nil {
		h.logger.Error("Test execution failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Test execution failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Tests executed", "success", result.Success, "summary", result.Summary)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleParseTestFailure parses test logs
func (h *Handler) handleParseTestFailure(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling parse_test_failure request")

	output, ok := request.Params.Arguments["output"].(string)
	if !ok || output == "" {
		return mcp.NewToolResultError("output argument is required"), nil
	}

	failures, err := h.logAnalyzer.Parse(output)
	if err != nil {
		h.logger.Error("Failed to parse output", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Parse failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(failures)
	if err != nil {
		h.logger.Error("Failed to marshal failures", "error", err)
		return mcp.NewToolResultError("Failed to format failures"), nil
	}

	h.logger.Info("Parsed test failures", "count", len(failures))
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handlePrioritizeTests ranks testing needs
func (h *Handler) handlePrioritizeTests(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling prioritize_tests request")

	metricsFile, ok := request.Params.Arguments["metricsFile"].(string)
	if !ok || metricsFile == "" {
		return mcp.NewToolResultError("metricsFile argument is required"), nil
	}

	// 1. Load Metrics
	provider := observability.NewFileMetricsProvider(metricsFile)
	metrics, err := provider.GetPageMetrics()
	if err != nil {
		h.logger.Error("Failed to read metrics", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to load metrics: %v", err)), nil
	}

	// 2. Analyze
	priorities := h.usageAnalyzer.Analyze(metrics)

	// 3. (Optional) Match with Codebase
	// projectPath, hasProject := request.Params.Arguments["projectPath"].(string)
	// if hasProject && projectPath != "" {
	//    files, _ := listFiles(projectPath)
	//    priorities = h.usageAnalyzer.MatchWithCodebase(priorities, files)
	// }

	resultJSON, err := json.Marshal(priorities)
	if err != nil {
		h.logger.Error("Failed to marshal priorities", "error", err)
		return mcp.NewToolResultError("Failed to format priorities"), nil
	}

	h.logger.Info("Prioritized tests", "count", len(priorities))
	return mcp.NewToolResultText(string(resultJSON)), nil
}



// handleAnalyzeImpact analyzes file modification impact
func (h *Handler) handleAnalyzeImpact(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling analyze_impact request")

	// Parse arguments
	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)
	targetFile, _ := args["targetFile"].(string)

	if projectPath == "" || targetFile == "" {
		return mcp.NewToolResultError("projectPath and targetFile are required"), nil
	}

	// 1. Build Graph
	graph, err := h.graphAnalyzer.Build(projectPath)
	if err != nil {
		h.logger.Error("Graph build failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build dependency graph: %v", err)), nil
	}

	// 2. Analyze Impact
	impactedFiles, err := h.graphAnalyzer.AnalyzeImpact(graph, targetFile)
	if err != nil {
		h.logger.Error("Impact analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	// Format response
	result := map[string]interface{}{
		"target":   targetFile,
		"impacted": impactedFiles,
		"count":    len(impactedFiles),
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Impact analysis completed successfully", "target", targetFile, "impacted_count", len(impactedFiles))
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// RegisterResources registers all available resources with the MCP server
func (h *Handler) RegisterResources(s *server.MCPServer) error {
	h.logger.Info("Registering resources")

	for _, r := range h.resourceProvider.List() {
		s.AddResource(r, func(ctx context.Context, request mcp.ReadResourceRequest) ([]interface{}, error) {
			return h.resourceProvider.Read(request.Params.URI)
		})
	}

	h.logger.Info("Resources registered successfully")
	return nil
}


// RegisterPrompts registers all available prompts with the MCP server
func (h *Handler) RegisterPrompts(s *server.MCPServer) error {
	h.logger.Info("Registering prompts")

	for _, p := range h.promptProvider.List() {
		// Capture p for closure
		pName := p.Name
		s.AddPrompt(p, func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			messages, err := h.promptProvider.Get(pName, req.Params.Arguments)
			if err != nil {
				return nil, err
			}
			return &mcp.GetPromptResult{
				Messages: messages,
			}, nil
		})
	}

	h.logger.Info("Prompts registered successfully")
	return nil
}

// handleListTools returns a list of all available tools
func (h *Handler) handleListTools(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling list_tools request")

	tools := []map[string]string{
		{
			"name":        "list_tools",
			"description": "List all available Saturday framework tools",
			"status":      "implemented",
		},
		{
			"name":        "generate_site",
			"description": "Generate a new Site class with page and flow registration",
			"status":      "implemented",
		},
		{
			"name":        "generate_page",
			"description": "Generate a new Page class with element registration",
			"status":      "implemented",
		},
		{
			"name":        "generate_flow",
			"description": "Generate a new Flow class for multi-step user journeys",
			"status":      "implemented",
		},
		{
			"name":        "generate_steps",
			"description": "Generate Cucumber step definitions from Gherkin",
			"status":      "implemented",
		},
		{
			"name":        "analyze_framework",
			"description": "Analyze existing framework structure and patterns",
			"status":      "implemented",
		},
		{
			"name":        "validate_patterns",
			"description": "Validate code against Saturday framework patterns",
			"status":      "implemented",
		},
	}

	toolsJSON, err := json.Marshal(tools)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tools: %w", err)
	}

	return mcp.NewToolResultText(string(toolsJSON)), nil
}

// handleGenerateSite generates a Site class
func (h *Handler) handleGenerateSite(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_site request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Extract file writing parameters
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)
	
	// Parse site generation request
	var req models.SiteGenerationRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	// Generate site
	resp, err := h.generator.SiteGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Site generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file if requested
	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		
		// Determine output directory (e.g., lib/sites/)
		relativePath := filepath.Join("lib", "sites", resp.FileName)
		
		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			h.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		h.logger.Info("File written successfully", "path", fullPath)
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Site generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGeneratePage generates a Page class
func (h *Handler) handleGeneratePage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_page request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Extract file writing parameters
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)
	
	// Parse page generation request
	var req models.PageGenerationRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	// Generate page
	resp, err := h.generator.PageGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Page generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file if requested
	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		
		// Determine output directory (e.g., lib/pages/)
		relativePath := filepath.Join("lib", "pages", resp.FileName)
		
		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			h.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		h.logger.Info("File written successfully", "path", fullPath)
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Page generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGenerateFlow generates a Flow class
func (h *Handler) handleGenerateFlow(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_flow request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Extract file writing parameters
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)
	
	// Parse flow generation request
	var req models.FlowGenerationRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	// Generate flow
	resp, err := h.generator.FlowGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Flow generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file if requested
	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		
		// Determine output directory (e.g., lib/flows/)
		relativePath := filepath.Join("lib", "flows", resp.FileName)
		
		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			h.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		h.logger.Info("File written successfully", "path", fullPath)
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Flow generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGenerateSteps generates Cucumber step definitions
func (h *Handler) handleGenerateSteps(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_steps request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Extract file writing parameters
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)
	
	// Parse step definition request
	var req models.StepDefinitionRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	// Generate steps
	resp, err := h.generator.StepGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Step generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file if requested
	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		
		// Determine output directory (e.g., tests/steps/)
		relativePath := filepath.Join("tests", "steps", resp.FileName)
		
		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			h.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		h.logger.Info("File written successfully", "path", fullPath)
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Steps generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGenerateElement generates an Element class
func (h *Handler) handleGenerateElement(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_element request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Extract file writing parameters
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)
	
	// Parse generation request
	var req models.ElementGenerationRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	// Generate code
	resp, err := h.generator.ElementGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Element generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file if requested
	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		
		// Determine output directory (e.g., lib/components/)
		relativePath := filepath.Join("lib", "components", resp.FileName)
		
		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			h.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		h.logger.Info("File written successfully", "path", fullPath)
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Element generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGenerateService generates a Service class
func (h *Handler) handleGenerateService(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_service request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Extract file writing parameters
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)
	
	// Parse generation request
	var req models.ServiceGenerationRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		h.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		h.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	// Generate code
	resp, err := h.generator.ServiceGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Service generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file if requested
	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		
		// Determine output directory (e.g., lib/services/)
		relativePath := filepath.Join("lib", "services", resp.FileName)
		
		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			h.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		h.logger.Info("File written successfully", "path", fullPath)
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Service generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleMigrateCode migrates legacy code
func (h *Handler) handleMigrateCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling migrate_code request")

	// Parse arguments
	args := request.Params.Arguments
	
	// Parse request manually since sourceCode might be complex
	sourceCode, _ := args["sourceCode"].(string)
	migrationType, ok := args["type"].(string)
	if !ok || migrationType == "" {
		migrationType = "page"
	}

	req := models.MigrationRequest{
		SourceCode: sourceCode,
		Type:       migrationType,
	}

	// Generate migrated code
	resp, err := h.generator.MigrationGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Migration failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Migration failed: %v", err)), nil
	}

	// Format response
	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Code migrated successfully")
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleAnalyzeFramework analyzes existing framework structure
func (h *Handler) handleAnalyzeFramework(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling analyze_framework request")

	// Parse arguments
	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	// Run analysis
	result, err := h.analyzer.Analyze(projectPath)
	if err != nil {
		h.logger.Error("Analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Analysis completed successfully", "path", projectPath)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleValidatePatterns validates code patterns
func (h *Handler) handleValidatePatterns(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling validate_patterns request")

	// Parse arguments
	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	// Validate patterns
	result, err := h.validator.Validate(projectPath)
	if err != nil {
		h.logger.Error("Validation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Validation failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Validation completed successfully", "path", projectPath, "valid", result.Valid)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleSuggestImprovements suggests code improvements
func (h *Handler) handleSuggestImprovements(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling suggest_improvements request")

	// Parse arguments
	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	// Analyze for improvements
	result, err := h.improvementAnalyzer.Analyze(projectPath)
	if err != nil {
		h.logger.Error("Improvement analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Improvement suggestions generated successfully", "path", projectPath, "count", len(result.Suggestions))
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleAnalyzePerformance analyzes performance
func (h *Handler) handleAnalyzePerformance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling analyze_performance request")

	// Parse arguments
	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)

	if projectPath == "" {
		return mcp.NewToolResultError("projectPath is required"), nil
	}

	// Analyze for performance
	result, err := h.performanceAnalyzer.Analyze(projectPath)
	if err != nil {
		h.logger.Error("Performance analysis failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Analysis failed: %v", err)), nil
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		h.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	h.logger.Info("Performance analysis generated successfully", "path", projectPath, "issues", len(result.Suggestions))
	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGenerateDocumentation generates project documentation
func (h *Handler) handleGenerateDocumentation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.logger.Info("Handling generate_documentation request")

	// Parse arguments
	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)
	outputPath, _ := args["outputPath"].(string)

	if projectPath == "" || outputPath == "" {
		return mcp.NewToolResultError("projectPath and outputPath are required"), nil
	}

	req := models.DocumentationRequest{
		ProjectPath: projectPath,
		OutputPath:  outputPath,
	}

	// Generate documentation
	resp, err := h.generator.DocumentationGenerator.Generate(req)
	if err != nil {
		h.logger.Error("Documentation generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	// Write to file
	if err := os.WriteFile(outputPath, []byte(resp.Code), 0644); err != nil {
		h.logger.Error("Failed to write documentation file", "path", outputPath, "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	h.logger.Info("Documentation generated successfully", "path", outputPath)
	
	// Return success
	result := map[string]interface{}{
		"success": true,
		"path":    outputPath,
		"pages":   resp.Metadata["pageCount"],
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

