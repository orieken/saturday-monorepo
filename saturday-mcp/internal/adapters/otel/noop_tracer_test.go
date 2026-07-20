package otel

import (
	"context"
	"errors"
	"testing"

	"github.com/orieken/saturday-mcp/internal/domain"
)

func TestNoopTracer_StartSpan_ReturnsSameContext(t *testing.T) {
	tracer := NewNoopTracer()

	type ctxKey string
	parent := context.WithValue(context.Background(), ctxKey("id"), "abc")

	got, end := tracer.StartSpan(parent, "any.span", domain.String("key", "value"))
	if got == nil {
		t.Fatal("expected non-nil context")
	}
	if got.Value(ctxKey("id")) != "abc" {
		t.Error("noop tracer should not strip parent context values")
	}
	if end == nil {
		t.Fatal("expected non-nil EndSpan")
	}
}

func TestNoopTracer_EndSpan_HandlesNilError(t *testing.T) {
	tracer := NewNoopTracer()

	_, end := tracer.StartSpan(context.Background(), "s")

	// The whole point of the noop is that these calls do nothing —
	// this test's assertion is "does not panic".
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("EndSpan(nil) panicked: %v", r)
		}
	}()
	end(nil)
}

func TestNoopTracer_EndSpan_HandlesErrorAndExtraAttrs(t *testing.T) {
	tracer := NewNoopTracer()

	_, end := tracer.StartSpan(context.Background(), "s")

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("EndSpan with err+attrs panicked: %v", r)
		}
	}()
	end(errors.New("boom"),
		domain.Bool("tool.success", false),
		domain.String("tool.error_class", "test_error"),
	)
}

func TestNewNoopTracer_ReturnsNonNil(t *testing.T) {
	if NewNoopTracer() == nil {
		t.Fatal("NewNoopTracer returned nil")
	}
}
