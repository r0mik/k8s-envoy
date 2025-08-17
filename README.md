# VPNaaS - VPN as a Service for Kubernetes

A complete VPN-as-a-Service solution that provides per-user VPN pods in Kubernetes with a modern web UI for management.

## Architecture

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

## Features

- **Per-user VPN pods**: Each user gets their own isolated VPN pod
- **WireGuard VPN**: Fast, secure, and modern VPN protocol
- **Web UI**: React-based dashboard for user management
- **API Gateway**: Envoy-based access point with authentication
- **Metrics**: Prometheus metrics and Grafana dashboards
- **Kubernetes Native**: Full K8s integration with CRDs

## Components

### 1. Backend API (Go)
- User management (CRUD operations)
- VPN pod lifecycle management
- Kubernetes pod orchestration
- Metrics collection

### 2. Frontend UI (React)
- User dashboard
- Add/remove users
- VPN configuration download
- Statistics and metrics visualization

### 3. API Gateway (Envoy)
- Authentication and authorization
- Rate limiting
- Load balancing
- SSL termination

### 4. VPN Pods (WireGuard)
- Per-user isolated pods
- Automatic configuration generation
- Health monitoring

## Quick Start

1. **Deploy to Kubernetes**:
   ```bash
   kubectl apply -f k8s/
   ```

2. **Access the UI**:
   ```bash
   kubectl port-forward svc/vpnaas-ui 3000:80
   ```

3. **Add a user**:
   - Open http://localhost:3000
   - Click "Add User"
   - Download VPN configuration

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker
- Kubernetes cluster
- kubectl

### Local Development
```bash
# Backend
cd backend && go run main.go

# Frontend
cd frontend && npm start

# Deploy to K8s
kubectl apply -f k8s/
```

## Configuration

See `config/` directory for configuration files and examples.

## License

MIT
