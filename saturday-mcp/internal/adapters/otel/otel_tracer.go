package otel

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/orieken/saturday-mcp/internal/domain"
)

// tracerName is the instrumentation library name registered with OTel.
// Kept stable across releases so backend queries can filter on it.
const tracerName = "github.com/orieken/saturday-mcp"

// OTelTracer is the OpenTelemetry-backed implementation of domain.Tracer.
// The only place in the codebase that imports go.opentelemetry.io/otel.
// Callers depend on domain.Tracer; the concrete type is injected only
// from cmd/saturday-mcp/main.go once, at startup.
type OTelTracer struct {
	tracer   oteltrace.Tracer
	shutdown func(context.Context) error
}

// Config carries the knobs main.go passes when building a real tracer.
// Endpoint mirrors the OTEL_EXPORTER_OTLP_ENDPOINT env var (already the
// OTel-conventional name), ServiceName tags every span with the source
// service, and Insecure toggles TLS for local collectors.
type Config struct {
	Endpoint    string
	ServiceName string
	Insecure    bool
}

// NewOTelTracer builds a tracer that ships spans to an OTLP gRPC
// endpoint. Chosen over OTLP HTTP because collectors in the Saturday
// stack default to gRPC and it's the more common production path; the
// noop implementation covers the "no collector configured" case so this
// path is only exercised when OTEL_EXPORTER_OTLP_ENDPOINT is set.
//
// Returns the tracer plus a shutdown func the caller MUST defer at
// process exit so the span exporter flushes queued batches.
func NewOTelTracer(ctx context.Context, cfg Config) (*OTelTracer, error) {
	exp, err := buildExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("build OTLP exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(cfg.ServiceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("build OTel resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return &OTelTracer{
		tracer:   tp.Tracer(tracerName),
		shutdown: tp.Shutdown,
	}, nil
}

// buildExporter isolates the gRPC-exporter construction so the ctor
// stays readable. Insecure controls plaintext transport for local
// collectors (Docker Compose, dev clusters).
func buildExporter(ctx context.Context, cfg Config) (*otlptrace.Exporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	return otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
}

// Shutdown flushes and closes the underlying provider. Safe to call with
// a bounded ctx (5s or so at process exit is typical).
func (t *OTelTracer) Shutdown(ctx context.Context) error {
	if t.shutdown == nil {
		return nil
	}
	return t.shutdown(ctx)
}

// StartSpan opens an OTel span and returns the context + an EndSpan
// closure that records the terminal error (if any) as span status before
// finishing. Attribute translation is confined to spanAttributes below —
// the domain layer's untyped Value gets mapped to OTel attribute.KeyValue
// here and nowhere else.
func (t *OTelTracer) StartSpan(ctx context.Context, name string, attrs ...domain.SpanAttribute) (context.Context, domain.EndSpan) {
	ctx, span := t.tracer.Start(ctx, name, oteltrace.WithAttributes(spanAttributes(attrs)...))
	start := time.Now()
	return ctx, func(err error) {
		span.SetAttributes(attribute.Int64("duration_ms", time.Since(start).Milliseconds()))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	}
}

// spanAttributes maps the domain-layer untyped attribute values into
// OTel attribute.KeyValue instances. Unknown types fall through as
// fmt.Sprintf("%v") strings so an unhandled Value never drops the span
// silently — the attribute just shows up as a stringified fallback.
func spanAttributes(attrs []domain.SpanAttribute) []attribute.KeyValue {
	out := make([]attribute.KeyValue, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, mapAttribute(a))
	}
	return out
}

func mapAttribute(a domain.SpanAttribute) attribute.KeyValue {
	switch v := a.Value.(type) {
	case string:
		return attribute.String(a.Key, v)
	case int:
		return attribute.Int(a.Key, v)
	case int64:
		return attribute.Int64(a.Key, v)
	case bool:
		return attribute.Bool(a.Key, v)
	default:
		return attribute.String(a.Key, fmt.Sprintf("%v", v))
	}
}
