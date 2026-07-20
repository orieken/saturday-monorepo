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

// GeneratePageTool generates a Saturday Page class with element registration.
// See mcp-add-plan Phase C op 7.
type GeneratePageTool struct {
	logger    *logging.Logger
	generator *generators.PageGenerator
}

// NewGeneratePageTool wires the tool with its dependencies.
func NewGeneratePageTool(logger *logging.Logger, generator *generators.PageGenerator) *GeneratePageTool {
	return &GeneratePageTool{logger: logger, generator: generator}
}

func (t *GeneratePageTool) Name() string { return "generate_page" }

func (t *GeneratePageTool) Description() string {
	return "Generate a Page class with element registration"
}

func (t *GeneratePageTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"name", "path", "elements"},
		Properties: map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the page class",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "URL path for the page",
			},
			"elements": map[string]interface{}{
				"type":        "array",
				"description": "List of elements on the page",
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"name", "selector"},
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Element name",
						},
						"selector": map[string]interface{}{
							"type":        "string",
							"description": "CSS selector for the element",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Element type (button, input, link, select, checkbox, radio)",
							"enum":        []string{"button", "input", "link", "select", "checkbox", "radio"},
						},
					},
				},
			},
			"description": descriptionProperty("Optional description of the page"),
			"writeToFile": writeToFileProperty(),
			"outputPath":  outputPathProperty(),
		},
	}
}

// OutputSchema advertises the shared GenerationResult response shape.
func (t *GeneratePageTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&GenerationResult{})
}

// Execute runs the page generator and optionally writes the output.
// Response shape preserved verbatim from server.Handler.handleGeneratePage.
func (t *GeneratePageTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_page request")

	args := request.GetArguments()

	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)

	var req models.PageGenerationRequest
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
		t.logger.Error("Page generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		relativePath := filepath.Join("lib", "pages", resp.FileName)

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

	t.logger.Info("Page generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
