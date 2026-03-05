#!/usr/bin/env bash
set -euo pipefail

# Writes the current kubeconfig into the project-local kubeconfig file
# Usage: ./export-kubeconfig.sh [--from-kind <cluster-name>]

OUT_FILE="$(pwd)/kubeconfig"

if [[ ${1-} == "--from-kind" ]]; then
  CLUSTER_NAME=${2:-kind}
  echo "Exporting kubeconfig from kind cluster '${CLUSTER_NAME}' to ${OUT_FILE}"
  kind get kubeconfig --name "${CLUSTER_NAME}" > "${OUT_FILE}"
else
  echo "Exporting current kubectl config to ${OUT_FILE}"
  kubectl config view --raw > "${OUT_FILE}"
fi

chmod 600 "${OUT_FILE}"

echo "Wrote kubeconfig to ${OUT_FILE}"

