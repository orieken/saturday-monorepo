// Package analyzers - walkutil.go
//
// Shared filepath.Walk helpers used by every analyzer that scans a
// project tree. Extracted after the third analyzer in this package
// (dependency_boundary_analyzer.go) grew its own copy of the same skip
// logic — accessibility_analyzer.go and ubiquitous_language_analyzer.go
// were the first two, this file makes them one implementation instead of
// three. Op 4 handoff notes explicitly flagged this extraction as the
// right call once a third copy appeared.
//
// The skip set is deliberately conservative: node_modules, .git, dist,
// build, vendor, .venv, plus any hidden dot-directory below the root.
// Individual analyzers can still add extra Extension filters on top of
// the walker's directory-level skipping — this file only owns the
// directory decisions every analyzer wants to share.
package analyzers

import (
	"path/filepath"
	"strings"
)

// SkippedDirNames lists directory basenames the walk should never descend
// into. Consolidates the historically-per-file `skippedDirs` maps into a
// single package-level set so a change here immediately affects every
// analyzer. Kept as a set (map[T]struct{}) rather than a slice so
// membership checks are O(1) and read as a set at the call site.
//
// Exported so cross-package callers (see internal/tools/bm25_retriever.go,
// which walks an installed project's docs corpus and needs the same skip
// set) get one source of truth rather than a second copy that can drift.
var SkippedDirNames = map[string]struct{}{
	"node_modules": {},
	".git":         {},
	"dist":         {},
	"build":        {},
	"vendor":       {},
	".venv":        {},
}

// SkipUninterestingDir returns filepath.SkipDir for directory names the
// walk should skip: any name in SkippedDirNames, plus any hidden
// dot-directory below the root. The root itself is exempt from the
// dot-directory rule so a caller can scan a project whose absolute path
// happens to include a hidden ancestor (e.g. `~/.dotfiles/project`).
//
// Exported (from `skipUninterestingDir`) so internal/tools/ callers can
// reuse the same walk-skip discipline the analyzers use — see
// bm25_retriever.go for the second consumer.
func SkipUninterestingDir(root, current, name string) error {
	if _, skip := SkippedDirNames[name]; skip {
		return filepath.SkipDir
	}
	if strings.HasPrefix(name, ".") && current != root {
		return filepath.SkipDir
	}
	return nil
}
