// Package tools contains one file per registered MCP tool. Each file
// exposes a *Tool type that satisfies domain.Tool. This is the target
// structure of mcp-add-plan Phase C — the server package stops being a
// God Object once every handleFoo lives here.
//
// Constructor injection uses concrete types today (e.g. *generators.SiteGenerator,
// *filewriter.FileWriter). Phase E (Milestone 2) will replace those with
// interfaces defined in the domain layer as part of adapter isolation.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/filewriter"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// GenerateSiteTool generates a Saturday Site class with page and flow
// registration. See mcp-add-plan Phase C op 6 for the extraction context.
type GenerateSiteTool struct {
	logger    *logging.Logger
	generator *generators.SiteGenerator
}

// NewGenerateSiteTool wires the tool with its dependencies.
func NewGenerateSiteTool(logger *logging.Logger, generator *generators.SiteGenerator) *GenerateSiteTool {
	return &GenerateSiteTool{
		logger:    logger,
		generator: generator,
	}
}

// Name is the MCP tool identifier — stable, part of the public contract.
func (t *GenerateSiteTool) Name() string { return "generate_site" }

// Description is the human-readable tool description sent to MCP clients.
func (t *GenerateSiteTool) Description() string {
	return "Generate a Site class with page and flow registration"
}

// InputSchema is the JSON Schema advertised via MCP tool discovery.
func (t *GenerateSiteTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"name", "baseUrl", "pages"},
		Properties: map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the site class",
			},
			"baseUrl": map[string]interface{}{
				"type":        "string",
				"description": "Base URL for the site",
			},
			"pages": map[string]interface{}{
				"type":        "array",
				"description": "List of page names to register",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"flows": map[string]interface{}{
				"type":        "array",
				"description": "Optional list of flow names",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"description": descriptionProperty("Optional description of the site"),
			"writeToFile": writeToFileProperty(),
			"outputPath":  outputPathProperty(),
		},
	}
}

// OutputSchema advertises the GenerationResult response shape. Shared
// across every code-generation tool; see internal/tools/responses.go.
func (t *GenerateSiteTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&GenerationResult{})
}

// Execute runs the site generator and optionally writes the output to disk.
// Response shape preserved verbatim from the old server.Handler.handleGenerateSite
// so the integration tests in internal/integration/e2e_test.go continue to pass
// without modification.
func (t *GenerateSiteTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_site request")

	args := request.GetArguments()

	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)

	var req models.SiteGenerationRequest
	argsJSON, err := json.Marshal(args)
	if err != nil {
		t.logger.Error("Failed to marshal arguments", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid arguments: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &req); err != nil {
		t.logger.Error("Failed to unmarshal request", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Invalid request format: %v", err)), nil
	}

	resp, err := t.generator.Generate(req)
	if err != nil {
		t.logger.Error("Site generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)

		relativePath := filepath.Join("lib", "sites", resp.FileName)

		if err := writer.WriteFile(relativePath, resp.Code); err != nil {
			t.logger.Error("Failed to write file", "error", err, "path", relativePath)
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
		}

		fullPath, _ := writer.GetFullPath(relativePath)
		filePath = fullPath
		t.logger.Info("File written successfully", "path", fullPath)
	}

	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	if writeToFile {
		result["filePath"] = filePath
		result["written"] = true
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Site generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
