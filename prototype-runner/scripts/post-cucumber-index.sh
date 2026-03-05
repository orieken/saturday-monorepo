#!/usr/bin/env bash
set -euo pipefail

# Usage: post-cucumber-index.sh PATH_TO_JSON
# Posts a CucumberIndex JSON to the running test-runner-service at localhost:9001

JSON_PATH=${1:-/tmp/cuke-index.json}
if [ ! -f "$JSON_PATH" ]; then
  echo "JSON file not found: $JSON_PATH"
  exit 1
fi

URL=${TEST_RUNNER_URL:-http://localhost:9001}
API_PATH="$URL/api/cucumber/index"

echo "Posting $JSON_PATH to $API_PATH"

# Do a verbose POST and print headers+body
curl -i -X POST -H "Content-Type: application/json" --data-binary "@$JSON_PATH" "$API_PATH"

echo "\nDone"

