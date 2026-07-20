// Package otel contains the OpenTelemetry-backed adapter that implements
// domain.Tracer, plus a noop implementation used by default and by tests.
// It is the only package in the codebase that may import
// go.opentelemetry.io/otel — architecture-guardrails #8 (Observability
// Boundaries) forbids OTel imports from domain, tools, workflows, or
// generators.
package otel

import (
	"context"

	"github.com/orieken/saturday-mcp/internal/domain"
)

// NoopTracer is the do-nothing implementation of domain.Tracer. It is
// the default whenever OTEL_EXPORTER_OTLP_ENDPOINT is unset — dev runs,
// unit tests, and installs without a collector all end up here and never
// pay the cost of the OTel SDK.
type NoopTracer struct{}

// NewNoopTracer returns the singleton-ish noop tracer. Constructor exists
// for symmetry with NewOTelTracer and so callers can inject it via the
// same code path that would inject a real one.
func NewNoopTracer() *NoopTracer {
	return &NoopTracer{}
}

// StartSpan returns the ctx unchanged and an EndSpan that does nothing.
// Callers still defer end(err) unconditionally; the shape matches the
// real tracer so switching between them is invisible to call sites.
func (t *NoopTracer) StartSpan(ctx context.Context, _ string, _ ...domain.SpanAttribute) (context.Context, domain.EndSpan) {
	return ctx, func(error) {}
}
