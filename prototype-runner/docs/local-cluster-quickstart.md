# Local cluster quickstart

A minimal quickstart to rebuild the demo images, load them into a local kind cluster, apply manifests, and run a smoke test.

This file intentionally keeps a single "do-everything" command for experienced users and a short expanded set of commands for clarity.

Prerequisites
- Docker running locally
- Homebrew (optional, for installing tools)
- `kind` on PATH
- `kubectl` on PATH and a working kubeconfig (ensure `~/.kube/config` is a file, not a directory)

Single-line quickstart (from repo root)

```bash
# Rebuild images, load into kind, apply manifests and wait for the service to be ready
cd /Users/oscarrieken/Projects/Rieken/tools/demo && ./final/scripts/kind-prepare.sh kind \
  && kubectl -n test-runner rollout status deployment/test-runner-service --timeout=120s \
  && kubectl -n test-runner get pods -l app=test-runner-service -o wide \
  && kubectl -n test-runner logs -l app=test-runner-service --tail=200
```

Expanded steps (what the one-liner does)

1. Build images and tag as `test-runner-service:local` and `cucumber-project:local`:

```bash
cd /Users/oscarrieken/Projects/Rieken/tools/demo/final/test-runner-service
DOCKER_BUILDKIT=1 docker build -t test-runner-service:local .
cd ../cucumber-project || true
DOCKER_BUILDKIT=1 docker build -t cucumber-project:local . || true
```

2. Load images into kind (replace cluster name if not `kind`):

```bash
kind load docker-image test-runner-service:local --name kind
kind load docker-image cucumber-project:local --name kind
```

3. Apply namespace, RBAC and deployment manifests (helper script does this):

```bash
kubectl apply -f final/local-cluster/k8s-demo/namespace.yaml
kubectl apply -f final/local-cluster/test-runner-rbac.yaml
kubectl apply -f final/local-cluster/k8s-demo/service-deployment.yaml
kubectl -n test-runner rollout restart deployment/test-runner-service
kubectl -n test-runner rollout status deployment/test-runner-service --timeout=120s
```

4. Quick smoke test (port-forward + HTTP check):

```bash
kubectl -n test-runner port-forward svc/test-runner-service 9001:9001 &
# then locally
curl -sS http://localhost:9001/ || curl -sS http://localhost:9001/health || true
```

Troubleshooting (common gotchas)
- kind not found: install with `brew install kind` or visit https://kind.sigs.k8s.io/ for platform-specific instructions.
- kubectl kubeconfig: if `kubectl` errors that `~/.kube/config` is a directory, back it up and create a proper file:

```bash
mv ~/.kube/config ~/.kube/config.bak-$(date +%Y%m%d%H%M%S)
touch ~/.kube/config && chmod 600 ~/.kube/config
```

- Images present locally but not in kind: run `kind load docker-image <image>:local --name <cluster>`.

- If the deployment fails to start, inspect events and pod logs:

```bash
kubectl -n test-runner describe deployment test-runner-service
kubectl -n test-runner get events --sort-by='.lastTimestamp' | tail -n 50
kubectl -n test-runner logs -l app=test-runner-service --tail=200
```

If you'd like, I can also:
- Add a tiny `scripts/quickstart.sh` that wraps the one-liner and performs safety checks (kind/kubectl present, kubeconfig sanity), or
- Run the quickstart for you now (I'll first validate `kind` and `kubectl` are installed and back up your kubeconfig if needed) — say "please run quickstart" to proceed.

