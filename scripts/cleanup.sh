#!/bin/bash
set -e

NAMESPACE=${1:-goforge}

echo "Cleaning up GoForge deployment from namespace: $NAMESPACE"
echo "WARNING: This will delete all resources in the $NAMESPACE namespace!"
read -p "Are you sure? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Cleanup cancelled"
    exit 0
fi

# Delete Istio resources
echo "Deleting Istio configurations..."
kubectl delete -f deployments/k8s/istio/ -n $NAMESPACE --ignore-not-found=true

# Delete Ingress
echo "Deleting Ingress..."
kubectl delete -f deployments/k8s/ingress/ -n $NAMESPACE --ignore-not-found=true

# Delete Deployments
echo "Deleting deployments..."
kubectl delete -f deployments/k8s/deployments/ -n $NAMESPACE --ignore-not-found=true

# Delete ConfigMaps
echo "Deleting ConfigMaps..."
kubectl delete -f deployments/k8s/configmaps/ -n $NAMESPACE --ignore-not-found=true

# Delete namespace
echo "Deleting namespace..."
kubectl delete namespace $NAMESPACE --ignore-not-found=true

echo "✓ Cleanup complete!"
