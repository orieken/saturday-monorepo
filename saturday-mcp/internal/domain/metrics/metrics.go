// Package metrics defines the domain concepts for page-usage metrics
// used by test prioritization (see analyzers.UsageAnalyzer). This is a
// DOMAIN concept — it models what the framework knows about production
// usage of pages, so the risk-based ranking that surfaces "which page
// needs test coverage next" has real data to work with.
//
// It is NOT OpenTelemetry infrastructure. The prior name
// (internal/observability) was misleading because it collided with the
// framework-wide "observability" term that usually means OTel spans and
// telemetry pipelines. Phase E op 13 (mcp-add-plan) renamed the package
// to eliminate that confusion.
package metrics

// PageMetric represents usage data for a specific page path in the
// system under test — the input to test prioritization.
type PageMetric struct {
	Path      string  `json:"path"`
	Visits    int     `json:"visits"`
	ErrorRate float64 `json:"errorRate"` // Percentage 0.0 to 100.0
	LatencyMs int     `json:"latencyMs"`
}

// Reader is the domain seam for loading PageMetric records from any
// backing source (a local JSON file, a metrics API, an in-memory fake,
// …). Renamed from the prior "MetricsProvider" because "Reader"
// describes intent — the callers only need to READ metrics — and stays
// short in `metrics.Reader` idiom.
//
// Contract:
//   - source is an adapter-specific locator string: the file-backed
//     adapter interprets it as a path; a future HTTP-backed adapter
//     would interpret it as a URL. The domain layer takes no position
//     on its shape beyond "the string the tool caller passed in".
//   - Moving source from constructor arg (as the prior FileMetricsProvider
//     did) to per-call arg is the design change that lets workflows
//     inject one Reader at wiring time and still handle a different
//     path per Run invocation — the metricsFile field on the tool input
//     changes with every request.
type Reader interface {
	ReadPageMetrics(source string) ([]PageMetric, error)
}
