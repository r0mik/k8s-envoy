# VPNaaS Setup Guide

This guide will help you set up and deploy the VPNaaS (VPN as a Service) system on Kubernetes.

## Prerequisites

### Required Software
- **Docker** (version 20.10+)
- **kubectl** (version 1.24+)
- **Kubernetes cluster** (minikube, kind, or cloud provider)
- **Go** (version 1.21+) - for local development
- **Node.js** (version 18+) - for local development

### System Requirements
- **CPU**: 2+ cores
- **Memory**: 4GB+ RAM
- **Storage**: 10GB+ free space
- **Network**: Internet access for pulling Docker images

## Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd vpnsaas
```

### 2. Build the System
```bash
./scripts/build.sh
```

### 3. Deploy to Kubernetes
```bash
./scripts/deploy.sh
```

### 4. Access the Application
```bash
# Access the web UI
kubectl port-forward svc/vpnaas-frontend 3000:80 -n vpnaas

# Access Grafana (admin/admin)
kubectl port-forward svc/grafana 3001:3000 -n vpnaas

# Access Prometheus
kubectl port-forward svc/prometheus 9090:9090 -n vpnaas
```

Open your browser and navigate to:
- **Web UI**: http://localhost:3000
- **Grafana**: http://localhost:3001
- **Prometheus**: http://localhost:9090

## Manual Setup

### 1. Configuration

#### Update VPN Endpoint
Edit `k8s/configmap.yaml` and update the VPN endpoint:
```yaml
vpn:
  endpoint: "your-vpn-endpoint.com"  # Replace with your actual endpoint
```

#### Environment Variables
You can customize the deployment by setting environment variables:
```bash
export VPNAAS_DEBUG=false
export VPNAAS_SERVER_PORT=8080
export VPNAAS_K8S_NAMESPACE=vpnaas
```

### 2. Build Components

#### Backend
```bash
cd backend
go mod download
docker build -t vpnaas-backend:latest .
cd ..
```

#### Frontend
```bash
cd frontend
npm install
docker build -t vpnaas-frontend:latest .
cd ..
```

### 3. Deploy to Kubernetes

#### Create Namespace
```bash
kubectl apply -f k8s/namespace.yaml
```

#### Deploy Components
```bash
# ConfigMap
kubectl apply -f k8s/configmap.yaml

# Backend
kubectl apply -f k8s/backend-deployment.yaml

# Frontend
kubectl apply -f k8s/frontend-deployment.yaml

# API Gateway
kubectl apply -f k8s/envoy-gateway.yaml

# Monitoring
kubectl apply -f k8s/prometheus.yaml
kubectl apply -f k8s/grafana.yaml
```

### 4. Verify Deployment
```bash
# Check pods
kubectl get pods -n vpnaas

# Check services
kubectl get services -n vpnaas

# Check logs
kubectl logs -f deployment/vpnaas-backend -n vpnaas
kubectl logs -f deployment/vpnaas-frontend -n vpnaas
```

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web UI        │    │   API Gateway   │    │   VPN Pods      │
│   (React)       │◄──►│   (Envoy)       │◄──►│   (WireGuard)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                       ┌─────────────────┐
                       │   Backend API   │
                       │   (Go)          │
                       └─────────────────┘
                                │
                       ┌─────────────────┐
                       │   Kubernetes    │
                       │   API           │
                       └─────────────────┘
```

## Components

### 1. Backend API (Go)
- **Purpose**: Manages VPN users and pod lifecycle
- **Features**:
  - User CRUD operations
  - VPN pod orchestration
  - WireGuard configuration generation
  - Prometheus metrics collection
- **Port**: 8080

### 2. Frontend UI (React)
- **Purpose**: Web-based user management interface
- **Features**:
  - User dashboard
  - Add/remove users
  - VPN configuration download
  - Statistics visualization
- **Port**: 80 (via Nginx)

### 3. API Gateway (Envoy)
- **Purpose**: Load balancing and routing
- **Features**:
  - Request routing
  - Load balancing
  - SSL termination
- **Port**: 80

### 4. VPN Pods (WireGuard)
- **Purpose**: Per-user VPN instances
- **Features**:
  - Isolated VPN per user
  - Automatic configuration
  - Health monitoring
- **Port**: 51820 (UDP)

### 5. Monitoring Stack
- **Prometheus**: Metrics collection
- **Grafana**: Metrics visualization
- **Custom Dashboards**: VPNaaS-specific metrics

## Usage

### 1. Adding a User
1. Open the web UI
2. Click "Add User"
3. Enter username and email
4. Click "Create User"
5. Download the VPN configuration

### 2. Managing Users
- **View Users**: Dashboard shows all users with status
- **Delete User**: Click the trash icon to remove a user
- **Download Config**: Click the download icon to get VPN config

### 3. Monitoring
- **Dashboard**: Overview of system metrics
- **Grafana**: Detailed metrics and custom dashboards
- **Prometheus**: Raw metrics data

## Troubleshooting

### Common Issues

#### 1. Pods Not Starting
```bash
# Check pod status
kubectl get pods -n vpnaas

# Check pod logs
kubectl logs <pod-name> -n vpnaas

# Check events
kubectl get events -n vpnaas --sort-by='.lastTimestamp'
```

#### 2. Backend API Issues
```bash
# Check backend logs
kubectl logs -f deployment/vpnaas-backend -n vpnaas

# Test API endpoint
kubectl port-forward svc/vpnaas-backend 8080:8080 -n vpnaas
curl http://localhost:8080/health
```

#### 3. Frontend Issues
```bash
# Check frontend logs
kubectl logs -f deployment/vpnaas-frontend -n vpnaas

# Check nginx configuration
kubectl exec -it <frontend-pod> -n vpnaas -- nginx -t
```

#### 4. VPN Pod Issues
```bash
# Check VPN pod status
kubectl get pods -l app=vpnaas,component=vpn -n vpnaas

# Check WireGuard logs
kubectl logs <vpn-pod-name> -n vpnaas
```

### Performance Tuning

#### Resource Limits
Adjust resource limits in the deployment files:
```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"
```

#### Scaling
```bash
# Scale backend
kubectl scale deployment vpnaas-backend --replicas=3 -n vpnaas

# Scale frontend
kubectl scale deployment vpnaas-frontend --replicas=3 -n vpnaas
```

## Security Considerations

### 1. Network Security
- Use Network Policies to restrict pod communication
- Implement proper RBAC for Kubernetes access
- Use TLS for API communication

### 2. Authentication
- Implement API authentication (JWT, OAuth)
- Use Kubernetes service accounts
- Enable audit logging

### 3. Data Protection
- Encrypt VPN configurations
- Use secrets for sensitive data
- Implement data retention policies

## Development

### Local Development
```bash
# Backend
cd backend
go run main.go

# Frontend
cd frontend
npm start
```

### Testing
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Building for Production
```bash
# Build optimized images
docker build --no-cache -t vpnaas-backend:latest backend/
docker build --no-cache -t vpnaas-frontend:latest frontend/
```

## Support

For issues and questions:
1. Check the troubleshooting section
2. Review the logs and metrics
3. Check the GitHub issues
4. Create a new issue with detailed information

## License

This project is licensed under the MIT License.
