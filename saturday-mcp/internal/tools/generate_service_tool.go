package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/filewriter"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// GenerateServiceTool generates a Saturday API Service class.
// See mcp-add-plan Phase C op 11.
type GenerateServiceTool struct {
	logger    *logging.Logger
	generator *generators.ServiceGenerator
}

// NewGenerateServiceTool wires the tool with its dependencies.
func NewGenerateServiceTool(logger *logging.Logger, generator *generators.ServiceGenerator) *GenerateServiceTool {
	return &GenerateServiceTool{logger: logger, generator: generator}
}

func (t *GenerateServiceTool) Name() string { return "generate_service" }

func (t *GenerateServiceTool) Description() string {
	return "Generate an API Service class"
}

func (t *GenerateServiceTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"name", "baseUrl", "endpoints"},
		Properties: map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the service",
			},
			"baseUrl": map[string]interface{}{
				"type":        "string",
				"description": "Base URL for the service",
			},
			"endpoints": map[string]interface{}{
				"type":        "array",
				"description": "List of API endpoints",
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"name", "method", "path"},
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "string",
						},
						"method": map[string]interface{}{
							"type": "string",
							"enum": []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
						},
						"path": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Optional description",
			},
			"writeToFile": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether to write the generated code to a file (default: false)",
				"default":     false,
			},
			"outputPath": map[string]interface{}{
				"type":        "string",
				"description": "Base directory for output files (required if writeToFile is true)",
			},
		},
	}
}

// Execute runs the service generator and optionally writes the output.
// Response shape preserved verbatim from server.Handler.handleGenerateService.
func (t *GenerateServiceTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_service request")

	args := request.Params.Arguments

	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)

	var req models.ServiceGenerationRequest
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
		t.logger.Error("Service generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}

		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)

		relativePath := filepath.Join("lib", "services", resp.FileName)

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

	t.logger.Info("Service generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
