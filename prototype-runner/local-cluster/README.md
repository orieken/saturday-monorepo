# Local Kubernetes Cluster Setup

This directory contains configurations and scripts for running the Test Runner system in a local Kubernetes cluster.

## 🎯 Quick Start

The fastest way to get everything running:

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

Then access:
- **UI**: http://localhost:30000
- **API**: http://localhost:30001
- **Web App**: http://localhost:30003

## 📚 Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get up and running in 5 minutes
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Visual diagrams and architecture overview
- **[REPORTS-FIX-SUMMARY.md](REPORTS-FIX-SUMMARY.md)** - Details on the HTML reports bug fix
- **[k8s-unified/README.md](k8s-unified/README.md)** - Detailed K8s configuration docs

## 🗂️ Directory Structure

```
local-cluster/
├── k8s-unified/          # ⭐ Unified K8s configuration (RECOMMENDED)
│   ├── 00-namespace.yaml
│   ├── rbac.yaml
│   ├── pvc.yaml
│   ├── service-deployment.yaml
│   ├── ui-deployment.yaml
│   ├── mock-api-deployment.yaml
│   ├── web-app-deployment.yaml
│   ├── nodeport-services.yaml
│   ├── demo-job.yaml
│   └── README.md
├── k8s-demo/             # Legacy: Local development config
├── k8s-remote/           # Legacy: Remote cluster config
├── run-local-k8s.sh      # Helper script for deployment
├── Makefile              # Convenient make targets
└── README.md             # This file
```

## 🚀 Usage Options

### Option 1: Helper Script (Recommended)

```bash
# Full deployment
./run-local-k8s.sh --all

# Individual steps
./run-local-k8s.sh --build           # Build images
./run-local-k8s.sh --load            # Load into kind
./run-local-k8s.sh --deploy          # Deploy manifests
./run-local-k8s.sh --port-forward    # Port forward services
./run-local-k8s.sh --run             # Trigger a test run
./run-local-k8s.sh --copy-reports    # Copy reports locally
```

### Option 2: Make Targets

```bash
make help          # Show all available targets
make all           # Build, load, and deploy everything
make build-all     # Build all Docker images
make deploy        # Deploy to cluster
make test-run      # Trigger a test run
make logs          # View service logs
make status        # Check cluster status
make clean         # Delete everything
```

### Option 3: Manual kubectl

```bash
# Build images
cd ../../apps/console && docker build -t test-runner-service:local .
cd ../../apps/cartridge && docker build -t test-runner-ui:local .
# ... etc

# Load into kind
kind load docker-image test-runner-service:local --name kind
# ... etc

# Deploy
kubectl apply -f k8s-unified/
```

## 🏗️ Architecture

The unified configuration uses a shared PVC for reports:

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
│  │ (RWO)        │         mounted by service                │
│  └──────────────┘         to serve reports                  │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed diagrams.

## 🧪 Running Tests

### Via UI
1. Open http://localhost:30000
2. Browse scenarios
3. Click "Run"
4. View HTML report when complete

### Via API
```bash
# Trigger a run
curl -X POST http://localhost:30001/api/runs \
  -H "Content-Type: application/json" \
  -d '{
    "framework": "cucumber",
    "suiteId": "final-cucumber-project",
    "scenarioId": "inventory-list",
    "executor": "k8s"
  }'

# Check status
curl http://localhost:30001/api/runs/<run-id> | jq .

# View report
# http://localhost:30001/reports/<suite>/<run-id>/index.html
```

## 🔍 Verifying Reports

```bash
# Check PVC is mounted
kubectl describe pod -l app=test-runner-service -n test-runner | grep -A 5 "Mounts:"

# List reports on PVC
kubectl exec -it deployment/test-runner-service -n test-runner -- ls -la /app/reports

# Find HTML reports
kubectl exec -it deployment/test-runner-service -n test-runner -- find /app/reports -name "index.html"
```

## 🛠️ Troubleshooting

### Reports Not Showing

1. Check job completed:
   ```bash
   kubectl get jobs -n test-runner
   kubectl logs job/<job-name> -n test-runner
   ```

2. Check reports exist:
   ```bash
   kubectl exec -it deployment/test-runner-service -n test-runner -- ls -la /app/reports
   ```

3. Check service logs:
   ```bash
   kubectl logs -f deployment/test-runner-service -n test-runner
   ```

### PVC Issues

```bash
# Check PVC status
kubectl get pvc -n test-runner

# Check pod events
kubectl describe pod -l app=test-runner-service -n test-runner
```

### Images Not Found

```bash
# List images in kind
docker exec -it kind-control-plane crictl images | grep local

# Reload if needed
make load-all
```

See [QUICKSTART.md](QUICKSTART.md) for more troubleshooting tips.

## 📦 What's Included

### Services
- **test-runner-service** - Go backend that orchestrates test runs
- **test-runner-ui** - Vue.js frontend (Cartridge)
- **web-app** - Demo application under test
- **mock-api** - Mock API for the web app

### Infrastructure
- **Namespace** - Isolated `test-runner` namespace
- **RBAC** - Service account with job creation permissions
- **PVC** - Shared persistent volume for reports
- **NodePort Services** - Easy local access

## 🔄 Differences from Legacy Configs

### vs k8s-demo
- ✅ Uses PVC instead of emptyDir
- ✅ Includes RBAC configuration
- ✅ Reports persist across pod restarts

### vs k8s-remote
- ✅ Uses local images (no registry needed)
- ✅ Includes all supporting services
- ✅ Optimized for local development

## 🧹 Cleanup

```bash
# Delete everything
kubectl delete namespace test-runner

# Or use make
make clean

# Delete kind cluster
kind delete cluster --name kind
```

## 📖 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [kind Documentation](https://kind.sigs.k8s.io/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)

## 🤝 Contributing

When making changes:
1. Test with `make all`
2. Verify reports work end-to-end
3. Update documentation
4. Check all services are healthy

## 📝 Notes

- The unified configuration uses `ReadWriteOnce` PVC which works in kind because all pods are on the same node
- For production, consider using `ReadWriteMany` with NFS or similar
- NodePort services use ports 30000-30003 for easy local access
- All images use the `:local` tag for local development

## 🎓 Learning Path

1. Start with [QUICKSTART.md](QUICKSTART.md)
2. Review [ARCHITECTURE.md](ARCHITECTURE.md) for understanding
3. Read [k8s-unified/README.md](k8s-unified/README.md) for details
4. Check [REPORTS-FIX-SUMMARY.md](REPORTS-FIX-SUMMARY.md) for context

---

**Need help?** Check the troubleshooting sections in the docs or review the logs with `make logs`.
