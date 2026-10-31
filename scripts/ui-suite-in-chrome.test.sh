#!/usr/bin/env bash
# Shell-level tests for scripts/ui-suite-in-chrome.sh's exit-code contract:
# a failed image build or an app that never becomes healthy must exit
# non-zero with a distinct, named cause, separate from a Playwright test
# failure's own exit code (which must pass through unchanged).
#
# Runs the real, unmodified script against a stubbed `docker` placed first on
# PATH, so every branch (no docker / build failed / app unhealthy /
# runner-check failed / suite failure passthrough / full success) is
# exercised without a real Docker daemon. `sleep` is also stubbed for the
# "app unhealthy" case so the real 180s bound is exercised without a 180s
# wall-clock wait.
#
# Usage: bash scripts/ui-suite-in-chrome.test.sh   (run from the repo root)
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SCRIPT_UNDER_TEST="$REPO_ROOT/scripts/ui-suite-in-chrome.sh"
FAKE_BIN_DIR="$(mktemp -d)"
FAILURES=0

cleanup() {
  rm -rf "$FAKE_BIN_DIR"
  rm -rf "$REPO_ROOT/artifacts/ui-suite"
}
trap cleanup EXIT

write_fake_docker() {
  # Split into small single-purpose functions (rather than one deeply nested
  # case) so each piece of this stub stays under the same complexity-6
  # budget as production code (design-principles.md / architecture
  # guardrails), even though it's test infrastructure.
  cat >"$FAKE_BIN_DIR/docker" <<'FAKE_DOCKER'
#!/usr/bin/env bash
dispatch_exec_command() {
  case "$1" in
    *"curl"*) exit "${FAKE_DOCKER_HEALTH_EXIT:-0}" ;;
    *"playwright test"*) exit "${FAKE_DOCKER_SUITE_EXIT:-0}" ;;
    *"pnpm install"*) exit "${FAKE_DOCKER_INSTALL_EXIT:-0}" ;;
    *) exit 0 ;;
  esac
}

dispatch_exec() {
  shift # drop "exec"
  while [[ "$1" == "-d" || "$1" == "-e" ]]; do
    if [[ "$1" == "-d" ]]; then
      exit "${FAKE_DOCKER_EXEC_DETACHED_EXIT:-0}"
    fi
    shift 2 # drop "-e" and its VAR=value
  done
  shift # drop container name
  if [[ "$1" == "runner-check" ]]; then
    exit "${FAKE_DOCKER_RUNNER_CHECK_EXIT:-0}"
  fi
  dispatch_exec_command "$*"
}

case "$1" in
  build) exit "${FAKE_DOCKER_BUILD_EXIT:-0}" ;;
  run) exit "${FAKE_DOCKER_RUN_EXIT:-0}" ;;
  exec) dispatch_exec "$@" ;;
  stop) exit 0 ;;
  *) exit 0 ;;
esac
FAKE_DOCKER
  chmod +x "$FAKE_BIN_DIR/docker"
}

write_fake_sleep() {
  cat >"$FAKE_BIN_DIR/sleep" <<'FAKE_SLEEP'
#!/usr/bin/env bash
exit 0
FAKE_SLEEP
  chmod +x "$FAKE_BIN_DIR/sleep"
}

assert_exit_code() {
  local description="$1" expected="$2" actual="$3"
  if [[ "$actual" -eq "$expected" ]]; then
    echo "PASS: $description (exit $actual)"
  else
    echo "FAIL: $description (expected exit $expected, got $actual)" >&2
    FAILURES=$((FAILURES + 1))
  fi
}

assert_contains() {
  local description="$1" needle="$2" haystack="$3"
  if [[ "$haystack" == *"$needle"* ]]; then
    echo "PASS: $description"
  else
    echo "FAIL: $description -- stderr did not contain: $needle" >&2
    echo "--- captured stderr ---" >&2
    echo "$haystack" >&2
    FAILURES=$((FAILURES + 1))
  fi
}

test_no_docker() {
  local stderr exit_code
  stderr="$(PATH="/usr/bin:/bin" bash "$SCRIPT_UNDER_TEST" 2>&1 1>/dev/null)"
  exit_code=$?
  assert_exit_code "no docker: distinct exit code" 2 "$exit_code"
  assert_contains "no docker: stderr names the cause" "docker is not installed" "$stderr"
}

test_build_failed() {
  local stderr exit_code
  stderr="$(PATH="$FAKE_BIN_DIR:/usr/bin:/bin" FAKE_DOCKER_BUILD_EXIT=1 \
    bash "$SCRIPT_UNDER_TEST" 2>&1 1>/dev/null)"
  exit_code=$?
  assert_exit_code "build failed: distinct exit code" 3 "$exit_code"
  assert_contains "build failed: stderr names the cause" "image build failed" "$stderr"
}

test_app_unhealthy() {
  local stderr exit_code
  stderr="$(PATH="$FAKE_BIN_DIR:/usr/bin:/bin" FAKE_DOCKER_HEALTH_EXIT=1 \
    bash "$SCRIPT_UNDER_TEST" 2>&1 1>/dev/null)"
  exit_code=$?
  assert_exit_code "app unhealthy: distinct exit code" 4 "$exit_code"
  assert_contains "app unhealthy: stderr names the cause" "did not become healthy within 180s" "$stderr"
}

test_runner_check_failed() {
  local stderr exit_code
  stderr="$(PATH="$FAKE_BIN_DIR:/usr/bin:/bin" FAKE_DOCKER_RUNNER_CHECK_EXIT=1 \
    bash "$SCRIPT_UNDER_TEST" 2>&1 1>/dev/null)"
  exit_code=$?
  assert_exit_code "runner-check failed: distinct exit code" 5 "$exit_code"
  assert_contains "runner-check failed: stderr names the cause" "runner-check failed" "$stderr"
}

test_suite_failure_passthrough() {
  local exit_code
  PATH="$FAKE_BIN_DIR:/usr/bin:/bin" FAKE_DOCKER_SUITE_EXIT=1 \
    bash "$SCRIPT_UNDER_TEST" >/dev/null 2>&1
  exit_code=$?
  assert_exit_code "playwright test failure: suite exit code passes through unchanged (not remapped to 2/3/4/5)" 1 "$exit_code"
}

test_suite_success() {
  local exit_code
  PATH="$FAKE_BIN_DIR:/usr/bin:/bin" \
    bash "$SCRIPT_UNDER_TEST" >/dev/null 2>&1
  exit_code=$?
  assert_exit_code "full happy path: exits 0" 0 "$exit_code"
}

main() {
  write_fake_docker
  write_fake_sleep

  test_no_docker
  test_build_failed
  test_app_unhealthy
  test_runner_check_failed
  test_suite_failure_passthrough
  test_suite_success

  echo ""
  if [[ "$FAILURES" -eq 0 ]]; then
    echo "All exit-code contract cases passed."
    exit 0
  else
    echo "$FAILURES case(s) failed." >&2
    exit 1
  fi
}

main "$@"
