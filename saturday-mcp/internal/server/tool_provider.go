package server

import (
	"github.com/orieken/saturday-mcp/internal/adapters/filesystem"
	"github.com/orieken/saturday-mcp/internal/adapters/metricsfile"
	"github.com/orieken/saturday-mcp/internal/adapters/testrunner"
	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/domain"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/tools"
	"github.com/orieken/saturday-mcp/internal/validators"
	"github.com/orieken/saturday-mcp/internal/workflows"
)

// This file owns Tool/Workflow construction. Handler stopped needing
// per-tool collaborator fields (analyzer, generator, testExecutor, …)
// as soon as Phase C+D extracted every handleFoo body out to
// internal/tools/ and internal/workflows/ — the collaborators exist only
// to feed constructors. Isolating that construction here lets Handler
// shrink to a thin composite (see mcp-add-plan Phase I op 22).

// buildTools wires up every domain.Tool the server exposes, in the order
// they should appear in the MCP tool-list advertisement. Adding a new
// tool means: write it under internal/tools/, then append it to this
// slice — nothing else in this package changes.
func buildTools(logger *logging.Logger, processor *templates.Processor) []domain.Tool {
	validator := validators.NewValidator()

	siteGen := generators.NewSiteGenerator(processor, validator)
	pageGen := generators.NewPageGenerator(processor, validator)
	flowGen := generators.NewFlowGenerator(processor, validator)
	stepGen := generators.NewStepGenerator(processor, validator)
	elementGen := generators.NewElementGenerator(processor, validator)
	serviceGen := generators.NewServiceGenerator(processor, validator)
	migrationGen := generators.NewMigrationGenerator(processor, validator)
	docGen := generators.NewDocumentationGenerator(processor, validator)

	complexityAnalyzer := analyzers.NewComplexityAnalyzer()
	accessibilityAnalyzer := analyzers.NewAccessibilityAnalyzer()
	frameworkAnalyzer := analyzers.NewFrameworkAnalyzer(logger)
	patternValidator := analyzers.NewPatternValidator(logger)
	improvementAnalyzer := analyzers.NewImprovementAnalyzer(logger)
	performanceAnalyzer := analyzers.NewPerformanceAnalyzer(logger)
	graphAnalyzer := analyzers.NewGraphAnalyzer(logger)
	logAnalyzer := analyzers.NewTestLogAnalyzer(logger)
	usageAnalyzer := analyzers.NewUsageAnalyzer(logger)
	testRunner := testrunner.NewExecRunner(logger)
	fs := filesystem.NewOSFileSystem()
	metricsReader := metricsfile.NewFileReader()

	return []domain.Tool{
		tools.NewGenerateSiteTool(logger, siteGen),
		tools.NewGeneratePageTool(logger, pageGen),
		tools.NewGenerateFlowTool(logger, flowGen),
		tools.NewGenerateStepsTool(logger, stepGen),
		tools.NewGenerateElementTool(logger, elementGen),
		tools.NewGenerateServiceTool(logger, serviceGen),
		tools.NewMigrateCodeTool(logger, migrationGen),
		tools.NewAnalyzePerformanceTool(logger, performanceAnalyzer),
		tools.NewGenerateDocumentationTool(logger, docGen, fs),
		tools.NewAnalyzeFrameworkTool(logger, frameworkAnalyzer),
		tools.NewValidatePatternsTool(logger, patternValidator),
		tools.NewSuggestImprovementsTool(logger, improvementAnalyzer),
		tools.NewAnalyzeImpactTool(logger, graphAnalyzer),
		tools.NewWorkflowTool(workflows.NewRunTestsWorkflow(logger, testRunner)),
		tools.NewParseTestFailureTool(logger, logAnalyzer),
		tools.NewWorkflowTool(workflows.NewPrioritizeTestsWorkflow(logger, usageAnalyzer, metricsReader)),
		tools.NewAnalyzeComplexityTool(logger, complexityAnalyzer),
		tools.NewCheckAccessibilityTool(logger, accessibilityAnalyzer),
	}
}

// indexByName builds the lookup map testing.go uses to route the exported
// HandleFoo wrappers to their tools by MCP name.
func indexByName(list []domain.Tool) map[string]domain.Tool {
	byName := make(map[string]domain.Tool, len(list))
	for _, t := range list {
		byName[t.Name()] = t
	}
	return byName
}
