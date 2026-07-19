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
// mcp-go v0.8.0's mcp.Tool struct does not currently carry an
// outputSchema field, so nothing wires these schemas into the MCP
// advertisement automatically. They exist in code today as a declarative
// contract — future mcp-go upgrades (or a middleware that inspects
// OutputSchema()) can surface them.

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

// reflectSchema wraps jsonschema.Reflect so every OutputSchema()
// implementation reads the same one-liner. Using a package-level
// jsonschema.Reflector with default settings keeps schema output
// stable and testable.
func reflectSchema(v interface{}) *jsonschema.Schema {
	return jsonschema.Reflect(v)
}
