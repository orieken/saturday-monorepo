#!/usr/bin/env bash
# Shell-level tests for scripts/update-playwright.sh.
#
# Each case builds a small fake monorepo in a temp directory (two workspace
# packages, runner-check, the real lock-versions.sh) and runs the real,
# unmodified script against it. `pnpm` and `npm` are stubs placed first on
# PATH: the pnpm stub writes a root pnpm-lock.yaml that resolves Playwright
# to whatever FAKE_LOCK_VERSIONS says, so the lockfile check can be made to
# pass or fail without a registry.
#
# Usage: bash scripts/update-playwright.test.sh   (run from the repo root)
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORK_DIR="$(mktemp -d)"
FAILURES=0
trap 'rm -rf "$WORK_DIR"' EXIT

write_fake_tools() {
  mkdir -p "$WORK_DIR/bin"
  cat >"$WORK_DIR/bin/npm" <<'FAKE_NPM'
#!/usr/bin/env bash
echo "${FAKE_NPM_LATEST:-9.9.9}"
FAKE_NPM
  cat >"$WORK_DIR/bin/pnpm" <<'FAKE_PNPM'
#!/usr/bin/env bash
# Only the root `pnpm install` writes the root lockfile; dedupe and the
# runner-check install do nothing here.
if [[ "$1" == "install" && "$*" != *"--ignore-workspace"* ]]; then
  {
    echo "packages:"
    for v in ${FAKE_LOCK_VERSIONS:-}; do
      echo "  '@playwright/test@$v':"
      echo "  playwright@$v:"
    done
  } >pnpm-lock.yaml
fi
FAKE_PNPM
  chmod +x "$WORK_DIR/bin/npm" "$WORK_DIR/bin/pnpm"
}

write_fixture_repo() {
  local root="$WORK_DIR/repo"
  rm -rf "$root"
  mkdir -p "$root/scripts" "$root/packages/lib" "$root/apps/app" \
    "$root/containers/saturday-runner/runner-check"
  cp "$REPO_ROOT/scripts/update-playwright.sh" "$root/scripts/"
  cp "$REPO_ROOT/containers/saturday-runner/lock-versions.sh" "$root/containers/saturday-runner/"
  echo '{ "name": "root", "packageManager": "pnpm@10.12.1" }' >"$root/package.json"
  cat >"$root/packages/lib/package.json" <<'JSON'
{
  "name": "lib",
  "peerDependencies": { "@playwright/test": ">=1.40.0" },
  "devDependencies": { "@playwright/test": "^1.40.0", "playwright": "^1.49.0" }
}
JSON
  cat >"$root/apps/app/package.json" <<'JSON'
{
    "name": "app",
    "devDependencies": { "@playwright/test": "^1.57.0", "vite": "^5.0.0" }
}
JSON
  echo '{ "name": "runner-check", "dependencies": { "playwright-core": "1.57.0" } }' \
    >"$root/containers/saturday-runner/runner-check/package.json"
  echo "$root"
}

run_script() {
  PATH="$WORK_DIR/bin:$PATH" bash "$1/scripts/update-playwright.sh" "${@:2}" >"$WORK_DIR/out" 2>&1
}

assert_equals() {
  local description="$1" expected="$2" actual="$3"
  if [[ "$expected" == "$actual" ]]; then
    echo "ok   - $description"
  else
    echo "FAIL - $description (expected '$expected', got '$actual')"
    FAILURES=$((FAILURES + 1))
  fi
}

json_field() {
  node -p "require('$1')$2"
}

test_bumps_every_pin_and_leaves_peers_alone() {
  local root
  root="$(write_fixture_repo)"
  FAKE_LOCK_VERSIONS="1.65.1" run_script "$root" --version 1.65.1
  assert_equals "bump exits 0 when the lockfile resolves one version" 0 $?
  assert_equals "devDependency @playwright/test becomes a caret range" "^1.65.1" \
    "$(json_field "$root/packages/lib/package.json" ".devDependencies['@playwright/test']")"
  assert_equals "devDependency playwright becomes a caret range" "^1.65.1" \
    "$(json_field "$root/packages/lib/package.json" ".devDependencies.playwright")"
  assert_equals "peerDependency range is left alone" ">=1.40.0" \
    "$(json_field "$root/packages/lib/package.json" ".peerDependencies['@playwright/test']")"
  assert_equals "unrelated dependency is left alone" "^5.0.0" \
    "$(json_field "$root/apps/app/package.json" ".devDependencies.vite")"
  assert_equals "runner-check playwright-core becomes an exact pin" "1.65.1" \
    "$(json_field "$root/containers/saturday-runner/runner-check/package.json" ".dependencies['playwright-core']")"
  assert_equals "versions.lock is regenerated on the target" "PLAYWRIGHT_VERSION=1.65.1" \
    "$(grep PLAYWRIGHT_VERSION "$root/containers/saturday-runner/versions.lock")"
  assert_equals "a four-space manifest keeps its indentation" '    "name": "app",' \
    "$(sed -n 2p "$root/apps/app/package.json")"
}

test_defaults_to_npm_latest() {
  local root
  root="$(write_fixture_repo)"
  FAKE_NPM_LATEST="1.70.0" FAKE_LOCK_VERSIONS="1.70.0" run_script "$root"
  assert_equals "no --version takes npm's latest dist-tag" "^1.70.0" \
    "$(json_field "$root/apps/app/package.json" ".devDependencies['@playwright/test']")"
}

test_fails_when_lockfile_keeps_two_versions() {
  local root
  root="$(write_fixture_repo)"
  FAKE_LOCK_VERSIONS="1.57.0 1.65.1" run_script "$root" --version 1.65.1
  assert_equals "bump exits 1 when an old version survives in the lockfile" 1 $?
  # lock-versions.sh refuses first, naming both versions on their own lines.
  assert_equals "the failure names the surviving old version" 1 \
    "$(grep -cx '1.57.0' "$WORK_DIR/out")"
}

test_rejects_a_prerelease_target() {
  local root
  root="$(write_fixture_repo)"
  run_script "$root" --version 1.66.0-alpha-2026-11-19
  assert_equals "a prerelease target exits 2" 2 $?
  assert_equals "a rejected target changes no manifest" "^1.57.0" \
    "$(json_field "$root/apps/app/package.json" ".devDependencies['@playwright/test']")"
}

test_check_mode_reports_without_writing() {
  local root
  root="$(write_fixture_repo)"
  run_script "$root" --check --version 1.65.1
  assert_equals "--check exits 1 when pins differ" 1 $?
  assert_equals "--check writes nothing" "^1.57.0" \
    "$(json_field "$root/apps/app/package.json" ".devDependencies['@playwright/test']")"
  assert_equals "--check lists each differing pin" 4 "$(grep -c -- ' -> ' "$WORK_DIR/out")"
}

write_fake_tools
test_bumps_every_pin_and_leaves_peers_alone
test_defaults_to_npm_latest
test_fails_when_lockfile_keeps_two_versions
test_rejects_a_prerelease_target
test_check_mode_reports_without_writing

if [[ "$FAILURES" -gt 0 ]]; then
  echo "$FAILURES failure(s)"
  exit 1
fi
echo "all update-playwright tests passed"
