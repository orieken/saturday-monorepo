package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/orieken/saturday-mcp/internal/generators"
)

// newMigrateCodeTool wires the tool with the real migration generator and
// template registry, matching the shape used by the other generator-tool
// tests. Migrate is not a file-writing tool — it only produces a code
// string — so there is no filewriter/outputPath axis to cover.
func newMigrateCodeTool(t *testing.T) *MigrateCodeTool {
	t.Helper()
	gen := generators.NewMigrationGenerator(newProcessor(t), newValidator())
	return NewMigrateCodeTool(silentLogger(), gen)
}

func TestMigrateCodeTool_Metadata(t *testing.T) {
	tool := newMigrateCodeTool(t)

	if tool.Name() != "migrate_code" {
		t.Errorf("Name: got %q, want %q", tool.Name(), "migrate_code")
	}
	if tool.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	// Only sourceCode is required at the transport layer — the tool
	// defaults type to "page" when the caller omits it.
	wantRequired := map[string]bool{"sourceCode": true}
	if len(schema.Required) != len(wantRequired) {
		t.Errorf("Required len: got %d, want %d", len(schema.Required), len(wantRequired))
	}
	for _, r := range schema.Required {
		if !wantRequired[r] {
			t.Errorf("unexpected required field: %s", r)
		}
	}
	if tool.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestMigrateCodeTool_Execute_Success_DefaultsToPage(t *testing.T) {
	tool := newMigrateCodeTool(t)

	// Two distinct selectors (button + input) drive two element rows in
	// the generated draft page. The regex in migrateGenerator dedupes on
	// selector, so a duplicate 3rd selector would not add a row.
	src := `
test('t', async ({ page }) => {
  await page.click('#login-btn');
  await page.fill('input[name="user"]', 'admin');
});
`
	req := buildRequest(map[string]any{
		"sourceCode": src,
		// type omitted — tool defaults to "page"
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got GenerationResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	// The real page.tmpl renders {{pascalCase .Name}}Page — the generator
	// hard-codes Name="DraftMigrated", so the rendered class is
	// DraftMigratedPage extending BasePage.
	if !strings.Contains(got.Code, "DraftMigratedPage") {
		t.Errorf("expected code to contain DraftMigratedPage, got:\n%s", got.Code)
	}
	if !strings.Contains(got.Code, "extends BasePage") {
		t.Errorf("expected code to extend BasePage, got:\n%s", got.Code)
	}
	if got.FileName != "draft-migrated-page.ts" {
		t.Errorf("FileName: got %q, want %q", got.FileName, "draft-migrated-page.ts")
	}
	// Metadata surfaced through the response — assert both keys exist
	// with the expected shape (elementCount reflects the deduped selector
	// count, type is always "page" for the supported branch).
	if got.Metadata["type"] != "page" {
		t.Errorf("Metadata[type]: got %q, want %q", got.Metadata["type"], "page")
	}
	if got.Metadata["elementCount"] != "2" {
		t.Errorf("Metadata[elementCount]: got %q, want %q", got.Metadata["elementCount"], "2")
	}
	if got.Metadata["migrated"] != "true" {
		t.Errorf("Metadata[migrated]: got %q, want %q", got.Metadata["migrated"], "true")
	}
}

func TestMigrateCodeTool_Execute_ExplicitPageType(t *testing.T) {
	tool := newMigrateCodeTool(t)

	// Same as the defaults test, but with type explicitly passed so the
	// non-default branch of the type-normalization block is exercised.
	req := buildRequest(map[string]any{
		"sourceCode": `await page.click('button.submit');`,
		"type":       "page",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute err: %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got IsError=true; content=%q", extractText(t, result))
	}

	var got GenerationResult
	if err := json.Unmarshal([]byte(extractText(t, result)), &got); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if !got.Success {
		t.Error("expected success=true")
	}
	if got.Metadata["elementCount"] != "1" {
		t.Errorf("Metadata[elementCount]: got %q, want %q", got.Metadata["elementCount"], "1")
	}
}

func TestMigrateCodeTool_Execute_UnsupportedTypeSurfacesAsToolError(t *testing.T) {
	tool := newMigrateCodeTool(t)

	// Type="test" passes the oneof validator but has no code path in
	// MigrationGenerator.Generate — it falls through to the
	// "unsupported migration type" error, which the tool must surface
	// as an IsError result rather than a transport error.
	req := buildRequest(map[string]any{
		"sourceCode": `await page.click('x');`,
		"type":       "test",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for unsupported type")
	}
	msg := extractText(t, result)
	if !strings.Contains(msg, "Migration failed") {
		t.Errorf("expected 'Migration failed' prefix, got %q", msg)
	}
	if !strings.Contains(msg, "unsupported migration type") {
		t.Errorf("expected 'unsupported migration type' in message, got %q", msg)
	}
}

func TestMigrateCodeTool_Execute_InvalidRequestSurfacesAsToolError(t *testing.T) {
	tool := newMigrateCodeTool(t)

	// Empty sourceCode fails the required validator on
	// models.MigrationRequest.SourceCode — surfaced as a tool result
	// error rather than a transport error. sourceCode key is present
	// (else the transport-layer schema required check would trigger
	// upstream of the tool), but the value is empty.
	req := buildRequest(map[string]any{
		"sourceCode": "",
		"type":       "page",
	})

	result, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true for invalid request")
	}
	if !strings.Contains(extractText(t, result), "Migration failed") {
		t.Errorf("expected 'Migration failed' prefix, got %q", extractText(t, result))
	}
}
