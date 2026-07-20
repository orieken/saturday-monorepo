package metricsfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/domain/metrics"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("test setup: WriteFile failed: %v", err)
	}
	return path
}

func TestFileReader_ReadPageMetrics_HappyPath(t *testing.T) {
	reader := NewFileReader()

	sample := []metrics.PageMetric{
		{Path: "/login", Visits: 500, ErrorRate: 2.5, LatencyMs: 120},
		{Path: "/checkout", Visits: 1200, ErrorRate: 0.8, LatencyMs: 340},
	}
	raw, err := json.Marshal(sample)
	if err != nil {
		t.Fatalf("test setup: marshal failed: %v", err)
	}
	path := writeTemp(t, "metrics.json", string(raw))

	got, err := reader.ReadPageMetrics(path)
	if err != nil {
		t.Fatalf("ReadPageMetrics failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
	if got[0].Path != "/login" || got[0].Visits != 500 {
		t.Errorf("first record mismatch: %+v", got[0])
	}
	if got[1].ErrorRate != 0.8 {
		t.Errorf("second record ErrorRate mismatch: %+v", got[1])
	}
}

func TestFileReader_ReadPageMetrics_EmptyArray(t *testing.T) {
	reader := NewFileReader()

	path := writeTemp(t, "empty_array.json", "[]")

	got, err := reader.ReadPageMetrics(path)
	if err != nil {
		t.Fatalf("empty array should not error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 records for empty array, got %d", len(got))
	}
}

func TestFileReader_ReadPageMetrics_MalformedJSON(t *testing.T) {
	reader := NewFileReader()

	path := writeTemp(t, "bad.json", "{not json}")

	got, err := reader.ReadPageMetrics(path)
	if err == nil {
		t.Fatalf("expected error for malformed JSON, got: %+v", got)
	}
	if got != nil {
		t.Errorf("expected nil result on error, got %+v", got)
	}
}

func TestFileReader_ReadPageMetrics_MissingFile(t *testing.T) {
	reader := NewFileReader()

	missing := filepath.Join(t.TempDir(), "does_not_exist.json")

	got, err := reader.ReadPageMetrics(missing)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected os.IsNotExist error, got %v", err)
	}
	if got != nil {
		t.Errorf("expected nil result on error, got %+v", got)
	}
}

func TestFileReader_ReadPageMetrics_EmptyFile(t *testing.T) {
	reader := NewFileReader()

	// An empty file is not valid JSON — json.Unmarshal on "" returns
	// an unexpected-end-of-input error. Adapter should surface this
	// verbatim rather than swallow it.
	path := writeTemp(t, "empty.json", "")

	got, err := reader.ReadPageMetrics(path)
	if err == nil {
		t.Fatal("expected error for empty file")
	}
	if got != nil {
		t.Errorf("expected nil result on error, got %+v", got)
	}
}

func TestFileReader_ReadPageMetrics_MultipleCallsSameReader(t *testing.T) {
	// Reader is stateless per the doc comment — a single instance must
	// handle sequential calls with different source paths cleanly.
	reader := NewFileReader()

	a := writeTemp(t, "a.json", `[{"path":"/a","visits":10,"errorRate":0,"latencyMs":50}]`)
	b := writeTemp(t, "b.json", `[{"path":"/b","visits":20,"errorRate":0,"latencyMs":60}]`)

	gotA, err := reader.ReadPageMetrics(a)
	if err != nil {
		t.Fatalf("read a: %v", err)
	}
	gotB, err := reader.ReadPageMetrics(b)
	if err != nil {
		t.Fatalf("read b: %v", err)
	}
	if gotA[0].Path != "/a" {
		t.Errorf("first read leaked into second: %+v", gotA[0])
	}
	if gotB[0].Path != "/b" {
		t.Errorf("second read did not overwrite: %+v", gotB[0])
	}
}

func TestNewFileReader_ReturnsNonNil(t *testing.T) {
	if NewFileReader() == nil {
		t.Fatal("NewFileReader returned nil")
	}
}
