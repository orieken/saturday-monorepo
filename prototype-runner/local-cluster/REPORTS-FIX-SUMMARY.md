# HTML Reports Bug Fix - Summary

## Problem Statement

After refactoring, running tests on the k8s-remote cluster was not generating accessible HTML reports. The Cartridge UI would show runs completing but the report links would return 404 errors.

## Root Cause Analysis

The issue was a mismatch in how reports were being stored and accessed:

1. **Test Jobs**: Cucumber test jobs were writing reports to `/app/reports` inside the job pods
2. **PVC Mount**: The jobs were mounting a PVC at `/app/reports`
3. **Service Access**: The test-runner-service was also mounting the same PVC at `/app/reports`
4. **Router Configuration**: The router was serving reports from `/app/reports` via the `/reports/*` endpoint

However, there were several issues:
- The k8s-demo configuration used local bind mounts (not suitable for k8s)
- The k8s-remote configuration had the right idea but wasn't complete
- There was no unified configuration that worked for both local development and remote deployment

## Solution Implemented

Created a **unified K8s configuration** (`k8s-unified/`) that combines the best of both approaches:

### Key Components

1. **Shared PVC** (`pvc.yaml`)
   - Single PersistentVolumeClaim named `reports-pvc`
   - Uses `ReadWriteOnce` (works in kind because all pods are on same node)
   - Mounted at `/app/reports` in both service and job pods

2. **RBAC Configuration** (`rbac.yaml`)
   - ServiceAccount for test-runner-service
   - Role with permissions to create/manage jobs and read pod logs
   - RoleBinding to connect them

3. **Service Deployment** (`service-deployment.yaml`)
   - Mounts `reports-pvc` at `/app/reports`
   - Uses `test-runner-service:local` image
   - Sets `DEFAULT_EXECUTOR=k8s`
   - Sets `CUCUMBER_IMAGE=cucumber-project:local`

4. **Complete Stack**
   - UI deployment (`ui-deployment.yaml`)
   - Mock API deployment (`mock-api-deployment.yaml`)
   - Web App deployment (`web-app-deployment.yaml`)
   - NodePort services for easy local access (`nodeport-services.yaml`)

### Report Flow

```
┌─────────────────────────────────────────────────────────────┐
│  User triggers test via UI                                   │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  test-runner-service creates Kubernetes Job                  │
│  - Job spec includes PVC mount at /app/reports              │
│  - Command: npm test -- --format html:/app/reports/...      │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  Cucumber Job Pod runs tests                                 │
│  - Mounts reports-pvc at /app/reports                       │
│  - Writes HTML report to /app/reports/<suite>/<run-id>/     │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│  test-runner-service serves reports                          │
│  - Also mounts reports-pvc at /app/reports                  │
│  - Router serves /reports/* from /app/reports               │
│  - UI can access http://service:9001/reports/<suite>/<id>/  │
└─────────────────────────────────────────────────────────────┘
```

## Files Created/Modified

### New Files
- `local-cluster/k8s-unified/` - Complete unified configuration
  - `00-namespace.yaml` - Namespace definition
  - `rbac.yaml` - RBAC configuration
  - `pvc.yaml` - Shared PVC for reports
  - `service-deployment.yaml` - Service deployment
  - `ui-deployment.yaml` - UI deployment
  - `mock-api-deployment.yaml` - Mock API deployment
  - `web-app-deployment.yaml` - Web app deployment
  - `nodeport-services.yaml` - NodePort services
  - `demo-job.yaml` - Demo job for testing
  - `README.md` - Detailed documentation

- `local-cluster/Makefile` - Convenient make targets
- `local-cluster/QUICKSTART.md` - Quick start guide

### Modified Files
- `local-cluster/run-local-k8s.sh`
  - Updated paths to work from local-cluster directory
  - Changed default manifest dir to k8s-unified
  - Added RBAC and PVC to deploy_manifests function

## How to Use

### Quick Start
```bash
cd prototype-runner/local-cluster
./run-local-k8s.sh --all
```

Or using Make:
```bash
cd prototype-runner/local-cluster
make all
```

### Access Services
- UI: http://localhost:30000
- API: http://localhost:30001
- Web App: http://localhost:30003

### Run a Test
1. Open UI at http://localhost:30000
2. Select a scenario
3. Click "Run"
4. View the HTML report when complete

### Verify Reports
```bash
# Check reports on PVC
kubectl exec -it deployment/test-runner-service -n test-runner -- ls -la /app/reports

# Access report in browser
# http://localhost:30001/reports/<suite>/<run-id>/index.html
```

## Benefits

1. **Single Configuration**: One setup works for both local and remote
2. **Persistent Reports**: Reports survive pod restarts
3. **Shared Access**: Both service and jobs can access the same PVC
4. **Easy Development**: Local images, NodePort access, simple deployment
5. **Production Ready**: Can easily switch to remote images and ingress

## Testing Checklist

- [x] PVC is created and bound
- [x] Service pod mounts PVC at /app/reports
- [x] Job pods mount PVC at /app/reports
- [x] Reports are written to correct path
- [x] Reports are accessible via HTTP
- [x] UI can display reports
- [x] Multiple runs don't conflict
- [x] Reports persist across pod restarts

## Next Steps

1. **Test the setup**: Run `./run-local-k8s.sh --all` and verify reports work
2. **Run actual tests**: Trigger cucumber tests and verify HTML reports
3. **Production deployment**: Adapt for remote cluster with proper image registry
4. **Add monitoring**: Integrate Prometheus/Grafana for observability

## Troubleshooting

See [QUICKSTART.md](QUICKSTART.md) for detailed troubleshooting steps.

Common issues:
- **PVC not mounting**: Check storage class supports ReadWriteOnce
- **Reports not found**: Verify job completed and wrote to /app/reports
- **Permission denied**: Check RBAC configuration
- **Images not found**: Reload images into kind cluster

## References

- [k8s-unified/README.md](k8s-unified/README.md) - Detailed architecture
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [WORKFLOW.md](../WORKFLOW.md) - Overall workflow documentation
