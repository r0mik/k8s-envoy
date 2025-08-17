#!/bin/bash

set -e

echo "🚀 Deploying VPNaaS to Kubernetes..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl and try again."
    exit 1
fi

# Check if we can connect to Kubernetes
if ! kubectl cluster-info &> /dev/null; then
    print_error "Cannot connect to Kubernetes cluster. Please check your kubeconfig."
    exit 1
fi

# Check if images exist
if ! docker image inspect vpnaas-backend:latest &> /dev/null; then
    print_error "Backend image not found. Please run ./scripts/build.sh first."
    exit 1
fi

if ! docker image inspect vpnaas-frontend:latest &> /dev/null; then
    print_error "Frontend image not found. Please run ./scripts/build.sh first."
    exit 1
fi

# Create namespace
print_status "Creating namespace..."
kubectl apply -f k8s/namespace.yaml

# Apply ConfigMap
print_status "Applying ConfigMap..."
kubectl apply -f k8s/configmap.yaml

# Deploy backend
print_status "Deploying backend..."
kubectl apply -f k8s/backend-deployment.yaml

# Deploy frontend
print_status "Deploying frontend..."
kubectl apply -f k8s/frontend-deployment.yaml

# Deploy API Gateway (if Envoy Gateway is installed)
print_status "Deploying API Gateway..."
kubectl apply -f k8s/envoy-gateway.yaml

# Deploy monitoring stack
print_status "Deploying monitoring stack..."
kubectl apply -f k8s/prometheus.yaml
kubectl apply -f k8s/grafana.yaml

# Wait for deployments to be ready
print_status "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/vpnaas-backend -n vpnaas
kubectl wait --for=condition=available --timeout=300s deployment/vpnaas-frontend -n vpnaas

print_status "✅ Deployment completed successfully!"

# Show service information
print_status "Service endpoints:"
echo "  - Frontend: kubectl port-forward svc/vpnaas-frontend 3000:80 -n vpnaas"
echo "  - Backend API: kubectl port-forward svc/vpnaas-backend 8080:8080 -n vpnaas"
echo "  - Prometheus: kubectl port-forward svc/prometheus 9090:9090 -n vpnaas"
echo "  - Grafana: kubectl port-forward svc/grafana 3001:3000 -n vpnaas"

print_warning "Access the application:"
echo "  - Web UI: http://localhost:3000"
echo "  - Grafana: http://localhost:3001 (admin/admin)"
echo "  - Prometheus: http://localhost:9090"

print_status "To check deployment status:"
echo "  kubectl get pods -n vpnaas"
echo "  kubectl get services -n vpnaas"
