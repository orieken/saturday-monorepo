package tools

import "github.com/mark3labs/mcp-go/mcp"

// This file collects the schema fragments that repeat across tool input
// schemas. It is an Extract Function pass — the external MCP contract is
// unchanged, only the construction is deduplicated. See mcp-add-plan
// Phase F op 15.
//
// Every helper returns a fresh map on each call. Callers may compose them
// into their Properties map without worrying about aliasing across tools.

// projectPathProperty returns the standard "Absolute path to the project
// root" string property. Six tools use this field verbatim:
// analyze_performance, suggest_improvements, analyze_framework,
// validate_patterns, analyze_impact, generate_documentation.
func projectPathProperty() map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": "Absolute path to the project root",
	}
}

// projectPathOnlySchema returns the full ToolInputSchema for the tools
// whose entire input is a single required projectPath field.
// analyze_performance and suggest_improvements match verbatim.
func projectPathOnlySchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath"},
		Properties: map[string]interface{}{
			"projectPath": projectPathProperty(),
		},
	}
}

// writeToFileProperty returns the boolean toggle property shared by every
// generator tool that can optionally emit its output to disk.
func writeToFileProperty() map[string]interface{} {
	return map[string]interface{}{
		"type":        "boolean",
		"description": "Whether to write the generated code to a file (default: false)",
		"default":     false,
	}
}

// outputPathProperty returns the "Base directory for output files"
// string property used by every generator tool alongside writeToFile.
// Note: generate_documentation uses outputPath differently (it's the
// full file path, not a base directory, and is required unconditionally)
// so it does NOT use this helper.
func outputPathProperty() map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": "Base directory for output files (required if writeToFile is true)",
	}
}

// descriptionProperty returns an optional free-form string description
// field. Callers supply the specific description text since it varies
// slightly per tool ("Optional description of the site" vs "Optional
// description of the flow" etc).
func descriptionProperty(description string) map[string]interface{} {
	return map[string]interface{}{
		"type":        "string",
		"description": description,
	}
}
