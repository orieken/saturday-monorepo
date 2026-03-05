#!/usr/bin/env bash
set -euo pipefail

# run-local-k8s.sh
# Helper to build local images, load them into a kind cluster, deploy manifests,
# port-forward the test-runner service, trigger a run (executor=k8s), stream logs,
# and copy reports locally.
#
# Usage examples:
#   ./run-local-k8s.sh --all
#   ./run-local-k8s.sh --build --load --deploy --port-forward --run --stream
#   ./run-local-k8s.sh --copy-reports
#
# The script expects to be placed in `final/local-cluster` and will resolve paths
# relative to that location.

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
K8S_MANIFEST_DIR="${K8S_MANIFEST_DIR:-$ROOT_DIR/local-cluster/k8s-remote}"
SERVICE_DIR="$ROOT_DIR/../apps/console"
UI_DIR="$ROOT_DIR/../apps/cartridge"
MOCK_API_DIR="$ROOT_DIR/../apps/mock-api"
WEB_APP_DIR="$ROOT_DIR/../apps/ye-olde-magic-shop"
RUNNER_DIR="$ROOT_DIR/../apps/ye-olde-magic-shop"

KIND_CLUSTER="${KIND_CLUSTER:-kind}"
IMAGE_TAG="${IMAGE_TAG:-local}"
NAMESPACE="${NAMESPACE:-test-runner}"
API_PORT_LOCAL="${API_PORT_LOCAL:-9001}"

PORT_FORWARD_PID_FILE=".test-runner-port-forward.pid"

# Resolve tool command paths (prefer explicit env overrides, then PATH, then ~/bin)
KIND_CMD="${KIND_CMD:-}"
KUBECTL_CMD="${KUBECTL_CMD:-}"

if [ -z "$KIND_CMD" ]; then
  if command -v kind >/dev/null 2>&1; then
    KIND_CMD=$(command -v kind)
  elif [ -x "$HOME/bin/kind" ]; then
    KIND_CMD="$HOME/bin/kind"
  else
    KIND_CMD=""
  fi
fi

if [ -z "$KUBECTL_CMD" ]; then
  if command -v kubectl >/dev/null 2>&1; then
    KUBECTL_CMD=$(command -v kubectl)
  elif [ -x "$HOME/bin/kubectl" ]; then
    KUBECTL_CMD="$HOME/bin/kubectl"
  else
    KUBECTL_CMD=""
  fi
fi

# Ensure a binary by downloading to ~/bin if not present
ensure_kind() {
  if [ -n "$KIND_CMD" ] && [ -x "$KIND_CMD" ]; then
    return 0
  fi
  echo "kind not found on PATH; downloading to $HOME/bin/kind"
  mkdir -p "$HOME/bin"
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH_RAW=$(uname -m)
  case "$ARCH_RAW" in
    x86_64) ARCH=amd64 ;; amd64) ARCH=amd64 ;; aarch64) ARCH=arm64 ;; arm64) ARCH=arm64 ;; *) ARCH=$ARCH_RAW;;
  esac
  KIND_VERSION="v0.20.0"
  URL="https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-${OS}-${ARCH}"
  curl -fsSL -o "$HOME/bin/kind" "$URL"
  chmod +x "$HOME/bin/kind"
  KIND_CMD="$HOME/bin/kind"
}

ensure_kubectl() {
  if [ -n "$KUBECTL_CMD" ] && [ -x "$KUBECTL_CMD" ]; then
    return 0
  fi
  echo "kubectl not found on PATH; downloading to $HOME/bin/kubectl"
  mkdir -p "$HOME/bin"
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH_RAW=$(uname -m)
  case "$ARCH_RAW" in
    x86_64) ARCH=amd64 ;; amd64) ARCH=amd64 ;; aarch64) ARCH=arm64 ;; arm64) ARCH=arm64 ;; *) ARCH=$ARCH_RAW;;
  esac
  STABLE=$(curl -fsSL https://dl.k8s.io/release/stable.txt)
  if [ "$OS" = "darwin" ]; then
    KURL="https://dl.k8s.io/release/${STABLE}/bin/darwin/${ARCH}/kubectl"
  else
    KURL="https://dl.k8s.io/release/${STABLE}/bin/linux/${ARCH}/kubectl"
  fi
  curl -fsSL -o "$HOME/bin/kubectl" "$KURL"
  chmod +x "$HOME/bin/kubectl"
  KUBECTL_CMD="$HOME/bin/kubectl"
}

# Make sure ~/bin is first on PATH so any downloaded binaries are found by 'kind' and 'kubectl'
export PATH="$HOME/bin:$PATH"

# Ensure kind/kubectl are present (download into ~/bin if needed)
ensure_kind
ensure_kubectl

require_cmd() {
  # Accepts either command name or full path; returns non-zero if missing
  cmd="$1"
  if [ -z "$cmd" ]; then
    return 1
  fi
  if [[ "$cmd" == */* ]]; then
    [ -x "$cmd" ] || return 1
  else
    command -v "$cmd" >/dev/null 2>&1 || return 1
  fi
  return 0
}

usage() {
  cat <<EOF
Usage: $0 [options]
Options:
  --all                 Run the full flow: build -> kind create (if missing) -> load -> deploy -> port-forward (includes mock-api + web-app in cluster)
  --all-in-kind         Same as --all (explicit)
  --build               Build local Docker images (service + ui + mock-api + web-app)
			--kind-create         Create kind cluster if missing
  --load                Load built images into kind
  --deploy              kubectl apply manifests (namespace, deployments, services)
  --export-kubeconfig   Export kubeconfig into k8s-demo/kubeconfig (uses current context or kind)
  --port-forward        Port-forward test-runner-service to localhost:${API_PORT_LOCAL}
  --run                 Trigger a run (uses default payload if not provided; see --payload)
  --payload '<json>'    JSON payload for POST /api/runs (overrides defaults)
  --stream              Open SSE stream for the created run (short timeout)
  --copy-reports        Copy /app/reports from service pod to ./tmp-reports
  --help                Show this help

Examples:
  $0 --all
  $0 --build --load --deploy
  $0 --port-forward --run --stream
EOF
  exit 1
}

build_images() {
  echo "[build] Building test-runner-service image..."
  if [ -f "$SERVICE_DIR/Dockerfile" ]; then
    docker build -t test-runner-service:${IMAGE_TAG} "$SERVICE_DIR"
  else
    echo "[build] No Dockerfile found for service at $SERVICE_DIR/Dockerfile" >&2
    return 1
  fi

  echo "[build] Building test-runner-ui image..."
  if [ -f "$UI_DIR/Dockerfile" ]; then
    docker build --build-arg VITE_API_BASE=http://localhost:9001 -t test-runner-ui:${IMAGE_TAG} "$UI_DIR"
  else
    echo "[build] No Dockerfile found for UI at $UI_DIR/Dockerfile" >&2
    return 1
  fi

  echo "[build] Building mock-api image if present..."
  if [ -f "$MOCK_API_DIR/Dockerfile" ]; then
    docker build -t mock-api:${IMAGE_TAG} "$MOCK_API_DIR"
  else
    echo "[build] No mock-api Dockerfile at $MOCK_API_DIR/Dockerfile; skipping mock-api build"
  fi

  echo "[build] Building web-app image if present..."
  if [ -f "$WEB_APP_DIR/Dockerfile" ]; then
    docker build -t web-app:${IMAGE_TAG} "$WEB_APP_DIR"
  else
    echo "[build] No web-app Dockerfile at $WEB_APP_DIR/Dockerfile; skipping web-app build"
  fi

  echo "[build] Building cucumber-project image..."
  if [ -f "$RUNNER_DIR/Dockerfile.runner" ]; then
    docker build -f "$RUNNER_DIR/Dockerfile.runner" -t cucumber-project:${IMAGE_TAG} "$RUNNER_DIR"
  else
    echo "[build] No Dockerfile.runner found at $RUNNER_DIR/Dockerfile.runner" >&2
    # This is critical for running tests, so we should allow validaiton
  fi
}

# Path for kind-generated kubeconfig; use a temp file to avoid touching a broken ~/.kube/config
KCFG="/tmp/kind-kubeconfig.yaml"

create_kind_if_missing() {
  # ensure kind binary is available
  ensure_kind
  if ! "$KIND_CMD" get clusters | grep -q "^${KIND_CLUSTER}$"; then
    echo "[kind] Creating cluster '${KIND_CLUSTER}'..."

    # If the user's kube config path is a directory (broken), avoid merge error by
    # writing kind's kubeconfig to a temp file and pointing KUBECONFIG at it while creating.
    TEMP_KCFG="/tmp/kind-kubeconfig-${KIND_CLUSTER}.yaml"
    if [ -d "$HOME/.kube/config" ]; then
      echo "[kind] Detected $HOME/.kube/config is a directory; using temporary kubeconfig $TEMP_KCFG for creation"
      export KUBECONFIG="$TEMP_KCFG"
    fi

    # create the cluster (kind will write into $KUBECONFIG if set)
    "$KIND_CMD" create cluster --name "${KIND_CLUSTER}"
  else
    echo "[kind] Cluster '${KIND_CLUSTER}' already exists"
  fi

  # After creation or if it already existed, always write a usable kubeconfig to KCFG and export it
  "$KIND_CMD" get kubeconfig --name "${KIND_CLUSTER}" > "$KCFG" || true
  chmod 600 "$KCFG" || true
  export KUBECONFIG="$KCFG"
  echo "[kind] Using KUBECONFIG=$KCFG"
}

# Ensure we have a usable kubeconfig exported for kubectl to use (call before any kubectl actions)
ensure_kubeconfig() {
  # If KUBECONFIG already set and points to file, keep it
  if [ -n "${KUBECONFIG:-}" ] && [ -f "${KUBECONFIG}" ]; then
    KCFG="$KUBECONFIG"
    return 0
  fi
  # Try to get kubeconfig from kind
  if [ -n "$KIND_CMD" ] && [ -x "$KIND_CMD" ]; then
    if "$KIND_CMD" get clusters | grep -q "^${KIND_CLUSTER}$"; then
      "$KIND_CMD" get kubeconfig --name "${KIND_CLUSTER}" > "$KCFG" || true
      chmod 600 "$KCFG" || true
      export KUBECONFIG="$KCFG"
      echo "[kubeconfig] Exported kind kubeconfig to $KCFG"
      return 0
    fi
  fi
  # Fall back to existing kubectl config if available (don't overwrite)
  if [ -f "$HOME/.kube/config" ]; then
    KCFG="$HOME/.kube/config"
    export KUBECONFIG="$KCFG"
    echo "[kubeconfig] Using existing KUBECONFIG=$HOME/.kube/config"
    return 0
  fi
  echo "[kubeconfig] Warning: no kubeconfig available"
}

load_images_into_kind() {
  ensure_kind
  echo "[kind] Loading images into kind (${KIND_CLUSTER})..."
  "$KIND_CMD" load docker-image "test-runner-service:${IMAGE_TAG}" --name "${KIND_CLUSTER}"
  "$KIND_CMD" load docker-image "test-runner-ui:${IMAGE_TAG}" --name "${KIND_CLUSTER}"
  "$KIND_CMD" load docker-image "cucumber-project:${IMAGE_TAG}" --name "${KIND_CLUSTER}"
  # load optional images
  if docker image inspect "mock-api:${IMAGE_TAG}" >/dev/null 2>&1; then
    "$KIND_CMD" load docker-image "mock-api:${IMAGE_TAG}" --name "${KIND_CLUSTER}"
  fi
  if docker image inspect "web-app:${IMAGE_TAG}" >/dev/null 2>&1; then
    "$KIND_CMD" load docker-image "web-app:${IMAGE_TAG}" --name "${KIND_CLUSTER}"
  fi
}

deploy_manifests() {
  require_cmd "$KUBECTL_CMD"
  ensure_kubeconfig
  echo "[k8s] Applying manifests from $K8S_MANIFEST_DIR"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/00-namespace.yaml"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/rbac.yaml"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/pvc.yaml"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/service-deployment.yaml"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/ui-deployment.yaml"

  # apply optional manifests if they exist
  if [ -f "$K8S_MANIFEST_DIR/mock-api-deployment.yaml" ]; then
    "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/mock-api-deployment.yaml"
  fi
  if [ -f "$K8S_MANIFEST_DIR/web-app-deployment.yaml" ]; then
    "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/web-app-deployment.yaml"
  fi
  if [ -f "$K8S_MANIFEST_DIR/nodeport-services.yaml" ]; then
    "$KUBECTL_CMD" --kubeconfig "$KCFG" apply -f "$K8S_MANIFEST_DIR/nodeport-services.yaml"
  fi

  echo "[k8s] Waiting for test-runner-service pod (timeout 120s)"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" wait --for=condition=ready pod -l app=test-runner-service -n ${NAMESPACE} --timeout=120s || true
  echo "[k8s] Waiting for test-runner-ui pod (timeout 120s)"
  "$KUBECTL_CMD" --kubeconfig "$KCFG" wait --for=condition=ready pod -l app=test-runner-ui -n ${NAMESPACE} --timeout=120s || true
  if [ -f "$K8S_MANIFEST_DIR/mock-api-deployment.yaml" ]; then
    "$KUBECTL_CMD" --kubeconfig "$KCFG" wait --for=condition=ready pod -l app=mock-api -n ${NAMESPACE} --timeout=120s || true
  fi
  if [ -f "$K8S_MANIFEST_DIR/web-app-deployment.yaml" ]; then
    "$KUBECTL_CMD" --kubeconfig "$KCFG" wait --for=condition=ready pod -l app=web-app -n ${NAMESPACE} --timeout=120s || true
  fi
}

export_kubeconfig_local() {
  OUT_FILE="$K8S_MANIFEST_DIR/kubeconfig"
  if [ -n "$KIND_CMD" ] && [ -x "$KIND_CMD" ] && "$KIND_CMD" get clusters | grep -q "^${KIND_CLUSTER}$"; then
    echo "[kubeconfig] Exporting kubeconfig from kind cluster '${KIND_CLUSTER}' to ${OUT_FILE}"
    "$KIND_CMD" get kubeconfig --name "${KIND_CLUSTER}" > "${OUT_FILE}"
  else
    echo "[kubeconfig] Exporting current kubectl config to ${OUT_FILE}"
    "$KUBECTL_CMD" config view --raw > "${OUT_FILE}"
  fi
  chmod 600 "${OUT_FILE}"
  echo "[kubeconfig] Wrote ${OUT_FILE}"
}

start_port_forward() {
  require_cmd "$KUBECTL_CMD"
  ensure_kubeconfig
  # avoid multiple port-forwards
  if [ -f "$PORT_FORWARD_PID_FILE" ]; then
    oldpid=$(cat "$PORT_FORWARD_PID_FILE") || true
    if [ -n "${oldpid:-}" ] && kill -0 "$oldpid" 2>/dev/null; then
      echo "[port-forward] Existing port-forward running with PID $oldpid"
      return
    else
      rm -f "$PORT_FORWARD_PID_FILE" || true
    fi
  fi

  echo "[port-forward] Starting port-forwards (Namespace: ${NAMESPACE})"

  # Service
  "$KUBECTL_CMD" --kubeconfig "$KCFG" port-forward svc/test-runner-service ${API_PORT_LOCAL}:9001 -n ${NAMESPACE} >>/tmp/port-forward.log 2>&1 &
  pids="$!"
  echo "[port-forward] API: http://localhost:${API_PORT_LOCAL} (PID $!)"

  # UI
  "$KUBECTL_CMD" --kubeconfig "$KCFG" port-forward svc/test-runner-ui 9000:9000 -n ${NAMESPACE} >>/tmp/port-forward.log 2>&1 &
  pids="$pids $!"
  echo "[port-forward] UI:  http://localhost:9000 (PID $!)"

  # Mock API
  "$KUBECTL_CMD" --kubeconfig "$KCFG" port-forward svc/mock-api 8001:8001 -n ${NAMESPACE} >>/tmp/port-forward.log 2>&1 &
  pids="$pids $!"
  echo "[port-forward] Mock: http://localhost:8001 (PID $!)"

  # Web App
  "$KUBECTL_CMD" --kubeconfig "$KCFG" port-forward svc/web-app 8000:8000 -n ${NAMESPACE} >>/tmp/port-forward.log 2>&1 &
  pids="$pids $!"
  echo "[port-forward] App:  http://localhost:8000 (PID $!)"

  echo "$pids" > "$PORT_FORWARD_PID_FILE"
  echo "[port-forward] Logs at /tmp/port-forward.log"
  sleep 1
}

stop_port_forward() {
  if [ -f "$PORT_FORWARD_PID_FILE" ]; then
    pids=$(cat "$PORT_FORWARD_PID_FILE") || true
    for pfpid in $pids; do
      if [ -n "${pfpid:-}" ] && kill -0 "$pfpid" 2>/dev/null; then
        echo "[port-forward] Stopping PID $pfpid"
        kill "$pfpid" || true
      fi
    done
    rm -f "$PORT_FORWARD_PID_FILE" || true
  else
    echo "[port-forward] No PID file found"
  fi
}

trigger_run() {
  require_cmd curl
  payload="${PAYLOAD:-}"
  if [ -z "$payload" ]; then
    # default payload
    payload='{"framework":"cucumber","suiteId":"final-cucumber-project","scenarioId":"inventory-list","executor":"k8s"}'
  fi
  # Diagnostics go to stderr so the function prints only the run id to stdout
  echo "[run] POST /api/runs -> $payload" >&2
  resp=$(curl -sS -X POST "http://localhost:${API_PORT_LOCAL}/api/runs" -H "Content-Type: application/json" -d "$payload")
  echo "[run] response: $resp" >&2
  # extract id
  run_id=$(echo "$resp" | python3 -c 'import sys,json
try:
    d=json.loads(sys.stdin.read())
    print(d.get("id",""))
except Exception:
    sys.stdout.write("")') || true
  if [ -z "$run_id" ]; then
    echo "[run] failed to get run id from response" >&2
    return 1
  fi
  # Only emit the id on stdout
  echo "$run_id"
}

stream_run() {
  rid="$1"
  require_cmd curl
  echo "[stream] Streaming SSE for run $rid (will timeout after 120s)"
  # Use timeout so it doesn't hang forever
  timeout 120s curl -sS "http://localhost:${API_PORT_LOCAL}/api/runs/${rid}/stream" || true
}

copy_reports() {
  require_cmd "$KUBECTL_CMD"
  ensure_kubeconfig
  POD=$("$KUBECTL_CMD" --kubeconfig "$KCFG" get pods -n ${NAMESPACE} -l app=test-runner-service -o jsonpath='{.items[0].metadata.name}') || true
  if [ -z "$POD" ]; then
    echo "[reports] Could not find test-runner-service pod in namespace ${NAMESPACE}" >&2
    return 1
  fi
  echo "[reports] Copying /app/reports from pod $POD to ./tmp-reports"
  rm -rf ./tmp-reports || true
  "$KUBECTL_CMD" --kubeconfig "$KCFG" cp ${NAMESPACE}/${POD}:/app/reports ./tmp-reports -n ${NAMESPACE} || true
  echo "[reports] Local ./tmp-reports content:"
  ls -la ./tmp-reports || true
}

# Parse args
if [ $# -eq 0 ]; then
  usage
fi

DO_BUILD=0
DO_KIND_CREATE=0
DO_LOAD=0
DO_DEPLOY=0
DO_EXPORT_KUBECONFIG=0
DO_PORT_FORWARD=0
DO_RUN=0
DO_STREAM=0
DO_COPY_REPORTS=0

PAYLOAD=""

while [ $# -gt 0 ]; do
  case "$1" in
    --all)
      DO_BUILD=1; DO_KIND_CREATE=1; DO_LOAD=1; DO_DEPLOY=1; DO_PORT_FORWARD=1
      shift
      ;;
    --all-in-kind)
      DO_BUILD=1; DO_KIND_CREATE=1; DO_LOAD=1; DO_DEPLOY=1; DO_PORT_FORWARD=1
      shift
      ;;
    --build)
      DO_BUILD=1; shift ;;
    --kind-create)
      DO_KIND_CREATE=1; shift ;;
    --load)
      DO_LOAD=1; shift ;;
    --deploy)
      DO_DEPLOY=1; shift ;;
    --export-kubeconfig)
      DO_EXPORT_KUBECONFIG=1; shift ;;
    --port-forward)
      DO_PORT_FORWARD=1; shift ;;
    --run)
      DO_RUN=1; shift ;;
    --payload)
      PAYLOAD="$2"; shift 2 ;;
    --stream)
      DO_STREAM=1; shift ;;
    --copy-reports)
      DO_COPY_REPORTS=1; shift ;;
    --help)
      usage ;;
    *)
      echo "Unknown arg: $1" >&2; usage ;;
  esac
done

# Run requested steps
if [ $DO_BUILD -eq 1 ]; then
  require_cmd docker
  build_images
fi

if [ $DO_KIND_CREATE -eq 1 ]; then
  create_kind_if_missing
fi

# Ensure kubeconfig is available before any kubectl operations
ensure_kubeconfig

if [ $DO_LOAD -eq 1 ]; then
  load_images_into_kind
fi

if [ $DO_DEPLOY -eq 1 ]; then
  deploy_manifests
fi

if [ $DO_EXPORT_KUBECONFIG -eq 1 ]; then
  export_kubeconfig_local
fi

if [ $DO_PORT_FORWARD -eq 1 ]; then
  start_port_forward
fi

run_id=""
if [ $DO_RUN -eq 1 ]; then
  run_id=$(trigger_run) || true
fi

if [ $DO_STREAM -eq 1 ] && [ -n "$run_id" ]; then
  stream_run "$run_id"
fi

if [ $DO_COPY_REPORTS -eq 1 ]; then
  copy_reports
fi

echo "[done]"
