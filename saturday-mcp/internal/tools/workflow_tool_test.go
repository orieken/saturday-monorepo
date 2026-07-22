package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
)

type fakeWorkflow struct {
	name         string
	description  string
	inputSchema  mcp.ToolInputSchema
	outputSchema *jsonschema.Schema
	result       *mcp.CallToolResult
	err          error

	called  bool
	lastCtx context.Context
	lastReq mcp.CallToolRequest
}

func (f *fakeWorkflow) Name() string { return f.name }

func (f *fakeWorkflow) Description() string { return f.description }

func (f *fakeWorkflow) InputSchema() mcp.ToolInputSchema { return f.inputSchema }

func (f *fakeWorkflow) OutputSchema() *jsonschema.Schema { return f.outputSchema }

func (f *fakeWorkflow) Run(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	f.called = true
	f.lastCtx = ctx
	f.lastReq = request
	return f.result, f.err
}

func TestWorkflowTool_Metadata(t *testing.T) {
	expectedSchema := &jsonschema.Schema{Title: "TestOutput"}
	expectedInput := mcp.ToolInputSchema{Type: "object"}

	fw := &fakeWorkflow{
		name:         "test_workflow",
		description:  "A test workflow description",
		inputSchema:  expectedInput,
		outputSchema: expectedSchema,
	}

	tool := NewWorkflowTool(fw)

	if got := tool.Name(); got != "test_workflow" {
		t.Errorf("Name(): got %q, want %q", got, "test_workflow")
	}

	if got := tool.Description(); got != "A test workflow description" {
		t.Errorf("Description(): got %q, want %q", got, "A test workflow description")
	}

	if got := tool.InputSchema(); got.Type != expectedInput.Type {
		t.Errorf("InputSchema(): got %v, want %v", got, expectedInput)
	}

	if got := tool.OutputSchema(); got != expectedSchema {
		t.Errorf("OutputSchema(): got %v, want %v", got, expectedSchema)
	}
}

func TestWorkflowTool_Execute_Success(t *testing.T) {
	expectedResult := mcp.NewToolResultText("workflow output")
	fw := &fakeWorkflow{
		result: expectedResult,
	}

	tool := NewWorkflowTool(fw)
	ctx := context.WithValue(context.Background(), "testKey", "testValue")
	req := buildRequest(map[string]any{"arg1": "val1"})

	res, err := tool.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Execute() unexpected error: %v", err)
	}

	if !fw.called {
		t.Error("Execute() did not call workflow.Run")
	}

	if fw.lastCtx != ctx {
		t.Error("Execute() did not pass context to workflow.Run")
	}

	argsMap, ok := fw.lastReq.Params.Arguments.(map[string]any)
	if !ok {
		t.Fatalf("Execute() Arguments not a map[string]any: %T", fw.lastReq.Params.Arguments)
	}
	if val, ok := argsMap["arg1"].(string); !ok || val != "val1" {
		t.Errorf("Execute() did not pass request correctly, got: %+v", fw.lastReq)
	}

	if res != expectedResult {
		t.Errorf("Execute(): got %v, want %v", res, expectedResult)
	}

	gotText := extractText(t, res)
	if gotText != "workflow output" {
		t.Errorf("extractText(): got %q, want %q", gotText, "workflow output")
	}
}

func TestWorkflowTool_Execute_Error(t *testing.T) {
	expectedErr := errors.New("pipeline error")
	fw := &fakeWorkflow{
		err: expectedErr,
	}

	tool := NewWorkflowTool(fw)
	ctx := context.Background()
	req := buildRequest(map[string]any{})

	res, err := tool.Execute(ctx, req)
	if !errors.Is(err, expectedErr) {
		t.Errorf("Execute(): got error %v, want %v", err, expectedErr)
	}

	if res != nil {
		t.Errorf("Execute(): got result %v, want nil on error", res)
	}

	if !fw.called {
		t.Error("Execute() did not call workflow.Run")
	}
}
