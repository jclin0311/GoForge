#!/bin/bash
set -e

# Deployment script for GoForge platform

NAMESPACE=${1:-goforge}
ENVIRONMENT=${2:-development}

echo "Deploying GoForge to namespace: $NAMESPACE (environment: $ENVIRONMENT)"

# Create namespace if it doesn't exist
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Label namespace for Istio injection
kubectl label namespace $NAMESPACE istio-injection=enabled --overwrite

# Apply ConfigMaps
echo "Applying ConfigMaps..."
kubectl apply -f deployments/k8s/configmaps/ -n $NAMESPACE

# Apply Deployments
echo "Deploying services..."
kubectl apply -f deployments/k8s/deployments/ -n $NAMESPACE

# Apply Ingress
echo "Configuring Ingress..."
kubectl apply -f deployments/k8s/ingress/ -n $NAMESPACE

# Apply Istio resources
if kubectl get namespace istio-system &> /dev/null; then
    echo "Applying Istio configurations..."
    kubectl apply -f deployments/k8s/istio/ -n $NAMESPACE
fi

# Wait for deployments
echo "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment --all -n $NAMESPACE

echo "✓ Deployment complete!"
echo ""
echo "Access services with:"
echo "  kubectl port-forward -n $NAMESPACE svc/api-gateway 50050:50050"
echo "  kubectl port-forward -n $NAMESPACE svc/model-service 50051:50051"
