package server

import (
	"context"
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
//
// tracer is the observability seam RegisterTools uses to wrap every
// tool invocation in a span. It's a domain.Tracer (interface), never a
// concrete OTel type — the concrete implementation is chosen in
// cmd/saturday-mcp/main.go and injected here at construction time.
type Handler struct {
	logger           *logging.Logger
	tools            []domain.Tool
	byName           map[string]domain.Tool
	tracer           domain.Tracer
	resourceProvider *resources.Provider
	promptProvider   *prompts.Provider
}

// Option configures a Handler at construction time. Functional-option
// pattern chosen over a wide constructor signature so future collaborators
// (e.g. custom persona provider, resource provider) can be layered in
// without churning every NewHandler call site.
type Option func(*Handler)

// WithTracer injects a domain.Tracer for span emission. When omitted,
// Handler defaults to an internal no-op tracer so tests and dev runs
// with no collector configured pay zero span-emission cost.
// cmd/saturday-mcp/main.go passes WithTracer(otel.NewOTelTracer(...))
// when OTEL_EXPORTER_OTLP_ENDPOINT is set.
func WithTracer(t domain.Tracer) Option {
	return func(h *Handler) { h.tracer = t }
}

// NewHandler wires the template system, the tool provider, and the
// resource/prompt providers, then returns the composite Handler. Pass
// WithTracer to inject a real domain.Tracer; otherwise every tool
// invocation goes through a private noopTracer that discards spans.
func NewHandler(logger *logging.Logger, opts ...Option) (*Handler, error) {
	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	toolList := buildTools(logger, processor)

	h := &Handler{
		logger:           logger,
		tools:            toolList,
		byName:           indexByName(toolList),
		tracer:           noopTracer{},
		resourceProvider: resources.NewProvider(logger, loader),
		promptProvider:   prompts.NewProvider(logger),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h, nil
}

// noopTracer is the private default so Handler never carries a nil
// tracer. Kept inside the server package to avoid an inward dependency
// on internal/adapters/otel — server is inner-layer, adapters are
// outer-layer, and Clean Architecture wants that arrow only pointing
// outward. The public adapters/otel.NoopTracer covers the "explicitly
// pass a noop" case main.go uses; this one covers the "caller passed
// nothing" case tests use.
type noopTracer struct{}

func (noopTracer) StartSpan(ctx context.Context, _ string, _ ...domain.SpanAttribute) (context.Context, domain.EndSpan) {
	return ctx, func(error, ...domain.SpanAttribute) {}
}
