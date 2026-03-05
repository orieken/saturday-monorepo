# Quick Start Guide - Unified K8s Test Runner

This guide will help you get the unified K8s test runner up and running quickly.

## Problem We're Solving

After refactoring, the k8s-remote cluster wasn't generating HTML reports. The issue was that:
- Test jobs were writing reports to a PVC
- But the configuration wasn't properly set up for the service to access those reports
- We needed a unified configuration that works for both local and remote scenarios

## Solution: Unified Configuration

The `k8s-unified` directory contains a single configuration that:
1. ✅ Mounts a shared PVC at `/app/reports` in both the service and test job pods
2. ✅ Uses local images for easy development
3. ✅ Includes all necessary components (service, UI, mock-api, web-app)
4. ✅ Provides NodePort services for easy local access

## Quick Start (5 minutes)

### Option 1: Using the Helper Script

```bash
cd prototype-runner/local-cluster
./run-local-k8s.sh --all
```

This will:
- Create a kind cluster (if needed)
- Build all Docker images
- Load images into kind
- Deploy all manifests
- Set up port forwarding

### Option 2: Using Make

```bash
cd prototype-runner/local-cluster
make all
```

This will build, load, and deploy everything.

### Option 3: Manual Steps

```bash
# 1. Build images
cd apps/console && docker build -t test-runner-service:local .
cd ../cartridge && docker build -t test-runner-ui:local .
cd ../../mock-api && docker build -t mock-api:local .
cd ../../web-app && docker build -t web-app:local .
cd ../../cucumber-project && docker build -t cucumber-project:local .

# 2. Create kind cluster
kind create cluster --name kind

# 3. Load images
kind load docker-image test-runner-service:local --name kind
kind load docker-image test-runner-ui:local --name kind
kind load docker-image mock-api:local --name kind
kind load docker-image web-app:local --name kind
kind load docker-image cucumber-project:local --name kind

# 4. Deploy
cd prototype-runner/local-cluster
kubectl apply -f k8s-unified/

# 5. Wait for pods
kubectl wait --for=condition=ready pod -l app=test-runner-service -n test-runner --timeout=120s
```

## Accessing the Services

### Via NodePort (Easiest)
- UI: http://localhost:30000
- API: http://localhost:30001
- Web App: http://localhost:30003

### Via Port Forwarding
```bash
kubectl port-forward svc/test-runner-service 9001:9001 -n test-runner &
kubectl port-forward svc/test-runner-ui 9000:9000 -n test-runner &
```

Then access:
- UI: http://localhost:9000
- API: http://localhost:9001

## Running a Test

### Via UI
1. Open http://localhost:30000 (or http://localhost:9000)
2. Browse scenarios
3. Click "Run" on any scenario
4. Watch the logs in real-time
5. View the HTML report when complete

### Via API
```bash
# Trigger a run
RUN_ID=$(curl -X POST http://localhost:30001/api/runs \
  -H "Content-Type: application/json" \
  -d '{"framework":"cucumber","suiteId":"final-cucumber-project","scenarioId":"inventory-list","executor":"k8s"}' \
  | jq -r '.id')

# Check status
curl http://localhost:30001/api/runs/$RUN_ID | jq .

# View report (once status is "passed" or "failed")
# Open in browser: http://localhost:30001/reports/<suite>/<run-id>/index.html
```

## Verifying Reports Are Working

### 1. Check the PVC is mounted
```bash
kubectl describe pod -l app=test-runner-service -n test-runner | grep -A 5 "Mounts:"
```

You should see `/app/reports` mounted.

### 2. Run a test and check the filesystem
```bash
# After running a test
kubectl exec -it deployment/test-runner-service -n test-runner -- ls -la /app/reports
```

You should see directories like `final-cucumber-project/<run-id>/`.

### 3. Access a report
```bash
# List reports
kubectl exec -it deployment/test-runner-service -n test-runner -- find /app/reports -name "index.html"

# View a report in your browser
# http://localhost:30001/reports/<suite>/<run-id>/index.html
```

## Troubleshooting

### Reports not showing up?

1. **Check if the job completed successfully:**
   ```bash
   kubectl get jobs -n test-runner
   kubectl logs job/<job-name> -n test-runner
   ```

2. **Check if reports were written:**
   ```bash
   kubectl exec -it deployment/test-runner-service -n test-runner -- ls -la /app/reports
   ```

3. **Check service logs:**
   ```bash
   kubectl logs -f deployment/test-runner-service -n test-runner
   ```

### PVC not mounting?

```bash
# Check PVC status
kubectl get pvc -n test-runner

# Check pod events
kubectl describe pod -l app=test-runner-service -n test-runner
```

### Images not found?

```bash
# List images in kind
docker exec -it kind-control-plane crictl images | grep local

# Reload images if needed
make load-all
```

## Key Differences from Previous Setups

### vs k8s-demo
- Now uses a PVC instead of emptyDir for reports
- Reports persist across pod restarts
- Service can access reports from jobs

### vs k8s-remote
- Uses local images instead of remote registry
- Includes all supporting services
- Optimized for local development
- No need to push to Docker Hub

## Next Steps

1. **Run your first test** via the UI
2. **Verify the HTML report** is accessible
3. **Check the logs** to see the test execution
4. **Explore the API** endpoints

## Clean Up

```bash
# Delete everything
kubectl delete namespace test-runner

# Or use make
make clean

# Delete the kind cluster
kind delete cluster --name kind
```

## Getting Help

- Check the [README](k8s-unified/README.md) for detailed documentation
- Review the [WORKFLOW](../WORKFLOW.md) for architecture details
- Look at the Makefile for available commands: `make help`
