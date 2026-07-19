package server

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/executor"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
	"github.com/orieken/saturday-mcp/internal/observability"
	"github.com/orieken/saturday-mcp/internal/prompts"
	"github.com/orieken/saturday-mcp/internal/resources"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/tools"
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

	// Extracted tools (Phase C of mcp-add-plan). Each field is a *tools.FooTool
	// implementing domain.Tool. Handler currently composes them; Phase I moves
	// registration + composition into a dedicated provider.
	generateSiteTool    *tools.GenerateSiteTool
	generatePageTool    *tools.GeneratePageTool
	generateFlowTool    *tools.GenerateFlowTool
	generateStepsTool   *tools.GenerateStepsTool
	generateElementTool *tools.GenerateElementTool
	generateServiceTool       *tools.GenerateServiceTool
	migrateCodeTool           *tools.MigrateCodeTool
	analyzePerformanceTool    *tools.AnalyzePerformanceTool
	generateDocumentationTool *tools.GenerateDocumentationTool
	analyzeFrameworkTool      *tools.AnalyzeFrameworkTool
	validatePatternsTool      *tools.ValidatePatternsTool
	suggestImprovementsTool   *tools.SuggestImprovementsTool
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

		generateSiteTool:    tools.NewGenerateSiteTool(logger, siteGen),
		generatePageTool:    tools.NewGeneratePageTool(logger, pageGen),
		generateFlowTool:    tools.NewGenerateFlowTool(logger, flowGen),
		generateStepsTool:   tools.NewGenerateStepsTool(logger, stepGen),
		generateElementTool: tools.NewGenerateElementTool(logger, elementGen),
		generateServiceTool:       tools.NewGenerateServiceTool(logger, serviceGen),
		migrateCodeTool:           tools.NewMigrateCodeTool(logger, migrationGen),
		analyzePerformanceTool:    tools.NewAnalyzePerformanceTool(logger, performanceAnalyzer),
		generateDocumentationTool: tools.NewGenerateDocumentationTool(logger, docGen),
		analyzeFrameworkTool:      tools.NewAnalyzeFrameworkTool(logger, analyzer),
		validatePatternsTool:      tools.NewValidatePatternsTool(logger, validatorTool),
		suggestImprovementsTool:   tools.NewSuggestImprovementsTool(logger, improvementAnalyzer),
	}, nil
}





// RegisterTools registers all available tools with the MCP server
func (h *Handler) RegisterTools(s *server.MCPServer) error {
	h.logger.Info("Registering tools")

	// Register generate_site tool (extracted — Phase C op 6)
	s.AddTool(mcp.Tool{
		Name:        h.generateSiteTool.Name(),
		Description: h.generateSiteTool.Description(),
		InputSchema: h.generateSiteTool.InputSchema(),
	}, h.generateSiteTool.Execute)

	// Register generate_page tool (extracted — Phase C op 7)
	s.AddTool(mcp.Tool{
		Name:        h.generatePageTool.Name(),
		Description: h.generatePageTool.Description(),
		InputSchema: h.generatePageTool.InputSchema(),
	}, h.generatePageTool.Execute)

	// Register generate_flow tool (extracted — Phase C op 8)
	s.AddTool(mcp.Tool{
		Name:        h.generateFlowTool.Name(),
		Description: h.generateFlowTool.Description(),
		InputSchema: h.generateFlowTool.InputSchema(),
	}, h.generateFlowTool.Execute)

	// Register generate_steps tool (extracted — Phase C op 9)
	s.AddTool(mcp.Tool{
		Name:        h.generateStepsTool.Name(),
		Description: h.generateStepsTool.Description(),
		InputSchema: h.generateStepsTool.InputSchema(),
	}, h.generateStepsTool.Execute)

	// Register generate_element tool (extracted — Phase C op 10)
	s.AddTool(mcp.Tool{
		Name:        h.generateElementTool.Name(),
		Description: h.generateElementTool.Description(),
		InputSchema: h.generateElementTool.InputSchema(),
	}, h.generateElementTool.Execute)

	// Register generate_service tool (extracted — Phase C op 11)
	s.AddTool(mcp.Tool{
		Name:        h.generateServiceTool.Name(),
		Description: h.generateServiceTool.Description(),
		InputSchema: h.generateServiceTool.InputSchema(),
	}, h.generateServiceTool.Execute)

	// Register migrate_code tool (extracted — Phase C op 12)
	s.AddTool(mcp.Tool{
		Name:        h.migrateCodeTool.Name(),
		Description: h.migrateCodeTool.Description(),
		InputSchema: h.migrateCodeTool.InputSchema(),
	}, h.migrateCodeTool.Execute)

	// Register analyze_performance tool (extracted — Phase C op 13)
	s.AddTool(mcp.Tool{
		Name:        h.analyzePerformanceTool.Name(),
		Description: h.analyzePerformanceTool.Description(),
		InputSchema: h.analyzePerformanceTool.InputSchema(),
	}, h.analyzePerformanceTool.Execute)

	// Register generate_documentation tool (extracted — Phase C op 14)
	s.AddTool(mcp.Tool{
		Name:        h.generateDocumentationTool.Name(),
		Description: h.generateDocumentationTool.Description(),
		InputSchema: h.generateDocumentationTool.InputSchema(),
	}, h.generateDocumentationTool.Execute)

	// Register analyze_framework tool (extracted — Phase C op 15)
	s.AddTool(mcp.Tool{
		Name:        h.analyzeFrameworkTool.Name(),
		Description: h.analyzeFrameworkTool.Description(),
		InputSchema: h.analyzeFrameworkTool.InputSchema(),
	}, h.analyzeFrameworkTool.Execute)

	// Register validate_patterns tool (extracted — Phase C op 16)
	s.AddTool(mcp.Tool{
		Name:        h.validatePatternsTool.Name(),
		Description: h.validatePatternsTool.Description(),
		InputSchema: h.validatePatternsTool.InputSchema(),
	}, h.validatePatternsTool.Execute)

	// Register suggest_improvements tool (extracted — Phase C op 17)
	s.AddTool(mcp.Tool{
		Name:        h.suggestImprovementsTool.Name(),
		Description: h.suggestImprovementsTool.Description(),
		InputSchema: h.suggestImprovementsTool.InputSchema(),
	}, h.suggestImprovementsTool.Execute)

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



