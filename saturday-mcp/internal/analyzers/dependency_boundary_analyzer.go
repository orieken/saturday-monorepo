// Package analyzers - dependency_boundary_analyzer.go
//
// DependencyBoundaryAnalyzer is the concrete analyzer behind the
// verify_dependencies MCP tool (mcp-expand M1 Op 5). It walks a project
// tree, classifies each Go / TypeScript source file into a Clean
// Architecture layer based on its path segments, parses its import
// statements, and flags every import that crosses a layer boundary in
// the wrong direction (inner layer importing an outer layer).
//
// Layer ordering (inner → outer, per shared/rules/architecture-
// guardrails.md #1): domain < usecases < adapters < frameworks.
// Anything not classified by a recognized path segment is layer
// "unknown" and skipped — neither flagged nor followed. Import paths are
// classified by the same segment-scan; an external import like
// `github.com/spf13/cobra` has no layer segment and is skipped, matching
// the "intra-project boundary check" scope this tool commits to.
//
// Parsing is deliberately regex-based, not AST-based, matching the style
// of the neighbouring analyzers (accessibility_analyzer.go,
// ubiquitous_language_analyzer.go). The Go `go/parser` package would give
// richer accuracy but pulling it in for two import shapes would be
// heavier than the tool needs — Op 5 spec calls for the regex approach.
// A future enhancement can swap the parser out behind readImports without
// touching the layer-classification or boundary-rule logic.
package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DependencyBoundaryViolation captures a single Clean Architecture layer
// violation. Field shape mirrors internal/tools/responses.go's
// DependencyBoundaryViolation so the tool layer can marshal analyzer
// output directly without re-mapping fields.
type DependencyBoundaryViolation struct {
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	FromLayer  string `json:"fromLayer"`
	ToLayer    string `json:"toLayer"`
	ImportPath string `json:"importPath"`
}

// DependencyVerificationResult is the full analyzer response, kept in
// lockstep with tools.DependencyVerificationResult (identical JSON tags).
type DependencyVerificationResult struct {
	Success         bool                          `json:"success"`
	ProjectPath     string                        `json:"projectPath"`
	ViolationsCount int                           `json:"violationsCount"`
	Violations      []DependencyBoundaryViolation `json:"violations,omitempty"`
	Summary         string                        `json:"summary"`
}

// DependencyBoundaryAnalyzer scans a project tree for Clean Architecture
// boundary violations. Stateless — one instance can be reused across many
// Analyze calls, including concurrent ones.
type DependencyBoundaryAnalyzer struct{}

// NewDependencyBoundaryAnalyzer constructs the analyzer.
func NewDependencyBoundaryAnalyzer() *DependencyBoundaryAnalyzer {
	return &DependencyBoundaryAnalyzer{}
}

// Layer identifiers. Kept as package-level string constants so tests can
// pin against stable IDs instead of prose. layerUnknown is the sentinel
// for "path did not match any layer segment" — such files are skipped in
// both directions (never a source of violations, never a target).
const (
	layerUnknown    = "unknown"
	layerDomain     = "domain"
	layerUsecases   = "usecases"
	layerAdapters   = "adapters"
	layerFrameworks = "frameworks"
)

// layerRank orders the layers inner → outer. A violation is defined as
// an import whose source layer has a strictly lower rank than the target
// layer's — inner layers importing outer layers.
var layerRank = map[string]int{
	layerDomain:     0,
	layerUsecases:   1,
	layerAdapters:   2,
	layerFrameworks: 3,
}

// layerSegments maps directory-name synonyms to their canonical layer.
// Aliases are the common Clean Architecture / DDD spellings seen across
// Go and TS projects. Matching is case-insensitive at the call site.
var layerSegments = map[string]string{
	"domain":         layerDomain,
	"entities":       layerDomain,
	"entity":         layerDomain,
	"usecases":       layerUsecases,
	"use-cases":      layerUsecases,
	"use_cases":      layerUsecases,
	"usecase":        layerUsecases,
	"application":    layerUsecases,
	"adapters":       layerAdapters,
	"adapter":        layerAdapters,
	"interfaces":     layerAdapters,
	"infrastructure": layerAdapters,
	"frameworks":     layerFrameworks,
}

// scannedImportExtensions is the set of source-file extensions this
// analyzer parses for imports. Op 5 spec scopes support to Go and
// TypeScript — other languages get skipped in collectSourceFiles.
var scannedImportExtensions = map[string]struct{}{
	".go":  {},
	".ts":  {},
	".tsx": {},
}

// Pre-compiled import regexes. Reused across every Analyze call.
//
// Go supports two import shapes, both handled here:
//   - `import "path"` — single-line, optionally preceded by an alias.
//   - `import ( ... )` — multi-line block; a stateful scan tracks
//     whether the current line is inside such a block (see scanGoFile).
//
// TypeScript imports come in two shapes that appear on one line each:
//   - `import ... from '<path>'` — the common ES-module form.
//   - `import '<path>'` — side-effect-only imports.
var (
	goSingleImportRe    = regexp.MustCompile(`^\s*import\s+(?:[A-Za-z_.][\w]*\s+)?"([^"]+)"`)
	goImportBlockOpenRe = regexp.MustCompile(`^\s*import\s*\(\s*$`)
	goImportBlockLineRe = regexp.MustCompile(`^\s*(?:[A-Za-z_.][\w]*\s+)?"([^"]+)"\s*(?://.*)?$`)
	tsFromImportRe      = regexp.MustCompile(`from\s+['"]([^'"]+)['"]`)
	tsSideImportRe      = regexp.MustCompile(`^\s*import\s+['"]([^'"]+)['"]`)
)

// importRef pairs a single import path with the line it was read from,
// so scanFile can attach an accurate line number to every violation.
type importRef struct {
	Path string
	Line int
}

// Analyze walks projectPath and returns a boundary-violation report.
// A missing projectPath returns an error; unreadable files inside the
// walk are silently skipped so one malformed file cannot fail the run.
func (a *DependencyBoundaryAnalyzer) Analyze(projectPath string) (*DependencyVerificationResult, error) {
	result := &DependencyVerificationResult{
		Success:     true,
		ProjectPath: projectPath,
		Violations:  []DependencyBoundaryViolation{},
	}

	files, err := a.collectSourceFiles(projectPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		a.scanFile(f, result)
	}

	result.ViolationsCount = len(result.Violations)
	result.Summary = summarizeDependencies(result.ViolationsCount)
	return result, nil
}

// summarizeDependencies returns the human-readable summary. Named to
// avoid colliding with the other analyzers' own summarize helpers
// (summarize in accessibility, summarizeLanguage in ubiquitous).
func summarizeDependencies(count int) string {
	if count == 0 {
		return "No dependency boundary violations found"
	}
	return "Dependency boundary violations found"
}

// classifyLayer returns the Clean Architecture layer of path based on
// its directory segments. The first (leftmost) recognized segment wins,
// so `internal/adapters/db/postgres.go` classifies as adapters even
// though "db" happens not to be a layer name. Segments are matched
// case-insensitively; both filesystem paths (`/`, `\`) and Go module
// import paths (always `/`) split the same way.
func classifyLayer(path string) string {
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	for _, p := range parts {
		if l, ok := layerSegments[strings.ToLower(p)]; ok {
			return l
		}
	}
	return layerUnknown
}

// isBoundaryViolation returns true when an import from → to crosses a
// Clean Architecture boundary in the wrong direction. Same-layer and
// outer→inner directions are allowed; unknown-layer participants are
// never flagged.
func isBoundaryViolation(from, to string) bool {
	fromRank, ok1 := layerRank[from]
	toRank, ok2 := layerRank[to]
	if !ok1 || !ok2 {
		return false
	}
	return fromRank < toRank
}

// collectSourceFiles returns every .go / .ts / .tsx file under root,
// applying the shared walkutil.go skip set for node_modules, .git,
// dist, build, vendor, .venv, and hidden dot-directories below the root.
// Single-file paths short-circuit the walk.
func (a *DependencyBoundaryAnalyzer) collectSourceFiles(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{root}, nil
	}
	var files []string
	err = filepath.Walk(root, func(p string, entry os.FileInfo, walkErr error) error {
		if walkErr != nil || entry == nil {
			return nil
		}
		if entry.IsDir() {
			return skipUninterestingDir(root, p, entry.Name())
		}
		if _, ok := scannedImportExtensions[strings.ToLower(filepath.Ext(p))]; ok {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// scanFile classifies file's layer, parses its imports, and appends a
// DependencyBoundaryViolation for every import that crosses a boundary
// inward-to-outward. Files whose own layer is unknown are skipped
// entirely — the tool commits to "flag violations among classified
// files" rather than guessing at unclassified callers.
func (a *DependencyBoundaryAnalyzer) scanFile(file string, result *DependencyVerificationResult) {
	fromLayer := classifyLayer(file)
	if fromLayer == layerUnknown {
		return
	}
	imports, err := readImports(file)
	if err != nil {
		return
	}
	for _, imp := range imports {
		appendIfViolation(file, fromLayer, imp, result)
	}
}

// appendIfViolation records one violation when imp crosses a boundary.
// Extracted so scanFile stays under the 30-LOC / complexity-7 limits and
// so the violation-shape mapping lives in one place.
func appendIfViolation(file, fromLayer string, imp importRef, result *DependencyVerificationResult) {
	toLayer := classifyLayer(imp.Path)
	if !isBoundaryViolation(fromLayer, toLayer) {
		return
	}
	result.Violations = append(result.Violations, DependencyBoundaryViolation{
		File:       file,
		LineNumber: imp.Line,
		FromLayer:  fromLayer,
		ToLayer:    toLayer,
		ImportPath: imp.Path,
	})
}

// readImports dispatches to the per-language import parser based on
// file extension. Unknown extensions return nil, nil — the walker
// already filters non-source extensions, so this branch is a safety net.
func readImports(file string) ([]importRef, error) {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".go":
		return readGoImports(file)
	case ".ts", ".tsx":
		return readTSImports(file)
	}
	return nil, nil
}

// readGoImports scans a Go source file and returns every import path,
// handling both single-line `import "x"` statements and multi-line
// `import ( ... )` blocks. State (inside/outside a block) is threaded
// through processGoImportLine so this loop stays a single-purpose scan.
func readGoImports(file string) ([]importRef, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var imports []importRef
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNum := 0
	inBlock := false
	for scanner.Scan() {
		lineNum++
		inBlock = processGoImportLine(scanner.Text(), lineNum, inBlock, &imports)
	}
	return imports, nil
}

// processGoImportLine consumes one line of a Go source file and returns
// the updated inBlock state. Extracted from readGoImports so the state
// machine's branches stay independently reviewable and readGoImports
// stays within the clean-code complexity budget.
func processGoImportLine(line string, lineNum int, inBlock bool, imports *[]importRef) bool {
	if inBlock {
		return processGoImportBlockLine(line, lineNum, imports)
	}
	if goImportBlockOpenRe.MatchString(line) {
		return true
	}
	if m := goSingleImportRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
	}
	return false
}

// processGoImportBlockLine consumes a line while we're inside an
// `import ( ... )` block. Returns false when the block closes, true to
// keep scanning inside it. Comment-only and blank lines pass through
// harmlessly — goImportBlockLineRe does not match them.
func processGoImportBlockLine(line string, lineNum int, imports *[]importRef) bool {
	if strings.TrimSpace(line) == ")" {
		return false
	}
	if m := goImportBlockLineRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
	}
	return true
}

// readTSImports scans a TypeScript source file for ES-module imports.
// Handles both `import ... from '<path>'` and side-effect-only
// `import '<path>'` shapes. A single line can match at most one shape
// here — the from-form is checked first because it is by far the more
// common in real projects.
func readTSImports(file string) ([]importRef, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var imports []importRef
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		appendTSImportsFromLine(scanner.Text(), lineNum, &imports)
	}
	return imports, nil
}

// appendTSImportsFromLine extracts an import from one TS source line, if
// any. Extracted so readTSImports stays a straight-line loop and the two
// shape regexes are dispatched in the same order everywhere.
func appendTSImportsFromLine(line string, lineNum int, imports *[]importRef) {
	if m := tsFromImportRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
		return
	}
	if m := tsSideImportRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
	}
}
