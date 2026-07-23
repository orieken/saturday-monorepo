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

// skippedDirNames lists directory basenames the walk should never descend
// into. Consolidates the historically-per-file `skippedDirs` maps into a
// single package-level set so a change here immediately affects every
// analyzer. Kept as a set (map[T]struct{}) rather than a slice so
// membership checks are O(1) and read as a set at the call site.
var skippedDirNames = map[string]struct{}{
	"node_modules": {},
	".git":         {},
	"dist":         {},
	"build":        {},
	"vendor":       {},
	".venv":        {},
}

// skipUninterestingDir returns filepath.SkipDir for directory names the
// walk should skip: any name in skippedDirNames, plus any hidden
// dot-directory below the root. The root itself is exempt from the
// dot-directory rule so a caller can scan a project whose absolute path
// happens to include a hidden ancestor (e.g. `~/.dotfiles/project`).
//
// Signature intentionally matches the closure the two prior analyzers
// each used before this extraction, so the callsite change is a straight
// swap — no arg-shape churn at the call sites.
func skipUninterestingDir(root, current, name string) error {
	if _, skip := skippedDirNames[name]; skip {
		return filepath.SkipDir
	}
	if strings.HasPrefix(name, ".") && current != root {
		return filepath.SkipDir
	}
	return nil
}
