package workflows

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// fakeRunner is a hand-rolled domain.TestRunner double. Keeps parity
// with the repo's existing test style — no gomock, no testify/mock.
type fakeRunner struct {
	called   bool
	lastCtx  context.Context
	lastReq  models.TestExecutionRequest
	result   *models.TestExecutionResult
	err      error
	blockFor time.Duration // when set, Run blocks until ctx expires or blockFor passes
}

func (f *fakeRunner) Run(ctx context.Context, req models.TestExecutionRequest) (*models.TestExecutionResult, error) {
	f.called = true
	f.lastCtx = ctx
	f.lastReq = req

	if f.blockFor > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(f.blockFor):
		}
	}

	if f.err != nil {
		return nil, f.err
	}
	if f.result != nil {
		return f.result, nil
	}
	// Default happy result so timeout/wiring tests that only care about
	// ctx propagation don't have to spell one out.
	return &models.TestExecutionResult{
		Success: true, Summary: "Tests Passed", Duration: "0s",
	}, nil
}

func newWorkflow(runner *fakeRunner) *RunTestsWorkflow {
	return NewRunTestsWorkflow(logging.NewLogger(&bytes.Buffer{}), runner)
}

func buildRequest(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func TestRunTestsWorkflow_Metadata(t *testing.T) {
	w := newWorkflow(&fakeRunner{})

	if w.Name() != "run_tests" {
		t.Errorf("Name mismatch: %q", w.Name())
	}
	if w.Description() == "" {
		t.Error("Description should not be empty")
	}

	schema := w.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema.Type: got %q", schema.Type)
	}
	if len(schema.Required) != 1 || schema.Required[0] != "projectPath" {
		t.Errorf("expected Required=[projectPath], got %v", schema.Required)
	}
	if _, ok := schema.Properties["timeoutSeconds"]; !ok {
		t.Error("expected timeoutSeconds property")
	}

	if w.OutputSchema() == nil {
		t.Error("expected non-nil OutputSchema")
	}
}

func TestRunTestsWorkflow_Run_Success(t *testing.T) {
	runner := &fakeRunner{
		result: &models.TestExecutionResult{
			Success:  true,
			Output:   "PASS",
			Duration: "1s",
			Summary:  "Tests Passed",
		},
	}
	w := newWorkflow(runner)

	req := buildRequest(map[string]any{
		"projectPath": "/tmp/project",
		"command":     "echo ok",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if result.IsError {
		t.Errorf("expected success result, got IsError=true; content=%+v", result.Content)
	}
	if !runner.called {
		t.Error("expected runner.Run to be called")
	}
	if runner.lastReq.ProjectPath != "/tmp/project" {
		t.Errorf("projectPath not forwarded: %+v", runner.lastReq)
	}
	if runner.lastReq.Command != "echo ok" {
		t.Errorf("command not forwarded: %+v", runner.lastReq)
	}

	// The result content should be the marshaled TestExecutionResult.
	got := extractText(t, result)
	var round models.TestExecutionResult
	if err := json.Unmarshal([]byte(got), &round); err != nil {
		t.Fatalf("result text is not valid JSON: %v (%q)", err, got)
	}
	if !round.Success || round.Summary != "Tests Passed" {
		t.Errorf("round-tripped result mismatch: %+v", round)
	}
}

func TestRunTestsWorkflow_Run_MissingProjectPath(t *testing.T) {
	runner := &fakeRunner{}
	w := newWorkflow(runner)

	req := buildRequest(map[string]any{
		"command": "echo ok",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when projectPath is missing")
	}
	if runner.called {
		t.Error("runner should not be called when validation fails")
	}
	if !strings.Contains(extractText(t, result), "projectPath is required") {
		t.Errorf("expected projectPath error message, got %q", extractText(t, result))
	}
}

func TestRunTestsWorkflow_Run_RunnerErrorSurfacedAsToolResult(t *testing.T) {
	runner := &fakeRunner{
		err: errors.New("exec.CommandContext: no such file"),
	}
	w := newWorkflow(runner)

	req := buildRequest(map[string]any{
		"projectPath": "/tmp/x",
	})

	result, err := w.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run should not surface transport err, got: %v", err)
	}
	if !result.IsError {
		t.Error("expected IsError=true when runner errors")
	}
	if !strings.Contains(extractText(t, result), "Test execution failed") {
		t.Errorf("expected prefix in error message, got %q", extractText(t, result))
	}
}

func TestRunTestsWorkflow_Run_TimeoutOverrideFloat64(t *testing.T) {
	// json.Unmarshal defaults to float64 for numbers; verify that path.
	runner := &fakeRunner{
		blockFor: 5 * time.Second, // never fires — ctx must cancel first
	}
	w := newWorkflow(runner)

	req := buildRequest(map[string]any{
		"projectPath":    "/tmp/x",
		"timeoutSeconds": float64(0), // sentinel: 0 disables the override
	})
	// Ignore the sentinel — set a real one that triggers the branch.
	req.Params.Arguments.(map[string]any)["timeoutSeconds"] = float64(1)

	start := time.Now()
	// blockFor = 5s but timeoutSeconds override = 1s → runner ctx should
	// expire, blockFor loop returns ctx.Err.
	_, err := w.Run(context.Background(), req)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Run should surface via tool result, got: %v", err)
	}
	if elapsed > 3*time.Second {
		t.Errorf("timeoutSeconds override not applied (elapsed %s)", elapsed)
	}

	if _, ok := runner.lastCtx.Deadline(); !ok {
		t.Error("expected runner ctx to carry a deadline from the override")
	}
}

func TestRunTestsWorkflow_Run_TimeoutOverrideInt(t *testing.T) {
	// Some clients send integers as int, not float64. Verify both work.
	runner := &fakeRunner{}
	w := newWorkflow(runner)

	req := buildRequest(map[string]any{
		"projectPath":    "/tmp/x",
		"timeoutSeconds": 30, // int
	})

	if _, err := w.Run(context.Background(), req); err != nil {
		t.Fatalf("Run err: %v", err)
	}
	deadline, ok := runner.lastCtx.Deadline()
	if !ok {
		t.Fatal("expected runner ctx to carry a deadline")
	}
	remaining := time.Until(deadline)
	if remaining < 25*time.Second || remaining > 31*time.Second {
		t.Errorf("expected ~30s deadline, remaining %s", remaining)
	}
}

func TestRunTestsWorkflow_Run_TimeoutOverrideZeroIgnored(t *testing.T) {
	// Zero / negative / non-numeric all mean "no override, use adapter
	// default". Runner's ctx should have no deadline in that case (the
	// adapter's ensureDeadline will apply its own 300s ceiling — but
	// that's ITS job, not the workflow's).
	runner := &fakeRunner{}
	w := newWorkflow(runner)

	req := buildRequest(map[string]any{
		"projectPath":    "/tmp/x",
		"timeoutSeconds": float64(0),
	})

	if _, err := w.Run(context.Background(), req); err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if _, ok := runner.lastCtx.Deadline(); ok {
		t.Error("expected NO deadline when timeoutSeconds is 0")
	}
}

func TestExtractTimeoutSeconds(t *testing.T) {
	cases := []struct {
		name    string
		args    map[string]any
		wantVal int
		wantOK  bool
	}{
		{"absent", map[string]any{}, 0, false},
		{"nil value", map[string]any{"timeoutSeconds": nil}, 0, false},
		{"positive float", map[string]any{"timeoutSeconds": float64(45)}, 45, true},
		{"positive int", map[string]any{"timeoutSeconds": 60}, 60, true},
		{"zero float", map[string]any{"timeoutSeconds": float64(0)}, 0, false},
		{"zero int", map[string]any{"timeoutSeconds": 0}, 0, false},
		{"negative float", map[string]any{"timeoutSeconds": float64(-1)}, 0, false},
		{"negative int", map[string]any{"timeoutSeconds": -5}, 0, false},
		{"non-numeric string", map[string]any{"timeoutSeconds": "300"}, 0, false},
		{"non-numeric bool", map[string]any{"timeoutSeconds": true}, 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotVal, gotOK := extractTimeoutSeconds(tc.args)
			if gotVal != tc.wantVal || gotOK != tc.wantOK {
				t.Errorf("got (%d, %v) want (%d, %v)", gotVal, gotOK, tc.wantVal, tc.wantOK)
			}
		})
	}
}

// extractText pulls the text field off the first content item of a
// CallToolResult so assertions can read the workflow's serialized
// response without reaching into the mcp-go internals in every test.
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
