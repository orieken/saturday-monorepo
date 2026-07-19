package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/generators"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// MigrateCodeTool analyzes legacy Playwright source and migrates it to the
// Saturday Framework's Site-Centric pattern. See mcp-add-plan Phase C op 12.
type MigrateCodeTool struct {
	logger    *logging.Logger
	generator *generators.MigrationGenerator
}

// NewMigrateCodeTool wires the tool with its dependencies.
func NewMigrateCodeTool(logger *logging.Logger, generator *generators.MigrationGenerator) *MigrateCodeTool {
	return &MigrateCodeTool{logger: logger, generator: generator}
}

func (t *MigrateCodeTool) Name() string { return "migrate_code" }

func (t *MigrateCodeTool) Description() string {
	return "Analyze and migrate legacy code to Saturday Framework patterns"
}

func (t *MigrateCodeTool) InputSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type:     "object",
		Required: []string{"sourceCode"},
		Properties: map[string]interface{}{
			"sourceCode": map[string]interface{}{
				"type":        "string",
				"description": "Legacy Playwright source code to migrate",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"description": "Type of migration (page or test)",
				"enum":        []string{"page", "test"},
				"default":     "page",
			},
		},
	}
}

// OutputSchema advertises the shared GenerationResult response shape.
// migrate_code omits filePath/written since it never writes to disk, but
// the schema fields being optional means the shared type still fits.
func (t *MigrateCodeTool) OutputSchema() *jsonschema.Schema {
	return reflectSchema(&GenerationResult{})
}

// Execute runs the migration generator. Response shape preserved verbatim
// from server.Handler.handleMigrateCode.
func (t *MigrateCodeTool) Execute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t.logger.Info("Handling migrate_code request")

	args := request.Params.Arguments

	sourceCode, _ := args["sourceCode"].(string)
	migrationType, ok := args["type"].(string)
	if !ok || migrationType == "" {
		migrationType = "page"
	}

	req := models.MigrationRequest{
		SourceCode: sourceCode,
		Type:       migrationType,
	}

	resp, err := t.generator.Generate(req)
	if err != nil {
		t.logger.Error("Migration failed", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Migration failed: %v", err)), nil
	}

	result := map[string]interface{}{
		"success":  true,
		"code":     resp.Code,
		"fileName": resp.FileName,
		"metadata": resp.Metadata,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.logger.Error("Failed to marshal result", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format result: %v", err)), nil
	}

	t.logger.Info("Code migrated successfully")
	return mcp.NewToolResultText(string(resultJSON)), nil
}
