// Package testrunner contains the os/exec-backed adapter that implements
// domain.TestRunner. It is the only place in the codebase that touches
// exec.Command for test execution — every other layer depends on the
// domain.TestRunner interface.
package testrunner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/orieken/saturday-mcp/internal/logging"
	"github.com/orieken/saturday-mcp/internal/models"
)

// DefaultTimeout bounds any Run call whose context has no deadline. The
// 300s ceiling is intentionally generous — Playwright suites routinely
// take minutes — but no unbounded exec.Command survives past this
// adapter. Callers that want a tighter deadline should pass a
// context.WithTimeout of their own; Phase H op 20 wires tool-input
// overrides on top of this default at the workflow layer.
const DefaultTimeout = 300 * time.Second

// ExecRunner wraps os/exec.CommandContext to satisfy domain.TestRunner.
// The struct is intentionally tiny — behavior is a single Run method that
// mirrors the pre-refactor executor.TestExecutor.Run body, minus the
// direct exec.Command call (now CommandContext) and the missing timeout
// (now applied via ensureDeadline).
type ExecRunner struct {
	logger *logging.Logger
}

// NewExecRunner wires the adapter with its logger. Constructor is
// deliberately parameterless beyond logging — timeout policy lives on the
// context passed to Run so tests can override it.
func NewExecRunner(logger *logging.Logger) *ExecRunner {
	return &ExecRunner{logger: logger}
}

// Run executes the test command described by req under ctx. If ctx has no
// deadline the adapter applies DefaultTimeout so no unbounded process
// escapes the boundary. Output is combined stdout+stderr, preserved
// verbatim from the pre-refactor executor.TestExecutor shape so the MCP
// response contract stays byte-identical.
func (r *ExecRunner) Run(ctx context.Context, req models.TestExecutionRequest) (*models.TestExecutionResult, error) {
	r.logger.Info("Starting test execution", "path", req.ProjectPath, "cmd", req.Command)

	ctx, cancel := ensureDeadline(ctx, DefaultTimeout)
	defer cancel()

	start := time.Now()

	cmdStr := req.Command
	if cmdStr == "" {
		cmdStr = "npx playwright test"
	}

	if req.Filter != "" {
		cmdStr += fmt.Sprintf(" -g '%s'", req.Filter)
	}

	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	head := parts[0]
	args := parts[1:]

	cmd := exec.CommandContext(ctx, head, args...)
	cmd.Dir = req.ProjectPath

	cmd.Env = os.Environ()
	for k, v := range req.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start).String()

	output := stdout.String() + "\n" + stderr.String()
	success := err == nil

	summary := "Tests executed"
	if success {
		summary = "Tests Passed"
	} else {
		summary = "Tests Failed"
	}

	r.logger.Info("Test execution complete", "success", success, "duration", duration)

	return &models.TestExecutionResult{
		Success:  success,
		Output:   output,
		Duration: duration,
		Summary:  summary,
	}, nil
}

// ensureDeadline returns ctx unchanged if it already carries a deadline;
// otherwise it wraps ctx with the supplied default. The returned cancel
// is always safe to call — the noop branch returns a cancel that does
// nothing beyond satisfying the defer contract.
func ensureDeadline(ctx context.Context, def time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, def)
}
