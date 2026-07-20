package server

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/domain"
)

// recordedSpan captures one StartSpan/EndSpan cycle so tests can assert
// what the middleware told the tracer. Kept as a plain struct — the
// hand-rolled fake shape matches what the rest of this repo does
// (no gomock/testify/mock).
type recordedSpan struct {
	name       string
	startAttrs []domain.SpanAttribute
	endErr     error
	endAttrs   []domain.SpanAttribute
	ended      bool
}

// fakeTracer implements domain.Tracer for tests. Not thread-safe; the
// middleware only calls into it sequentially per invocation.
type fakeTracer struct {
	spans []*recordedSpan
}

func (f *fakeTracer) StartSpan(ctx context.Context, name string, attrs ...domain.SpanAttribute) (context.Context, domain.EndSpan) {
	span := &recordedSpan{
		name:       name,
		startAttrs: append([]domain.SpanAttribute(nil), attrs...),
	}
	f.spans = append(f.spans, span)
	end := func(err error, extra ...domain.SpanAttribute) {
		span.endErr = err
		span.endAttrs = append([]domain.SpanAttribute(nil), extra...)
		span.ended = true
	}
	return ctx, end
}

// findAttr returns the value for key or nil if absent. Match by exact
// string comparison on the untyped Value so tests can assert
// "tool.success=true" without caring about the concrete type.
func findAttr(attrs []domain.SpanAttribute, key string) any {
	for _, a := range attrs {
		if a.Key == key {
			return a.Value
		}
	}
	return nil
}

func TestWithTracing_SuccessRecordsOneSpan(t *testing.T) {
	tracer := &fakeTracer{}
	next := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("ok"), nil
	}

	wrapped := withTracing(tracer, "example_tool", next)

	_, err := wrapped(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("unexpected transport err: %v", err)
	}

	if len(tracer.spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(tracer.spans))
	}
	span := tracer.spans[0]
	if span.name != "mcp.tool.example_tool" {
		t.Errorf("span name mismatch: got %q", span.name)
	}
	if findAttr(span.startAttrs, "tool.name") != "example_tool" {
		t.Errorf("expected tool.name=example_tool at start, got %v", span.startAttrs)
	}
	if !span.ended {
		t.Error("expected EndSpan to be called")
	}
	if span.endErr != nil {
		t.Errorf("expected nil endErr on success, got %v", span.endErr)
	}
	if findAttr(span.endAttrs, "tool.success") != true {
		t.Errorf("expected tool.success=true, got %v", span.endAttrs)
	}
	if findAttr(span.endAttrs, "tool.error_class") != "" {
		t.Errorf("expected tool.error_class empty on success, got %v", span.endAttrs)
	}
}

func TestWithTracing_ToolResultErrorFlagsSpan(t *testing.T) {
	tracer := &fakeTracer{}
	next := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("business failure"), nil
	}

	wrapped := withTracing(tracer, "bad_tool", next)

	_, err := wrapped(context.Background(), mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("expected nil transport err (tool-result-error is IN-band), got: %v", err)
	}

	span := tracer.spans[0]
	if findAttr(span.endAttrs, "tool.success") != false {
		t.Errorf("expected tool.success=false, got %v", span.endAttrs)
	}
	if findAttr(span.endAttrs, "tool.error_class") != "tool_result_error" {
		t.Errorf("expected tool.error_class=tool_result_error, got %v", span.endAttrs)
	}
	if span.endErr != nil {
		t.Errorf("EndSpan err should be nil for tool-result errors, got %v", span.endErr)
	}
}

func TestWithTracing_TransportErrorRecordsErrClass(t *testing.T) {
	tracer := &fakeTracer{}
	sentinel := errors.New("network blew up")
	next := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return nil, sentinel
	}

	wrapped := withTracing(tracer, "boom_tool", next)

	_, err := wrapped(context.Background(), mcp.CallToolRequest{})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel propagation, got %v", err)
	}

	span := tracer.spans[0]
	if span.endErr != sentinel {
		t.Errorf("expected EndSpan err = sentinel, got %v", span.endErr)
	}
	if findAttr(span.endAttrs, "tool.success") != false {
		t.Errorf("expected tool.success=false, got %v", span.endAttrs)
	}
	// error_class is fmt.Sprintf("%T", err) which for errors.New yields
	// "*errors.errorString". Assert on prefix to stay robust to stdlib
	// implementation shifts.
	class, _ := findAttr(span.endAttrs, "tool.error_class").(string)
	if !strings.Contains(class, "errorString") && !strings.Contains(class, "errors.") {
		t.Errorf("expected error_class to reflect type, got %q", class)
	}
}

func TestWithTracing_ContextPassedToInner(t *testing.T) {
	tracer := &fakeTracer{}
	type ctxKey string
	seen := ""
	next := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if v, ok := ctx.Value(ctxKey("id")).(string); ok {
			seen = v
		}
		return mcp.NewToolResultText("ok"), nil
	}

	wrapped := withTracing(tracer, "t", next)
	parent := context.WithValue(context.Background(), ctxKey("id"), "abc")

	if _, err := wrapped(parent, mcp.CallToolRequest{}); err != nil {
		t.Fatalf("wrapped err: %v", err)
	}
	if seen != "abc" {
		t.Errorf("expected inner handler to receive parent ctx value, got %q", seen)
	}
}

func TestClassifyOutcome(t *testing.T) {
	cases := []struct {
		name         string
		result       *mcp.CallToolResult
		err          error
		wantSuccess  bool
		wantErrClass string // empty means "just don't panic — we only check success"
	}{
		{"nil err, ok result", mcp.NewToolResultText("ok"), nil, true, ""},
		{"nil err, nil result", nil, nil, true, ""},
		{"transport err", nil, errors.New("x"), false, "*errors.errorString"},
		{"tool_result_error", mcp.NewToolResultError("x"), nil, false, "tool_result_error"},
		{"err beats result", mcp.NewToolResultError("x"), errors.New("wire"), false, "*errors.errorString"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			success, errClass := classifyOutcome(tc.result, tc.err)
			if success != tc.wantSuccess {
				t.Errorf("success: got %v want %v", success, tc.wantSuccess)
			}
			if tc.wantErrClass != "" && errClass != tc.wantErrClass {
				t.Errorf("errClass: got %q want %q", errClass, tc.wantErrClass)
			}
		})
	}
}
