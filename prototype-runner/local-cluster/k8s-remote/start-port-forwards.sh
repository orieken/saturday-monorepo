#!/usr/bin/env bash
set -euo pipefail

KUBECONFIG=${KUBECONFIG:-/tmp/kind-kubeconfig.yaml}
NAMESPACE=test-runner

# Start port-forwards in background and save PIDs
# UI (cluster svc test-runner-ui -> localhost:9000)
kubectl --kubeconfig="$KUBECONFIG" -n "$NAMESPACE" port-forward svc/test-runner-ui 9000:9000 > /tmp/portfwd-ui.log 2>&1 &
echo $! > /tmp/portfwd-ui.pid

# web-app (svc/web-app -> localhost:8000)
kubectl --kubeconfig="$KUBECONFIG" -n "$NAMESPACE" port-forward svc/web-app 8000:8000 > /tmp/portfwd-webapp.log 2>&1 &
echo $! > /tmp/portfwd-webapp.pid

# backend service (svc/test-runner-service -> localhost:9001)
kubectl --kubeconfig="$KUBECONFIG" -n "$NAMESPACE" port-forward svc/test-runner-service 9001:9001 > /tmp/portfwd-backend.log 2>&1 &
echo $! > /tmp/portfwd-backend.pid

# mock-api (svc/mock-api -> localhost:8001)
kubectl --kubeconfig="$KUBECONFIG" -n "$NAMESPACE" port-forward svc/mock-api 8001:8001 > /tmp/portfwd-mockapi.log 2>&1 &
echo $! > /tmp/portfwd-mockapi.pid

cat <<EOF
Started port-forwards:
  UI      -> http://localhost:9000 (pid $(cat /tmp/portfwd-ui.pid))
  Web     -> http://localhost:8000 (pid $(cat /tmp/portfwd-webapp.pid))
  API     -> http://localhost:9001 (pid $(cat /tmp/portfwd-backend.pid))
  Mock API-> http://localhost:8001 (pid $(cat /tmp/portfwd-mockapi.pid))
Logs: /tmp/portfwd-*.log
To stop: ./stop-port-forwards.sh
EOF
