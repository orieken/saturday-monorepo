#!/usr/bin/env bash
set -euo pipefail

NS=test-runner

echo "Port-forwarding UI -> localhost:9000 and service -> localhost:9001"

kubectl -n $NS port-forward svc/test-runner-ui 9000:9000 &
PF1=$!
kubectl -n $NS port-forward svc/test-runner-service 9001:9001 &
PF2=$!

echo "Forwarding (PID $PF1, $PF2). Press Ctrl-C to stop."
wait $PF1 $PF2

