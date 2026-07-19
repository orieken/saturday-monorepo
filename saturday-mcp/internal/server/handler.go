package server

import (
	"fmt"
	"time"

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





