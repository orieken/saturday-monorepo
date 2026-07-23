package tools

import "github.com/invopop/jsonschema"

// This file declares the typed response structs that every extracted
// tool serializes into its mcp.NewToolResultText body. Response shapes
// are preserved verbatim from the pre-extraction handlers — the structs
// exist so OutputSchema() can be generated via jsonschema.Reflect
// rather than hand-authored, and so tests can round-trip the JSON
// against a real Go type instead of an untyped map. See mcp-add-plan
// Phase F op 14.
//
// Schemas are marshaled onto mcp.Tool.RawOutputSchema by RegisterTools,
// so they participate in MCP tool discovery. Response shapes stay in
// lockstep with the schema because they are the same Go types.

// GenerationResult is the response returned by every code-generation
// tool: the six Saturday scaffold generators and migrate_code. filePath
// and written are populated only when the caller passes writeToFile=true.
type GenerationResult struct {
	Success  bool              `json:"success"`
	Code     string            `json:"code"`
	FileName string            `json:"fileName"`
	Metadata map[string]string `json:"metadata,omitempty"`
	FilePath string            `json:"filePath,omitempty"`
	Written  bool              `json:"written,omitempty"`
}

// DocumentationResult is the response from generate_documentation.
// pages surfaces the generator's pageCount metadata value as a string
// because generators.DocumentationGenerator's Metadata map is
// map[string]string — preserved verbatim from the pre-extraction shape.
type DocumentationResult struct {
	Success bool   `json:"success"`
	Path    string `json:"path"`
	Pages   string `json:"pages,omitempty"`
}

// ImpactResult is the response from analyze_impact.
type ImpactResult struct {
	Target   string   `json:"target"`
	Impacted []string `json:"impacted"`
	Count    int      `json:"count"`
}

// FunctionComplexity captures complexity metrics for a single function.
type FunctionComplexity struct {
	File         string `json:"file"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
	Complexity   int    `json:"complexity"`
	LineCount    int    `json:"lineCount"`
	Status       string `json:"status"`
}

// ComplexityAnalysisResult is the response from analyze_complexity.
type ComplexityAnalysisResult struct {
	Success         bool                 `json:"success"`
	ProjectPath     string               `json:"projectPath"`
	TotalFiles      int                  `json:"totalFiles"`
	TotalFunctions  int                  `json:"totalFunctions"`
	ViolationsCount int                  `json:"violationsCount"`
	Violations      []FunctionComplexity `json:"violations,omitempty"`
	Summary         string               `json:"summary"`
}

// AccessibilityViolation details a single a11y defect.
type AccessibilityViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	Element     string `json:"element"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

// AccessibilityReportResult is the response from check_accessibility.
type AccessibilityReportResult struct {
	Success         bool                     `json:"success"`
	Path            string                   `json:"path"`
	TotalFiles      int                      `json:"totalFiles"`
	ViolationsCount int                      `json:"violationsCount"`
	Violations      []AccessibilityViolation `json:"violations,omitempty"`
	Summary         string                   `json:"summary"`
}

// LanguageViolation represents a domain term violation.
type LanguageViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	InvalidTerm string `json:"invalidTerm"`
	Suggested   string `json:"suggested"`
}

// UbiquitousLanguageResult is the response from check_ubiquitous_language.
type UbiquitousLanguageResult struct {
	Success         bool                `json:"success"`
	ProjectPath     string              `json:"projectPath"`
	ViolationsCount int                 `json:"violationsCount"`
	Violations      []LanguageViolation `json:"violations,omitempty"`
	Summary         string              `json:"summary"`
}

// DependencyBoundaryViolation represents a Clean Architecture layer violation.
type DependencyBoundaryViolation struct {
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	FromLayer  string `json:"fromLayer"`
	ToLayer    string `json:"toLayer"`
	ImportPath string `json:"importPath"`
}

// DependencyVerificationResult is the response from verify_dependencies.
type DependencyVerificationResult struct {
	Success         bool                          `json:"success"`
	ProjectPath     string                        `json:"projectPath"`
	ViolationsCount int                           `json:"violationsCount"`
	Violations      []DependencyBoundaryViolation `json:"violations,omitempty"`
	Summary         string                        `json:"summary"`
}

// KIMatch represents a Knowledge Item search hit.
type KIMatch struct {
	Title        string   `json:"title"`
	Path         string   `json:"path"`
	Summary      string   `json:"summary"`
	Tags         []string `json:"tags,omitempty"`
	Relevance    float64  `json:"relevance"`
}

// KISearchResult is the response from search_ki.
type KISearchResult struct {
	Success    bool      `json:"success"`
	Query      string    `json:"query"`
	TotalHits  int       `json:"totalHits"`
	Matches    []KIMatch `json:"matches,omitempty"`
}

// reflectSchema wraps jsonschema.Reflect so every OutputSchema()
// implementation reads the same one-liner. Using a package-level
// jsonschema.Reflector with default settings keeps schema output
// stable and testable.
func reflectSchema(v interface{}) *jsonschema.Schema {
	return jsonschema.Reflect(v)
}
