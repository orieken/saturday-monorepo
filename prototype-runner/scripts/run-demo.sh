#!/usr/bin/env bash
set -euo pipefail

# Demo automation script for the test-runner stack
# - Builds `cucumber-project:local` if not present
# - Starts compose services (mock-api, web-app, test-runner-service, test-runner-ui)
# - Registers cucumber index (if present)
# - Triggers a run for a default scenario (or $SCENARIO_ID)
# - Polls run status until finished and prints report URL
# - Optionally opens the report in the host OS browser if OPEN_REPORT=1

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT_DIR"

# Configuration
IMAGE_NAME="cucumber-project:local"
COMPOSE_SERVICES=(mock-api web-app test-runner-service test-runner-ui)
SCENARIO_ID="${SCENARIO_ID:-webapp_cart.feature:9}"
SUITE_ID="final-cucumber-project"
FRAMEWORK="cucumber"
MAX_WAIT_SECONDS=120
POLL_INTERVAL=2

# Open report automatically if set to 1
OPEN_REPORT="${OPEN_REPORT:-0}"

echo "[demo] root: $ROOT_DIR"

# Build the cucumber image if missing (skip with SKIP_BUILD=1)
if [ "${SKIP_BUILD:-0}" != "1" ]; then
  if [ -z "$(docker images -q $IMAGE_NAME 2>/dev/null || true)" ]; then
    echo "[demo] building $IMAGE_NAME..."
    (cd "$ROOT_DIR/cucumber-project" && docker build -t $IMAGE_NAME .)
  else
    echo "[demo] image $IMAGE_NAME already present, skipping build"
  fi
else
  echo "[demo] SKIP_BUILD=1 set; skipping image build"
fi

# Start compose services
echo "[demo] starting compose services: ${COMPOSE_SERVICES[*]}"
docker compose up --build -d ${COMPOSE_SERVICES[*]}

# Wait for test-runner-service to become healthy (simple HTTP check)
SERVICE_URL="http://localhost:9001"
WAITED=0
echo "[demo] waiting for test-runner-service at $SERVICE_URL..."
until curl -sS --fail "$SERVICE_URL/api/frameworks" >/dev/null 2>&1 || [ $WAITED -ge $MAX_WAIT_SECONDS ]; do
  sleep 1
  WAITED=$((WAITED+1))
  printf "."
done
if [ $WAITED -ge $MAX_WAIT_SECONDS ]; then
  echo "\n[demo] timeout waiting for test-runner-service; check docker compose logs"
  exit 1
fi

echo "\n[demo] test-runner-service is up"

# Register cucumber index if file exists
INDEX_FILE="$ROOT_DIR/test-runner-service/data/cucumber_index.json"
if [ -f "$INDEX_FILE" ]; then
  echo "[demo] registering cucumber index from $INDEX_FILE"
  curl -sS -X POST "$SERVICE_URL/api/cucumber/index" -H "Content-Type: application/json" --data-binary "@$INDEX_FILE" | jq . || true
else
  echo "[demo] no index file found at $INDEX_FILE; skipping registration"
fi

# Trigger a run
echo "[demo] requesting run: framework=$FRAMEWORK suite=$SUITE_ID scenario=$SCENARIO_ID"
RUN_JSON=$(curl -sS -X POST "$SERVICE_URL/api/runs" -H "Content-Type: application/json" -d "$(jq -n --arg f "$FRAMEWORK" --arg s "$SUITE_ID" --arg sc "$SCENARIO_ID" '{framework:$f, suiteId:$s, scenarioId:$sc}')")
RUN_ID=$(echo "$RUN_JSON" | jq -r '.id // .runId // empty')
if [ -z "$RUN_ID" ] || [ "$RUN_ID" = "null" ]; then
  echo "[demo] failed to start run, response:"
  echo "$RUN_JSON" | jq . || true
  exit 1
fi

echo "[demo] started run: $RUN_ID"

# Poll run status
while true; do
  sleep $POLL_INTERVAL
  STATUS_JSON=$(curl -sS "$SERVICE_URL/api/runs/$RUN_ID" || true)
  STATUS=$(echo "$STATUS_JSON" | jq -r '.status // empty')
  echo "[demo] run $RUN_ID status: ${STATUS:-unknown}"
  if [ "${STATUS:-}" = "running" ]; then
    continue
  fi
  # finished
  REPORT_URL=$(echo "$STATUS_JSON" | jq -r '.reportUrl // empty')
  echo "[demo] run finished with status: ${STATUS:-unknown}"
  if [ -n "$REPORT_URL" ]; then
    # If reportUrl is a relative path, make it absolute to localhost
    if [[ "$REPORT_URL" == /* ]]; then
      REPORT_URL="http://localhost:9001${REPORT_URL}"
    fi
    echo "[demo] report available at: $REPORT_URL"
    if [ "${OPEN_REPORT}" = "1" ]; then
      echo "[demo] opening report in default browser..."
      if command -v open >/dev/null 2>&1; then
        open "$REPORT_URL" || true
      elif command -v xdg-open >/dev/null 2>&1; then
        xdg-open "$REPORT_URL" || true
      else
        echo "[demo] no OS open command found (tried 'open' and 'xdg-open'), please open the URL manually"
      fi
    fi
  else
    echo "[demo] no report URL provided by service"
  fi
  break
done

echo "[demo] demo finished. Containers left running for inspection; use 'docker compose down' to stop them."
