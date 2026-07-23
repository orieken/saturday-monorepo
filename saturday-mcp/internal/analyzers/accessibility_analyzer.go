// Package analyzers - accessibility_analyzer.go
//
// AccessibilityAnalyzer is the concrete analyzer behind the
// check_accessibility MCP tool (mcp-expand M1 Op 3). It walks a project
// path (or scans a single file) matching common UI template extensions
// and reports semantic-HTML / ARIA violations line by line.
//
// The analyzer is intentionally regex-based rather than parsing HTML
// with a real DOM library — check_accessibility is a first-pass
// framework-wide check, not a full accessibility audit. Regexes keep the
// tool dependency-free, easy to reason about, and match the style of the
// existing ComplexityAnalyzer (see complexity_analyzer.go). Rules are
// kept small, single-purpose, and named with stable IDs so downstream
// consumers can filter or trend them over time.
package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AccessibilityViolation captures a single a11y defect the analyzer
// found. Field shape mirrors internal/tools/responses.go's
// AccessibilityViolation so the tool layer can marshal the analyzer
// output directly without re-mapping fields.
type AccessibilityViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	Element     string `json:"element"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

// AccessibilityReportResult is the full analyzer response, kept in
// lockstep with tools.AccessibilityReportResult (identical JSON tags).
type AccessibilityReportResult struct {
	Success         bool                     `json:"success"`
	Path            string                   `json:"path"`
	TotalFiles      int                      `json:"totalFiles"`
	ViolationsCount int                      `json:"violationsCount"`
	Violations      []AccessibilityViolation `json:"violations,omitempty"`
	Summary         string                   `json:"summary"`
}

// AccessibilityAnalyzer scans UI template files for semantic-HTML and
// ARIA violations. Stateless — one instance can be reused across many
// Analyze calls, including concurrent ones.
type AccessibilityAnalyzer struct{}

// NewAccessibilityAnalyzer constructs an analyzer with the default rule
// set. Rules are hard-wired for now; a configurable rule registry can
// be introduced later if downstream consumers need to opt individual
// rules out.
func NewAccessibilityAnalyzer() *AccessibilityAnalyzer {
	return &AccessibilityAnalyzer{}
}

// scannedExtensions is the set of file extensions the analyzer walks
// when given a directory. Single-file mode ignores the extension — the
// caller has already chosen the file explicitly.
var scannedExtensions = map[string]struct{}{
	".html":   {},
	".htm":    {},
	".vue":    {},
	".jsx":    {},
	".tsx":    {},
	".svelte": {},
}

// Rule IDs — kept as package-level constants so tests can pin against
// stable identifiers instead of prose descriptions.
const (
	ruleClickableNonInteractive = "clickable-non-interactive"
	ruleImgAlt                  = "img-alt"
	ruleAnchorName              = "anchor-name"
	ruleInputLabel              = "input-label"
	ruleHTMLLang                = "html-lang"
)

// Pre-compiled regexes. Every match is line-scoped, so the analyzer can
// report the exact line number without re-scanning the file. Patterns
// are lower-case-insensitive so <DIV ONCLICK=...> gets caught the same
// as <div onclick=...>.
var (
	clickableNonInteractiveRe = regexp.MustCompile(`(?i)<(div|span)\b[^>]*\bon[a-z]+\s*=`)
	imgTagRe                  = regexp.MustCompile(`(?i)<img\b[^>]*>`)
	imgHasAltRe               = regexp.MustCompile(`(?i)\balt\s*=`)
	anchorOpenRe              = regexp.MustCompile(`(?i)<a\b[^>]*>`)
	ariaLabelRe               = regexp.MustCompile(`(?i)\baria-label\s*=\s*["'][^"']+["']`)
	inputTagRe                = regexp.MustCompile(`(?i)<input\b[^>]*>`)
	inputHasIDRe              = regexp.MustCompile(`(?i)\bid\s*=\s*["']([^"']+)["']`)
	inputTypeHiddenRe         = regexp.MustCompile(`(?i)\btype\s*=\s*["']hidden["']`)
	labelForRe                = regexp.MustCompile(`(?i)<label\b[^>]*\bfor\s*=\s*["']([^"']+)["']`)
	htmlTagRe                 = regexp.MustCompile(`(?i)<html\b[^>]*>`)
	htmlHasLangRe             = regexp.MustCompile(`(?i)\blang\s*=`)
)

// Analyze inspects path and returns a report. path may be a single file
// or a directory; when it is a directory the walk skips node_modules,
// dot-directories, and any file whose extension is not in
// scannedExtensions. Errors reading a specific file do not fail the
// whole run — the file is skipped, the walk continues.
func (a *AccessibilityAnalyzer) Analyze(path string) (*AccessibilityReportResult, error) {
	result := &AccessibilityReportResult{
		Success:    true,
		Path:       path,
		Violations: []AccessibilityViolation{},
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	files, err := a.collectFiles(path, info.IsDir())
	if err != nil {
		return nil, err
	}
	result.TotalFiles = len(files)

	for _, f := range files {
		a.analyzeFile(f, result)
	}

	result.ViolationsCount = len(result.Violations)
	result.Summary = summarize(result.ViolationsCount)
	return result, nil
}

// summarize returns the human-readable summary string. Extracted so
// Analyze stays under the 30-LOC and complexity-7 clean-code limits.
func summarize(count int) string {
	if count == 0 {
		return "No accessibility violations found"
	}
	return "Accessibility violations found"
}

// collectFiles returns the concrete file list to scan. Extracted from
// Analyze so the walk logic stays a single-purpose helper.
func (a *AccessibilityAnalyzer) collectFiles(path string, isDir bool) ([]string, error) {
	if !isDir {
		return []string{path}, nil
	}
	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return skipUninteresting(path, p, info.Name())
		}
		if _, ok := scannedExtensions[strings.ToLower(filepath.Ext(p))]; ok {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// skipUninteresting returns filepath.SkipDir for directories we should
// never descend into (node_modules, hidden subdirectories) and nil to
// keep walking otherwise. Extracted from the filepath.Walk closure to
// keep that closure small.
func skipUninteresting(root, current, name string) error {
	if name == "node_modules" {
		return filepath.SkipDir
	}
	if strings.HasPrefix(name, ".") && current != root {
		return filepath.SkipDir
	}
	return nil
}

// analyzeFile runs every rule against one file. Labels are collected in
// a first pass and passed to the per-line rule scan so the input-label
// rule can resolve <label for="..."> references defined anywhere in the
// same file. Line-scoped rules avoid loading the entire file into
// memory as a single string.
func (a *AccessibilityAnalyzer) analyzeFile(file string, result *AccessibilityReportResult) {
	labels := collectLabelIDs(file)
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		runLineRules(file, lineNum, scanner.Text(), labels, result)
	}
}

// runLineRules dispatches every per-line rule for one line. Kept as a
// flat dispatch so analyzeFile does not blow past the complexity limit;
// each rule function stays independently testable.
func runLineRules(file string, lineNum int, line string, labels map[string]struct{}, result *AccessibilityReportResult) {
	checkHTMLLang(file, lineNum, line, result)
	checkClickableNonInteractive(file, lineNum, line, result)
	checkImgAlt(file, lineNum, line, result)
	checkAnchorName(file, lineNum, line, result)
	checkInputLabel(file, lineNum, line, labels, result)
}

// checkHTMLLang flags any <html> tag missing a lang attribute.
func checkHTMLLang(file string, lineNum int, line string, result *AccessibilityReportResult) {
	if !htmlTagRe.MatchString(line) || htmlHasLangRe.MatchString(line) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "html",
		Rule:        ruleHTMLLang,
		Description: "<html> element is missing a lang attribute",
	})
}

// checkClickableNonInteractive flags <div>/<span> that carry any
// on-something handler. Use a native button/link/input instead — that
// is the standard WAI-ARIA authoring-practices recommendation.
func checkClickableNonInteractive(file string, lineNum int, line string, result *AccessibilityReportResult) {
	m := clickableNonInteractiveRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	el := strings.ToLower(m[1])
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     el,
		Rule:        ruleClickableNonInteractive,
		Description: "Non-interactive element (<" + el + ">) has an event handler; use <button> or <a> instead",
	})
}

// checkImgAlt flags any <img> tag missing an alt attribute (an empty
// alt="" is allowed and expected for decorative images).
func checkImgAlt(file string, lineNum int, line string, result *AccessibilityReportResult) {
	tag := imgTagRe.FindString(line)
	if tag == "" || imgHasAltRe.MatchString(tag) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "img",
		Rule:        ruleImgAlt,
		Description: "<img> element is missing an alt attribute",
	})
}

// checkAnchorName flags <a> tags with neither visible inline text nor
// an aria-label. Uses a same-line heuristic: if the anchor opens and
// closes on this line with no text between, it is empty; if the anchor
// opens but does not close on this line, aria-label on the open tag is
// required.
func checkAnchorName(file string, lineNum int, line string, result *AccessibilityReportResult) {
	openTag := anchorOpenRe.FindString(line)
	if openTag == "" || ariaLabelRe.MatchString(openTag) {
		return
	}
	if hasVisibleAnchorText(line, openTag) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "a",
		Rule:        ruleAnchorName,
		Description: "<a> element has no visible text and no aria-label",
	})
}

// hasVisibleAnchorText returns true when the substring following the
// anchor open tag has non-whitespace content (either between the open
// tag and a same-line </a>, or trailing off the line).
func hasVisibleAnchorText(line, openTag string) bool {
	rest := line[strings.Index(line, openTag)+len(openTag):]
	closeIdx := strings.Index(strings.ToLower(rest), "</a>")
	if closeIdx < 0 {
		return strings.TrimSpace(rest) != ""
	}
	return strings.TrimSpace(rest[:closeIdx]) != ""
}

// checkInputLabel flags <input> tags lacking both aria-label and a
// matching <label for=...> anywhere in the same file. Hidden inputs
// are skipped — they have no visible affordance to label.
func checkInputLabel(file string, lineNum int, line string, labels map[string]struct{}, result *AccessibilityReportResult) {
	tag := inputTagRe.FindString(line)
	if tag == "" || inputTypeHiddenRe.MatchString(tag) {
		return
	}
	if ariaLabelRe.MatchString(tag) || hasMatchingLabel(tag, labels) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "input",
		Rule:        ruleInputLabel,
		Description: "<input> element has no associated <label for=...> and no aria-label",
	})
}

// hasMatchingLabel returns true when the input tag exposes an id that
// some label[for=...] in the same file references.
func hasMatchingLabel(tag string, labels map[string]struct{}) bool {
	idMatch := inputHasIDRe.FindStringSubmatch(tag)
	if idMatch == nil {
		return false
	}
	_, ok := labels[idMatch[1]]
	return ok
}

// collectLabelIDs scans file once and returns every id value referenced
// by a <label for="..."> tag. Called once per file before rule
// dispatch so checkInputLabel can resolve label references forward and
// backward relative to the input tag. Read errors return an empty set —
// the caller (analyzeFile) will hit the same error on its own open and
// skip the file entirely.
func collectLabelIDs(file string) map[string]struct{} {
	ids := map[string]struct{}{}
	f, err := os.Open(file)
	if err != nil {
		return ids
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		for _, m := range labelForRe.FindAllStringSubmatch(scanner.Text(), -1) {
			ids[m[1]] = struct{}{}
		}
	}
	return ids
}
