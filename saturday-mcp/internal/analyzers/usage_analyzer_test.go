package analyzers

import (
	"testing"

	"github.com/orieken/saturday-mcp/internal/domain/metrics"
)

func TestUsageAnalyzer_Analyze(t *testing.T) {
	analyzer := NewUsageAnalyzer(nil)

	metricsData := []metrics.PageMetric{
		{Path: "/login", Visits: 1000, ErrorRate: 0.0},        // Score: 1000
		{Path: "/checkout", Visits: 500, ErrorRate: 5.0},      // Score: 525
		{Path: "/home", Visits: 5000, ErrorRate: 0.1},         // Score: 5005
		{Path: "/broken", Visits: 100, ErrorRate: 50.0},       // Score: 150
	}

	priorities := analyzer.Analyze(metricsData)

	if len(priorities) != 4 {
		t.Fatalf("Expected 4 results, got %d", len(priorities))
	}

	// Expect /home to be first (highest mass)
	if priorities[0].Path != "/home" {
		t.Errorf("Expected top priority to be /home, got %s", priorities[0].Path)
	}

	// Expect /login (1000) to be higher than /checkout (525)
	if priorities[1].Path != "/login" {
		t.Errorf("Expected second priority to be /login, got %s", priorities[1].Path)
	}
}
