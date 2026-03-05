#!/usr/bin/env bash
set -euo pipefail

# Deploy test-runner demo stack to a local Kubernetes cluster (kind/minikube)
# - creates namespace
# - deploys test-runner-service and test-runner-ui (images must be available in the cluster)
# - deploys a ClusterIP service for both
# - (optional) loads local images into kind cluster: use 'kind load docker-image'

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NS=test-runner

echo "Applying namespace..."
kubectl apply -f "$DIR/namespace.yaml"

echo "Deploying service and UI..."
kubectl apply -f "$DIR/service-deployment.yaml"
kubectl apply -f "$DIR/ui-deployment.yaml"

echo "Waiting for pods to be ready..."
kubectl -n "$NS" wait --for=condition=available --timeout=120s deployment/test-runner-service || true
kubectl -n "$NS" wait --for=condition=available --timeout=120s deployment/test-runner-ui || true

echo "To trigger demo job run: kubectl apply -f $DIR/demo-job.yaml -n $NS"

echo "If using kind, load local images before deploying (example):"
cat <<'KIND'
# build images locally
cd /path/to/tools/demo/final/test-runner-ui
docker build -t test-runner-ui:local .
cd /path/to/tools/demo/final/test-runner-service
docker build -t test-runner-service:local .
# load into kind cluster
kind load docker-image test-runner-ui:local
kind load docker-image test-runner-service:local
KIND
