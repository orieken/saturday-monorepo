# Local Kubernetes cluster (kind) — where the configs live and how the demo uses them

This document centralizes where our local Kubernetes demo manifests, helper scripts, and workflows live so any agent or developer can find and use them quickly.

Summary
- All local k8s configs, manifests, and helper scripts live in the `local-cluster` directories in the repository. There are two common places you’ll see them depending on which demo surface you use:
  - `final/local-cluster/` — the manifests and scripts used by the final demo setup (used by `final/scripts/kind-prepare.sh` and the `final/test-runner-service` deployment)
  - `local-cluster/` — higher-level, repository-top local cluster helpers and additional manifests (alternate entrypoint for other demos)
- The demo namespace used by the `test-runner-service` deployment is `test-runner`.

Key files and what they do

- `final/scripts/kind-prepare.sh`
  - Builds local images (`test-runner-service:local` and `cucumber-project:local`), loads them into the named kind cluster, applies namespace + RBAC, applies the deployment manifest and restarts the deployment. This is the main helper script used to refresh images and apply the demo manifests.

- `final/local-cluster/k8s-demo/service-deployment.yaml`
  - The Deployment + Service for `test-runner-service` (namespace: `test-runner`). The Deployment uses image `test-runner-service:local` and mounts a `reports` emptyDir by default. It sets `serviceAccountName: test-runner-sa` in our repo manifests.

- `final/local-cluster/test-runner-rbac.yaml`
  - Role and RoleBinding (and ServiceAccount) needed so the test-runner can create Jobs, read pods/logs, and exec into pods. This file was added to unblock local demos; it should be hardened later.

- `final/local-cluster/k8s-demo/namespace.yaml`
  - Namespace manifest for `test-runner` (applied by the helper script if present).

- `final/test-runner-service/Dockerfile` and `final/test-runner-service/Makefile`
  - How the test-runner service is built and packaged. The Dockerfile builds a static Go binary then packages it on Alpine.

- `final/cucumber-project/Dockerfile` (if present)
  - The image used as the Job template. Jobs created by the service refer to `cucumber-project:local` in the demo — that image must be loaded into kind for Jobs to schedule.

- `final/local-cluster/k8s-demo/start-port-forwards.sh` and `port-forward.sh`
  - Scripts to port-forward services (useful for local UI and API access).

- `local-cluster/run-local-k8s.sh`, `local-cluster/README.md`
  - Project-level convenience scripts and documentation for starting a local k8s demo using the other manifests.

Common workflow (copy/paste commands)

1) Preconditions:
- Docker installed and running
- kind installed and on PATH
- kubectl installed and configured
- If you have a non-standard kubeconfig, set KUBECONFIG accordingly:

```bash
export KUBECONFIG=/path/to/your/kubeconfig
```

2) Build the images (service + cucumber project) and tag them for local kind usage

```bash
# from repo root
cd /Users/oscarrieken/Projects/Rieken/tools/demo
# Build test-runner-service image
cd final/test-runner-service
DOCKER_BUILDKIT=1 docker build -t test-runner-service:local .

# Build cucumber project image (if the Dockerfile exists)
cd ../cucumber-project || true
DOCKER_BUILDKIT=1 docker build -t cucumber-project:local . || true
```

3) Load images into kind (replace cluster name if not `kind`)

```bash
kind load docker-image test-runner-service:local --name kind
kind load docker-image cucumber-project:local --name kind
```

4) Apply namespace, RBAC and deployment manifests (or use the helper script)

```bash
# Apply namespace and RBAC
kubectl apply -f final/local-cluster/k8s-demo/namespace.yaml
kubectl apply -f final/local-cluster/test-runner-rbac.yaml

# Deploy the service (will create Deployment + Service)
kubectl apply -f final/local-cluster/k8s-demo/service-deployment.yaml

# Restart to pick up local images (imagePullPolicy: IfNotPresent helps)
kubectl -n test-runner rollout restart deployment/test-runner-service
kubectl -n test-runner rollout status deployment/test-runner-service --timeout=120s
```

Or run the helper which automates steps 2–4 (recommended):

```bash
# from repo root
./final/scripts/kind-prepare.sh kind
```

Verification and quick checks

- Check that the deployment exists and the image matches the local tag:

```bash
kubectl -n test-runner get deployment test-runner-service -o yaml | grep image -A1
kubectl -n test-runner get pods -l app=test-runner-service -o wide
```

- Inspect logs:

```bash
kubectl -n test-runner logs -l app=test-runner-service --tail=200
```

- Port-forward the service for local API/UI access:

```bash
kubectl -n test-runner port-forward svc/test-runner-service 9001:9001 &
curl -sS http://localhost:9001/ || curl -sS http://localhost:9001/health || true
```

Notes, gotchas and recommendations

- kubeconfig: ensure `~/.kube/config` is a file (not a directory). If kubectl errors about the config being a directory, back it up and create a proper file or set the `KUBECONFIG` env variable to a file path.

- RBAC: The repo currently includes `final/local-cluster/test-runner-rbac.yaml` which grants permissions used in development. For production or shared demos, replace the `default` SA usage with a dedicated `test-runner-sa` and tighten the Role to minimum required verbs.

- Image availability: Jobs referencing `cucumber-project:local` will fail to create containers in kind if the image hasn't been loaded into the cluster. Use `kind load docker-image` to make locally-built images available.

- Persistent reports: By default the Deployment uses an `emptyDir` for `/app/reports` which does not persist across restarts. For demo persistence mount a PVC into `/app/reports` and `/app/data` in the Deployment manifest.

- Helper scripts: `final/scripts/kind-prepare.sh` is the canonical helper for the final demo. It builds the images, loads them into kind, applies namespace & RBAC, and restarts the deployment. Use it as the first step when updating the service image.

Where agents should look first
- `final/scripts/kind-prepare.sh` — automation entrypoint for reloads
- `final/local-cluster/k8s-demo/service-deployment.yaml` — service deployment/service manifest
- `final/local-cluster/test-runner-rbac.yaml` — RBAC used in demos
- `final/test-runner-service/Dockerfile` and `Makefile` — how the service is built and packaged
- `final/cucumber-project` — job image used by the service

If you'd like, I can also:
- Add a small `docs/local-cluster-quickstart.md` that contains a single-line script to fully rebuild/load/apply everything and run a smoke test.
- Modify `final/scripts/kind-prepare.sh` to accept a `--kubeconfig` or export `KUBECONFIG` to avoid kubeconfig directory issues.

---
Document created at `docs/local-cluster.md`. Update or request additions and I will extend it.
