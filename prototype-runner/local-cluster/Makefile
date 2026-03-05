.PHONY: help build-all load-all deploy clean test-run logs

# Default target
help:
	@echo "Unified K8s Test Runner - Make targets:"
	@echo ""
	@echo "  make build-all    - Build all Docker images"
	@echo "  make load-all     - Load all images into kind cluster"
	@echo "  make deploy       - Deploy all manifests to cluster"
	@echo "  make all          - Build, load, and deploy everything"
	@echo "  make test-run     - Trigger a test run via API"
	@echo "  make logs         - Show logs from test-runner-service"
	@echo "  make clean        - Delete all resources from cluster"
	@echo "  make port-forward - Start port forwarding for local access"
	@echo ""

# Build all images
build-all:
	@echo "Building all images..."
	cd ../../apps/console && docker build -t test-runner-service:local .
	cd ../../apps/cartridge && docker build -t test-runner-ui:local .
	cd ../../mock-api && docker build -t mock-api:local .
	cd ../../web-app && docker build -t web-app:local .
	cd ../../cucumber-project && docker build -t cucumber-project:local .
	@echo "All images built successfully!"

# Load images into kind
load-all:
	@echo "Loading images into kind cluster..."
	kind load docker-image test-runner-service:local --name kind
	kind load docker-image test-runner-ui:local --name kind
	kind load docker-image mock-api:local --name kind
	kind load docker-image web-app:local --name kind
	kind load docker-image cucumber-project:local --name kind
	@echo "All images loaded successfully!"

# Deploy all manifests
deploy:
	@echo "Deploying manifests..."
	kubectl apply -f k8s-unified/00-namespace.yaml
	kubectl apply -f k8s-unified/rbac.yaml
	kubectl apply -f k8s-unified/pvc.yaml
	kubectl apply -f k8s-unified/service-deployment.yaml
	kubectl apply -f k8s-unified/ui-deployment.yaml
	kubectl apply -f k8s-unified/mock-api-deployment.yaml
	kubectl apply -f k8s-unified/web-app-deployment.yaml
	kubectl apply -f k8s-unified/nodeport-services.yaml
	@echo "Waiting for pods to be ready..."
	kubectl wait --for=condition=ready pod -l app=test-runner-service -n test-runner --timeout=120s || true
	kubectl wait --for=condition=ready pod -l app=test-runner-ui -n test-runner --timeout=120s || true
	kubectl wait --for=condition=ready pod -l app=mock-api -n test-runner --timeout=120s || true
	kubectl wait --for=condition=ready pod -l app=web-app -n test-runner --timeout=120s || true
	@echo "Deployment complete!"

# Build, load, and deploy everything
all: build-all load-all deploy
	@echo "Full deployment complete!"
	@echo ""
	@echo "Access the services:"
	@echo "  UI:      http://localhost:30000"
	@echo "  API:     http://localhost:30001"
	@echo "  Web App: http://localhost:30003"
	@echo ""
	@echo "Or use port-forwarding:"
	@echo "  make port-forward"

# Trigger a test run
test-run:
	@echo "Triggering test run..."
	@curl -X POST http://localhost:30001/api/runs \
		-H "Content-Type: application/json" \
		-d '{"framework":"cucumber","suiteId":"final-cucumber-project","scenarioId":"inventory-list","executor":"k8s"}' \
		| jq .

# Show logs
logs:
	kubectl logs -f deployment/test-runner-service -n test-runner

# Clean up
clean:
	@echo "Cleaning up resources..."
	kubectl delete namespace test-runner --ignore-not-found=true
	@echo "Cleanup complete!"

# Port forwarding
port-forward:
	@echo "Starting port forwarding..."
	@echo "UI will be available at http://localhost:9000"
	@echo "API will be available at http://localhost:9001"
	@echo "Press Ctrl+C to stop"
	kubectl port-forward svc/test-runner-service 9001:9001 -n test-runner & \
	kubectl port-forward svc/test-runner-ui 9000:9000 -n test-runner

# Check status
status:
	@echo "Checking cluster status..."
	@echo ""
	@echo "Pods:"
	kubectl get pods -n test-runner
	@echo ""
	@echo "Services:"
	kubectl get svc -n test-runner
	@echo ""
	@echo "PVC:"
	kubectl get pvc -n test-runner
	@echo ""
	@echo "Jobs:"
	kubectl get jobs -n test-runner
