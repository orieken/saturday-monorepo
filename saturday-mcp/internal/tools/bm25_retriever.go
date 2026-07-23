// Package tools - bm25_retriever.go
//
// BM25 retrieval over the installed project's docs corpus. Per framework
// ADR-002 (Corpus-Aware Retrieval Strategy), this is retrieval tier 2:
// prose retrieval via sqlite-fts5 with the built-in bm25() ranking
// function. The corpus is docs/features/, docs/adrs/, docs/patterns/,
// and docs/runbooks/ inside the installed project — typically small
// enough to lazy-build the index on first use.
//
// **CGO note.** The task spec named github.com/mattn/go-sqlite3, but
// saturday-mcp's release builds pin CGO_ENABLED=0 (see .goreleaser.yaml)
// and mattn/go-sqlite3 requires CGO. We use the pure-Go
// modernc.org/sqlite driver instead — same FTS5 support (compiled in by
// default), no CGO, no binary-size regression versus the original plan.
// The tradeoff is slower cold-start and marginally slower queries than
// the C driver, both immaterial for docs-corpus scale.
//
// The Retriever interface (see retriever.go) is shared with the tier-1
// KICorpusRetriever so the search_docs tool (Op 7b) can consume this
// through the same abstraction the search_ki tool already uses.
package tools

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orieken/saturday-mcp/internal/analyzers"

	// Register the pure-Go sqlite driver under the name "sqlite".
	_ "modernc.org/sqlite"
)

// BM25Retriever implements Retriever using fts5's built-in bm25()
// ranking. One instance owns one sqlite database file — safe for
// concurrent Retrieve calls (database/sql handles connection pooling),
// but EnsureIndex is not safe to call concurrently with itself.
type BM25Retriever struct {
	db     *sql.DB
	dbPath string
}

// bm25MaxResults caps the surfaced set so the calling LLM sees a
// manageable list rather than the entire matching corpus. Matches the
// tier-1 KICorpusRetriever's cap so both retrievers behave uniformly.
const bm25MaxResults = 25

// docsFTSSchema creates the fts5 virtual table used by every query.
// Path is UNINDEXED because we surface it in results but never want it
// to participate in MATCH scoring — a query like "docs" shouldn't rank
// files under a "docs/" directory higher purely because of the path.
// Tokenizer: porter unicode61 — stemming + unicode word segmentation,
// the fts5 default for prose corpora.
const docsFTSSchema = `CREATE VIRTUAL TABLE IF NOT EXISTS docs_fts USING fts5(
	path UNINDEXED,
	title,
	body,
	tokenize = 'porter unicode61'
)`

// NewBM25Retriever opens (or creates) the sqlite database at dbPath and
// ensures the docs_fts virtual table exists. Returns an error only on
// unrecoverable DB open/setup failure — a corrupt existing DB is fatal
// by design (per Op 7a spec); the caller's remedy is to delete the file
// and let the next construction rebuild.
func NewBM25Retriever(dbPath string) (*BM25Retriever, error) {
	if strings.TrimSpace(dbPath) == "" {
		return nil, fmt.Errorf("bm25 retriever: dbPath is required")
	}
	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("bm25 retriever: cannot create db dir %q: %w", dir, err)
		}
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("bm25 retriever: open %q failed: %w", dbPath, err)
	}
	if _, err := db.Exec(docsFTSSchema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("bm25 retriever: schema init failed: %w", err)
	}
	return &BM25Retriever{db: db, dbPath: dbPath}, nil
}

// Close releases the underlying database handle. Safe to call on a nil
// retriever or after a prior Close — both no-op.
func (r *BM25Retriever) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	return r.db.Close()
}

// EnsureIndex walks each corpus root and (re)indexes every .md / .mdx
// file found under it. Idempotent by design: for each file, an existing
// row at that path is DELETEd first, then the current file contents are
// INSERTed — so re-runs converge on the current disk state rather than
// accumulating duplicates. Directory-level skip logic delegates to
// analyzers.SkipUninterestingDir so this walker shares one implementation
// with every analyzer walker (node_modules, .git, dist, build, vendor,
// .venv, dot-directories).
func (r *BM25Retriever) EnsureIndex(corpusPaths []string) error {
	for _, root := range corpusPaths {
		if strings.TrimSpace(root) == "" {
			continue
		}
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			continue // corpus root may not exist on this install — skip silently
		}
		walkErr := filepath.Walk(root, func(p string, entry os.FileInfo, walkErr error) error {
			if walkErr != nil || entry == nil {
				return nil
			}
			if entry.IsDir() {
				return analyzers.SkipUninterestingDir(root, p, entry.Name())
			}
			if !isMarkdownExtension(p) {
				return nil
			}
			return r.indexFile(p)
		})
		if walkErr != nil {
			return fmt.Errorf("bm25 retriever: walk %q: %w", root, walkErr)
		}
	}
	return nil
}

// isMarkdownExtension is the corpus-file gate. Kept separate so the
// extension set is one obvious line to extend later (e.g. .rst) without
// touching the walker.
func isMarkdownExtension(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".mdx"
}

// indexFile upserts one file into docs_fts. Missing/malformed frontmatter
// is tolerated — falls back to a filename-derived title so retrieval
// still surfaces the file rather than skipping it entirely. Read errors
// silently skip the file (per ADR-002's "markdown stays canonical, don't
// fail retrieval because one file is transient/unreadable").
func (r *BM25Retriever) indexFile(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	title := extractFrontmatterName(body)
	if title == "" {
		title = titleFromFilename(path)
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM docs_fts WHERE path = ?", path); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(
		"INSERT INTO docs_fts (path, title, body) VALUES (?, ?, ?)",
		path, title, string(body),
	); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// extractFrontmatterName pulls the `name:` field from a YAML frontmatter
// block. Returns "" when no frontmatter is present or the name field is
// absent — the caller's fallback handles the empty case. Deliberately a
// tiny hand-parser rather than pulling in a YAML dependency; the docs
// corpus's frontmatter shape is stable and small (see
// shared/schemas/ki-frontmatter.schema.json).
func extractFrontmatterName(body []byte) string {
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	if !scanner.Scan() || scanner.Text() != "---" {
		return ""
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			return ""
		}
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		}
	}
	return ""
}

// Retrieve runs an fts5 MATCH query and returns up to bm25MaxResults
// hits ordered by BM25 rank (title-column-boosted 10x over body).
//
// tags and domain args are accepted for Retriever interface conformance
// with KICorpusRetriever, but BM25 doesn't consume structured metadata
// — that's tier 1's job. Documented as ignored rather than silently
// dropped so callers know why filtering by tag has no effect here.
//
// An empty query, an unindexed DB, or a query that trips fts5's MATCH
// parser all surface as zero hits rather than errors — the caller's LLM
// reformulates more usefully than the tool escalating.
func (r *BM25Retriever) Retrieve(query string, tags []string, domain string) ([]Reference, error) {
	_ = tags
	_ = domain
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, nil
	}
	escaped := escapeFTS5Query(trimmed)
	if escaped == "" {
		return nil, nil
	}
	rows, err := r.db.Query(
		`SELECT path, title,
		        snippet(docs_fts, 2, '[', ']', '...', 20) AS summary,
		        bm25(docs_fts, 1.0, 10.0, 1.0) AS relevance
		 FROM docs_fts
		 WHERE docs_fts MATCH ?
		 ORDER BY relevance
		 LIMIT ?`,
		escaped, bm25MaxResults,
	)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	var refs []Reference
	for rows.Next() {
		var (
			path, title, summary string
			rank                 float64
		)
		if err := rows.Scan(&path, &title, &summary, &rank); err != nil {
			continue
		}
		// fts5 bm25 rank is lower-is-better (negative for matches); invert
		// so callers can compare against KICorpusRetriever's positive scale.
		refs = append(refs, Reference{
			Path:      path,
			Title:     title,
			Summary:   summary,
			Relevance: -rank,
		})
	}
	return refs, nil
}

// escapeFTS5Query turns a natural-language query into an fts5 MATCH
// expression that won't error on punctuation. Each whitespace-separated
// token is wrapped in double-quotes (fts5 phrase syntax); embedded
// double-quotes are doubled per fts5 escape rules. Multiple quoted
// tokens implicit-AND together, matching intuitive keyword search
// semantics.
func escapeFTS5Query(q string) string {
	tokens := strings.Fields(q)
	if len(tokens) == 0 {
		return ""
	}
	quoted := make([]string, 0, len(tokens))
	for _, t := range tokens {
		escaped := strings.ReplaceAll(t, `"`, `""`)
		quoted = append(quoted, `"`+escaped+`"`)
	}
	return strings.Join(quoted, " ")
}
