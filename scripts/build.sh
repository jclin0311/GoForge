#!/bin/bash
set -e

echo "Building GoForge Services..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create bin directory
mkdir -p bin

# Services to build
SERVICES=("api-gateway" "model-service" "auth-service" "metrics-service")

for service in "${SERVICES[@]}"; do
    echo -e "${BLUE}Building $service...${NC}"
    cd cmd/$service
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../bin/$service .
    cd ../..
    echo -e "${GREEN}✓ $service built successfully${NC}"
done

echo -e "${GREEN}All services built successfully!${NC}"
