package tools

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/templates"
	"github.com/orieken/saturday-mcp/internal/validators"
)

// This file holds the small helpers every *_tool_test.go in this package
// re-uses. Kept in one place so the pattern established by Phase J Part 1
// (see internal/workflows/*_test.go) stays literal here — buildRequest,
// extractText, silent logger, template setup — rather than being copy-
// pasted into fourteen files. See mcp-add-plan Phase J.

// buildRequest turns an argument map into an mcp.CallToolRequest the way
// the mcp-go SDK's JSON transport would, so tests don't have to spell out
// the Params.Arguments plumbing on every call.
func buildRequest(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

// extractText pulls the text off the first Content entry of a tool
// result. Every tool in this package encodes its response as a single
// mcp.TextContent, so assertions can read it back without descending into
// the SDK's typed content union.
func extractText(t *testing.T, r *mcp.CallToolResult) string {
	t.Helper()
	if r == nil {
		t.Fatal("nil result")
	}
	if len(r.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := r.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", r.Content[0])
	}
	return tc.Text
}

// silentLogger returns a logger that discards its output. Every tool
// takes a *logging.Logger; the tests care about the return value, not
// the log spew, so we pipe it into a throwaway buffer.
func silentLogger() *logging.Logger {
	return logging.NewLogger(&bytes.Buffer{})
}

// newProcessor wires up the real template processor with every Saturday
// scaffolding template loaded from disk. The generator tools genuinely
// exercise these templates end-to-end — swapping in a mock processor
// would hide the code the tool actually produces, which is exactly the
// assertion path the tests need to see (per shared/rules/testing-
// conventions.md "moist tests" guidance).
func newProcessor(t *testing.T) *templates.Processor {
	t.Helper()

	registry := templates.NewRegistry()
	loader := templates.NewLoader(registry)
	cache := templates.NewCache(5 * time.Minute)
	processor := templates.NewProcessor(registry, cache)

	if err := loader.LoadAll(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}
	return processor
}

// newValidator returns the shared validators.Validator every generator
// requires. Extracted only so tool tests don't repeat the import path.
func newValidator() *validators.Validator {
	return validators.NewValidator()
}

// captureFileWriter is a domain.FileWriter double that records the last
// WriteFile call for assertion and lets tests inject a return error.
// Used by generate_documentation which takes the low-level one-shot
// writer interface (see internal/domain/filesystem.go).
type captureFileWriter struct {
	called   bool
	lastPath string
	lastBody []byte
	err      error
}

func (c *captureFileWriter) WriteFile(path string, content []byte) error {
	c.called = true
	c.lastPath = path
	// Copy — callers may mutate their buffer after WriteFile returns and
	// we want the recorded body to reflect what was passed in.
	c.lastBody = append([]byte(nil), content...)
	return c.err
}

// writeTSFile is a small helper for the analyzer tests that need a
// minimal fake Saturday project on disk. It writes a .ts file at a
// project-relative path, creating parent directories as needed.
func writeTSFile(t *testing.T, root, relPath, body string) {
	t.Helper()
	full := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(full, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", relPath, err)
	}
}

// writeFile writes body to an absolute path, creating parent directories
// as needed. Shared across every tool test that stands up a fake
// project on disk with more than a single flat file (the
// check_accessibility and check_ubiquitous_language tests, plus any
// future analyzer tool tests). Promoted here from
// check_accessibility_tool_test.go per mcp-expand M1 Op 4 handoff notes
// — cleaner than duplicating the helper into every new test file.
func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
