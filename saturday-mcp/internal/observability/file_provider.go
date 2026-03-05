package observability

import (
	"encoding/json"
	"os"
)

// FileMetricsProvider implements MetricsProvider using a JSON file
type FileMetricsProvider struct {
	filePath string
}

// NewFileMetricsProvider creates a new provider reading from the given path
func NewFileMetricsProvider(filePath string) *FileMetricsProvider {
	return &FileMetricsProvider{
		filePath: filePath,
	}
}

// GetPageMetrics reads and parses the JSON metrics
func (p *FileMetricsProvider) GetPageMetrics() ([]PageMetric, error) {
	file, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, err
	}

	var metrics []PageMetric
	if err := json.Unmarshal(file, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}
