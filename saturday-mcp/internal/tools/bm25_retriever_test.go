package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBM25Retriever_EnsureIndexAndRetrieveRoundtrip pins the retriever's
// happy path: index a temp corpus, MATCH-search it, receive a hit whose
// path points at the on-disk file. The single most important behaviour
// this whole file guards — everything else is an edge case around it.
func TestBM25Retriever_EnsureIndexAndRetrieveRoundtrip(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "auth-runbook.md"), `---
name: OAuth token refresh runbook
---

When the OAuth token refresh flow fails, follow these steps.
`)

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}

	refs, err := retriever.Retrieve("OAuth token refresh", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("want 1 hit, got %d", len(refs))
	}
	if !strings.HasSuffix(refs[0].Path, "auth-runbook.md") {
		t.Errorf("hit path = %q, want suffix auth-runbook.md", refs[0].Path)
	}
	if refs[0].Title != "OAuth token refresh runbook" {
		t.Errorf("title = %q, want frontmatter name", refs[0].Title)
	}
	if refs[0].Summary == "" {
		t.Error("summary should be non-empty (snippet of body)")
	}
}

// TestBM25Retriever_UnindexedDBReturnsEmpty confirms querying an empty
// database is a silent zero-hit, not an error. The tool binary lazy-
// builds on first use — if a caller happens to query before the first
// EnsureIndex, they get nothing back, not a failure.
func TestBM25Retriever_UnindexedDBReturnsEmpty(t *testing.T) {
	retriever := newBM25(t)
	refs, err := retriever.Retrieve("anything", nil, "")
	if err != nil {
		t.Fatalf("Retrieve on empty DB: %v", err)
	}
	if len(refs) != 0 {
		t.Errorf("empty DB should yield 0 hits, got %d", len(refs))
	}
}

// TestBM25Retriever_EnsureIndexIsIdempotent covers the DELETE-then-INSERT
// pattern in indexFile. Op 7b's tool will re-run EnsureIndex on every
// invocation (or on `/reindex-docs`), so re-indexing MUST NOT accumulate
// duplicate rows.
func TestBM25Retriever_EnsureIndexIsIdempotent(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "note.md"), "# Some note\n\nBody about widgets.\n")

	retriever := newBM25(t)
	for i := 0; i < 3; i++ {
		if err := retriever.EnsureIndex([]string{corpus}); err != nil {
			t.Fatalf("EnsureIndex iter %d: %v", i, err)
		}
	}

	if got := countRows(t, retriever); got != 1 {
		t.Errorf("row count after 3 EnsureIndex passes = %d, want 1", got)
	}
}

// TestBM25Retriever_TitleMatchOutranksBodyMatch pins the column-weight
// contract. `bm25(docs_fts, 1.0, 10.0, 1.0)` in Retrieve gives title 10x
// the weight of body — search_docs callers care most about title hits,
// and this test locks that ranking behaviour so a future weight tweak
// can't silently invert it.
func TestBM25Retriever_TitleMatchOutranksBodyMatch(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "body-only.md"), `---
name: Deployment checklist
---

Unrelated title, but body mentions telemetry once.
`)
	writeFile(t, filepath.Join(corpus, "title-hit.md"), `---
name: Telemetry setup guide
---

Body talks about something else entirely.
`)

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}

	refs, err := retriever.Retrieve("telemetry", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if len(refs) < 2 {
		t.Fatalf("expected both docs to match, got %d", len(refs))
	}
	if !strings.HasSuffix(refs[0].Path, "title-hit.md") {
		t.Errorf("title-hit.md should rank first; got order: %v", pathsOnly(refs))
	}
}

// TestBM25Retriever_MalformedMarkdownDoesNotBreakIndexing covers the
// "one broken file mustn't take down search" property from ADR-002.
// Empty file, no-frontmatter file, and a valid file all in the same
// corpus — the valid file MUST still surface.
func TestBM25Retriever_MalformedMarkdownDoesNotBreakIndexing(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "empty.md"), "")
	writeFile(t, filepath.Join(corpus, "no-frontmatter.md"), "Just some body text mentioning zebras.\n")
	writeFile(t, filepath.Join(corpus, "good.md"), `---
name: Zebra research notes
---

Zebras have stripes.
`)

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex must not fail on malformed siblings: %v", err)
	}

	refs, err := retriever.Retrieve("zebras", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if len(refs) < 1 {
		t.Fatalf("expected at least one hit, got 0")
	}
	// The no-frontmatter file should fall back to a filename-derived title
	// rather than being skipped. Find it in results to prove indexing.
	foundNoFM := false
	for _, r := range refs {
		if strings.HasSuffix(r.Path, "no-frontmatter.md") {
			foundNoFM = true
			if r.Title == "" {
				t.Error("no-frontmatter.md should have a filename-derived fallback title")
			}
		}
	}
	if !foundNoFM {
		t.Error("no-frontmatter.md should be indexed and retrievable via body match")
	}
}

// TestBM25Retriever_WalkSkipsUninterestingDirs proves EnsureIndex uses
// the shared analyzers.SkipUninterestingDir helper — files under
// node_modules/, .git/, and .hidden/ MUST NOT be indexed. Guards against
// the "everything in the tree got indexed" bug that a naive walker
// would introduce.
func TestBM25Retriever_WalkSkipsUninterestingDirs(t *testing.T) {
	corpus := t.TempDir()
	// One legitimately-visible file to prove the walk actually ran.
	writeFile(t, filepath.Join(corpus, "docs", "visible.md"), "# Visible\n\nContent about widgets.\n")
	// Files under skip-directories — all mention the same term.
	writeFile(t, filepath.Join(corpus, "node_modules", "pkg.md"), "widgets widgets widgets\n")
	writeFile(t, filepath.Join(corpus, ".git", "note.md"), "widgets widgets widgets\n")
	writeFile(t, filepath.Join(corpus, ".hidden", "note.md"), "widgets widgets widgets\n")
	writeFile(t, filepath.Join(corpus, "vendor", "readme.md"), "widgets widgets widgets\n")

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}

	refs, err := retriever.Retrieve("widgets", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	for _, r := range refs {
		for _, forbidden := range []string{"/node_modules/", "/.git/", "/.hidden/", "/vendor/"} {
			if strings.Contains(r.Path, forbidden) {
				t.Errorf("indexed file under skip-dir %s: %q", forbidden, r.Path)
			}
		}
	}
	if len(refs) != 1 {
		t.Errorf("want exactly the one visible file, got %d hits: %v", len(refs), pathsOnly(refs))
	}
}

// TestBM25Retriever_MissingCorpusPathSilentlySkipped covers the "install
// doesn't have a docs/ dir" scenario — the retriever must not fail if a
// configured corpus root doesn't exist on disk.
func TestBM25Retriever_MissingCorpusPathSilentlySkipped(t *testing.T) {
	retriever := newBM25(t)
	err := retriever.EnsureIndex([]string{
		"/nonexistent/absolutely/not/there",
		"", // empty string too — should be skipped
	})
	if err != nil {
		t.Errorf("EnsureIndex should tolerate missing/empty corpus paths, got: %v", err)
	}
}

// TestBM25Retriever_EmptyQueryReturnsNil pins the "empty query = no
// work" contract — matches KICorpusRetriever's behaviour so Op 7b's
// tool can rely on both retrievers behaving the same way.
func TestBM25Retriever_EmptyQueryReturnsNil(t *testing.T) {
	retriever := newBM25(t)
	for _, q := range []string{"", "   ", "\t\n"} {
		refs, err := retriever.Retrieve(q, nil, "")
		if err != nil {
			t.Errorf("Retrieve(%q): %v", q, err)
		}
		if refs != nil {
			t.Errorf("Retrieve(%q) = %v, want nil", q, refs)
		}
	}
}

// TestBM25Retriever_QueryWithSpecialCharsDoesNotError covers punctuation
// that would trip a raw fts5 MATCH parser (colons, quotes, parentheses).
// The escapeFTS5Query helper quotes each token; the retriever must
// therefore return silent zero-hits rather than propagating a MATCH
// syntax error.
func TestBM25Retriever_QueryWithSpecialCharsDoesNotError(t *testing.T) {
	retriever := newBM25(t)
	for _, q := range []string{
		`what is "auth"?`,
		`(retrieval)`,
		`token: refresh`,
		`a AND b`, // AND is fts5-reserved; quoting must neutralise it
	} {
		if _, err := retriever.Retrieve(q, nil, ""); err != nil {
			t.Errorf("Retrieve(%q) must not error, got: %v", q, err)
		}
	}
}

// TestBM25Retriever_NewBM25RetrieverRejectsBlankPath pins the constructor
// input validation — blank/whitespace-only paths return an error, not a
// broken retriever the caller would only discover on first use.
func TestBM25Retriever_NewBM25RetrieverRejectsBlankPath(t *testing.T) {
	for _, p := range []string{"", "   "} {
		if _, err := NewBM25Retriever(p); err == nil {
			t.Errorf("NewBM25Retriever(%q) should error, got nil", p)
		}
	}
}

// TestBM25Retriever_CreatesParentDir confirms the constructor mkdirs the
// db file's parent path — real callers pass paths like
// .claude/rag/docs-fts5.sqlite where the .claude/rag directory may not
// exist yet.
func TestBM25Retriever_CreatesParentDir(t *testing.T) {
	root := t.TempDir()
	dbPath := filepath.Join(root, "nested", "subdir", "fts.sqlite")

	retriever, err := NewBM25Retriever(dbPath)
	if err != nil {
		t.Fatalf("NewBM25Retriever: %v", err)
	}
	t.Cleanup(func() { _ = retriever.Close() })

	if _, err := os.Stat(dbPath); err != nil {
		t.Errorf("db file not created at %q: %v", dbPath, err)
	}
}

// TestBM25Retriever_MDXFilesIndexed covers the second markdown extension
// referenced in Op 7's spec — .mdx (React-flavoured markdown) is a real
// shape in the framework's docs, and skipping it would silently exclude
// pages from search.
func TestBM25Retriever_MDXFilesIndexed(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "component.mdx"), "# React thing\n\nMentions capacitors.\n")

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}

	refs, err := retriever.Retrieve("capacitors", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if len(refs) != 1 || !strings.HasSuffix(refs[0].Path, "component.mdx") {
		t.Errorf(".mdx file should be indexed and retrievable; got %v", pathsOnly(refs))
	}
}

// TestBM25Retriever_NonMarkdownFilesSkipped pins the extension gate — a
// .txt sibling under the corpus MUST NOT be indexed. Guards against the
// walker growing a bug where isMarkdownExtension gets bypassed.
func TestBM25Retriever_NonMarkdownFilesSkipped(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "note.md"), "# Real doc\n\nContent mentioning frobnicate.\n")
	writeFile(t, filepath.Join(corpus, "readme.txt"), "frobnicate frobnicate frobnicate\n")
	writeFile(t, filepath.Join(corpus, "config.yaml"), "frobnicate: yes\n")

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}
	if got := countRows(t, retriever); got != 1 {
		t.Errorf("only the .md file should be indexed; got %d rows", got)
	}
}

// TestBM25Retriever_FrontmatterWithoutNameFallsBackToFilename covers the
// case where a well-formed frontmatter block exists but has no `name:`
// key. extractFrontmatterName must return "" so indexFile falls back to
// the filename-derived title.
func TestBM25Retriever_FrontmatterWithoutNameFallsBackToFilename(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "misc-doc.md"), `---
domain: infrastructure
tags: [ops]
---

Body mentions grommets.
`)

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}

	refs, err := retriever.Retrieve("grommets", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if len(refs) != 1 {
		t.Fatalf("want 1 hit, got %d", len(refs))
	}
	if refs[0].Title != "misc doc" {
		t.Errorf("title = %q, want filename-fallback %q", refs[0].Title, "misc doc")
	}
}

// TestBM25Retriever_FrontmatterMissingCloseDelimiterHandled covers the
// scanner's "opened frontmatter but never saw the closing ---" edge —
// file must still be indexable via body match with a filename-derived
// title.
func TestBM25Retriever_FrontmatterMissingCloseDelimiterHandled(t *testing.T) {
	corpus := t.TempDir()
	// Note: no trailing --- ever
	writeFile(t, filepath.Join(corpus, "unclosed.md"), "---\nfoo: bar\nbaz: qux\ncontent about squirrels\n")

	retriever := newBM25(t)
	if err := retriever.EnsureIndex([]string{corpus}); err != nil {
		t.Fatalf("EnsureIndex: %v", err)
	}
	refs, err := retriever.Retrieve("squirrels", nil, "")
	if err != nil {
		t.Fatalf("Retrieve: %v", err)
	}
	if len(refs) != 1 {
		t.Errorf("unclosed-frontmatter file should still index; got %d hits", len(refs))
	}
}

// TestBM25Retriever_IndexFileErrorsWhenDBClosed exercises the transaction
// error path in indexFile — after Close, tx.Begin fails and the error
// must propagate up through EnsureIndex.
func TestBM25Retriever_IndexFileErrorsWhenDBClosed(t *testing.T) {
	corpus := t.TempDir()
	writeFile(t, filepath.Join(corpus, "doc.md"), "# Doc\n\nBody.\n")

	dbPath := filepath.Join(t.TempDir(), "closed.sqlite")
	retriever, err := NewBM25Retriever(dbPath)
	if err != nil {
		t.Fatalf("NewBM25Retriever: %v", err)
	}
	if err := retriever.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	err = retriever.EnsureIndex([]string{corpus})
	if err == nil {
		t.Error("EnsureIndex on closed DB should surface an error")
	}
}

// TestBM25Retriever_RetrieveOnClosedDBReturnsNil pins the "MATCH-parse
// or DB error surfaces as zero hits, not a tool failure" contract. A
// closed DB is the easy repro; the same silent-empty path covers fts5
// parser errors in production.
func TestBM25Retriever_RetrieveOnClosedDBReturnsNil(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "closed-retrieve.sqlite")
	retriever, err := NewBM25Retriever(dbPath)
	if err != nil {
		t.Fatalf("NewBM25Retriever: %v", err)
	}
	if err := retriever.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	refs, err := retriever.Retrieve("anything", nil, "")
	if err != nil {
		t.Errorf("Retrieve on closed DB should return nil, not error; got err=%v", err)
	}
	if refs != nil {
		t.Errorf("Retrieve on closed DB should return nil refs; got %v", refs)
	}
}

// TestBM25Retriever_ConstructorFailsWhenParentIsFile covers the MkdirAll
// error branch — if the intended db-file parent path already exists as a
// regular file, mkdir cannot create it and the constructor must error
// out cleanly.
func TestBM25Retriever_ConstructorFailsWhenParentIsFile(t *testing.T) {
	root := t.TempDir()
	// Create "collision" as a plain file, then try to open a db under it.
	collision := filepath.Join(root, "collision")
	if err := os.WriteFile(collision, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("prep: %v", err)
	}
	if _, err := NewBM25Retriever(filepath.Join(collision, "should-fail.sqlite")); err == nil {
		t.Error("NewBM25Retriever should error when parent path is a file, not a directory")
	}
}

// TestBM25Retriever_CloseIsNilSafe pins the "safe to Close idle/nil"
// contract — long-lived tools may Close on shutdown paths that also
// handle the never-constructed case.
func TestBM25Retriever_CloseIsNilSafe(t *testing.T) {
	var r *BM25Retriever
	if err := r.Close(); err != nil {
		t.Errorf("nil retriever Close: %v", err)
	}

	r2 := newBM25(t)
	if err := r2.Close(); err != nil {
		t.Errorf("first Close: %v", err)
	}
	// t.Cleanup would double-close via newBM25 — but Close-after-Close on
	// a database/sql handle is safe (returns sql.ErrConnDone or nil).
}

// newBM25 constructs a fresh retriever pointed at a temp-directory
// sqlite file and registers Close for the test's cleanup. Every test
// uses it — the one line of setup noise DRY'd away without hiding the
// interesting behaviour under test.
func newBM25(t *testing.T) *BM25Retriever {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test-fts5.sqlite")
	r, err := NewBM25Retriever(dbPath)
	if err != nil {
		t.Fatalf("NewBM25Retriever: %v", err)
	}
	t.Cleanup(func() { _ = r.Close() })
	return r
}

// countRows exposes the docs_fts row count for idempotency assertions.
// Tests live in-package so direct db access is fine and cheaper than
// exposing a Count method on the production surface.
func countRows(t *testing.T, r *BM25Retriever) int {
	t.Helper()
	var n int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM docs_fts").Scan(&n); err != nil {
		t.Fatalf("count rows: %v", err)
	}
	return n
}

// pathsOnly is a tiny formatting helper for readable ordering-assertion
// failure messages.
func pathsOnly(refs []Reference) []string {
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		out = append(out, filepath.Base(r.Path))
	}
	return out
}
