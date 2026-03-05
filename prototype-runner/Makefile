SHELL := /bin/bash

# Directories inside root/
API_DIR := ../apps/mock-api
WEB_DIR := ../apps/ye-olde-magic-shop
FRONTEND_DIR := ../apps/cartridge
SERVICE_DIR := ../apps/console

# PID and log files (placed in final/)
API_PID := mock-api.pid
WEB_PID := web-app.pid
FRONTEND_PID := frontend.pid
SERVICE_PID := service.pid
API_LOG := mock-api.log
WEB_LOG := web-app.log
FRONTEND_LOG := frontend.log
SERVICE_LOG := service.log

# Kind cluster and k8s helper vars
KIND_CLUSTER ?= kind-test-runner
K8S_MANIFEST_DIR := local-cluster/k8s-demo
KUBEFILE := local-cluster/kubeconfig

.PHONY: help start-api start-web start start-service dev stop logs status start-frontend start-frontend-backend start-test-runner-service start-cluster stop-cluster rebuild-web-app rebuild-mock-api rebuild-test-runner-ui rebuild-test-runner-server destroy-cluster restart-port-forwards

help:
	@echo "Makefile targets:"
	@echo "  make start-api   # start mock API in background"
	@echo "  make start-web   # start web-app in background"
	@echo "  make start       # start both (both background)"
	@echo "  make dev         # start API (background) and run web-app in foreground"
	@echo "  make stop        # stop services started by this Makefile"
	@echo "  make logs        # tail both logs"
	@echo "  make status      # show pid files if present"
	@echo "  make restart-port-forwards # restart local cluster port forwarding"

start-api: ## Start the mock API in the background and write PID to $(API_PID)
	@echo "Starting mock API (in $(API_DIR)) on port 8001..."
	@cd $(API_DIR) && PORT=8001 npm start > ../$(API_LOG) 2>&1 & echo $$! > ../$(API_PID)
	@echo "Mock API PID: $$(cat $(API_PID))"

start-web: ## Start the web-app in the background and write PID to $(WEB_PID)
	@echo "Starting web-app (in $(WEB_DIR)) on port 8000..."
	@cd $(WEB_DIR) && PORT=8000 npm run dev:local > ../$(WEB_LOG) 2>&1 & echo $$! > ../$(WEB_PID)
	@echo "Web-app PID: $$(cat $(WEB_PID))"

start-frontend: ## Start the test-runner frontend Vite dev server only (background) on port 9000
	@echo "Starting frontend (in $(FRONTEND_DIR)) on port 9000..."
	@cd $(FRONTEND_DIR) && npx vite --port 9000 > ../$(FRONTEND_LOG) 2>&1 & echo $$! > ../$(FRONTEND_PID)
	@echo "Frontend PID: $$(cat $(FRONTEND_PID))"

start-test-runner-service: ## Start the Go test-runner-service on port 9001 (if Go is available)
	@echo "Starting test-runner-service (Go) on port 9001..."
	@cd $(SERVICE_DIR) && PORT=9001 nohup go run ./cmd/server > ../$(SERVICE_LOG) 2>&1 & echo $$! > ../$(SERVICE_PID) || (echo "Go not available or failed to start; run 'cd $(SERVICE_DIR) && go run ./cmd/server' manually" && false)
	@echo "Service PID: $$(cat $(SERVICE_PID) 2>/dev/null || echo 'none')"

start-frontend-backend: start-test-runner-service start-frontend ## Start backend (service) and frontend (Vite) in background
	@echo "Started test-runner-service and frontend (background). Logs: $(SERVICE_LOG), $(FRONTEND_LOG)"

start: start-api start-web start-test-runner-service ## start mock-api, web-app, and test-runner-service
	@echo "Started mock-api (8001), web-app (8000), and test-runner-service (9001)."

dev: start-api ## Start API in background, run web-app in foreground (blocks)
	@echo "Starting web-app in foreground (will block). Use Ctrl-C to exit."
	cd $(WEB_DIR) && npm run dev:local

stop: ## Stop services if pid files exist
	@echo "Stopping services..."
	-@if [ -f $(API_PID) ]; then kill $$(cat $(API_PID)) 2>/dev/null || true; rm -f $(API_PID); fi
	-@if [ -f $(WEB_PID) ]; then kill $$(cat $(WEB_PID)) 2>/dev/null || true; rm -f $(WEB_PID); fi
	-@if [ -f $(FRONTEND_PID) ]; then kill $$(cat $(FRONTEND_PID)) 2>/dev/null || true; rm -f $(FRONTEND_PID); fi
	-@if [ -f $(SERVICE_PID) ]; then kill $$(cat $(SERVICE_PID)) 2>/dev/null || true; rm -f $(SERVICE_PID); fi
	@echo "Stopped."

logs: ## Tail logs for services
	@echo "Tailing logs: $(API_LOG) $(WEB_LOG) $(FRONTEND_LOG) $(SERVICE_LOG)"
	@tail -f --pid=$$(cat $(API_PID) 2>/dev/null || echo $$) $(API_LOG) $(WEB_LOG) $(FRONTEND_LOG) $(SERVICE_LOG)

status: ## Show pid files (if present)
	@echo "API PID: $$(cat $(API_PID) 2>/dev/null || echo "none")"
	@echo "Web PID: $$(cat $(WEB_PID) 2>/dev/null || echo "none")"
	@echo "Frontend PID: $$(cat $(FRONTEND_PID) 2>/dev/null || echo "none")"
	@echo "Service PID: $$(cat $(SERVICE_PID) 2>/dev/null || echo "none")"

start-cluster: ## Create a kind cluster (if missing), build/load images and apply manifests
	@echo "Creating kind cluster '$(KIND_CLUSTER)' if missing..."
	@kind get clusters | grep -q "^$(KIND_CLUSTER)$$" || kind create cluster --name "$(KIND_CLUSTER)"
	@echo "Exporting kubeconfig to $(KUBEFILE)"
	@kind get kubeconfig --name "$(KIND_CLUSTER)" > $(KUBEFILE) || true
	@echo "Building and loading images into kind (mock-api, test-runner-service, test-runner-ui, web-app)..."
	@$(MAKE) rebuild-mock-api
	@$(MAKE) rebuild-test-runner-server
	@$(MAKE) rebuild-test-runner-ui
	@$(MAKE) rebuild-web-app
	@echo "Applying k8s manifests from $(K8S_MANIFEST_DIR)"
	@kubectl --kubeconfig $(KUBEFILE) apply -f $(K8S_MANIFEST_DIR)/00-namespace.yaml || true
	@kubectl --kubeconfig $(KUBEFILE) apply -f $(K8S_MANIFEST_DIR) || true
	@echo "Cluster started and manifests applied."
	@echo "Starting port-forwards..."
	@KUBECONFIG=$(KUBEFILE) $(K8S_MANIFEST_DIR)/start-port-forwards.sh

restart-port-forwards: ## Restart the port forwarding scripts
	@echo "Restarting port-forwards..."
	@$(K8S_MANIFEST_DIR)/stop-port-forwards.sh || true
	@KUBECONFIG=$(KUBEFILE) $(K8S_MANIFEST_DIR)/start-port-forwards.sh

stop-cluster: ## Delete the kind cluster if it exists
	@echo "Stopping port-forwards..."
	@$(K8S_MANIFEST_DIR)/stop-port-forwards.sh || true
	@echo "Deleting kind cluster '$(KIND_CLUSTER)' (if exists)..."
	@kind get clusters | grep -q "^$(KIND_CLUSTER)$$" && kind delete cluster --name "$(KIND_CLUSTER)" || echo "Cluster '$(KIND_CLUSTER)' not found; nothing to do."

destroy-cluster: stop-cluster ## Alias to stop-cluster, also remove generated kubeconfigs and PID files
	@echo "Removing generated kubeconfig files and PID files..."
	@rm -f $(K8S_MANIFEST_DIR)/kubeconfig || true
	@rm -f local-cluster/kubeconfig || true
	@rm -f local-cluster/.test-runner-port-forward.pid || true
	@rm -f /tmp/portfwd-*.pid || true
	@rm -f /tmp/portfwd-*.log || true
	@echo "Destroyed cluster artifacts."

publish: ## Build and push images to registry, generate remote manifests. Usage: make publish REGISTRY=username
	@scripts/publish_images.sh $(REGISTRY)

# Rebuild and redeploy helpers
rebuild-mock-api: ## Build mock-api image, load into kind and restart deployment
	@echo "Building mock-api image..."
	@cd $(API_DIR) && docker build -t mock-api:local .
	@echo "Loading mock-api:local into kind cluster '$(KIND_CLUSTER)'..."
	@kind load docker-image mock-api:local --name "$(KIND_CLUSTER)" || kind load docker-image mock-api:local || true
	@echo "Restarting mock-api deployment in namespace test-runner (if present)..."
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout restart deployment/mock-api || true
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout status deployment/mock-api --timeout=60s || true

rebuild-web-app: ## Build web-app image, load into kind and restart web-app deployment
	@echo "Building web-app image..."
	@cd $(WEB_DIR) && docker build -t web-app:local .
	@echo "Loading web-app:local into kind cluster '$(KIND_CLUSTER)'..."
	@kind load docker-image web-app:local --name "$(KIND_CLUSTER)" || kind load docker-image web-app:local || true
	@echo "Restarting web-app deployment in namespace test-runner (if present)..."
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout restart deployment/web-app || true
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout status deployment/web-app --timeout=60s || true

rebuild-test-runner-ui: ## Build test-runner-ui image, load into kind and restart deployment
	@echo "Building test-runner-ui image..."
	@cd $(FRONTEND_DIR) && docker build -t test-runner-ui:local .
	@echo "Loading test-runner-ui:local into kind cluster '$(KIND_CLUSTER)'..."
	@kind load docker-image test-runner-ui:local --name "$(KIND_CLUSTER)" || kind load docker-image test-runner-ui:local || true
	@echo "Restarting test-runner-ui deployment in namespace test-runner (if present)..."
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout restart deployment/test-runner-ui || true
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout status deployment/test-runner-ui --timeout=60s || true

rebuild-test-runner-server: ## Build test-runner-service image, load into kind and restart deployment
	@echo "Building test-runner-service image..."
	@cd $(SERVICE_DIR) && docker build -t test-runner-service:local .
	@echo "Loading test-runner-service:local into kind cluster '$(KIND_CLUSTER)'..."
	@kind load docker-image test-runner-service:local --name "$(KIND_CLUSTER)" || kind load docker-image test-runner-service:local || true
	@echo "Restarting test-runner-service deployment in namespace test-runner (if present)..."
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout restart deployment/test-runner-service || true
	@kubectl --kubeconfig $(KUBEFILE) -n test-runner rollout status deployment/test-runner-service --timeout=60s || true
