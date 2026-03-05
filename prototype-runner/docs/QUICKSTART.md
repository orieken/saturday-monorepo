Quickstart — Local cluster and demo services

This repository includes a small local Kubernetes demo under `final/local-cluster`.
All Kubernetes manifests and helper scripts for the demo cluster live under:

  final/local-cluster/

Key files:

- `final/local-cluster/kubeconfig` — a kubeconfig pre-configured for the local kind cluster used by the demo.
- `final/local-cluster/k8s-demo/` — K8s manifests (namespace, deployments, services, demo job).
- `final/local-cluster/deploy.sh` and other helper scripts used by the k8s demo.

How the quickstart / rebuild flow works

1) Build the Docker image(s) locally from the `final/` subfolders (for example `final/mock-api`).
2) Load built images into the kind cluster (`kind load docker-image <image> --name <cluster>`).
3) Apply manifests from `final/local-cluster/k8s-demo` (or use `kubectl` with the `kubeconfig` stored in that directory).
4) Restart deployments to pick up newly loaded images (optional): `kubectl -n test-runner rollout restart deployment/<name>`.

Quick commands (use from repo root):

  # build mock-api and load into the cluster named in repo kubeconfig
  ./scripts/rebuild-mock-api.sh

  # run the repo quickstart (builds service + cucumber images and applies manifests)
  ./scripts/quickstart.sh [KIND_CLUSTER_NAME]

Notes:
- The scripts expect `docker`, `kind` and `kubectl` to be installed and Docker to be running.
- The scripts default to using `final/local-cluster/kubeconfig` as the kubeconfig file. You can override the kubeconfig by setting the KUBECONFIG environment variable.
- If you use a different cluster name, pass it as the first argument to the scripts.

If you'd like, I can also add CI tasks or Makefile targets to automate these steps.

