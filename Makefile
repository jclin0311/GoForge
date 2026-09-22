.PHONY: all build test clean proto docker k8s-deploy k8s-delete helm-install helm-uninstall

# Variables
SERVICES := api-gateway model-service auth-service metrics-service
DOCKER_REGISTRY := goforge
VERSION := v1.0.0
NAMESPACE := goforge

all: proto build

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@mkdir -p pkg/proto
	@for file in api/proto/*.proto; do \
		protoc --go_out=pkg/proto --go_opt=paths=source_relative \
		       --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
		       -I api/proto $$file; \
	done

# Build all services
build:
	@echo "Building services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd cmd/$$service && CGO_ENABLED=0 GOOS=linux go build -o ../../bin/$$service . && cd ../..; \
	done

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Build Docker images
docker:
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) -f deployments/docker/Dockerfile.$$service .; \
	done

# Deploy to Kubernetes
k8s-deploy:
	@echo "Deploying to Kubernetes..."
	@kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	@kubectl apply -f deployments/k8s/namespace.yaml
	@kubectl apply -f deployments/k8s/configmaps/
	@kubectl apply -f deployments/k8s/secrets/
	@kubectl apply -f deployments/k8s/services/
	@kubectl apply -f deployments/k8s/deployments/
	@kubectl apply -f deployments/k8s/ingress/

# Delete from Kubernetes
k8s-delete:
	@echo "Deleting from Kubernetes..."
	@kubectl delete -f deployments/k8s/ingress/ || true
	@kubectl delete -f deployments/k8s/deployments/ || true
	@kubectl delete -f deployments/k8s/services/ || true
	@kubectl delete namespace $(NAMESPACE) || true

# Install with Helm
helm-install:
	@echo "Installing with Helm..."
	@helm upgrade --install goforge deployments/helm/goforge \
		--namespace $(NAMESPACE) \
		--create-namespace \
		--values deployments/helm/goforge/values.yaml

# Uninstall Helm release
helm-uninstall:
	@echo "Uninstalling Helm release..."
	@helm uninstall goforge --namespace $(NAMESPACE)

# Install Istio
istio-install:
	@echo "Installing Istio..."
	@istioctl install --set profile=demo -y
	@kubectl label namespace $(NAMESPACE) istio-injection=enabled --overwrite

# Install monitoring stack
monitoring-install:
	@echo "Installing Prometheus and Grafana..."
	@helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
	@helm repo add grafana https://grafana.github.io/helm-charts
	@helm repo update
	@helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
		--namespace monitoring \
		--create-namespace \
		--values deployments/monitoring/prometheus-values.yaml

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf pkg/proto/*.pb.go
	@rm -f coverage.out

# Run locally (development)
run-api-gateway:
	@go run cmd/api-gateway/main.go

run-model-service:
	@go run cmd/model-service/main.go

run-auth-service:
	@go run cmd/auth-service/main.go

run-metrics-service:
	@go run cmd/metrics-service/main.go
