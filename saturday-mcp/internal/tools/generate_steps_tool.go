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

// GenerateStepsTool generates Cucumber step definitions from Gherkin patterns.
// See mcp-add-plan Phase C op 9.
type GenerateStepsTool struct {
	logger    *logging.Logger
	generator *generators.StepGenerator
}

// NewGenerateStepsTool wires the tool with its dependencies.
func NewGenerateStepsTool(logger *logging.Logger, generator *generators.StepGenerator) *GenerateStepsTool {
	return &GenerateStepsTool{logger: logger, generator: generator}
}

func (t *GenerateStepsTool) Name() string { return "generate_steps" }

func (t *GenerateStepsTool) Description() string {
	return "Generate Cucumber step definitions from Gherkin patterns"
}

func (t *GenerateStepsTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"steps"},
		Properties: map[string]interface{}{
			"steps": map[string]interface{}{
				"type":        "array",
				"description": "List of step definitions",
				"items": map[string]interface{}{
					"type":     "object",
					"required": []string{"type", "pattern"},
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Step type (Given, When, Then)",
							"enum":        []string{"Given", "When", "Then"},
						},
						"pattern": map[string]interface{}{
							"type":        "string",
							"description": "Gherkin step pattern",
						},
						"parameters": map[string]interface{}{
							"type":        "string",
							"description": "Optional parameter types",
						},
					},
				},
			},
			"language": map[string]interface{}{
				"type":        "string",
				"description": "Target language (typescript or javascript)",
				"enum":        []string{"typescript", "javascript"},
				"default":     "typescript",
			},
			"description": descriptionProperty("Optional description"),
			"writeToFile": writeToFileProperty(),
			"outputPath":  outputPathProperty(),
		},
	}
}

// OutputSchema advertises the shared GenerationResult response shape.
func (t *GenerateStepsTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&GenerationResult{})
}

// Execute runs the step generator and optionally writes the output.
// Response shape preserved verbatim from server.Handler.handleGenerateSteps.
func (t *GenerateStepsTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_steps request")

	args := request.Params.Arguments
	writeToFile, _ := args["writeToFile"].(bool)
	outputPath, _ := args["outputPath"].(string)

	var req models.StepDefinitionRequest
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
		t.logger.Error("Step generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	var filePath string
	if writeToFile {
		if outputPath == "" {
			return mcp.NewToolResultError("outputPath is required when writeToFile is true"), nil
		}
		writer := filewriter.NewFileWriter(outputPath, filewriter.WriteModeOverwrite, false)
		relativePath := filepath.Join("tests", "steps", resp.FileName)
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

	t.logger.Info("Steps generated successfully", "fileName", resp.FileName, "written", writeToFile)
	return mcp.NewToolResultText(string(resultJSON)), nil
}
