#!/bin/bash

set -e

echo "🚀 Building VPNaaS system..."

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

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Build backend
print_status "Building backend..."
cd backend
if [ ! -f go.mod ]; then
    print_error "go.mod not found in backend directory"
    exit 1
fi

# Download Go dependencies
print_status "Downloading Go dependencies..."
go mod download

# Build backend Docker image
print_status "Building backend Docker image..."
docker build -t vpnaas-backend:latest .

cd ..

# Build frontend
print_status "Building frontend..."
cd frontend
if [ ! -f package.json ]; then
    print_error "package.json not found in frontend directory"
    exit 1
fi

# Install Node.js dependencies
print_status "Installing Node.js dependencies..."
npm install

# Build frontend Docker image
print_status "Building frontend Docker image..."
docker build -t vpnaas-frontend:latest .

cd ..

print_status "✅ Build completed successfully!"
print_status "Images created:"
echo "  - vpnaas-backend:latest"
echo "  - vpnaas-frontend:latest"

print_warning "Next steps:"
echo "  1. Update the VPN endpoint in k8s/configmap.yaml"
echo "  2. Deploy to Kubernetes: kubectl apply -f k8s/"
echo "  3. Access the UI: kubectl port-forward svc/vpnaas-frontend 3000:80"
