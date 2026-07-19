package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// GenerateDocumentationTool renders project markdown documentation and
// writes it to disk. See mcp-add-plan Phase C op 14.
//
// NOTE: this tool still calls os.WriteFile directly rather than routing
// through the filewriter package. That direct-syscall path is called out in
// the plan (Pattern Conformance row #2) as a Phase E concern — introducing a
// FileSystem adapter interface — which is deliberately out of scope for the
// Phase C Extract-Class op. Leave the inline os.WriteFile in place here; it
// will be replaced when Phase E lands.
type GenerateDocumentationTool struct {
	logger    *logging.Logger
	generator *generators.DocumentationGenerator
}

// NewGenerateDocumentationTool wires the tool with its dependencies.
func NewGenerateDocumentationTool(logger *logging.Logger, generator *generators.DocumentationGenerator) *GenerateDocumentationTool {
	return &GenerateDocumentationTool{logger: logger, generator: generator}
}

func (t *GenerateDocumentationTool) Name() string { return "generate_documentation" }

func (t *GenerateDocumentationTool) Description() string {
	return "Generate markdown documentation for the project"
}

func (t *GenerateDocumentationTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"projectPath", "outputPath"},
		Properties: map[string]interface{}{
			"projectPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path to the project root",
			},
			"outputPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path for the output markdown file",
			},
		},
	}
}

// Execute runs the documentation generator and writes the result to disk.
// Response shape preserved verbatim from server.Handler.handleGenerateDocumentation.
func (t *GenerateDocumentationTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_documentation request")

	args := request.Params.Arguments
	projectPath, _ := args["projectPath"].(string)
	outputPath, _ := args["outputPath"].(string)

	if projectPath == "" || outputPath == "" {
		return mcp.NewToolResultError("projectPath and outputPath are required"), nil
	}

	req := models.DocumentationRequest{
		ProjectPath: projectPath,
		OutputPath:  outputPath,
	}

	resp, err := t.generator.Generate(req)
	if err != nil {
		t.logger.Error("Documentation generation failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Generation failed: %v", err)), nil
	}

	if err := os.WriteFile(outputPath, []byte(resp.Code), 0644); err != nil {
		t.logger.Error("Failed to write documentation file", "path", outputPath, "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to write file: %v", err)), nil
	}

	t.logger.Info("Documentation generated successfully", "path", outputPath)

	result := map[string]interface{}{
		"success": true,
		"path":    outputPath,
		"pages":   resp.Metadata["pageCount"],
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}
