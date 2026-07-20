package domain

import "context"

// SpanAttribute is a single key/value pair to tag on a span. Value is
// intentionally untyped — the adapter maps it to whatever the underlying
// tracer (OTel, noop, in-memory test double) expects. The domain layer
// takes no position on OTel semantic conventions; adapters do.
type SpanAttribute struct {
	Key   string
	Value any
}

// String builds a string-typed SpanAttribute. Convenience over exposing
// the struct literal — tool/workflow call sites read cleaner with
// `domain.String("tool.name", t.Name())` than a struct literal.
func String(key, value string) SpanAttribute {
	return SpanAttribute{Key: key, Value: value}
}

// Int builds an int-typed SpanAttribute.
func Int(key string, value int) SpanAttribute {
	return SpanAttribute{Key: key, Value: value}
}

// Int64 builds an int64-typed SpanAttribute — mostly for duration_ms
// values pulled from time.Since().Milliseconds().
func Int64(key string, value int64) SpanAttribute {
	return SpanAttribute{Key: key, Value: value}
}

// Bool builds a bool-typed SpanAttribute — mostly for success flags.
func Bool(key string, value bool) SpanAttribute {
	return SpanAttribute{Key: key, Value: value}
}

// EndSpan closes a span opened by Tracer.StartSpan. Callers use it via
// `ctx, end := tracer.StartSpan(...); defer end(err)` — passing the
// terminal error (nil on success) lets the adapter record it on the
// span before finishing. Extra attributes discovered during the span's
// life (tool.success, tool.error_class, per-step outcome) can be
// attached at close time via the variadic attrs. Idempotent: calling
// twice is a no-op.
type EndSpan func(err error, extraAttrs ...SpanAttribute)

// Tracer is the domain seam for emitting spans. All span emission goes
// through this interface; no code outside internal/adapters/otel/ ever
// imports go.opentelemetry.io/otel directly. This is what
// architecture-guardrails #8 (Observability Boundaries) requires — OTel
// is infrastructure, not a domain concern.
//
// Contract:
//   - StartSpan returns a derived context carrying the span; downstream
//     calls that receive that context can be linked into the trace.
//   - The returned EndSpan MUST be called exactly once, typically via
//     defer. Pass the terminal error (nil on success) so adapters can
//     record it as a span status.
//   - Attributes provided at StartSpan attach at span creation; adapters
//     MAY support attribute additions after the fact, but the domain
//     interface does not expose that surface.
//
// The noop implementation (internal/adapters/otel/noop_tracer.go) is the
// default when OTEL_EXPORTER_OTLP_ENDPOINT is unset, so dev runs and
// tests never require a running collector.
type Tracer interface {
	StartSpan(ctx context.Context, name string, attrs ...SpanAttribute) (context.Context, EndSpan)
}
