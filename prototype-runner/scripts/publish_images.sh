#!/usr/bin/env bash
set -e

REGISTRY=${1:-orieken}
echo "Building and pushing images to registry: $REGISTRY"

# Build and push mock-api
echo "--- Processing mock-api ---"
cd mock-api
docker build -t $REGISTRY/mock-api:latest .
docker push $REGISTRY/mock-api:latest
cd ..

# Build and push web-app
echo "--- Processing web-app ---"
cd web-app
docker build -t $REGISTRY/web-app:latest .
docker push $REGISTRY/web-app:latest
cd ..

# Build and push test-runner-ui
echo "--- Processing test-runner-ui ---"
cd test-runner-ui
docker build -t $REGISTRY/test-runner-ui:latest .
docker push $REGISTRY/test-runner-ui:latest
cd ..

# Build and push test-runner-service
echo "--- Processing test-runner-service ---"
cd test-runner-service
docker build -t $REGISTRY/test-runner-service:latest .
docker push $REGISTRY/test-runner-service:latest
cd ..

# Build and push cucumber-project
echo "--- Processing cucumber-project ---"
cd cucumber-project
docker build -t $REGISTRY/cucumber-project:latest .
docker push $REGISTRY/cucumber-project:latest
cd ..

echo "--- Generating Remote Kubernetes manifests ---"
REMOTE_K8S_DIR="local-cluster/k8s-remote"
rm -rf "$REMOTE_K8S_DIR"
cp -r "local-cluster/k8s-demo" "$REMOTE_K8S_DIR"

# Update image references in k8s-remote
# mock-api:local -> registry/mock-api:latest
# etc.

# Using sed to replace images
# We assume the local images are named:
# mock-api:local
# web-app:local
# test-runner-ui:local
# test-runner-service:local
# cucumber-project:local (though this doesn't have a k8s deployment usually, unless demo-job uses it? demo-job uses alpine per grep)

find "$REMOTE_K8S_DIR" -name "*.yaml" -print0 | xargs -0 sed -i '' "s|mock-api:local|$REGISTRY/mock-api:latest|g"
find "$REMOTE_K8S_DIR" -name "*.yaml" -print0 | xargs -0 sed -i '' "s|web-app:local|$REGISTRY/web-app:latest|g"
find "$REMOTE_K8S_DIR" -name "*.yaml" -print0 | xargs -0 sed -i '' "s|test-runner-ui:local|$REGISTRY/test-runner-ui:latest|g"
find "$REMOTE_K8S_DIR" -name "*.yaml" -print0 | xargs -0 sed -i '' "s|test-runner-service:local|$REGISTRY/test-runner-service:latest|g"

# Update image pull policy to Always for remote deployments to ensure freshness
find "$REMOTE_K8S_DIR" -name "*.yaml" -print0 | xargs -0 sed -i '' "s|imagePullPolicy: IfNotPresent|imagePullPolicy: Always|g"

echo "--- Configuring CUCUMBER_IMAGE env var for service ---"
# Add env var to service-deployment.yaml if not present, to tell service to use the remote cucumber image
# We'll look for the existing PORT env var block and append below it.
# Assuming standard structure in service-deployment.yaml

SERVICE_DEPLOY="$REMOTE_K8S_DIR/service-deployment.yaml"
# We inject the env var using sed. 
# Identifying the env block.
# We replace '- name: DEFAULT_EXECUTOR' with the new block
# This inserts it before DEFAULT_EXECUTOR.

sed -i '' "s|- name: DEFAULT_EXECUTOR|- name: CUCUMBER_IMAGE\\
              value: \"$REGISTRY/cucumber-project:latest\"\\
            - name: DEFAULT_EXECUTOR|g" "$SERVICE_DEPLOY"

echo "Done! Remote manifests generated in $REMOTE_K8S_DIR"
echo "You can deploy these to any cluster by running:"
echo "kubectl apply -f $REMOTE_K8S_DIR/00-namespace.yaml"
echo "kubectl apply -f $REMOTE_K8S_DIR"
