package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Service      ServiceConfig
	GRPC         GRPCConfig
	Kubernetes   K8sConfig
	CircuitBreaker CircuitBreakerConfig
	RateLimit    RateLimitConfig
	Observability ObservabilityConfig
}

type ServiceConfig struct {
	Name        string
	Version     string
	Environment string
	Port        int
	MetricsPort int
}

type GRPCConfig struct {
	Port                int
	MaxConnectionIdle   time.Duration
	MaxConnectionAge    time.Duration
	MaxConnectionAgeGrace time.Duration
	Time                time.Duration
	Timeout             time.Duration
}

type K8sConfig struct {
	Namespace       string
	ServiceDiscovery bool
	InCluster       bool
}

type CircuitBreakerConfig struct {
	MaxRequests     uint32
	Interval        time.Duration
	Timeout         time.Duration
	ReadyToTrip     func(counts uint64) bool
}

type RateLimitConfig struct {
	RequestsPerSecond float64
	BurstSize         int
}

type ObservabilityConfig struct {
	PrometheusEnabled bool
	PrometheusPort    int
	TracingEnabled    bool
	TracingEndpoint   string
	LogLevel          string
}

func LoadConfig() *Config {
	return &Config{
		Service: ServiceConfig{
			Name:        getEnv("SERVICE_NAME", "goforge-service"),
			Version:     getEnv("SERVICE_VERSION", "v1.0.0"),
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnvAsInt("SERVICE_PORT", 8080),
			MetricsPort: getEnvAsInt("METRICS_PORT", 9090),
		},
		GRPC: GRPCConfig{
			Port:                  getEnvAsInt("GRPC_PORT", 50051),
			MaxConnectionIdle:     getEnvAsDuration("GRPC_MAX_CONN_IDLE", 5*time.Minute),
			MaxConnectionAge:      getEnvAsDuration("GRPC_MAX_CONN_AGE", 10*time.Minute),
			MaxConnectionAgeGrace: getEnvAsDuration("GRPC_MAX_CONN_AGE_GRACE", 5*time.Minute),
			Time:                  getEnvAsDuration("GRPC_KEEPALIVE_TIME", 2*time.Minute),
			Timeout:               getEnvAsDuration("GRPC_KEEPALIVE_TIMEOUT", 20*time.Second),
		},
		Kubernetes: K8sConfig{
			Namespace:        getEnv("K8S_NAMESPACE", "goforge"),
			ServiceDiscovery: getEnvAsBool("K8S_SERVICE_DISCOVERY", true),
			InCluster:        getEnvAsBool("K8S_IN_CLUSTER", false),
		},
		CircuitBreaker: CircuitBreakerConfig{
			MaxRequests: uint32(getEnvAsInt("CB_MAX_REQUESTS", 10)),
			Interval:    getEnvAsDuration("CB_INTERVAL", 60*time.Second),
			Timeout:     getEnvAsDuration("CB_TIMEOUT", 30*time.Second),
		},
		RateLimit: RateLimitConfig{
			RequestsPerSecond: float64(getEnvAsInt("RATE_LIMIT_RPS", 100)),
			BurstSize:         getEnvAsInt("RATE_LIMIT_BURST", 200),
		},
		Observability: ObservabilityConfig{
			PrometheusEnabled: getEnvAsBool("PROMETHEUS_ENABLED", true),
			PrometheusPort:    getEnvAsInt("PROMETHEUS_PORT", 9090),
			TracingEnabled:    getEnvAsBool("TRACING_ENABLED", true),
			TracingEndpoint:   getEnv("TRACING_ENDPOINT", "jaeger:14268"),
			LogLevel:          getEnv("LOG_LEVEL", "info"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
