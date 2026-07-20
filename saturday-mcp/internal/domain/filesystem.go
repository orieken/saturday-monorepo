package domain

// FileWriter is the domain seam for writing a byte payload to an absolute
// path on disk. It exists so tools that need a single low-level write
// (e.g. GenerateDocumentationTool emitting rendered markdown to a
// user-supplied outputPath) don't have to reach for os.WriteFile
// directly — architecture-guardrails #1 (Clean Architecture dependency
// direction) forbids the use-case layer from depending on syscalls.
//
// The adapter lives at internal/adapters/filesystem/os_filesystem.go and
// is the only implementation that touches os.WriteFile for this shape.
//
// Note: this is intentionally separate from internal/filewriter, which is
// a higher-level relative-path-plus-write-mode abstraction the generator
// tools use. Both exist because they solve different problems — this
// interface is the minimal one-shot write; filewriter.FileWriter carries
// baseDir, WriteModeOverwrite/Skip/Fail, and dryRun. Merging them would
// force one caller to inherit machinery it doesn't want.
//
// Contract:
//   - path is absolute. The adapter does not resolve, join, or validate
//     path traversal — the caller owns that. (filewriter.FileWriter does
//     the traversal check because it owns the baseDir seam; this
//     interface leaves it to the caller.)
//   - content is written with 0644 permissions. Adapters MAY create
//     parent directories, but the current OS adapter does not — callers
//     that need mkdir semantics should compose that themselves.
type FileWriter interface {
	WriteFile(path string, content []byte) error
}
