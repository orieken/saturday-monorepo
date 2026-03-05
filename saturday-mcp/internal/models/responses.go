package models

// GenerationResponse represents a successful code generation response
type GenerationResponse struct {
	Code     string            `json:"code"`
	FileName string            `json:"fileName"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// AnalysisResponse represents framework analysis results
type AnalysisResponse struct {
	Sites    []string          `json:"sites"`
	Pages    []string          `json:"pages"`
	Flows    []string          `json:"flows"`
	Elements []string          `json:"elements"`
	Metrics  map[string]int    `json:"metrics"`
	Issues   []ValidationIssue `json:"issues,omitempty"`
}

// ValidationResponse represents validation results
type ValidationResponse struct {
	Valid  bool              `json:"valid"`
	Issues []ValidationIssue `json:"issues,omitempty"`
}

// ValidationIssue represents a validation or pattern issue
type ValidationIssue struct {
	Severity string `json:"severity"` // error, warning, info
	Message  string `json:"message"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Rule     string `json:"rule,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// ImprovementResponse represents code improvement suggestions
type ImprovementResponse struct {
	Suggestions []ValidationIssue      `json:"suggestions"`
	Summary     map[string]interface{} `json:"summary,omitempty"`
}

// TestExecutionResult represents the result of a test run
type TestExecutionResult struct {
	Success  bool   `json:"success"`
	Output   string `json:"output"` // stdout/stderr combined
	Duration string `json:"duration"`
	Summary  string `json:"summary"` // e.g. "Passed: 5, Failed: 1"
}
