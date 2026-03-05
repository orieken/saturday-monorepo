package analyzers

import (
	"sort"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/observability"
)

// UsageAnalyzer ranks pages based on production metrics
type UsageAnalyzer struct {
	logger *logging.Logger
}

// PagePriority represents a page ranked by importance
type PagePriority struct {
	Path        string  `json:"path"`
	Score       float64 `json:"score"`
	Visits      int     `json:"visits"`
	Explanation string  `json:"explanation"`
}

// NewUsageAnalyzer creates a new analyzer
func NewUsageAnalyzer(logger *logging.Logger) *UsageAnalyzer {
	return &UsageAnalyzer{
		logger: logger,
	}
}

// Analyze ranks pages from raw metrics
func (a *UsageAnalyzer) Analyze(metrics []observability.PageMetric) []PagePriority {
	var priorities []PagePriority

	// Normalize data for scoring
	// Simple formula: Score = Visits * (1 + ErrorRate)
	// High traffic and High errors = Highest Priority
	
	// Find max visits for normalization (optional, but let's keep it simple first)
	
	for _, m := range metrics {
		// Example scoring logic
		// Visits: 1000, ErrorRate: 0.05 (5%) -> Score: 1000 * 1.05 = 1050
		score := float64(m.Visits) * (1.0 + (m.ErrorRate / 100.0))
		
		explanation := "Moderate traffic"
		if m.Visits > 1000 {
			explanation = "High traffic"
		}
		if m.ErrorRate > 1.0 {
			explanation += ", High error rate"
		}

		priorities = append(priorities, PagePriority{
			Path:        m.Path,
			Score:       score,
			Visits:      m.Visits,
			Explanation: explanation,
		})
	}

	// Sort by Score descending
	sort.Slice(priorities, func(i, j int) bool {
		return priorities[i].Score > priorities[j].Score
	})

	return priorities
}

// MatchWithCodebase allows connecting generic metric paths to actual Project files
// This is a placeholder for logic that would use the GraphAnalyzer to map URLs to Files.
func (a *UsageAnalyzer) MatchWithCodebase(priorities []PagePriority, projectFiles []string) []PagePriority {
	// Simple keyword matching for MVP
	// If priority path is "/login" and we have "LoginPage.ts", it's a match.
	return priorities 
}
