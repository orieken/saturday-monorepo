#!/usr/bin/env bash
# Update this repo's installed loom framework to the version embedded in `loom`.
# Usage: scripts/update-framework.sh [--apply]    (default: report only, change nothing)
# Docs: docs/updating-loom.md
set -euo pipefail

LOOM="${LOOM:-loom}"
APPLY=false
[[ "${1:-}" == "--apply" ]] && APPLY=true
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

require_clean_tree() {
  if [[ -n "$(git status --porcelain)" ]]; then
    echo "working tree is not clean; commit or stash first" >&2
    exit 2
  fi
}

dry_run() { "$LOOM" install --target . --copy --dry-run > "$WORK/dry.log" 2>&1; }

detected_platforms() { grep '✓' "$WORK/dry.log" | awk '{print $3}'; }

unowned_paths() {
  grep 'not installed by loom' "$WORK/dry.log" \
    | sed -E 's/^\[dry-run\] +skipped (.*) \(not installed by loom.*/\1/'
}

# Files loom ships today, for every platform detected here, from a scratch install.
shipped_paths() {
  local scratch="$WORK/fresh"
  mkdir -p "$scratch" && git -C "$scratch" init -q
  local platform
  for platform in "$@"; do
    "$LOOM" install --target "$scratch" --copy --platform "$platform" > /dev/null 2>&1
  done
  (cd "$scratch" && find . -type f -not -path './.git/*' -not -name .loom-manifest.json | sed 's|^\./||' | sort)
}

# Tracked files in framework locations that loom no longer ships.
stale_paths() {
  git ls-files | grep -E '^(\.claude/(agents|skills|rules)/|\.agents/|\.cursor/(agents|skills|rules)/|\.github/instructions/|\.roo/|AGENTS\.md$|\.windsurfrules$|\.cursorrules$|\.openai\.md$|\.roomodes$|\.github/copilot-instructions\.md$)' \
    | sort | comm -23 - "$WORK/shipped.txt"
}

main() {
  require_clean_tree
  dry_run
  head -1 "$WORK/dry.log"
  mapfile -t platforms < <(detected_platforms)
  unowned_paths > "$WORK/unowned.txt"
  echo "platforms detected: ${platforms[*]}"
  echo "existing files loom would skip as not its own: $(wc -l < "$WORK/unowned.txt")"
  if [[ "$APPLY" == false ]]; then
    echo "report only; run with --apply to update"
    return 0
  fi
  # Remove exactly what loom would write but cannot claim, then install per platform so
  # detection does not depend on the files just removed. The manifest merges across runs.
  xargs -a "$WORK/unowned.txt" -r git rm -q --cached --ignore-unmatch --
  xargs -a "$WORK/unowned.txt" -r rm -f --
  local platform
  for platform in "${platforms[@]}"; do
    "$LOOM" install --target . --copy --platform "$platform" > "$WORK/install-$platform.log" 2>&1 \
      || { cat "$WORK/install-$platform.log" >&2; exit 1; }
  done
  git add -A
  dry_run
  echo "still not loom's after install (this repo's own files): $(grep -c 'not installed by loom' "$WORK/dry.log" || true)"
  shipped_paths "${platforms[@]}" > "$WORK/shipped.txt"
  stale_paths > stale-framework-files.txt
  echo "stale framework files (no longer shipped): $(wc -l < stale-framework-files.txt), listed in stale-framework-files.txt"
  git status --short | cut -c1-2 | sort | uniq -c
}

main
