# Unified K8s Configuration

This directory contains a unified Kubernetes configuration that combines the best of both `k8s-demo` (local) and `k8s-remote` (remote) setups.

## Key Features

1. **Shared PVC for Reports**: Both the test-runner-service and test job pods mount the same PVC at `/app/reports`, ensuring reports are accessible.
2. **Local Images**: Uses `test-runner-service:local`, `cucumber-project:local`, etc. for local development.
3. **Complete Stack**: Includes all necessary components (service, UI, mock-api, web-app).
4. **NodePort Access**: Provides NodePort services for easy local access.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                       │
│                                                               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │ test-runner  │    │ test-runner  │    │   web-app    │  │
│  │   service    │◄───┤      UI      │    │   :8000      │  │
│  │   :9001      │    │   :9000      │    └──────────────┘  │
│  └──────┬───────┘    └──────────────┘           ▲          │
│         │                                        │          │
│         │ creates                                │          │
│         ▼                                        │          │
│  ┌──────────────┐                                │          │
│  │ Cucumber Job │────────────────────────────────┘          │
│  │   (Pod)      │      runs tests against                   │
│  └──────┬───────┘                                           │
│         │                                                    │
│         │ writes                                             │
│         ▼                                                    │
│  ┌──────────────┐                                           │
│  │ reports-pvc  │◄──────────────────────────────────────────┤
│  │ (RWX)        │         mounted by service                │
│  └──────────────┘         to serve reports                  │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Quick Start

### Prerequisites

- Docker installed
- kind or another Kubernetes cluster
- kubectl configured

### Build Images

```bash
# From the monorepo root
cd apps/console
docker build -t test-runner-service:local .

cd ../cartridge
docker build -t test-runner-ui:local .

cd ../../mock-api
docker build -t mock-api:local .

cd ../web-app
docker build -t web-app:local .

cd ../cucumber-project
docker build -t cucumber-project:local .
```

### Deploy to Cluster

Using the helper script:

```bash
cd prototype-runner/local-cluster
./run-local-k8s.sh --all
```

Or manually:

```bash
# Create kind cluster (if needed)
kind create cluster --name kind

# Load images into kind
kind load docker-image test-runner-service:local --name kind
kind load docker-image test-runner-ui:local --name kind
kind load docker-image mock-api:local --name kind
kind load docker-image web-app:local --name kind
kind load docker-image cucumber-project:local --name kind

# Apply manifests
kubectl apply -f k8s-unified/00-namespace.yaml
kubectl apply -f k8s-unified/rbac.yaml
kubectl apply -f k8s-unified/pvc.yaml
kubectl apply -f k8s-unified/service-deployment.yaml
kubectl apply -f k8s-unified/ui-deployment.yaml
kubectl apply -f k8s-unified/mock-api-deployment.yaml
kubectl apply -f k8s-unified/web-app-deployment.yaml
kubectl apply -f k8s-unified/nodeport-services.yaml

# Wait for pods to be ready
kubectl wait --for=condition=ready pod -l app=test-runner-service -n test-runner --timeout=120s
```

### Access the Services

With NodePort services:
- UI: http://localhost:30000
- Service API: http://localhost:30001
- Web App: http://localhost:30003

Or use port-forwarding:
```bash
kubectl port-forward svc/test-runner-service 9001:9001 -n test-runner
kubectl port-forward svc/test-runner-ui 9000:9000 -n test-runner
```

## Running Tests

### Via the UI

1. Open http://localhost:30000 (or http://localhost:9000 if port-forwarding)
2. Select a scenario
3. Click "Run"
4. View the report when complete

### Via API

```bash
# Trigger a test run
curl -X POST http://localhost:30001/api/runs \
  -H "Content-Type: application/json" \
  -d '{
    "framework": "cucumber",
    "suiteId": "final-cucumber-project",
    "scenarioId": "inventory-list",
    "executor": "k8s"
  }'

# Check run status
curl http://localhost:30001/api/runs/<run-id>

# View report
# The reportUrl will be something like /reports/<suite>/<run-id>/index.html
# Access it at: http://localhost:30001/reports/<suite>/<run-id>/index.html
```

## Troubleshooting

### Reports Not Showing

1. **Check PVC is mounted**:
   ```bash
   kubectl describe pod -l app=test-runner-service -n test-runner
   # Look for volume mounts at /app/reports
   ```

2. **Check reports exist on PVC**:
   ```bash
   kubectl exec -it deployment/test-runner-service -n test-runner -- ls -la /app/reports
   ```

3. **Check job logs**:
   ```bash
   kubectl logs -l run-id=<run-id> -n test-runner
   ```

### PVC Access Mode Issues

If you see "multi-attach error" or similar, ensure your storage class supports `ReadWriteMany`. For kind/local development, you may need to:

1. Use `ReadWriteOnce` and ensure only one pod writes at a time
2. Or use a storage class that supports RWX (like NFS)

For local kind clusters, the default `standard` storage class typically only supports `ReadWriteOnce`. The current configuration uses `ReadWriteMany` which works in kind because all pods are on the same node.

### Service Can't Create Jobs

Check RBAC permissions:
```bash
kubectl get rolebinding -n test-runner
kubectl describe role test-runner-role -n test-runner
```

## Configuration Files

- `00-namespace.yaml` - Creates the test-runner namespace
- `rbac.yaml` - Service account and permissions for creating jobs
- `pvc.yaml` - Persistent volume claim for reports (ReadWriteMany)
- `service-deployment.yaml` - Test runner service deployment
- `ui-deployment.yaml` - UI deployment
- `mock-api-deployment.yaml` - Mock API deployment
- `web-app-deployment.yaml` - Demo web app deployment
- `nodeport-services.yaml` - NodePort services for local access
- `demo-job.yaml` - Example job for testing

## Differences from k8s-demo and k8s-remote

### vs k8s-demo
- ✅ Uses PVC instead of emptyDir for reports
- ✅ Includes RBAC configuration
- ✅ Uses NodePort for easier local access
- ✅ Service mounts the same PVC as jobs

### vs k8s-remote
- ✅ Uses local images instead of remote registry
- ✅ Includes all supporting services (mock-api, web-app)
- ✅ Optimized for local development
- ✅ No need to push images to registry

## Next Steps

1. **Production Setup**: For production, change `imagePullPolicy` to `Always` and use a proper image registry
2. **Storage**: Use a proper storage class with RWX support (NFS, Ceph, etc.)
3. **Ingress**: Add ingress rules instead of NodePort for production
4. **Monitoring**: Add Prometheus/Grafana for observability
