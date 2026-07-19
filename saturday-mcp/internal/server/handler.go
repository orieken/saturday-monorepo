package server

import (
	"fmt"
	"time"

	"github.com/orieken/saturday-mcp/internal/domain"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/prompts"
	"github.com/orieken/saturday-mcp/internal/resources"
	"github.com/orieken/saturday-mcp/internal/templates"
)

// Handler is the composite that wires the three MCP provider slices
// (tools, resources, prompts) to a single MCP server. Construction of
// individual tools/workflows lives in tool_provider.go; MCP registration
// lives in registration.go; testing.go exposes named public wrappers so
// internal/integration/e2e_test.go can invoke tools directly. This
// struct is intentionally shallow — every collaborator that fed a tool
// constructor (analyzer, generator, testExecutor, …) lives inside the
// closures buildTools created, not on Handler.
type Handler struct {
	logger           *logging.Logger
	tools            []domain.Tool
	byName           map[string]domain.Tool
	resourceProvider *resources.Provider
	promptProvider   *prompts.Provider
}

// NewHandler wires the template system, the tool provider, and the
// resource/prompt providers, then returns the composite Handler.
func NewHandler(logger *logging.Logger) (*Handler, error) {
	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	toolList := buildTools(logger, processor)

	return &Handler{
		logger:           logger,
		tools:            toolList,
		byName:           indexByName(toolList),
		resourceProvider: resources.NewProvider(logger, loader),
		promptProvider:   prompts.NewProvider(logger),
	}, nil
}
