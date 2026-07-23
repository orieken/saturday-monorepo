package server

import (
	"os"
	"path/filepath"

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
	ubiquitousLanguageAnalyzer := analyzers.NewUbiquitousLanguageAnalyzer()
	dependencyBoundaryAnalyzer := analyzers.NewDependencyBoundaryAnalyzer()
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
		tools.NewCheckUbiquitousLanguageTool(logger, ubiquitousLanguageAnalyzer),
		tools.NewVerifyDependenciesTool(logger, dependencyBoundaryAnalyzer),
		tools.NewSearchKITool(logger, tools.NewKICorpusRetriever(frameworkCorpusPaths())),
		newSearchDocsTool(logger),
	}
}

// newSearchDocsTool constructs the search_docs tool and its BM25
// backend. Split out from buildTools so the failure branch (sqlite
// setup failed on this install — permission denied, disk full, etc.)
// stays readable: the tool ALWAYS registers, so the MCP tool inventory
// stays stable across environments; a broken backend surfaces as a
// per-call "no docs retriever configured" diagnostic rather than a
// startup crash. Follow-up: retriever.Close() has no shutdown hook
// today — process exit closes the sqlite handle (see modernc.org/sqlite
// durability), and adding a proper Handler.Shutdown() is out of scope
// for M1 Op 7b.
func newSearchDocsTool(logger *logging.Logger) domain.Tool {
	dbPath, ok := docsFTSDBPath()
	if !ok {
		logger.Info("BM25 retriever not initialised — no installed-project marker in CWD; search_docs will return no-corpus diagnostic responses")
		return tools.NewSearchDocsTool(logger, nil, nil)
	}
	retriever, err := tools.NewBM25Retriever(dbPath)
	if err != nil {
		logger.Warn("BM25 retriever init failed — search_docs will return no-corpus diagnostic responses",
			"dbPath", dbPath, "error", err)
		return tools.NewSearchDocsTool(logger, nil, nil)
	}
	return tools.NewSearchDocsTool(logger, retriever, retriever)
}

// docsFTSDBPath resolves the sqlite-fts5 index location for the docs
// corpus. Precedence:
//  1. SATURDAY_MCP_DOCS_FTS_PATH env var — an explicit override; always
//     returned when set (so a caller who's opted in can put the DB
//     anywhere, including under a test's t.TempDir()).
//  2. Otherwise, .claude/rag/docs-fts5.sqlite relative to CWD — but only
//     when a .claude/ directory already exists at CWD (the marker that
//     this process is running inside an installed project). Without the
//     marker, returns ("", false) so newSearchDocsTool skips backend
//     construction — this keeps unit tests from spuriously creating
//     .claude/rag/ under every test-package directory.
//  3. If CWD is unresolvable at all, returns ("", false) — the tool
//     degrades to the no-corpus diagnostic path.
func docsFTSDBPath() (string, bool) {
	if override := os.Getenv("SATURDAY_MCP_DOCS_FTS_PATH"); override != "" {
		return override, true
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	if _, err := os.Stat(filepath.Join(cwd, ".claude")); err != nil {
		return "", false
	}
	return filepath.Join(cwd, ".claude", "rag", "docs-fts5.sqlite"), true
}

// frameworkCorpusPaths resolves the shared/knowledge/ and docs/adrs/
// roots of the ai-assistant-dot-files install so search_ki can walk
// them. Reads AI_ASSISTANT_DOTFILES_PATH env; returns nil (empty corpus)
// when unset — the tool then reports zero hits with a diagnostic note
// rather than failing, matching mcp-expand-plan §3 Op 6's stance that
// LLM-as-retriever must degrade gracefully when the corpus isn't
// configured. Follow-up: when AOS Phase 3 lands shared/rag/, discovery
// moves there and this helper becomes a thin adapter over it.
func frameworkCorpusPaths() []string {
	root := os.Getenv("AI_ASSISTANT_DOTFILES_PATH")
	if root == "" {
		return nil
	}
	return []string{
		filepath.Join(root, "shared", "knowledge"),
		filepath.Join(root, "docs", "adrs"),
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
