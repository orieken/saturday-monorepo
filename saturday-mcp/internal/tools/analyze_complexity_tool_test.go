package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
)

func TestAnalyzeComplexityTool_Metadata(t *testing.T) {
	tool := NewAnalyzeComplexityTool(silentLogger(), analyzers.NewComplexityAnalyzer())

	if tool.Name() != "analyze_complexity" {
		t.Errorf("Name(): got %q, want %q", tool.Name(), "analyze_complexity")
	}

	if tool.Description() == "" {
		t.Error("Description() should not be empty")
	}

	if tool.OutputSchema() == nil {
		t.Error("OutputSchema() should not be nil")
	}
}

func TestAnalyzeComplexityTool_Execute_Success(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "sample.go")
	body := `package sample

func Simple() {
	_ = 1
}

func Complex(x int) int {
	if x > 0 {
		if x > 10 {
			if x > 20 {
				return 3
			}
			return 2
		}
		return 1
	}
	return 0
}
`
	if err := os.WriteFile(filePath, []byte(body), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	tool := NewAnalyzeComplexityTool(silentLogger(), analyzers.NewComplexityAnalyzer())
	req := buildRequest(map[string]any{
		"projectPath":   filePath,
		"maxComplexity": float64(2),
		"maxLines":      float64(30),
	})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	text := extractText(t, res)
	var out ComplexityAnalysisResult
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v (%s)", err, text)
	}

	if !out.Success {
		t.Error("expected Success=true")
	}
	if out.ViolationsCount == 0 {
		t.Error("expected violations count > 0")
	}
}

func TestAnalyzeComplexityTool_Execute_MissingPath(t *testing.T) {
	tool := NewAnalyzeComplexityTool(silentLogger(), analyzers.NewComplexityAnalyzer())
	req := buildRequest(map[string]any{})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if !res.IsError {
		t.Error("expected IsError=true when projectPath is missing")
	}
}
