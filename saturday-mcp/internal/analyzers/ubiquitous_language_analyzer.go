// Package analyzers - ubiquitous_language_analyzer.go
//
// UbiquitousLanguageAnalyzer is the concrete analyzer behind the
// check_ubiquitous_language MCP tool (mcp-expand M1 Op 4). It parses a
// DOMAIN_DICTIONARY.md, extracts canonical terms plus their forbidden
// synonyms, walks a project tree, and reports every source-file line
// that mentions any synonym alongside the canonical replacement.
//
// Two dictionary shapes are supported:
//  1. Markdown table with `| **Term** | Domain | Description | Synonyms
//     to AVOID |` rows — the shape used by the framework's actual
//     DOMAIN_DICTIONARY.md.
//  2. `## Term` section headers followed by `**Synonyms to avoid**:` or
//     `**Not**:` lines — the shape described by mcp-expand M1 Op 4.
//
// Both share the "backtick-quoted or comma-separated synonym list"
// grammar. Parsing is line-scan / regex based rather than a real markdown
// AST walker; that matches the style of the neighbouring analyzers
// (complexity_analyzer.go, accessibility_analyzer.go) and keeps the tool
// dependency-free. A present-but-unparseable dictionary yields zero
// violations, not an error — the analyzer treats it as "no rules to
// enforce" per Op 4 spec.
package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// LanguageViolation captures a single unapproved-synonym usage. Field
// shape mirrors internal/tools/responses.go's LanguageViolation so the
// tool layer can marshal analyzer output directly without re-mapping.
type LanguageViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	InvalidTerm string `json:"invalidTerm"`
	Suggested   string `json:"suggested"`
}

// UbiquitousLanguageResult is the full analyzer response, kept in
// lockstep with tools.UbiquitousLanguageResult (identical JSON tags).
type UbiquitousLanguageResult struct {
	Success         bool                `json:"success"`
	ProjectPath     string              `json:"projectPath"`
	ViolationsCount int                 `json:"violationsCount"`
	Violations      []LanguageViolation `json:"violations,omitempty"`
	Summary         string              `json:"summary"`
}

// UbiquitousLanguageAnalyzer scans a project tree for uses of unapproved
// synonyms defined in DOMAIN_DICTIONARY.md. Stateless — one instance can
// be reused across many Analyze calls, including concurrent ones.
type UbiquitousLanguageAnalyzer struct{}

// NewUbiquitousLanguageAnalyzer constructs the analyzer.
func NewUbiquitousLanguageAnalyzer() *UbiquitousLanguageAnalyzer {
	return &UbiquitousLanguageAnalyzer{}
}

// scannedSourceExtensions is the set of source-code extensions the walker
// considers. Matches the eight-language list from Op 4 spec.
var scannedSourceExtensions = map[string]struct{}{
	".go":   {},
	".ts":   {},
	".tsx":  {},
	".js":   {},
	".jsx":  {},
	".py":   {},
	".java": {},
	".cs":   {},
}

// Pre-compiled regexes. Reused across every Analyze call.
var (
	backtickTermRe   = regexp.MustCompile("`([^`]+)`")
	sectionHeaderRe  = regexp.MustCompile(`^##\s+(.+?)\s*$`)
	numberedHeaderRe = regexp.MustCompile(`^\d+\.\s`)
	synonymLineRe    = regexp.MustCompile(`(?i)^\s*\*\*(?:synonyms to avoid|synonyms|avoid|not)\*\*\s*:\s*(.+)$`)
	tableRowRe       = regexp.MustCompile(`^\|\s*\*\*([^*]+)\*\*\s*\|(.+)$`)
)

// termMatcher pairs a compiled word-boundary regex with the canonical
// term that should replace it. Compiled once per Analyze call and reused
// across every scanned line.
type termMatcher struct {
	re        *regexp.Regexp
	synonym   string
	canonical string
}

// Analyze parses dictionaryPath, walks projectPath, and returns a
// violation report. Missing dictionary or project paths return an error;
// a present-but-unparseable dictionary yields zero violations.
func (a *UbiquitousLanguageAnalyzer) Analyze(projectPath, dictionaryPath string) (*UbiquitousLanguageResult, error) {
	result := &UbiquitousLanguageResult{
		Success:     true,
		ProjectPath: projectPath,
		Violations:  []LanguageViolation{},
	}

	synonyms, err := loadDictionary(dictionaryPath)
	if err != nil {
		return nil, err
	}
	if len(synonyms) == 0 {
		result.Summary = summarizeLanguage(0)
		return result, nil
	}

	files, err := a.collectSourceFiles(projectPath)
	if err != nil {
		return nil, err
	}

	matchers := compileMatchers(synonyms)
	for _, f := range files {
		a.scanFile(f, matchers, result)
	}

	result.ViolationsCount = len(result.Violations)
	result.Summary = summarizeLanguage(result.ViolationsCount)
	return result, nil
}

// summarizeLanguage returns the human-readable summary. Named to avoid
// colliding with accessibility_analyzer.go's own summarize helper.
func summarizeLanguage(count int) string {
	if count == 0 {
		return "No ubiquitous language violations found"
	}
	return "Ubiquitous language violations found"
}

// loadDictionary parses dictionaryPath into a synonym→canonical map.
// Returns an error only for I/O failures; a parseable-but-empty file
// returns an empty map, which the caller treats as "no rules to enforce".
func loadDictionary(dictionaryPath string) (map[string]string, error) {
	f, err := os.Open(dictionaryPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	synonyms := map[string]string{}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var currentTerm string
	for scanner.Scan() {
		line := scanner.Text()
		currentTerm = updateTermFromHeader(line, currentTerm)
		parseTableRow(line, synonyms)
		parseSynonymLine(line, currentTerm, synonyms)
	}
	return synonyms, nil
}

// updateTermFromHeader returns the term name when line is a `## Term`
// heading, otherwise returns the previously-active term. Numbered
// organizational headings like `## 1. Core Domains` reset the current
// term so their subsequent synonym lines are not misattributed.
func updateTermFromHeader(line, current string) string {
	m := sectionHeaderRe.FindStringSubmatch(line)
	if m == nil {
		return current
	}
	if numberedHeaderRe.MatchString(m[1]) {
		return ""
	}
	return m[1]
}

// parseTableRow extracts synonyms from a dictionary-table row like
// `| **Term** | Domain | Description | Synonyms to AVOID |` and appends
// them to the synonyms map.
func parseTableRow(line string, synonyms map[string]string) {
	m := tableRowRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	canonical := strings.TrimSpace(m[1])
	parts := splitTableCells(m[2])
	if len(parts) < 2 {
		return
	}
	addSynonyms(parts[len(parts)-1], canonical, synonyms)
}

// splitTableCells splits raw on `|` and drops the trailing empty cell
// caused by the row's closing pipe. Extracted so parseTableRow stays a
// straight-line function.
func splitTableCells(raw string) []string {
	parts := strings.Split(raw, "|")
	for len(parts) > 0 && strings.TrimSpace(parts[len(parts)-1]) == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// parseSynonymLine handles `**Synonyms to avoid**: ...` lines that follow
// a `## Term` section heading. No-op when there is no active term (before
// any header, or under a numbered organizational heading).
func parseSynonymLine(line, currentTerm string, synonyms map[string]string) {
	if currentTerm == "" {
		return
	}
	m := synonymLineRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	addSynonyms(m[1], currentTerm, synonyms)
}

// addSynonyms extracts synonym terms from raw and maps each to canonical.
// Prefers backtick-wrapped terms; if none, falls back to comma-splitting
// the raw text and trimming each piece.
func addSynonyms(raw, canonical string, synonyms map[string]string) {
	matches := backtickTermRe.FindAllStringSubmatch(raw, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			recordSynonym(m[1], canonical, synonyms)
		}
		return
	}
	for _, s := range strings.Split(raw, ",") {
		recordSynonym(s, canonical, synonyms)
	}
}

// recordSynonym normalizes term and, when non-empty, stores it in
// synonyms lowercase-keyed. Trailing parentheticals like `role (too
// generic)` are stripped to their bare term.
func recordSynonym(term, canonical string, synonyms map[string]string) {
	term = strings.TrimSpace(term)
	term = strings.Trim(term, "`*_ ")
	if idx := strings.Index(term, "("); idx > 0 {
		term = strings.TrimSpace(term[:idx])
	}
	if term == "" {
		return
	}
	synonyms[strings.ToLower(term)] = canonical
}

// compileMatchers pre-compiles a case-insensitive word-boundary regex
// per synonym so scanFile can iterate a small slice instead of
// re-compiling per line. Malformed patterns (should not happen because
// synonyms are already quote-escaped) are silently dropped.
func compileMatchers(synonyms map[string]string) []termMatcher {
	matchers := make([]termMatcher, 0, len(synonyms))
	for synonym, canonical := range synonyms {
		pattern := `(?i)\b` + regexp.QuoteMeta(synonym) + `\b`
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		matchers = append(matchers, termMatcher{re: re, synonym: synonym, canonical: canonical})
	}
	return matchers
}

// collectSourceFiles returns every scannable source file under root,
// skipping node_modules, .git, dist, build, vendor, .venv, and other
// hidden dot-directories. Single-file paths short-circuit the walk.
func (a *UbiquitousLanguageAnalyzer) collectSourceFiles(root string) ([]string, error) {
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
			return SkipUninterestingDir(root, p, entry.Name())
		}
		if _, ok := scannedSourceExtensions[strings.ToLower(filepath.Ext(p))]; ok {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// scanFile reads file line-by-line and appends a LanguageViolation for
// every synonym match. Read errors cause the file to be silently skipped
// so the walk continues — consistent with the accessibility analyzer.
func (a *UbiquitousLanguageAnalyzer) scanFile(file string, matchers []termMatcher, result *UbiquitousLanguageResult) {
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
		appendLineViolations(file, lineNum, scanner.Text(), matchers, result)
	}
}

// appendLineViolations runs every matcher against one line and appends
// any hits. Extracted from scanFile so scanFile stays under the
// complexity-7 and 30-LOC clean-code limits.
func appendLineViolations(file string, lineNum int, line string, matchers []termMatcher, result *UbiquitousLanguageResult) {
	for _, m := range matchers {
		if m.re.MatchString(line) {
			result.Violations = append(result.Violations, LanguageViolation{
				File:        file,
				LineNumber:  lineNum,
				InvalidTerm: m.synonym,
				Suggested:   m.canonical,
			})
		}
	}
}
