#!/bin/bash
set -e

# Test script for GoForge gRPC services

HOST=${1:-localhost}
API_GATEWAY_PORT=${2:-50050}
MODEL_SERVICE_PORT=${3:-50051}
AUTH_SERVICE_PORT=${4:-50052}

echo "Testing GoForge Services on $HOST"
echo "=================================="

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo "Error: grpcurl is not installed"
    echo "Install with: brew install grpcurl (macOS) or go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Test Auth Service
echo ""
echo "Testing Auth Service on port $AUTH_SERVICE_PORT..."
echo "---------------------------------------------------"
AUTH_RESPONSE=$(grpcurl -plaintext -d '{
  "username": "testuser",
  "password": "testpassword",
  "client_id": "test-client"
}' $HOST:$AUTH_SERVICE_PORT goforge.auth.AuthService/Authenticate)

echo "$AUTH_RESPONSE"
TOKEN=$(echo "$AUTH_RESPONSE" | grep -o '"access_token": "[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "✗ Auth Service test failed"
else
    echo "✓ Auth Service test passed (token: ${TOKEN:0:20}...)"
fi

# Test Model Service - Health Check
echo ""
echo "Testing Model Service Health Check on port $MODEL_SERVICE_PORT..."
echo "------------------------------------------------------------------"
HEALTH_RESPONSE=$(grpcurl -plaintext $HOST:$MODEL_SERVICE_PORT goforge.model.ModelService/HealthCheck)
echo "$HEALTH_RESPONSE"

if echo "$HEALTH_RESPONSE" | grep -q '"healthy": true'; then
    echo "✓ Model Service health check passed"
else
    echo "✗ Model Service health check failed"
fi

# Test Model Service - Predict
echo ""
echo "Testing Model Service Predict on port $MODEL_SERVICE_PORT..."
echo "-------------------------------------------------------------"
PREDICT_RESPONSE=$(grpcurl -plaintext -d '{
  "model_id": "gpt-neo-125m",
  "model_version": "1.0.0",
  "input_data": "SGVsbG8gV29ybGQh",
  "request_id": "test-req-001",
  "parameters": {
    "temperature": "0.7",
    "max_tokens": "100"
  }
}' $HOST:$MODEL_SERVICE_PORT goforge.model.ModelService/Predict)

echo "$PREDICT_RESPONSE"

if echo "$PREDICT_RESPONSE" | grep -q '"request_id": "test-req-001"'; then
    echo "✓ Model Service predict test passed"
else
    echo "✗ Model Service predict test failed"
fi

# Test Model Service - Get Model Info
echo ""
echo "Testing Model Service GetModelInfo on port $MODEL_SERVICE_PORT..."
echo "------------------------------------------------------------------"
INFO_RESPONSE=$(grpcurl -plaintext -d '{
  "model_id": "gpt-neo-125m",
  "model_version": "1.0.0"
}' $HOST:$MODEL_SERVICE_PORT goforge.model.ModelService/GetModelInfo)

echo "$INFO_RESPONSE"

if echo "$INFO_RESPONSE" | grep -q '"model_id": "gpt-neo-125m"'; then
    echo "✓ Model Service info test passed"
else
    echo "✗ Model Service info test failed"
fi

echo ""
echo "=================================="
echo "Service Testing Complete!"
