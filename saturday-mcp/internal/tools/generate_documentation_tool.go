package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/domain"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// GenerateDocumentationTool renders project markdown documentation and
// writes it to disk. See mcp-add-plan Phase C op 14.
//
// Phase E op 12 replaced the previous direct os.WriteFile call with an
// injected domain.FileWriter — the tool no longer imports the os package.
// The adapter (internal/adapters/filesystem.OSFileSystem) is wired in
// tool_provider.go; unit tests can substitute an in-memory fake.
type GenerateDocumentationTool struct {
	logger    *logging.Logger
	generator *generators.DocumentationGenerator
	fs        domain.FileWriter
}

// NewGenerateDocumentationTool wires the tool with its dependencies.
func NewGenerateDocumentationTool(logger *logging.Logger, generator *generators.DocumentationGenerator, fs domain.FileWriter) *GenerateDocumentationTool {
	return &GenerateDocumentationTool{logger: logger, generator: generator, fs: fs}
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
			"projectPath": projectPathProperty(),
			// NOTE: unlike other generator tools, generate_documentation
			// treats outputPath as an unconditional full file path, not a
			// base directory — so it does not share the outputPathProperty
			// helper the generator tools use alongside writeToFile.
			"outputPath": map[string]interface{}{
				"type":        "string",
				"description": "Absolute path for the output markdown file",
			},
		},
	}
}

// OutputSchema advertises the DocumentationResult response shape —
// distinct from GenerationResult because generate_documentation returns
// {success, path, pages} rather than the code-generator envelope.
func (t *GenerateDocumentationTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&DocumentationResult{})
}

// Execute runs the documentation generator and writes the result to disk.
// Response shape preserved verbatim from server.Handler.handleGenerateDocumentation.
func (t *GenerateDocumentationTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling generate_documentation request")

	args := request.GetArguments()
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

	if err := t.fs.WriteFile(outputPath, []byte(resp.Code)); err != nil {
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
