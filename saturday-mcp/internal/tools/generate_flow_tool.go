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

// GenerateFlowTool generates a Saturday Flow class for multi-step user
// journeys. See mcp-add-plan Phase C op 8.
type GenerateFlowTool struct {
	logger    *logging.Logger
	generator *generators.FlowGenerator
}

// NewGenerateFlowTool wires the tool with its dependencies.
func NewGenerateFlowTool(logger *logging.Logger, generator *generators.FlowGenerator) *GenerateFlowTool {
	return &GenerateFlowTool{logger: logger, generator: generator}
}

func (t *GenerateFlowTool) Name() string { return "generate_flow" }

func (t *GenerateFlowTool) Description() string {
	return "Generate a Flow class for multi-step user journeys"
}

func (t *GenerateFlowTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"name", "steps"},
		Properties: map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the flow class",
			},
			"steps": map[string]interface{}{
				"type":        "array",
				"description": "List of step method names in the flow",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
			"description": descriptionProperty("Optional description of the flow"),
			"writeToFile": writeToFileProperty(),
			"outputPath":  outputPathProperty(),
		},
	}
}

// OutputSchema advertises the shared GenerationResult response shape.
func (t *GenerateFlowTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&GenerationResult{})
}

// Execute runs the flow generator and optionally writes the output.
// Response shape preserved verbatim from server.Handler.handleGenerateFlow.
func (t *GenerateFlowTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_flow request")

	args := request.Params.Arguments
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)

	var req models.FlowGenerationRequest
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
		t.logger.Error("Flow generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}
		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		relativePath := filepath.Join("lib", "flows", resp.FileName)
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

	t.logger.Info("Flow generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
