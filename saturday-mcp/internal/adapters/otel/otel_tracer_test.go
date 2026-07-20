package otel

import (
	"context"
	"errors"
	"testing"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/orieken/saturday-mcp/internal/domain"
)

// newTestTracer builds an OTelTracer wired to an in-memory SpanRecorder
// so tests observe span emission without a running collector. Bypasses
// NewOTelTracer, which constructs an OTLP gRPC exporter that would try
// to dial a real endpoint at ctor time.
func newTestTracer(t *testing.T) (*OTelTracer, *tracetest.SpanRecorder) {
	t.Helper()
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	return &OTelTracer{tracer: tp.Tracer(tracerName)}, sr
}

func TestOTelTracer_StartSpan_RecordsNameAndStartAttrs(t *testing.T) {
	tracer, sr := newTestTracer(t)

	_, end := tracer.StartSpan(context.Background(), "mcp.tool.example",
		domain.String("tool.name", "example"),
	)
	end(nil)

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected exactly 1 span, got %d", len(spans))
	}
	if spans[0].Name() != "mcp.tool.example" {
		t.Errorf("span name mismatch: got %q", spans[0].Name())
	}

	if !hasAttr(spans[0].Attributes(), "tool.name", "example") {
		t.Errorf("expected tool.name=example attribute, got %+v", spans[0].Attributes())
	}
}

func TestOTelTracer_EndSpan_RecordsSuccessStatus(t *testing.T) {
	tracer, sr := newTestTracer(t)

	_, end := tracer.StartSpan(context.Background(), "s")
	end(nil)

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Status().Code != codes.Ok {
		t.Errorf("expected codes.Ok, got %v", spans[0].Status().Code)
	}
}

func TestOTelTracer_EndSpan_RecordsErrorStatus(t *testing.T) {
	tracer, sr := newTestTracer(t)

	_, end := tracer.StartSpan(context.Background(), "s")
	end(errors.New("boom"))

	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}
	if spans[0].Status().Code != codes.Error {
		t.Errorf("expected codes.Error, got %v", spans[0].Status().Code)
	}
	if spans[0].Status().Description != "boom" {
		t.Errorf("expected description=boom, got %q", spans[0].Status().Description)
	}
	if len(spans[0].Events()) == 0 {
		t.Error("expected recorded exception event")
	}
}

func TestOTelTracer_EndSpan_MergesExtraAttrs(t *testing.T) {
	tracer, sr := newTestTracer(t)

	_, end := tracer.StartSpan(context.Background(), "s",
		domain.String("tool.name", "foo"),
	)
	end(nil,
		domain.Bool("tool.success", true),
		domain.String("tool.error_class", ""),
	)

	spans := sr.Ended()
	attrs := spans[0].Attributes()
	if !hasAttr(attrs, "tool.name", "foo") {
		t.Error("expected start attr preserved")
	}
	if !hasAttr(attrs, "tool.success", true) {
		t.Error("expected extra bool attr recorded")
	}
	if !hasAttr(attrs, "tool.error_class", "") {
		t.Error("expected extra string attr recorded")
	}
}

func TestOTelTracer_EndSpan_AddsDurationMs(t *testing.T) {
	tracer, sr := newTestTracer(t)

	_, end := tracer.StartSpan(context.Background(), "s")
	end(nil)

	attrs := sr.Ended()[0].Attributes()
	found := false
	for _, kv := range attrs {
		if string(kv.Key) == "duration_ms" && kv.Value.Type() == attribute.INT64 {
			found = true
			if kv.Value.AsInt64() < 0 {
				t.Errorf("negative duration_ms: %d", kv.Value.AsInt64())
			}
		}
	}
	if !found {
		t.Errorf("expected duration_ms attribute on ended span, got %+v", attrs)
	}
}

func TestNewOTelTracer_ConstructorWiresExporter(t *testing.T) {
	// The gRPC OTLP exporter is lazy — it does not dial the endpoint at
	// construction time, so we can hand it an unreachable localhost port
	// without the test hanging or requiring a running collector.
	// Insecure=true skips TLS setup; ServiceName tags the resource.
	ctx, cancel := context.WithTimeout(context.Background(), 5*testTimeout())
	defer cancel()

	tracer, err := NewOTelTracer(ctx, Config{
		Endpoint:    "localhost:0", // unroutable, but not dialed until export
		ServiceName: "saturday-mcp-test",
		Insecure:    true,
	})
	if err != nil {
		t.Fatalf("NewOTelTracer failed: %v", err)
	}
	if tracer == nil {
		t.Fatal("expected non-nil tracer")
	}
	if tracer.tracer == nil {
		t.Error("expected inner tracer to be wired")
	}
	if tracer.shutdown == nil {
		t.Error("expected shutdown closure to be wired")
	}

	// Shutdown with a bounded ctx — the exporter's batcher flush should
	// return promptly since no spans were emitted.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*testTimeout())
	defer shutdownCancel()
	if err := tracer.Shutdown(shutdownCtx); err != nil {
		// Shutdown may error trying to flush to the dead endpoint; the
		// point of this test is that construction wires everything.
		t.Logf("shutdown returned (expected — dead endpoint): %v", err)
	}
}

// testTimeout returns a small time unit used to bound the ctx passed to
// OTel setup — kept as a helper so a slow CI can bump it in one place.
func testTimeout() time.Duration { return 200 * time.Millisecond }

func TestOTelTracer_Shutdown_NilIsNoop(t *testing.T) {
	tracer := &OTelTracer{} // no shutdown closure set
	if err := tracer.Shutdown(context.Background()); err != nil {
		t.Errorf("Shutdown with nil closure must be no-op, got: %v", err)
	}
}

func TestOTelTracer_Shutdown_DelegatesToInner(t *testing.T) {
	called := false
	tracer := &OTelTracer{shutdown: func(context.Context) error {
		called = true
		return nil
	}}
	if err := tracer.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown err: %v", err)
	}
	if !called {
		t.Error("expected inner shutdown to be called")
	}
}

func TestMapAttribute_TypeMapping(t *testing.T) {
	// The mapping surface has to keep parity with the domain.SpanAttribute
	// constructors (String/Int/Int64/Bool). Anything falls back to a
	// stringified default rather than dropping the value.
	cases := []struct {
		name     string
		in       domain.SpanAttribute
		wantType attribute.Type
		wantStr  string // for INVALID/STRING checks
	}{
		{"string", domain.String("k", "v"), attribute.STRING, "v"},
		{"int", domain.Int("k", 42), attribute.INT64, ""},
		{"int64", domain.Int64("k", int64(99)), attribute.INT64, ""},
		{"bool", domain.Bool("k", true), attribute.BOOL, ""},
		{"fallback float", domain.SpanAttribute{Key: "k", Value: 1.5}, attribute.STRING, "1.5"},
		{"fallback struct", domain.SpanAttribute{Key: "k", Value: struct{ X int }{7}}, attribute.STRING, "{7}"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mapAttribute(tc.in)
			if got.Value.Type() != tc.wantType {
				t.Errorf("type mismatch: got %v want %v", got.Value.Type(), tc.wantType)
			}
			if tc.wantStr != "" && got.Value.AsString() != tc.wantStr {
				t.Errorf("string value mismatch: got %q want %q", got.Value.AsString(), tc.wantStr)
			}
		})
	}
}

func TestSpanAttributes_PreservesOrder(t *testing.T) {
	in := []domain.SpanAttribute{
		domain.String("a", "1"),
		domain.String("b", "2"),
		domain.String("c", "3"),
	}
	got := spanAttributes(in)
	if len(got) != 3 {
		t.Fatalf("expected 3 attributes, got %d", len(got))
	}
	for i, want := range []string{"a", "b", "c"} {
		if string(got[i].Key) != want {
			t.Errorf("order broken at %d: got %s want %s", i, got[i].Key, want)
		}
	}
}

// hasAttr searches the attribute list for a key whose value matches
// want. Kept generic so callers pass string/bool/int64 without wrapping.
func hasAttr(attrs []attribute.KeyValue, key string, want any) bool {
	for _, kv := range attrs {
		if string(kv.Key) != key {
			continue
		}
		switch v := want.(type) {
		case string:
			return kv.Value.AsString() == v
		case bool:
			return kv.Value.AsBool() == v
		case int64:
			return kv.Value.AsInt64() == v
		case int:
			return kv.Value.AsInt64() == int64(v)
		}
	}
	return false
}
