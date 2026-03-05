# Remote Deployment Guide

This guide explains how to deploy the "Cartridge" Test Runner stack on a fresh machine (e.g., a colleague's laptop or a cloud VM) using the remote Kubernetes manifests.

## Prerequisites

On the target machine, ensure you have the following installed:

1.  **Docker** (Desktop or Engine)
2.  **Kind** (Kubernetes in Docker) - *optional if you have another k8s cluster*
    *   Install via Homebrew: `brew install kind`
    *   Or binary: [https://kind.sigs.k8s.io/docs/user/quick-start/](https://kind.sigs.k8s.io/docs/user/quick-start/)
3.  **Kubectl**
    *   Install via Homebrew: `brew install kubectl`

## Podman Configuration

If you are using **Podman** instead of Docker, you must configure `kind` to use the Podman provider.

1.  **System Setup:** Ensure your user belongs to a group allowed to use Podman (usually default on modern installs), and that the Podman socket is active if needed (though Kind typically runs "rootless" or via the CLI wrapper).
2.  **Kind Configuration:**
    `kind` supports Podman experimentally. You need to set the `KIND_EXPERIMENTAL_PROVIDER` environment variable before creating the cluster.

    ```bash
    # Create the cluster using podman
    KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster --name test-runner
    ```

    *Note: Ensure you have `podman` version 3.0+.*

## Deployment Steps

1.  **Copy the Manifests**
    Copy the `local-cluster/k8s-remote` folder from your development machine to the target machine.

2.  **Create a Cluster** (if expected)
    ```bash
    kind create cluster --name test-runner
    ```

3.  **Apply Manifests**
    Go into the `k8s-remote` directory:
    ```bash
    cd k8s-remote
    ```

    Apply the namespace first, then the rest of the resources:
    ```bash
    kubectl apply -f 00-namespace.yaml
    kubectl apply -f .
    ```

4.  **Verify Pods**
    Check if the pods are pulling the images (from Docker Hub) and starting up:
    ```bash
    kubectl -n test-runner get pods -w
    ```
    Wait until they are `Running`.

5.  **Access the Application**
    To access the services from your host machine (browser), you need to forward the ports. A script is provided for convenience, but check the `KUBECONFIG` variable.

    If you used `kind create cluster`, your config is likely at `~/.kube/config`.

    ```bash
    # Ensure it uses your default config, overriding the Makefile-specific default
    export KUBECONFIG=~/.kube/config
    
    ./start-port-forwards.sh
    ```

    You should now be able to access:
    *   **UI:** [http://localhost:9000](http://localhost:9000)
    *   **Service Root:** [http://localhost:9001](http://localhost:9001)
    *   **App:** [http://localhost:8000](http://localhost:8000)

## Troubleshooting

*   **Image Pull Errors:** Ensure the images were successfully pushed to Docker Hub and the repository is public (or you have configured image pull secrets).
    *   **Note:** The `test-runner-demo-job` uses `alpine:3.18`. If this fails with `ImagePullBackOff`, check if the node has internet access to pull from Docker Hub.
*   **Port Conflicts:** Ensure ports 8000, 8001, 9000, and 9001 are free on the target machine.
*   **Service Access:** If `NodePort` services are not accessible, rely on the `start-port-forwards.sh` script (which uses `kubectl port-forward`). This is the most reliable way to access services in Kind across different OSs (Mac/Linux/Windows).

## Teardown

To stop the application and remove the cluster:

1.  **Stop Port Forwarding:**
    Run the provided stop script to clean up background port-forward processes:
    ```bash
    ./stop-port-forwards.sh
    ```

2.  **Delete the Cluster:**
    If you created a specific kind cluster named `test-runner`:
    ```bash
    kind delete cluster --name test-runner
    ```
    This completely removes the cluster and all installed resources.
