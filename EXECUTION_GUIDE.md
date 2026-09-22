# GoForge - Complete Execution Guide

This guide provides step-by-step instructions for deploying, operating, and managing the GoForge AI Infrastructure Platform.

## Table of Contents

1. [Initial Setup](#initial-setup)
2. [Local Development](#local-development)
3. [Kubernetes Deployment](#kubernetes-deployment)
4. [Production Deployment](#production-deployment)
5. [Day-2 Operations](#day-2-operations)
6. [Performance Tuning](#performance-tuning)
7. [Troubleshooting](#troubleshooting)

---

## Initial Setup

### Prerequisites Installation

#### 1. Install Docker

**macOS:**
```bash
brew install --cask docker
```

**Linux (Ubuntu/Debian):**
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

**Verify:**
```bash
docker --version
# Docker version 24.0.0 or higher
```

#### 2. Install Kubernetes (Minikube for Local)

**macOS:**
```bash
brew install minikube
```

**Linux:**
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
```

**Start Minikube:**
```bash
minikube start --cpus=4 --memory=8192 --driver=docker
minikube addons enable ingress
minikube addons enable metrics-server
```

**Verify:**
```bash
kubectl get nodes
# NAME       STATUS   ROLES           AGE   VERSION
# minikube   Ready    control-plane   1m    v1.29.0
```

#### 3. Install Helm

**macOS:**
```bash
brew install helm
```

**Linux:**
```bash
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

**Verify:**
```bash
helm version
# version.BuildInfo{Version:"v3.14.0"...}
```

#### 4. Install Istio

```bash
# Download Istio
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.20.0
export PATH=$PWD/bin:$PATH

# Install Istio
istioctl install --set profile=demo -y

# Verify installation
kubectl get pods -n istio-system
```

#### 5. Install Go (for development)

**macOS:**
```bash
brew install go@1.22
```

**Linux:**
```bash
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

**Verify:**
```bash
go version
# go version go1.22.0
```

#### 6. Install Protocol Buffers Compiler

**macOS:**
```bash
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Linux:**
```bash
sudo apt-get install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

**Add to PATH:**

macOS (zsh):
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

Linux (bash):
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

**Verify:**
```bash
protoc --version
# libprotoc 3.x.x or higher

protoc-gen-go --version
# protoc-gen-go v1.x.x

protoc-gen-go-grpc --version
# protoc-gen-go-grpc 1.x.x
```

---

## Local Development

### Step 1: Clone and Setup

```bash
# Clone repository
git clone https://github.com/goforge/ai-platform.git
cd ai-platform

# Install Go dependencies
go mod download
```

> **Before generating protobuf code**, ensure the protoc plugins installed in Step 6 above are on your PATH in the current shell session:
> ```bash
> export PATH=$PATH:$(go env GOPATH)/bin
> ```

```bash
# Generate protobuf Go code into pkg/proto/
# This resolves "could not import github.com/goforge/ai-platform/pkg/proto" errors
make proto
```

### Step 2: Build Services

```bash
# Build all services
make build

# Or build individually
cd cmd/model-service && go build -o ../../bin/model-service
cd cmd/auth-service && go build -o ../../bin/auth-service
cd cmd/api-gateway && go build -o ../../bin/api-gateway
cd cmd/metrics-service && go build -o ../../bin/metrics-service
```

### Step 3: Run Services Locally

**Option A: Using Make (Recommended)**

Open 4 terminal windows:

```bash
# Terminal 1: Auth Service
make run-auth-service

# Terminal 2: Model Service
make run-model-service

# Terminal 3: Metrics Service
make run-metrics-service

# Terminal 4: API Gateway
make run-api-gateway
```

**Option B: Manual Execution**

```bash
# Set environment variables
export SERVICE_NAME=model-service
export GRPC_PORT=50051
export LOG_LEVEL=debug

# Run service
./bin/model-service
```

### Step 4: Test Locally

Install grpcurl:
```bash
brew install grpcurl  # macOS
# or
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Test endpoints:
```bash
# List services
grpcurl -plaintext localhost:50051 list

# Test Model Service - Predict
grpcurl -plaintext -d '{
  "model_id": "gpt-neo",
  "model_version": "1.0",
  "input_data": "SGVsbG8gV29ybGQ=",
  "request_id": "req-001"
}' localhost:50051 goforge.model.ModelService/Predict

# Test Auth Service - Authenticate
grpcurl -plaintext -d '{
  "username": "admin",
  "password": "secret",
  "client_id": "web-app"
}' localhost:50052 goforge.auth.AuthService/Authenticate

# Test Health Check
grpcurl -plaintext localhost:50051 goforge.model.ModelService/HealthCheck
```

---

## Kubernetes Deployment

### Step 1: Prepare Kubernetes Cluster

```bash
# Verify cluster is running
kubectl cluster-info

# Create namespace with Istio injection
kubectl create namespace goforge
kubectl label namespace goforge istio-injection=enabled

# Verify Istio injection label
kubectl get namespace goforge --show-labels
```

### Step 2: Build and Push Docker Images

**Option A: Local Minikube (No Registry Needed)**

```bash
# Use Minikube's Docker daemon
eval $(minikube docker-env)

# Build images
make docker

# Verify images
docker images | grep goforge
```

**Option B: Push to Container Registry**

```bash
# Tag images for your registry
export DOCKER_REGISTRY=your-registry.io
docker tag goforge/api-gateway:v1.0.0 $DOCKER_REGISTRY/api-gateway:v1.0.0
docker tag goforge/model-service:v1.0.0 $DOCKER_REGISTRY/model-service:v1.0.0
docker tag goforge/auth-service:v1.0.0 $DOCKER_REGISTRY/auth-service:v1.0.0
docker tag goforge/metrics-service:v1.0.0 $DOCKER_REGISTRY/metrics-service:v1.0.0

# Push to registry
docker push $DOCKER_REGISTRY/api-gateway:v1.0.0
docker push $DOCKER_REGISTRY/model-service:v1.0.0
docker push $DOCKER_REGISTRY/auth-service:v1.0.0
docker push $DOCKER_REGISTRY/metrics-service:v1.0.0
```

### Step 3: Deploy with Helm

```bash
# Install GoForge platform
make helm-install

# Or manually with custom values
helm upgrade --install goforge deployments/helm/goforge \
  --namespace goforge \
  --create-namespace \
  --values deployments/helm/goforge/values.yaml \
  --set global.imageRegistry=$DOCKER_REGISTRY \
  --wait

# Verify deployment
kubectl get pods -n goforge
kubectl get svc -n goforge
```

**Expected Output:**
```
NAME                              READY   STATUS    RESTARTS   AGE
api-gateway-7d8f5c8b9d-abcde      2/2     Running   0          2m
api-gateway-7d8f5c8b9d-fghij      2/2     Running   0          2m
model-service-6c9d8f7b6c-klmno    2/2     Running   0          2m
model-service-6c9d8f7b6c-pqrst    2/2     Running   0          2m
auth-service-5b8c7d6a5b-uvwxy     2/2     Running   0          2m
metrics-service-4a7b6c5d4a-zzabc  2/2     Running   0          2m
```

### Step 4: Deploy Istio Resources

```bash
# Apply Istio Gateway
kubectl apply -f deployments/k8s/istio/gateway.yaml

# Apply Virtual Services
kubectl apply -f deployments/k8s/istio/destination-rules.yaml

# Enable mTLS
kubectl apply -f deployments/k8s/istio/peer-authentication.yaml

# Verify Istio configuration
istioctl analyze -n goforge
```

### Step 5: Install Monitoring Stack

```bash
# Add Helm repositories
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Install Prometheus and Grafana
make monitoring-install

# Or manually
helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --create-namespace \
  --values deployments/monitoring/prometheus-values.yaml \
  --wait

# Deploy ServiceMonitors
kubectl apply -f deployments/monitoring/servicemonitor.yaml

# Deploy Prometheus Rules
kubectl apply -f deployments/monitoring/prometheus-rules.yaml
```

### Step 6: Access Services

#### Port Forwarding (Development)

```bash
# API Gateway
kubectl port-forward -n goforge svc/api-gateway 50050:50050 &

# Model Service
kubectl port-forward -n goforge svc/model-service 50051:50051 &

# Grafana
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80 &

# Prometheus
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090 &
```

Access in browser:
- Grafana: http://localhost:3000 (admin/admin123)
- Prometheus: http://localhost:9090

#### Using Minikube Tunnel (Recommended for Local)

```bash
# Start tunnel (requires sudo)
minikube tunnel

# Get Istio Ingress Gateway external IP
kubectl get svc -n istio-system istio-ingressgateway

# Add to /etc/hosts
echo "$(minikube ip) api.goforge.local model.goforge.local" | sudo tee -a /etc/hosts
```

Access via:
- http://api.goforge.local
- http://model.goforge.local

---

## Production Deployment

### Step 1: Pre-Production Checklist

- [ ] Production Kubernetes cluster provisioned (GKE/EKS/AKS)
- [ ] Container registry configured
- [ ] DNS records configured
- [ ] TLS certificates obtained
- [ ] Secrets management solution deployed (Vault/Sealed Secrets)
- [ ] Backup strategy defined
- [ ] Monitoring alerts configured
- [ ] Disaster recovery plan documented

### Step 2: Configure Production Values

Create `values-production.yaml`:

```yaml
global:
  namespace: goforge-prod
  imageRegistry: your-registry.io
  imageTag: v1.0.0
  imagePullPolicy: Always

replicaCount:
  apiGateway: 5
  modelService: 10
  authService: 3
  metricsService: 3

resources:
  modelService:
    requests:
      memory: "512Mi"
      cpu: "500m"
    limits:
      memory: "2Gi"
      cpu: "2000m"

autoscaling:
  enabled: true
  modelService:
    minReplicas: 10
    maxReplicas: 50
    targetCPUUtilization: 60

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: api.goforge.io
      tls:
        - secretName: goforge-tls
          hosts:
            - api.goforge.io

monitoring:
  prometheus:
    enabled: true
    retention: 90d
    storageSize: 200Gi
```

### Step 3: Deploy to Production

```bash
# Create production namespace
kubectl create namespace goforge-prod
kubectl label namespace goforge-prod istio-injection=enabled

# Install with production values
helm upgrade --install goforge-prod deployments/helm/goforge \
  --namespace goforge-prod \
  --values values-production.yaml \
  --atomic \
  --timeout 10m

# Verify deployment
kubectl get pods -n goforge-prod
kubectl get hpa -n goforge-prod
```

### Step 4: Configure TLS

```bash
# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Create ClusterIssuer
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@goforge.io
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF
```

### Step 5: Implement GitOps (Optional but Recommended)

#### Using ArgoCD

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Create Application
cat <<EOF | kubectl apply -f -
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: goforge
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/goforge/ai-platform.git
    targetRevision: main
    path: deployments/helm/goforge
    helm:
      valueFiles:
      - values-production.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: goforge-prod
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
EOF
```

---

## Day-2 Operations

### Scaling

#### Manual Scaling

```bash
# Scale deployment
kubectl scale deployment model-service -n goforge --replicas=10

# Update HPA
kubectl patch hpa model-service-hpa -n goforge --patch '{"spec":{"maxReplicas":20}}'
```

#### Automatic Scaling

HPA automatically scales based on metrics:
```bash
# Check HPA status
kubectl get hpa -n goforge

# Describe HPA for details
kubectl describe hpa model-service-hpa -n goforge
```

### Updates and Rollouts

#### Rolling Update

```bash
# Update image version
kubectl set image deployment/model-service -n goforge \
  model-service=goforge/model-service:v1.1.0

# Check rollout status
kubectl rollout status deployment/model-service -n goforge

# View rollout history
kubectl rollout history deployment/model-service -n goforge
```

#### Canary Deployment

```bash
# Deploy canary version with Istio
cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: model-service
  namespace: goforge
spec:
  hosts:
  - model-service
  http:
  - route:
    - destination:
        host: model-service
        subset: v1
      weight: 90
    - destination:
        host: model-service
        subset: v2
      weight: 10
EOF

# Monitor canary metrics
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80

# Gradually increase traffic to v2
# Update weight to 50/50, then 0/100
```

#### Rollback

```bash
# Rollback to previous version
kubectl rollout undo deployment/model-service -n goforge

# Rollback to specific revision
kubectl rollout undo deployment/model-service -n goforge --to-revision=2

# Verify rollback
kubectl rollout status deployment/model-service -n goforge
```

### Configuration Updates

```bash
# Edit ConfigMap
kubectl edit configmap app-config -n goforge

# Restart deployments to pick up changes
kubectl rollout restart deployment/model-service -n goforge
kubectl rollout restart deployment/api-gateway -n goforge
```

### Backup and Restore

#### Backup with Velero

```bash
# Install Velero
velero install --provider aws --bucket goforge-backups

# Create backup
velero backup create goforge-backup-$(date +%Y%m%d) \
  --include-namespaces goforge-prod

# Schedule daily backups
velero schedule create goforge-daily \
  --schedule="0 2 * * *" \
  --include-namespaces goforge-prod
```

#### Restore

```bash
# List backups
velero backup get

# Restore from backup
velero restore create --from-backup goforge-backup-20260601
```

---

## Performance Tuning

### Optimize Resource Allocation

Monitor resource usage:
```bash
# Get resource consumption
kubectl top nodes
kubectl top pods -n goforge

# Adjust resource requests/limits
kubectl set resources deployment model-service -n goforge \
  --requests=cpu=500m,memory=512Mi \
  --limits=cpu=2000m,memory=2Gi
```

### Tune gRPC Connection Pool

Edit `pkg/config/config.go`:
```go
GRPC: GRPCConfig{
    MaxConnectionIdle:     5 * time.Minute,
    MaxConnectionAge:      10 * time.Minute,
    MaxConnectionAgeGrace: 5 * time.Minute,
    Time:                  2 * time.Minute,
    Timeout:               20 * time.Second,
},
```

### Optimize Circuit Breaker

Edit `pkg/config/config.go`:
```go
CircuitBreaker: CircuitBreakerConfig{
    MaxRequests: 10,           // Half-open state max requests
    Interval:    60 * time.Second,  // Closed state counter reset
    Timeout:     30 * time.Second,  // Open state duration
},
```

### Tune Rate Limiting

```go
RateLimit: RateLimitConfig{
    RequestsPerSecond: 100,
    BurstSize:         200,
},
```

### Database Connection Pooling (Future Enhancement)

```go
// Add to config when database is integrated
DB: DBConfig{
    MaxOpenConns:    25,
    MaxIdleConns:    25,
    ConnMaxLifetime: 5 * time.Minute,
},
```

---

## Troubleshooting

### Pod Crashes

```bash
# Check pod status
kubectl get pods -n goforge

# Describe pod for events
kubectl describe pod <pod-name> -n goforge

# View logs
kubectl logs <pod-name> -n goforge -c model-service

# View previous container logs (if crashed)
kubectl logs <pod-name> -n goforge -c model-service --previous

# Check Istio sidecar logs
kubectl logs <pod-name> -n goforge -c istio-proxy
```

### High Latency

```bash
# Check service mesh metrics
istioctl dashboard kiali

# Analyze request flow
kubectl port-forward -n istio-system svc/kiali 20001:20001

# Check Prometheus metrics
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
# Query: histogram_quantile(0.95, rate(grpc_request_duration_seconds_bucket[5m]))
```

### Circuit Breaker Issues

```bash
# Check logs for circuit breaker state changes
kubectl logs deployment/model-service -n goforge | grep "Circuit breaker"

# Verify destination rules
kubectl get destinationrule -n goforge
kubectl describe destinationrule model-service -n goforge
```

### Service Discovery Problems

```bash
# Check service endpoints
kubectl get endpoints -n goforge

# Test DNS resolution from a pod
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- /bin/bash
nslookup model-service.goforge.svc.cluster.local

# Verify service mesh configuration
istioctl analyze -n goforge
```

### Memory Leaks

```bash
# Monitor memory usage over time
kubectl top pod -n goforge --containers

# Enable Go profiling (add to main.go)
import _ "net/http/pprof"
go func() {
    http.ListenAndServe(":6060", nil)
}()

# Access pprof
kubectl port-forward <pod-name> 6060:6060
go tool pprof http://localhost:6060/debug/pprof/heap
```

### Istio Sidecar Issues

```bash
# Check sidecar injection
kubectl get pod <pod-name> -n goforge -o jsonpath='{.spec.containers[*].name}'

# Verify Istio configuration
istioctl proxy-status

# Check sidecar logs
kubectl logs <pod-name> -n goforge -c istio-proxy
```

---

## Maintenance Tasks

### Log Rotation

```bash
# Configure log retention in Kubernetes
kubectl edit deployment model-service -n goforge
# Add to container spec:
# args:
#   - --log-max-size=100
#   - --log-max-backups=3
```

### Certificate Renewal

Cert-manager handles this automatically, but verify:
```bash
# Check certificate status
kubectl get certificate -n goforge

# Force renewal
kubectl delete certificaterequest <name> -n goforge
```

### Database Maintenance (When Integrated)

```bash
# Backup database
kubectl exec -it postgres-0 -n goforge -- pg_dump -U postgres goforge > backup.sql

# Vacuum and analyze
kubectl exec -it postgres-0 -n goforge -- psql -U postgres -c "VACUUM ANALYZE;"
```

---

## Appendix

### Useful Commands Reference

```bash
# Quick status check
kubectl get all -n goforge

# Watch pod status
watch kubectl get pods -n goforge

# Tail logs from all model-service pods
kubectl logs -f -l app=model-service -n goforge --all-containers=true

# Execute command in pod
kubectl exec -it <pod-name> -n goforge -- /bin/sh

# Copy files from pod
kubectl cp goforge/<pod-name>:/tmp/file.log ./file.log

# Get resource usage
kubectl top nodes
kubectl top pods -n goforge --containers

# Debug networking
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- /bin/bash
```

### Performance Benchmarking

```bash
# Install ghz (gRPC benchmarking tool)
go install github.com/bojand/ghz/cmd/ghz@latest

# Benchmark model service
ghz --insecure \
  --proto api/proto/model.proto \
  --call goforge.model.ModelService/Predict \
  -d '{"model_id":"test","model_version":"1.0","input_data":"dGVzdA==","request_id":"bench-1"}' \
  -c 50 \
  -n 10000 \
  localhost:50051
```

---

**End of Execution Guide**

For additional support, refer to the main README.md or contact the GoForge team.
