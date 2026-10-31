#!/usr/bin/env bash
# Runs the demo app's Playwright suite in real Chrome, inside the Saturday
# runner image, under Xvfb. One command for the factory's pre-ship check —
# see docs/runner-image.md and containers/saturday-runner/.
#
# Usage: scripts/ui-suite-in-chrome.sh   (run from the repository root)
set -euo pipefail

EXIT_NO_DOCKER=2
EXIT_BUILD_FAILED=3
EXIT_APP_UNHEALTHY=4
EXIT_RUNNER_CHECK_FAILED=5

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE_TAG="saturday-runner:debian"
CONTAINER_NAME="ui-suite-in-chrome-$$"
PNPM_STORE_VOLUME="ui-suite-pnpm-store"
APP_HEALTH_TIMEOUT_SECONDS=180
ARTIFACTS_DIR="$REPO_ROOT/artifacts/ui-suite"

require_docker() {
  if ! command -v docker >/dev/null 2>&1; then
    echo "ui-suite-in-chrome: docker is not installed — cannot run" >&2
    exit "$EXIT_NO_DOCKER"
  fi
}

build_image() {
  if ! docker build --platform linux/amd64 --target debian \
      -t "$IMAGE_TAG" \
      -f "$REPO_ROOT/containers/saturday-runner/Containerfile" \
      "$REPO_ROOT/containers/saturday-runner"; then
    echo "ui-suite-in-chrome: image build failed" >&2
    exit "$EXIT_BUILD_FAILED"
  fi
}

start_container() {
  docker run -d --rm --platform linux/amd64 --name "$CONTAINER_NAME" \
    -v "$REPO_ROOT:/work" \
    -v "$PNPM_STORE_VOLUME:/home/runner/.pnpm-store" \
    -e "RUNNER_SLOW_MO=${RUNNER_SLOW_MO:-}" \
    -e "RUNNER_VIDEO=${RUNNER_VIDEO:-}" \
    "$IMAGE_TAG" \
    bash -lc 'sleep infinity' >/dev/null
}

stop_container() {
  docker stop "$CONTAINER_NAME" >/dev/null 2>&1 || true
}

# CI=true makes pnpm auto-answer its "recreate node_modules from scratch?"
# prompt instead of waiting on a TTY that docker exec never attaches —
# without it, a bind-mounted node_modules built for a different platform
# (e.g. the host's arch instead of the container's) hangs this step
# indefinitely, silently skipping the 180s app-health budget entirely since
# that budget starts only after this step returns.
install_dependencies() {
  docker exec -e CI=true "$CONTAINER_NAME" bash -lc '
    pnpm config set store-dir /home/runner/.pnpm-store >/dev/null
    cd /work && pnpm install --frozen-lockfile
    # The suite imports workspace packages from their dist/, which a fresh
    # checkout does not have. Build only what the demo app depends on.
    pnpm --filter "@orieken/ye-olde-magic-shop^..." -r build
  '
}

start_app_and_mock_api() {
  docker exec -d "$CONTAINER_NAME" bash -lc '
    cd /work/apps/ye-olde-magic-shop && npm run dev >/tmp/app.log 2>&1
  '
}

# PRECONDITION: this is a bounded, one-shot readiness poll of a process this
# script itself started (the dev server launched by start_app_and_mock_api),
# not resilience against an unreliable remote dependency — same shape as
# containers/saturday-runner/entrypoint.sh's wait_for_display.
# ENFORCED-BY: judgment-only, per architecture-notes.md's Stability Design
# decision; enforcement is the 180s bound (APP_HEALTH_TIMEOUT_SECONDS) plus
# the distinct EXIT_APP_UNHEALTHY exit code, not a CircuitBreaker or
# ExponentialBackoffStrategy.
wait_for_app_healthy() {
  local waited=0
  until docker exec "$CONTAINER_NAME" bash -lc \
      'curl -sf --max-time 5 http://localhost:8000/ >/dev/null && curl -sf --max-time 5 http://localhost:8001/api/health >/dev/null' \
      >/dev/null 2>&1; do
    waited=$((waited + 1))
    if [ "$waited" -ge "$APP_HEALTH_TIMEOUT_SECONDS" ]; then
      echo "ui-suite-in-chrome: app/mock-api did not become healthy within ${APP_HEALTH_TIMEOUT_SECONDS}s" >&2
      docker exec "$CONTAINER_NAME" bash -lc 'tail -n 50 /tmp/app.log' >&2 || true
      exit "$EXIT_APP_UNHEALTHY"
    fi
    sleep 1
  done
}

run_runner_check() {
  if ! docker exec "$CONTAINER_NAME" runner-check; then
    echo "ui-suite-in-chrome: runner-check failed — WebGL unavailable in the runner image" >&2
    exit "$EXIT_RUNNER_CHECK_FAILED"
  fi
}

run_chrome_suite() {
  docker exec "$CONTAINER_NAME" bash -lc \
    'cd /work/apps/ye-olde-magic-shop && npx playwright test --project=chrome'
}

copy_report() {
  mkdir -p "$ARTIFACTS_DIR"
  rm -rf "$ARTIFACTS_DIR/playwright-report"
  cp -r "$REPO_ROOT/apps/ye-olde-magic-shop/playwright-report" "$ARTIFACTS_DIR/playwright-report" \
    || echo "ui-suite-in-chrome: warning: failed to copy playwright report" >&2
}

main() {
  require_docker
  build_image
  trap stop_container EXIT
  start_container
  install_dependencies
  start_app_and_mock_api
  wait_for_app_healthy
  run_runner_check
  local suite_exit=0
  run_chrome_suite || suite_exit=$?
  copy_report
  exit "$suite_exit"
}

main "$@"
