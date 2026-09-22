# GoForge - Project Structure

```
GoForge/
│
├── README.md                          # Main project documentation
├── QUICKSTART.md                      # 10-minute quick start guide
├── EXECUTION_GUIDE.md                 # Comprehensive deployment guide
├── PROJECT_SUMMARY.md                 # Technical summary and achievements
├── CONTRIBUTING.md                    # Contribution guidelines
├── LICENSE                            # MIT License
├── .gitignore                         # Git ignore patterns
├── .dockerignore                      # Docker ignore patterns
├── go.mod                             # Go module dependencies
├── go.sum                             # Go dependency checksums
├── Makefile                           # Build and deployment automation
│
├── api/                               # API definitions
│   └── proto/                         # Protocol Buffer definitions
│       ├── model.proto                # Model service API (Predict, BatchPredict, GetModelInfo)
│       ├── auth.proto                 # Auth service API (Authenticate, ValidateToken)
│       ├── gateway.proto              # Gateway service API (RouteRequest, GetServiceHealth)
│       └── metrics.proto              # Metrics service API (RecordMetric, GetMetrics)
│
├── cmd/                               # Service entrypoints (main packages)
│   ├── api-gateway/
│   │   └── main.go                    # API Gateway server (port 50050)
│   ├── model-service/
│   │   └── main.go                    # Model Service server (port 50051)
│   ├── auth-service/
│   │   └── main.go                    # Auth Service server (port 50052)
│   └── metrics-service/
│       └── main.go                    # Metrics Service server (port 50053)
│
├── internal/                          # Private application code
│   ├── gateway/
│   │   └── server.go                  # Gateway service implementation
│   ├── model/
│   │   └── server.go                  # Model service implementation
│   ├── auth/
│   │   └── server.go                  # Auth service implementation
│   └── metrics/
│       └── server.go                  # Metrics service implementation
│
├── pkg/                               # Shared libraries (importable)
│   ├── config/
│   │   └── config.go                  # Configuration management (env vars, defaults)
│   ├── logger/
│   │   └── logger.go                  # Structured logging with Zap
│   ├── middleware/
│   │   ├── auth.go                    # Authentication interceptor
│   │   └── metrics.go                 # Metrics collection interceptor
│   ├── circuitbreaker/
│   │   └── circuitbreaker.go          # Circuit breaker implementation (gobreaker)
│   ├── ratelimit/
│   │   └── ratelimit.go               # Rate limiting implementation (token bucket)
│   ├── discovery/
│   │   └── discovery.go               # Kubernetes service discovery
│   └── proto/                         # Generated protobuf code (auto-generated)
│       ├── *.pb.go                    # Protobuf message definitions
│       └── *_grpc.pb.go               # gRPC service stubs
│
├── deployments/                       # Deployment configurations
│   ├── docker/                        # Dockerfiles
│   │   ├── Dockerfile.api-gateway     # Multi-stage build for API Gateway
│   │   ├── Dockerfile.model-service   # Multi-stage build for Model Service
│   │   ├── Dockerfile.auth-service    # Multi-stage build for Auth Service
│   │   └── Dockerfile.metrics-service # Multi-stage build for Metrics Service
│   │
│   ├── k8s/                           # Kubernetes manifests
│   │   ├── namespace.yaml             # Namespace with Istio injection enabled
│   │   │
│   │   ├── configmaps/                # Configuration
│   │   │   └── app-config.yaml        # Environment variables, settings
│   │   │
│   │   ├── secrets/                   # Sensitive data (not committed)
│   │   │
│   │   ├── deployments/               # Service deployments
│   │   │   ├── api-gateway.yaml       # Deployment + Service + HPA
│   │   │   ├── model-service.yaml     # Deployment + Service + HPA
│   │   │   ├── auth-service.yaml      # Deployment + Service + HPA
│   │   │   └── metrics-service.yaml   # Deployment + Service
│   │   │
│   │   ├── services/                  # (Included in deployment files)
│   │   │
│   │   ├── ingress/                   # Ingress configurations
│   │   │   └── ingress.yaml           # NGINX Ingress with rate limiting
│   │   │
│   │   └── istio/                     # Istio service mesh configurations
│   │       ├── gateway.yaml           # Istio Gateway + VirtualService
│   │       ├── destination-rules.yaml # Load balancing, circuit breaking
│   │       └── peer-authentication.yaml # mTLS enforcement
│   │
│   ├── helm/                          # Helm charts
│   │   └── goforge/                   # GoForge Helm chart
│   │       ├── Chart.yaml             # Chart metadata
│   │       ├── values.yaml            # Default configuration values
│   │       └── templates/             # Kubernetes manifest templates
│   │           └── deployment.yaml    # Templated deployments
│   │
│   └── monitoring/                    # Observability stack
│       ├── prometheus-values.yaml     # Prometheus configuration
│       ├── prometheus-rules.yaml      # Alert rules
│       ├── servicemonitor.yaml        # Prometheus ServiceMonitor CRDs
│       └── grafana-dashboard.json     # Pre-built Grafana dashboards
│
├── scripts/                           # Utility scripts
│   ├── build.sh                       # Build all services
│   ├── deploy.sh                      # Deploy to Kubernetes
│   ├── test-services.sh               # Test gRPC endpoints with grpcurl
│   └── cleanup.sh                     # Clean up deployments
│
└── bin/                               # Compiled binaries (generated, not committed)
    ├── api-gateway
    ├── model-service
    ├── auth-service
    └── metrics-service
```

## File Count Summary

| Category | Count | Description |
|----------|-------|-------------|
| Go Source Files | 17 | Service implementations and libraries |
| Protobuf Definitions | 4 | gRPC service definitions |
| Dockerfiles | 4 | Container images |
| Kubernetes Manifests | 10+ | Deployments, services, config |
| Helm Templates | 3 | Chart, values, templates |
| Monitoring Configs | 4 | Prometheus, Grafana, alerts |
| Scripts | 4 | Automation scripts |
| Documentation | 6 | README, guides, contributing |
| **Total** | **50+** | Production-ready files |

## Key Directories Explained

### `/api/proto`
Protocol Buffer definitions that define the gRPC service contracts. These are language-agnostic and generate code for Go (and potentially other languages).

### `/cmd`
Entry points for each microservice. Each subdirectory contains a `main.go` that:
- Initializes configuration
- Sets up logging
- Creates gRPC server with interceptors
- Starts HTTP metrics server
- Handles graceful shutdown

### `/internal`
Business logic implementations. These are private to this module and can't be imported by external packages. Contains the actual gRPC service handlers.

### `/pkg`
Reusable libraries that could be imported by other Go modules. Includes:
- Configuration management
- Logging utilities
- Middleware (auth, metrics, tracing)
- Circuit breaker wrapper
- Rate limiter wrapper
- Service discovery client

### `/deployments`
Everything needed to deploy the platform:
- **docker/**: Dockerfiles using multi-stage builds
- **k8s/**: Raw Kubernetes manifests
- **helm/**: Templated Kubernetes deployments
- **monitoring/**: Prometheus and Grafana configurations

### `/scripts`
Shell scripts for common operations:
- Building services
- Deploying to Kubernetes
- Testing endpoints
- Cleaning up resources

## Generated Files (Not Committed)

These files are generated during build and not tracked in Git:

```
├── bin/                    # Compiled Go binaries
│   ├── api-gateway
│   ├── model-service
│   ├── auth-service
│   └── metrics-service
│
├── pkg/proto/              # Generated protobuf code
│   ├── model.pb.go
│   ├── model_grpc.pb.go
│   ├── auth.pb.go
│   ├── auth_grpc.pb.go
│   ├── gateway.pb.go
│   ├── gateway_grpc.pb.go
│   ├── metrics.pb.go
│   └── metrics_grpc.pb.go
│
└── vendor/                 # Go dependencies (if using vendoring)
```

## Configuration Flow

```
Environment Variables
         ↓
   pkg/config/config.go
         ↓
   ConfigMaps (K8s)
         ↓
   Service Initialization
         ↓
   Runtime Configuration
```

## Build Flow

```
1. make proto          → Generate protobuf code
2. make build          → Compile Go binaries
3. make docker         → Build container images
4. make helm-install   → Deploy to Kubernetes
5. make monitoring-install → Set up observability
```

## Request Flow

```
Client Request
      ↓
NGINX Ingress (TLS termination, rate limiting)
      ↓
Istio Gateway (routing, retries)
      ↓
Envoy Sidecar (mTLS, circuit breaking, metrics)
      ↓
API Gateway (authentication, request routing)
      ↓
Backend Service (Model/Auth/Metrics)
      ↓
Response (with metrics, traces, logs)
```

## Dependency Graph

```
cmd/[service]/main.go
    ├── internal/[service]/server.go  (business logic)
    ├── pkg/config/config.go          (configuration)
    ├── pkg/logger/logger.go          (logging)
    ├── pkg/middleware/*              (interceptors)
    ├── pkg/circuitbreaker/*          (resilience)
    ├── pkg/ratelimit/*               (throttling)
    └── pkg/proto/*.pb.go             (generated gRPC code)
```

## Port Allocation

| Service | gRPC Port | Metrics Port |
|---------|-----------|--------------|
| API Gateway | 50050 | 9090 |
| Model Service | 50051 | 9090 |
| Auth Service | 50052 | 9090 |
| Metrics Service | 50053 | 9090 |
| Prometheus | - | 9090 |
| Grafana | - | 3000 |

## Technology Stack by Directory

| Directory | Technologies |
|-----------|--------------|
| `/api` | Protocol Buffers, gRPC |
| `/cmd` | Go 1.22, gRPC-Go |
| `/internal` | Go, gRPC server implementations |
| `/pkg` | Go, Zap, gobreaker, rate limiter |
| `/deployments/docker` | Docker, multi-stage builds |
| `/deployments/k8s` | Kubernetes 1.29, Istio 1.20 |
| `/deployments/helm` | Helm 3 |
| `/deployments/monitoring` | Prometheus, Grafana, AlertManager |

---

**This structure follows Go best practices and cloud-native patterns for maximum maintainability and scalability.**
