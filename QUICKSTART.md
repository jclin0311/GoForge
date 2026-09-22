# GoForge - Quick Start Guide

Get GoForge up and running in 10 minutes!

## Prerequisites

- Docker installed
- Kubernetes cluster (Minikube for local)
- kubectl configured
- 8GB RAM minimum

## Option 1: Local Development (Fastest)

### 1. Install Dependencies

```bash
# macOS
brew install go protobuf grpcurl

# Linux
sudo apt-get install golang-go protobuf-compiler
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### 2. Clone and Build

```bash
git clone https://github.com/goforge/ai-platform.git
cd ai-platform
go mod download
make proto
make build
```

### 3. Run Services

Open 4 terminals:

```bash
# Terminal 1
make run-auth-service

# Terminal 2
make run-model-service

# Terminal 3
make run-metrics-service

# Terminal 4
make run-api-gateway
```

### 4. Test

```bash
# Test model service
grpcurl -plaintext -d '{
  "model_id": "gpt-neo",
  "model_version": "1.0",
  "input_data": "SGVsbG8=",
  "request_id": "test-1"
}' localhost:50051 goforge.model.ModelService/Predict

# Test auth service
grpcurl -plaintext -d '{
  "username": "admin",
  "password": "secret",
  "client_id": "web"
}' localhost:50052 goforge.auth.AuthService/Authenticate
```

**Done!** Services are running locally.

---

## Option 2: Kubernetes (Minikube)

### 1. Start Minikube

```bash
minikube start --cpus=4 --memory=8192
minikube addons enable ingress
minikube addons enable metrics-server
```

### 2. Build Images

```bash
eval $(minikube docker-env)
make docker
```

### 3. Deploy

```bash
make helm-install
```

### 4. Wait for Pods

```bash
kubectl get pods -n goforge -w
# Wait until all pods are Running (2/2 READY)
```

### 5. Access Services

```bash
# Port forward
kubectl port-forward -n goforge svc/model-service 50051:50051 &
kubectl port-forward -n goforge svc/auth-service 50052:50052 &

# Test
./scripts/test-services.sh
```

**Done!** Services are running in Kubernetes.

---

## Option 3: Full Stack (with Monitoring)

### 1. Start Kubernetes

```bash
minikube start --cpus=4 --memory=8192
```

### 2. Install Istio

```bash
make istio-install
```

### 3. Build and Deploy

```bash
eval $(minikube docker-env)
make docker
make helm-install
```

### 4. Install Monitoring

```bash
make monitoring-install
```

### 5. Access Grafana

```bash
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80
```

Open browser: http://localhost:3000
- Username: `admin`
- Password: `admin123`

**Done!** Full platform with monitoring is running.

---

## Next Steps

### View Metrics
```bash
kubectl port-forward -n monitoring svc/prometheus-grafana 3000:80
# Open http://localhost:3000
```

### View Logs
```bash
kubectl logs -f deployment/model-service -n goforge
```

### Scale Services
```bash
kubectl scale deployment model-service -n goforge --replicas=5
```

### Check Health
```bash
grpcurl -plaintext localhost:50051 goforge.model.ModelService/HealthCheck
```

---

## Troubleshooting

### Pods not starting?
```bash
kubectl describe pod <pod-name> -n goforge
kubectl logs <pod-name> -n goforge
```

### Can't connect to services?
```bash
kubectl get svc -n goforge
kubectl get endpoints -n goforge
```

### Minikube issues?
```bash
minikube delete
minikube start --cpus=4 --memory=8192
```

---

## Clean Up

### Stop Local Services
```bash
# Press Ctrl+C in each terminal
```

### Delete Kubernetes Deployment
```bash
make helm-uninstall
# or
./scripts/cleanup.sh goforge
```

### Stop Minikube
```bash
minikube stop
# or delete completely
minikube delete
```

---

## Architecture at a Glance

```
User → NGINX Ingress → Istio Gateway → Services
                                         ├─ API Gateway
                                         ├─ Model Service
                                         ├─ Auth Service
                                         └─ Metrics Service
                                              │
                                              ↓
                                         Prometheus → Grafana
```

---

## Key Endpoints

| Service | Port | Endpoint |
|---------|------|----------|
| API Gateway | 50050 | localhost:50050 |
| Model Service | 50051 | localhost:50051 |
| Auth Service | 50052 | localhost:50052 |
| Metrics Service | 50053 | localhost:50053 |
| Prometheus | 9090 | localhost:9090 |
| Grafana | 3000 | localhost:3000 |

---

## Useful Commands

```bash
# Generate protobuf
make proto

# Build all services
make build

# Run tests
make test

# Build Docker images
make docker

# Deploy to Kubernetes
make k8s-deploy

# Deploy with Helm
make helm-install

# Install monitoring
make monitoring-install

# Clean up
make clean
```

---

## Learn More

- **Full Documentation**: See [README.md](README.md)
- **Deployment Guide**: See [EXECUTION_GUIDE.md](EXECUTION_GUIDE.md)
- **Project Details**: See [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md)

---

**Questions?** Open an issue on GitHub or join our Slack channel.

**Happy coding!** 🚀
