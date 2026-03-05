K8s demo for `final/local-cluster`

This folder contains a small Kubernetes demo to run the test-runner stack on a local cluster (kind / minikube) and an example Job that simulates a test run and writes a simple HTML report.

Files
- `namespace.yaml` — namespace for the demo (`test-runner`).
- `service-deployment.yaml` — Deployment + Service for `test-runner-service` (ClusterIP).
- `ui-deployment.yaml` — Deployment + Service for `test-runner-ui` (ClusterIP).
- `demo-job.yaml` — Job manifest that simulates a run and writes `/reports/index.html` inside the pod.
- `deploy.sh` — helper to apply namespace + deployments. (Remember to load local images into the cluster first if you're using built images.)
- `port-forward.sh` — helper script to forward the UI and API to localhost for easy testing.

Quick start

1) Build local images (optional, if you want to use code-built images):

```bash
# from repo root
cd final/test-runner-ui
docker build -t test-runner-ui:local .
cd ../test-runner-service
docker build -t test-runner-service:local .
# if using kind, load images into cluster
kind load docker-image test-runner-ui:local
kind load docker-image test-runner-service:local
```

2) Deploy stack

```bash
cd final/local-cluster/k8s-demo
./deploy.sh
```

3) Expose UI and API locally (port-forward)

```bash
./port-forward.sh
# UI -> http://localhost:9000
# Service -> http://localhost:9001
```

4) Trigger the demo job (manually)

```bash
# in this folder
kubectl apply -f demo-job.yaml -n test-runner
# find job pod
kubectl get pods -n test-runner --selector=job-name=test-runner-demo-job
# stream logs
kubectl logs -f <POD> -n test-runner
# copy report
kubectl cp -n test-runner <POD>:/reports/index.html ./demo-report-demo-1.html
```

Trigger via UI

- If you want the UI to trigger jobs directly, the backend must support creating Kubernetes Jobs (i.e., an executor that posts Job manifests to the cluster when `POST /api/runs` is called with `executor: "k8s"`).
- The current demo provides the k8s deploy + Job manifest as an example; wiring the UI to create jobs requires server-side changes described in `test-runner-ui/TODO-run-execution.md`.

Notes
- The demo uses `emptyDir` to hold the report inside the pod; in production you'd write artifacts to object storage or a PVC and expose a `reportUrl` from the backend.
- `port-forward.sh` runs foreground `kubectl port-forward` commands — keep the terminal open while testing, or run in separate terminals.

Troubleshooting
- If port-forward fails, ensure kubectl context points to the cluster and the pods are ready.
- If images are not found, load them into the cluster (kind) or push them to a registry accessible by the cluster.

