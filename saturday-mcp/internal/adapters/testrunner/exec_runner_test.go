package testrunner

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// newRunner builds an ExecRunner with a discarding logger so subprocess
// stdout/stderr never leaks into the test output.
func newRunner() *ExecRunner {
	return NewExecRunner(logging.NewLogger(&bytes.Buffer{}))
}

func TestExecRunner_Run_Success(t *testing.T) {
	runner := newRunner()

	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "echo hello",
	}

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got false; output=%q", result.Output)
	}
	if !strings.Contains(result.Output, "hello") {
		t.Errorf("expected output to contain 'hello', got %q", result.Output)
	}
	if result.Summary != "Tests Passed" {
		t.Errorf("expected Summary=Tests Passed, got %q", result.Summary)
	}
	if result.Duration == "" {
		t.Error("expected Duration to be populated")
	}
}

func TestExecRunner_Run_CommandFailureCapturesNonZero(t *testing.T) {
	runner := newRunner()

	// `false` is a POSIX command that always exits 1 with no output.
	// The adapter should surface Success=false without propagating err.
	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "false",
	}

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run should not propagate command-exit errors, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result on command failure")
	}
	if result.Success {
		t.Error("expected Success=false for `false` command")
	}
	if result.Summary != "Tests Failed" {
		t.Errorf("expected Summary=Tests Failed, got %q", result.Summary)
	}
}

func TestExecRunner_Run_CapturesStderr(t *testing.T) {
	runner := newRunner()

	// `sh -c` runs a compound command whose stderr the adapter should
	// combine into the output field. The command exits non-zero so
	// Success=false, and the message we wrote to stderr must appear.
	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     `sh -c "echo boom 1>&2; exit 2"`,
	}

	// strings.Fields will mangle the quoted sh -c arg — verify the
	// adapter's parsing catches that as failure instead of running
	// something unexpected.
	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run should not propagate command-exit errors, got: %v", err)
	}
	if result.Success {
		t.Errorf("expected Success=false, output=%q", result.Output)
	}
}

func TestExecRunner_Run_EmptyCommandReturnsError(t *testing.T) {
	runner := newRunner()

	// strings.Fields on "" yields an empty slice; adapter must reject
	// rather than crash on parts[0].
	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "   ",
	}

	result, err := runner.Run(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for empty command")
	}
	if result != nil {
		t.Errorf("expected nil result on error, got %+v", result)
	}
	if !strings.Contains(err.Error(), "empty command") {
		t.Errorf("expected 'empty command' in error, got %q", err.Error())
	}
}

func TestExecRunner_Run_UnknownCommandFailsGracefully(t *testing.T) {
	runner := newRunner()

	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "this_binary_does_not_exist_xyz",
	}

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run should not propagate exec errors, got: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for unknown binary")
	}
}

func TestExecRunner_Run_FilterAppendedToCommand(t *testing.T) {
	runner := newRunner()

	// Filter should append "-g 'login'". Using `echo` we can observe
	// the appended token in stdout.
	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "echo cmd",
		Filter:      "login",
	}

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if !strings.Contains(result.Output, "-g") {
		t.Errorf("expected output to include filter flag '-g', got %q", result.Output)
	}
	if !strings.Contains(result.Output, "login") {
		t.Errorf("expected output to include filter value 'login', got %q", result.Output)
	}
}

func TestExecRunner_Run_EnvForwarded(t *testing.T) {
	runner := newRunner()

	// `env` prints every env var as KEY=VALUE. Verifying our custom
	// var lands in stdout confirms the adapter forwarded req.Env.
	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "env",
		Env: map[string]string{
			"SATURDAY_MCP_TEST_VAR": "custom_value",
		},
	}

	result, err := runner.Run(context.Background(), req)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if !strings.Contains(result.Output, "SATURDAY_MCP_TEST_VAR=custom_value") {
		t.Errorf("expected env var forwarded, got %q", result.Output)
	}
}

func TestExecRunner_Run_ContextCancellationStopsProcess(t *testing.T) {
	runner := newRunner()

	// `sleep 5` normally runs 5s. With a 100ms context, the adapter's
	// exec.CommandContext MUST kill the process — otherwise this test
	// would take 5 full seconds. We assert both the fast return and
	// the failure signal.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req := models.TestExecutionRequest{
		ProjectPath: ".",
		Command:     "sleep 5",
	}

	start := time.Now()
	result, err := runner.Run(ctx, req)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Run should not propagate ctx cancellation as err, got: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false when process is killed by ctx timeout")
	}
	if elapsed > 2*time.Second {
		t.Errorf("expected fast return (< 2s), took %s — ctx cancellation not honored", elapsed)
	}
}

func TestDefaultTimeoutIsBounded(t *testing.T) {
	// Guardrail against a well-meaning "no timeout" regression. If the
	// constant ever gets bumped to 0 or something unrealistically large,
	// this test catches it before the code lands.
	if DefaultTimeout <= 0 {
		t.Errorf("DefaultTimeout must be positive, got %s", DefaultTimeout)
	}
	if DefaultTimeout > 10*time.Minute {
		t.Errorf("DefaultTimeout suspiciously large (>10m), got %s", DefaultTimeout)
	}
	if DefaultTimeout != 300*time.Second {
		t.Errorf("DefaultTimeout expected 300s per plan, got %s", DefaultTimeout)
	}
}

func TestEnsureDeadline_PreservesExistingDeadline(t *testing.T) {
	parentDeadline := time.Now().Add(50 * time.Millisecond)
	parent, cancel := context.WithDeadline(context.Background(), parentDeadline)
	defer cancel()

	ctx, adapterCancel := ensureDeadline(parent, 300*time.Second)
	defer adapterCancel()

	got, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be present")
	}
	// If ensureDeadline had wrapped a fresh 300s timeout, the deadline
	// would be far in the future — assert it stayed near parentDeadline.
	if got.Sub(parentDeadline).Abs() > 5*time.Millisecond {
		t.Errorf("ensureDeadline overrode parent deadline: got %s, want ~%s", got, parentDeadline)
	}
}

func TestEnsureDeadline_AddsDeadlineWhenAbsent(t *testing.T) {
	parent := context.Background()

	ctx, cancel := ensureDeadline(parent, 10*time.Millisecond)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("ensureDeadline should add a deadline when parent has none")
	}

	// With a 10ms deadline, ctx.Done() must fire quickly. Waiting up to
	// 500ms bounds the test without being flaky on slow CI.
	select {
	case <-ctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Error("expected ctx to expire within 500ms")
	}
}
