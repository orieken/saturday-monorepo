#!/usr/bin/env bash
set -euo pipefail

# kind-prepare.sh
# Builds demo images, loads them into the kind cluster, and (re)applies basic k8s assets (NS, RBAC, Deployment).
# Usage: ./kind-prepare.sh [KIND_CLUSTER_NAME]

KIND_CLUSTER_NAME=${1:-kind}
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
CUKE_PROJECT_DIR="$ROOT_DIR/cucumber-project"
SERVICE_DIR="$ROOT_DIR/test-runner-service"
RBAC_FILE="$ROOT_DIR/local-cluster/test-runner-rbac.yaml"
NAMESPACE_FILE="$ROOT_DIR/local-cluster/k8s-demo/namespace.yaml"
SERVICE_DEPLOY_FILE="$ROOT_DIR/local-cluster/k8s-demo/service-deployment.yaml"

echo "Using kind cluster: $KIND_CLUSTER_NAME"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind not found in PATH. Install kind and retry." >&2
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "docker not found in PATH. Install Docker and retry." >&2
  exit 1
fi

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl not found in PATH. Install kubectl and retry." >&2
  exit 1
fi

# Build images
if [ -d "$CUKE_PROJECT_DIR" ] && [ -f "$CUKE_PROJECT_DIR/Dockerfile" ]; then
  echo "Building cucumber-project:local image..."
  (cd "$CUKE_PROJECT_DIR" && DOCKER_BUILDKIT=1 docker build -t cucumber-project:local .)
else
  echo "Skipping cucumber-project image build (directory or Dockerfile missing at $CUKE_PROJECT_DIR)" >&2
fi

if [ -d "$SERVICE_DIR" ] && [ -f "$SERVICE_DIR/Dockerfile" ]; then
  echo "Building test-runner-service:local image..."
  (cd "$SERVICE_DIR" && DOCKER_BUILDKIT=1 docker build -t test-runner-service:local .)
else
  echo "Skipping test-runner-service image build (directory or Dockerfile missing at $SERVICE_DIR)" >&2
fi

# Try loading into the named kind cluster
echo "Loading images into kind cluster '$KIND_CLUSTER_NAME'..."
if kind get clusters | grep -q "^$KIND_CLUSTER_NAME$"; then
  # Load images that may exist locally
  docker image inspect cucumber-project:local >/dev/null 2>&1 && kind load docker-image cucumber-project:local --name "$KIND_CLUSTER_NAME" || true
  docker image inspect test-runner-service:local >/dev/null 2>&1 && kind load docker-image test-runner-service:local --name "$KIND_CLUSTER_NAME" || true
else
  echo "Cluster '$KIND_CLUSTER_NAME' not found. Listing available kind clusters:" >&2
  kind get clusters || true
  echo "Trying default 'kind load docker-image' (no --name) as fallback"
  docker image inspect cucumber-project:local >/dev/null 2>&1 && kind load docker-image cucumber-project:local || true
  docker image inspect test-runner-service:local >/dev/null 2>&1 && kind load docker-image test-runner-service:local || true
fi

# Ensure namespace exists (apply manifest if present)
if [ -f "$NAMESPACE_FILE" ]; then
  echo "Ensuring namespace exists: $NAMESPACE_FILE"
  kubectl apply -f "$NAMESPACE_FILE"
else
  echo "Namespace manifest not found at $NAMESPACE_FILE; make sure namespace 'test-runner' exists." >&2
fi

# Apply RBAC (if file exists)
if [ -f "$RBAC_FILE" ]; then
  echo "Applying RBAC: $RBAC_FILE"
  kubectl apply -f "$RBAC_FILE"
else
  echo "RBAC file not found at $RBAC_FILE; create it first (already present in repo under local-cluster/test-runner-rbac.yaml)" >&2
fi

# Ensure Deployment is applied and restart to pick latest image
if [ -f "$SERVICE_DEPLOY_FILE" ]; then
  echo "Applying Deployment: $SERVICE_DEPLOY_FILE"
  kubectl apply -f "$SERVICE_DEPLOY_FILE"
  echo "Restarting Deployment to pick up freshly loaded image..."
  kubectl -n test-runner rollout restart deployment/test-runner-service || true
  echo "Waiting for rollout to complete..."
  kubectl -n test-runner rollout status deployment/test-runner-service --timeout=90s || true
  echo "Current pods:"
  kubectl -n test-runner get pods -l app=test-runner-service -o wide || true
else
  echo "Deployment manifest not found at $SERVICE_DEPLOY_FILE" >&2
fi

echo "Done. Next:"
echo "- Ensure port-forward is active to the service (see final/local-cluster/k8s-demo/start-port-forwards.sh)."
echo "- POST a run to http://localhost:9001/api/runs with executor: \"k8s\" and then use 'kubectl -n test-runner get pods -l run-id=<id>' to watch the job pod start."

