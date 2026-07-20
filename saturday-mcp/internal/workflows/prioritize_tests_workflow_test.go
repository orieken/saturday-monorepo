package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
	"github.com/orieken/saturday-mcp/internal/domain/metrics"
	"github.com/orieken/saturday-mcp/internal/logging"
)

// fakeReader is a hand-rolled metrics.Reader double. Mirrors fakeRunner's
// shape from run_tests_workflow_test.go so both files stay symmetric.
type fakeReader struct {
	data       []metrics.PageMetric
	err        error
	called     bool
	lastSource string
}

func (f *fakeReader) ReadPageMetrics(source string) ([]metrics.PageMetric, error) {
	f.called = true
	f.lastSource = source
	return f.data, f.err
}

func newPrioritizeWorkflow(reader *fakeReader) *PrioritizeTestsWorkflow {
	logger := logging.NewLogger(&bytes.Buffer{})
	return NewPrioritizeTestsWorkflow(logger, analyzers.NewUsageAnalyzer(logger), reader)
}

func TestPrioritizeTestsWorkflow_Metadata(t *testing.T) {
	w := newPrioritizeWorkflow(&fakeReader{})

	if w.Name() != "prioritize_tests" {
		t.Errorf("Name mismatch: %q", w.Name())
	}
	if w.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := w.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	if len(schema.Required) != 1 || schema.Required[0] != "metricsFile" {
		t.Errorf("expected Required=[metricsFile], got %v", schema.Required)
	}

	if w.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestPrioritizeTestsWorkflow_Run_Success(t *testing.T) {
	reader := &fakeReader{
		data: []metrics.PageMetric{
			{Path: "/checkout", Visits: 1200, ErrorRate: 2.5, LatencyMs: 340},
			{Path: "/home", Visits: 300, ErrorRate: 0.1, LatencyMs: 80},
		},
	}
	w := newPrioritizeWorkflow(reader)

	req := buildRequest(map[string]any{
		"metricsFile": "/tmp/metrics.json",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if result.IsError {
		t.Errorf("expected success, got IsError=true; text=%q", extractText(t, result))
	}
	if !reader.called {
		t.Error("expected reader.ReadPageMetrics to be called")
	}
	if reader.lastSource != "/tmp/metrics.json" {
		t.Errorf("expected metricsFile forwarded, got %q", reader.lastSource)
	}

	var priorities []analyzers.PagePriority
	if err := json.Unmarshal([]byte(extractText(t, result)), &priorities); err != nil {
		t.Fatalf("result text is not valid JSON: %v", err)
	}
	if len(priorities) != 2 {
		t.Fatalf("expected 2 priorities, got %d", len(priorities))
	}
	// Ranking: /checkout has higher visits & higher error rate → first.
	if priorities[0].Path != "/checkout" {
		t.Errorf("expected /checkout first, got %+v", priorities[0])
	}
}

func TestPrioritizeTestsWorkflow_Run_MissingMetricsFile(t *testing.T) {
	reader := &fakeReader{}
	w := newPrioritizeWorkflow(reader)

	result, err := w.Run(context.Background(), buildRequest(map[string]any{}))
	if err != nil {
		t.Fatalf("Run should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when metricsFile is missing")
	}
	if reader.called {
		t.Error("reader should not be called when validation fails")
	}
	if !strings.Contains(extractText(t, result), "metricsFile argument is required") {
		t.Errorf("expected validation message, got %q", extractText(t, result))
	}
}

func TestPrioritizeTestsWorkflow_Run_EmptyMetricsFileTreatedAsMissing(t *testing.T) {
	reader := &fakeReader{}
	w := newPrioritizeWorkflow(reader)

	req := buildRequest(map[string]any{
		"metricsFile": "",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for empty metricsFile")
	}
	if reader.called {
		t.Error("reader should not be called for empty file arg")
	}
}

func TestPrioritizeTestsWorkflow_Run_NonStringMetricsFileRejected(t *testing.T) {
	// The type-assertion path in Run must reject non-string values
	// rather than panic on a bad cast.
	reader := &fakeReader{}
	w := newPrioritizeWorkflow(reader)

	req := buildRequest(map[string]any{
		"metricsFile": 42,
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for non-string metricsFile")
	}
	if reader.called {
		t.Error("reader should not be called for non-string arg")
	}
}

func TestPrioritizeTestsWorkflow_Run_ReaderErrorSurfacedAsToolResult(t *testing.T) {
	reader := &fakeReader{
		err: errors.New("open metrics.json: permission denied"),
	}
	w := newPrioritizeWorkflow(reader)

	req := buildRequest(map[string]any{
		"metricsFile": "/tmp/metrics.json",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when reader errors")
	}
	if !strings.Contains(extractText(t, result), "Failed to load metrics") {
		t.Errorf("expected error prefix, got %q", extractText(t, result))
	}
}

func TestPrioritizeTestsWorkflow_Run_EmptyMetricsProducesEmptyArray(t *testing.T) {
	reader := &fakeReader{
		data: []metrics.PageMetric{},
	}
	w := newPrioritizeWorkflow(reader)

	req := buildRequest(map[string]any{
		"metricsFile": "/tmp/metrics.json",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if result.IsError {
		t.Errorf("expected success even with empty metrics, got IsError=true")
	}

	text := extractText(t, result)
	// Empty priorities marshal to "null" (nil slice) — assert on either
	// possibility so the test does not lock in an incidental JSON shape.
	if text != "null" && text != "[]" {
		var out []analyzers.PagePriority
		if err := json.Unmarshal([]byte(text), &out); err != nil {
			t.Fatalf("unexpected result shape %q: %v", text, err)
		}
		if len(out) != 0 {
			t.Errorf("expected 0 priorities, got %d", len(out))
		}
	}
}
