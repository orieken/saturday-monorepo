package server

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/executor"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/prompts"
	"github.com/orieken/saturday-mcp/internal/resources"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/tools"
	"github.com/orieken/saturday-mcp/internal/validators"
	"github.com/orieken/saturday-mcp/internal/workflows"
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
	analyzeImpactTool         *tools.AnalyzeImpactTool
	parseTestFailureTool      *tools.ParseTestFailureTool

	// Extracted workflows adapted to domain.Tool via tools.WorkflowTool
	// (Phase D op 3). The Workflow types themselves live under
	// internal/workflows/; the *WorkflowTool wrapper here exists purely
	// so RegisterTools can treat them identically to single-step tools.
	runTestsTool        *tools.WorkflowTool
	prioritizeTestsTool *tools.WorkflowTool
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
		analyzeImpactTool:         tools.NewAnalyzeImpactTool(logger, graphAnalyzer),
		parseTestFailureTool:      tools.NewParseTestFailureTool(logger, logAnalyzer),

		runTestsTool:        tools.NewWorkflowTool(workflows.NewRunTestsWorkflow(logger, testExecutor)),
		prioritizeTestsTool: tools.NewWorkflowTool(workflows.NewPrioritizeTestsWorkflow(logger, usageAnalyzer)),
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

	// Register analyze_impact tool (extracted — Phase C op 18)
	s.AddTool(mcp.Tool{
		Name:        h.analyzeImpactTool.Name(),
		Description: h.analyzeImpactTool.Description(),
		InputSchema: h.analyzeImpactTool.InputSchema(),
	}, h.analyzeImpactTool.Execute)

	// Register run_tests workflow (extracted — Phase D op 3, adapted via
	// tools.WorkflowTool so registration is identical to a single-step tool)
	s.AddTool(mcp.Tool{
		Name:        h.runTestsTool.Name(),
		Description: h.runTestsTool.Description(),
		InputSchema: h.runTestsTool.InputSchema(),
	}, h.runTestsTool.Execute)

	// Register parse_test_failure tool (extracted — Phase C op 19)
	s.AddTool(mcp.Tool{
		Name:        h.parseTestFailureTool.Name(),
		Description: h.parseTestFailureTool.Description(),
		InputSchema: h.parseTestFailureTool.InputSchema(),
	}, h.parseTestFailureTool.Execute)

	// Register prioritize_tests workflow (extracted — Phase D op 3)
	s.AddTool(mcp.Tool{
		Name:        h.prioritizeTestsTool.Name(),
		Description: h.prioritizeTestsTool.Description(),
		InputSchema: h.prioritizeTestsTool.InputSchema(),
	}, h.prioritizeTestsTool.Execute)

	h.logger.Info("Tools registered successfully")
	return nil
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



