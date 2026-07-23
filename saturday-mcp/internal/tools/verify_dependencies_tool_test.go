package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/orieken/saturday-mcp/internal/analyzers"
)

// Tests for VerifyDependenciesTool follow the pattern established by
// check_ubiquitous_language_tool_test.go: metadata assertions, an
// Execute success path against a t.TempDir() fixture, layer-classification
// coverage across all four layers, TS + Go coverage, dir-skipping
// coverage, and error paths. Coverage target ≥ 85% per mcp-expand M1
// Op 5 spec. Shared helpers (writeFile, silentLogger, buildRequest,
// extractText) come from testfixtures_test.go.

func TestVerifyDependenciesTool_Metadata(t *testing.T) {
	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())

	if tool.Name() != "verify_dependencies" {
		t.Errorf("Name(): got %q, want %q", tool.Name(), "verify_dependencies")
	}
	if tool.Description() == "" {
		t.Error("Description() should not be empty")
	}
	if tool.OutputSchema() == nil {
		t.Error("OutputSchema() should not be nil")
	}
	schema := tool.InputSchema()
	if schema.Type != "object" {
		t.Errorf("InputSchema Type: got %q, want object", schema.Type)
	}
	if len(schema.Required) != 1 || schema.Required[0] != "projectPath" {
		t.Errorf("Required: got %v, want [projectPath]", schema.Required)
	}
	if _, ok := schema.Properties["projectPath"]; !ok {
		t.Error("InputSchema is missing projectPath property")
	}
}

func TestVerifyDependenciesTool_Execute_DomainImportsAdaptersIsViolation(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package domain\n\n"+
			"import (\n"+
			"\t\"fmt\"\n"+
			"\t\"github.com/example/project/adapters/db\"\n"+
			")\n\n"+
			"func New() { _ = fmt.Println; _ = db.Open }\n")
	writeFile(t, filepath.Join(tmp, "adapters/db/postgres.go"),
		"package db\n\nfunc Open() {}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	req := buildRequest(map[string]any{"projectPath": tmp})

	res, err := tool.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("Execute returned error result: %s", extractText(t, res))
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if !out.Success {
		t.Error("expected Success=true")
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	v := out.Violations[0]
	if v.FromLayer != "domain" {
		t.Errorf("FromLayer: got %q, want domain", v.FromLayer)
	}
	if v.ToLayer != "adapters" {
		t.Errorf("ToLayer: got %q, want adapters", v.ToLayer)
	}
	if v.ImportPath != "github.com/example/project/adapters/db" {
		t.Errorf("ImportPath: got %q", v.ImportPath)
	}
	if v.LineNumber != 5 {
		t.Errorf("LineNumber: got %d, want 5", v.LineNumber)
	}
}

func TestVerifyDependenciesTool_Execute_AdaptersImportsDomainIsAllowed(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "adapters/db/postgres.go"),
		"package db\n\n"+
			"import \"github.com/example/project/domain/user\"\n\n"+
			"func New() { _ = user.New }\n")
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package user\n\nfunc New() {}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Summary != "No dependency boundary violations found" {
		t.Errorf("Summary: got %q, want clean summary", out.Summary)
	}
}

func TestVerifyDependenciesTool_Execute_UsecasesImportsFrameworksIsViolation(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "usecases/create_user.go"),
		"package usecases\n\n"+
			"import \"github.com/example/project/frameworks/http\"\n\n"+
			"func Do() { _ = http.Listen }\n")
	writeFile(t, filepath.Join(tmp, "frameworks/http/server.go"),
		"package http\n\nfunc Listen() {}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	v := out.Violations[0]
	if v.FromLayer != "usecases" || v.ToLayer != "frameworks" {
		t.Errorf("layers: got %q → %q, want usecases → frameworks", v.FromLayer, v.ToLayer)
	}
}

func TestVerifyDependenciesTool_Execute_SameLayerImportIsAllowed(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package domain\n\n"+
			"import \"github.com/example/project/domain/other\"\n\n"+
			"var _ = other.X\n")
	writeFile(t, filepath.Join(tmp, "domain/other/other.go"),
		"package other\n\nvar X = 1\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0 for same-layer import", out.ViolationsCount)
	}
}

func TestVerifyDependenciesTool_Execute_UnclassifiedFilesAndImportsAreSkipped(t *testing.T) {
	tmp := t.TempDir()
	// A file whose path has no layer segment should not be classified,
	// so it produces no violations even when it imports across layers.
	writeFile(t, filepath.Join(tmp, "internal/helpers/util.go"),
		"package helpers\n\nimport \"github.com/example/project/frameworks/x\"\nvar _ = x.Y\n")
	writeFile(t, filepath.Join(tmp, "frameworks/x/x.go"), "package x\nvar Y = 1\n")
	// A classified file importing an unclassified external dep also
	// produces no violation — the import target has no layer.
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package domain\nimport \"github.com/spf13/cobra\"\nvar _ = cobra.Command{}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0; violations=%+v", out.ViolationsCount, out.Violations)
	}
}

func TestVerifyDependenciesTool_Execute_TypeScriptRelativeImportViolation(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "src/domain/user.ts"),
		"import { Db } from '../adapters/db';\n\nexport class User { db: Db }\n")
	writeFile(t, filepath.Join(tmp, "src/adapters/db.ts"),
		"export class Db {}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	v := out.Violations[0]
	if v.FromLayer != "domain" || v.ToLayer != "adapters" {
		t.Errorf("layers: got %q → %q, want domain → adapters", v.FromLayer, v.ToLayer)
	}
	if v.ImportPath != "../adapters/db" {
		t.Errorf("ImportPath: got %q, want ../adapters/db", v.ImportPath)
	}
	if v.LineNumber != 1 {
		t.Errorf("LineNumber: got %d, want 1", v.LineNumber)
	}
}

func TestVerifyDependenciesTool_Execute_TypeScriptSideEffectImport(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "src/domain/bootstrap.tsx"),
		"import '../frameworks/init';\n\nexport const X = 1;\n")
	writeFile(t, filepath.Join(tmp, "src/frameworks/init.ts"), "export {};\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Violations[0].ToLayer != "frameworks" {
		t.Errorf("ToLayer: got %q, want frameworks", out.Violations[0].ToLayer)
	}
}

func TestVerifyDependenciesTool_Execute_EntitiesAndInfrastructureAliases(t *testing.T) {
	tmp := t.TempDir()
	// "entities" is a domain alias; "infrastructure" is an adapters alias.
	writeFile(t, filepath.Join(tmp, "src/entities/user.ts"),
		"import { Db } from '../infrastructure/db';\nexport class User {}\n")
	writeFile(t, filepath.Join(tmp, "src/infrastructure/db.ts"), "export class Db {}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Violations[0].FromLayer != "domain" || out.Violations[0].ToLayer != "adapters" {
		t.Errorf("aliases not resolved: got %q → %q, want domain → adapters",
			out.Violations[0].FromLayer, out.Violations[0].ToLayer)
	}
}

func TestVerifyDependenciesTool_Execute_GoSingleLineImportAlias(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package domain\n\n"+
			"import db \"github.com/example/project/adapters/db\"\n\n"+
			"var _ = db.Open\n")
	writeFile(t, filepath.Join(tmp, "adapters/db/db.go"), "package db\nfunc Open() {}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1 (aliased import); violations=%+v", out.ViolationsCount, out.Violations)
	}
}

func TestVerifyDependenciesTool_Execute_SkipsExcludedDirsAndUnknownExtensions(t *testing.T) {
	tmp := t.TempDir()
	// One real source file that should be scanned.
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package domain\nimport \"github.com/example/project/adapters/db\"\nvar _ = db.X\n")
	writeFile(t, filepath.Join(tmp, "adapters/db/db.go"), "package db\nvar X = 1\n")
	// Excluded directories that should NOT be scanned even though they
	// contain the exact same violation shape.
	writeFile(t, filepath.Join(tmp, "node_modules/pkg/domain/x.go"),
		"package domain\nimport \"github.com/example/project/adapters/db\"\n")
	writeFile(t, filepath.Join(tmp, "vendor/x/domain/y.go"),
		"package domain\nimport \"github.com/example/project/adapters/db\"\n")
	writeFile(t, filepath.Join(tmp, "dist/domain/out.go"),
		"package domain\nimport \"github.com/example/project/adapters/db\"\n")
	writeFile(t, filepath.Join(tmp, ".git/domain/hooks.go"),
		"package domain\nimport \"github.com/example/project/adapters/db\"\n")
	// Non-source extensions should be ignored regardless of layer path.
	writeFile(t, filepath.Join(tmp, "domain/README.md"), "domain adapters\n")
	writeFile(t, filepath.Join(tmp, "domain/config.json"), "{}\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Errorf("ViolationsCount: got %d, want 1 (only domain/user.go should be scanned); violations=%+v", out.ViolationsCount, out.Violations)
	}
}

func TestVerifyDependenciesTool_Execute_CleanProjectHasNoViolations(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "domain/user.go"), "package domain\n\ntype User struct{}\n")
	writeFile(t, filepath.Join(tmp, "usecases/create_user.go"),
		"package usecases\nimport \"github.com/example/project/domain\"\n"+
			"var _ = domain.User{}\n")
	writeFile(t, filepath.Join(tmp, "adapters/db/user_repo.go"),
		"package db\nimport \"github.com/example/project/usecases\"\nvar _ = usecases.X\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 0 {
		t.Errorf("ViolationsCount: got %d, want 0 for clean project; violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Summary != "No dependency boundary violations found" {
		t.Errorf("Summary: got %q, want clean summary", out.Summary)
	}
}

func TestVerifyDependenciesTool_Execute_GoImportBlockCommentsIgnored(t *testing.T) {
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "domain/user.go"),
		"package domain\n\n"+
			"import (\n"+
			"\t// standard library\n"+
			"\t\"fmt\"\n"+
			"\n"+
			"\t// project imports\n"+
			"\t\"github.com/example/project/adapters/db\"\n"+
			")\n\n"+
			"var _ = fmt.Println\nvar _ = db.X\n")
	writeFile(t, filepath.Join(tmp, "adapters/db/db.go"), "package db\nvar X = 1\n")

	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{"projectPath": tmp}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	var out DependencyVerificationResult
	if err := json.Unmarshal([]byte(extractText(t, res)), &out); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.ViolationsCount != 1 {
		t.Fatalf("ViolationsCount: got %d, want 1 (only db import is a violation); violations=%+v", out.ViolationsCount, out.Violations)
	}
	if out.Violations[0].ImportPath != "github.com/example/project/adapters/db" {
		t.Errorf("ImportPath: got %q", out.Violations[0].ImportPath)
	}
}

func TestVerifyDependenciesTool_Execute_MissingProjectPath(t *testing.T) {
	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true when projectPath is missing")
	}
}

func TestVerifyDependenciesTool_Execute_NonexistentProjectPath(t *testing.T) {
	tool := NewVerifyDependenciesTool(silentLogger(), analyzers.NewDependencyBoundaryAnalyzer())
	res, err := tool.Execute(context.Background(), buildRequest(map[string]any{
		"projectPath": filepath.Join(t.TempDir(), "does-not-exist"),
	}))
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true for nonexistent project path")
	}
}
