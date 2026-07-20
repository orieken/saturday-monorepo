// Package metricsfile is the JSON-file-backed adapter that implements
// domain/metrics.Reader. It is the only place in the codebase that reads
// PageMetric data from disk — every use-case depends on the metrics.Reader
// interface, never on this concrete type.
package metricsfile

import (
	"encoding/json"
	"os"

	"github.com/orieken/saturday-mcp/internal/domain/metrics"
)

// FileReader reads PageMetric records from a JSON file. The reader is
// stateless: the caller supplies the source path per-call so a single
// FileReader instance can serve many workflow invocations that each pass
// a different metricsFile argument.
type FileReader struct{}

// NewFileReader returns the adapter. Constructor takes no path — the
// path lives on the ReadPageMetrics call arg. This is the shape change
// Phase E op 13 introduced: the prior FileMetricsProvider captured its
// path at construction time, which forced PrioritizeTestsWorkflow.Run to
// instantiate a new provider inline for every request. Moving the path
// to the call site is what makes constructor-time injection possible.
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadPageMetrics reads and parses the JSON file at source. Preserves
// the exact "raw os error, no wrapping" behavior the prior
// FileMetricsProvider.GetPageMetrics had, so the error message MCP
// clients see does not change.
func (r *FileReader) ReadPageMetrics(source string) ([]metrics.PageMetric, error) {
	file, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}

	var out []metrics.PageMetric
	if err := json.Unmarshal(file, &out); err != nil {
		return nil, err
	}

	return out, nil
}
