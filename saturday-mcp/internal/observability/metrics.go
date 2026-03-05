package observability

// PageMetric represents usage data for a specific page path
type PageMetric struct {
	Path      string  `json:"path"`
	Visits    int     `json:"visits"`
	ErrorRate float64 `json:"errorRate"` // Percentage 0.0 to 100.0
	LatencyMs int     `json:"latencyMs"`
}

// MetricsProvider defines how to fetch page metrics
type MetricsProvider interface {
	GetPageMetrics() ([]PageMetric, error)
}
