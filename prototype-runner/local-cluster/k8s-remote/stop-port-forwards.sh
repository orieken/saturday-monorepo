#!/usr/bin/env bash
set -euo pipefail

pids=(/tmp/portfwd-ui.pid /tmp/portfwd-webapp.pid /tmp/portfwd-backend.pid /tmp/portfwd-mockapi.pid)
for p in "${pids[@]}"; do
  if [ -f "$p" ]; then
    pid=$(cat "$p")
    # validate pid is numeric
    if ! [[ "$pid" =~ ^[0-9]+$ ]]; then
      echo "Invalid pid in $p: '$pid' — removing stale pid file"
      rm -f "$p"
      continue
    fi

    if ps -p "$pid" > /dev/null 2>&1; then
      echo "Stopping pid $pid from $p"
      kill "$pid" || true
      # wait up to 5 seconds for process to exit
      for i in {1..10}; do
        if ! ps -p "$pid" > /dev/null 2>&1; then
          break
        fi
        sleep 0.5
      done
      if ps -p "$pid" > /dev/null 2>&1; then
        echo "Process $pid did not exit gracefully, sending SIGKILL"
        kill -9 "$pid" || true
      fi
    else
      echo "No process $pid found for $p, removing pid file"
    fi

    rm -f "$p"
  fi
done

# remove any leftover logs
rm -f /tmp/portfwd-*.log || true

echo "Port-forwards stopped"
