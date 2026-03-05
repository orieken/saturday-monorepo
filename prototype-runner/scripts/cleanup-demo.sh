#!/usr/bin/env bash
set -euo pipefail

# Cleanup script for the test-runner demo
# - Stops and removes containers started by docker compose
# - Removes volumes and orphan containers
# - Optionally removes the cucumber-project image if REMOVE_IMAGE=1
# - Optionally prunes Docker system (PRUNE_DOCKER=1) and builder cache (PRUNE_BUILDER=1)
# - Optionally removes generated reports (REMOVE_REPORTS=1)

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

REMOVE_IMAGE="${REMOVE_IMAGE:-0}"
PRUNE_DOCKER="${PRUNE_DOCKER:-0}"
PRUNE_BUILDER="${PRUNE_BUILDER:-0}"
REMOVE_REPORTS="${REMOVE_REPORTS:-0}"
COMPOSE_SERVICES=(mock-api web-app test-runner-service test-runner-ui)

echo "[cleanup] stopping and removing compose services"
docker compose down --volumes --remove-orphans

# Optionally prune Docker system (images, containers, networks, build cache)
if [ "$PRUNE_DOCKER" = "1" ]; then
  echo "[cleanup] running 'docker system prune -af --volumes' (this will remove unused images/containers/networks and volumes)"
  docker system prune -af --volumes || true
fi

# Optionally prune the builder cache (docker build cache)
if [ "$PRUNE_BUILDER" = "1" ]; then
  echo "[cleanup] running 'docker builder prune -af' (this will clear build cache)"
  docker builder prune -af || true
fi

# Optionally remove the cucumber image
if [ "$REMOVE_IMAGE" = "1" ]; then
  echo "[cleanup] removing image: cucumber-project:local"
  docker image rm -f cucumber-project:local || true
fi

# Optionally remove generated reports directory so fresh runs regenerate them
if [ "$REMOVE_REPORTS" = "1" ]; then
  REPORT_DIR="$ROOT_DIR/test-runner-service/reports"
  if [ -d "$REPORT_DIR" ]; then
    echo "[cleanup] removing reports directory: $REPORT_DIR"
    rm -rf "$REPORT_DIR"
  else
    echo "[cleanup] no reports directory found at $REPORT_DIR"
  fi
fi

echo "[cleanup] cleanup complete."
