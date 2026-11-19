#!/usr/bin/env bash
# Moves every Playwright pin in the monorepo to one release, together, so the
# workspace packages, the runner image's runner-check and the image's
# versions.lock can never disagree. Run by .github/workflows/playwright-update.yml
# every week, and by hand when needed.
#
# What it changes:
#   - packages/*/package.json and apps/*/package.json: the dependencies and
#     devDependencies entries for @playwright/test, playwright and
#     playwright-core become ^<version>. peerDependencies are left alone: they
#     state what consumers may use, not what this repo tests with.
#   - containers/saturday-runner/runner-check/package.json: playwright-core
#     becomes exactly <version> (the Containerfile checks this against
#     versions.lock and fails the image build on a mismatch).
#   - pnpm-lock.yaml, runner-check/pnpm-lock.yaml and versions.lock are
#     regenerated, never hand-edited.
#
# It does not touch apt.lock or rpm.lock. The image runs the real Chrome and
# Edge stable channels, which move on their own schedule; refresh those with
# containers/saturday-runner/lock-packages.sh.
#
# Usage: scripts/update-playwright.sh [--version <x.y.z>] [--check]
#   --version  target release (default: the npm "latest" dist-tag)
#   --check    change nothing; exit 1 if any pin differs from the target
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUNNER_DIR="$REPO_ROOT/containers/saturday-runner"
RUNNER_CHECK_DIR="$RUNNER_DIR/runner-check"
PLAYWRIGHT_PACKAGES='["@playwright/test","playwright","playwright-core"]'

TARGET_VERSION=""
IS_CHECK_ONLY=false

parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --version) TARGET_VERSION="${2:?--version needs a value}"; shift 2 ;;
      --check) IS_CHECK_ONLY=true; shift ;;
      *) echo "update-playwright: unknown argument $1" >&2; exit 2 ;;
    esac
  done
}

resolve_target_version() {
  if [[ -z "$TARGET_VERSION" ]]; then
    TARGET_VERSION="$(npm view @playwright/test dist-tags.latest)"
  fi
  if [[ ! "$TARGET_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "update-playwright: '$TARGET_VERSION' is not a stable x.y.z release" >&2
    exit 2
  fi
}

workspace_manifests() {
  find "$REPO_ROOT/packages" "$REPO_ROOT/apps" -mindepth 2 -maxdepth 2 \
    -name package.json -not -path '*/node_modules/*' | sort
}

# Prints "<file> <section> <package> <old> -> <new>" for each pin that differs.
# Writes the file unless CHECK_ONLY=true. Keeps the file's own indentation.
rewrite_manifest() {
  local manifest="$1" range="$2"
  CHECK_ONLY="$IS_CHECK_ONLY" node -e '
    const fs = require("fs");
    const [file, range, names] = process.argv.slice(1);
    const raw = fs.readFileSync(file, "utf8");
    const pkg = JSON.parse(raw);
    let isChanged = false;
    for (const section of ["dependencies", "devDependencies"]) {
      for (const name of JSON.parse(names)) {
        const current = pkg[section] && pkg[section][name];
        if (current === undefined || current === range) continue;
        console.log(`${file} ${section} ${name} ${current} -> ${range}`);
        pkg[section][name] = range;
        isChanged = true;
      }
    }
    if (isChanged && process.env.CHECK_ONLY !== "true") {
      const indent = (raw.match(/^( +)"/m) || [, "  "])[1];
      fs.writeFileSync(file, JSON.stringify(pkg, null, indent) + (raw.endsWith("\n") ? "\n" : ""));
    }
  ' "$manifest" "$range" "$PLAYWRIGHT_PACKAGES"
}

rewrite_all_manifests() {
  local manifest
  while IFS= read -r manifest; do
    rewrite_manifest "$manifest" "^$TARGET_VERSION"
  done < <(workspace_manifests)
  rewrite_manifest "$RUNNER_CHECK_DIR/package.json" "$TARGET_VERSION"
}

regenerate_lockfiles() {
  (cd "$REPO_ROOT" && pnpm install --no-frozen-lockfile && pnpm dedupe)
  (cd "$RUNNER_CHECK_DIR" && pnpm install --lockfile-only --ignore-workspace)
  "$RUNNER_DIR/lock-versions.sh"
}

# Fails when the root lockfile resolves any Playwright package to more than
# one version or to a version other than the target, or when versions.lock
# (which the runner image builds from) is on another version.
verify_single_version() {
  local resolved
  resolved="$(grep -oE "^  '?(@playwright/test|playwright|playwright-core)@[0-9][^':(]*" \
    "$REPO_ROOT/pnpm-lock.yaml" | sed -E 's/.*@//' | sort -u)"
  if [[ "$resolved" != "$TARGET_VERSION" ]]; then
    echo "update-playwright: pnpm-lock.yaml resolves Playwright to: $(echo $resolved)" >&2
    echo "update-playwright: expected only $TARGET_VERSION" >&2
    exit 1
  fi
  if ! grep -qx "PLAYWRIGHT_VERSION=$TARGET_VERSION" "$RUNNER_DIR/versions.lock"; then
    echo "update-playwright: versions.lock is not on $TARGET_VERSION; rerun lock-versions.sh" >&2
    exit 1
  fi
}

main() {
  parse_args "$@"
  resolve_target_version
  local differences
  differences="$(rewrite_all_manifests)"
  if [[ -n "$differences" ]]; then echo "$differences"; fi
  if [[ "$IS_CHECK_ONLY" == true && -n "$differences" ]]; then exit 1; fi
  if [[ "$IS_CHECK_ONLY" == false ]]; then regenerate_lockfiles; fi
  verify_single_version
  echo "update-playwright: every pin is on $TARGET_VERSION"
}

main "$@"
