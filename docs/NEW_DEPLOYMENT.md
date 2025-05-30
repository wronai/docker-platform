# 🚀 Media Vault Deployment Guide

This guide provides comprehensive instructions for deploying the Media Vault platform across different environments, from local development to production.

## 📋 Table of Contents

### Getting Started
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)

### Deployment Types
- [Local Development](#-local-development)
- [Production Deployment](#-production-deployment)
- [Kubernetes](#-kubernetes-deployment)

### Configuration
- [Environment Variables](#-environment-variables)
- [Service Ports](#-service-ports)
- [Storage Configuration](#-storage-configuration)
- [Networking](#-networking)

### Operations
- [Scaling](#-scaling)
- [Monitoring](#-monitoring)
- [Backup & Recovery](#-backup--recovery)
- [Upgrading](#-upgrading)
- [Troubleshooting](#-troubleshooting)

## 📋 Prerequisites

### System Requirements

#### Development Environment
- **Minimum**:
  - 2 vCPUs
  - 4GB RAM
  - 20GB available storage
  - Docker 20.10+
  - Docker Compose 2.0+
  - Git
  - Make (recommended)

#### Production Environment
- **Minimum**:
  - 4 vCPUs
  - 8GB RAM
  - 100GB+ SSD storage
  - Linux-based OS (Ubuntu 20.04 LTS recommended)
  - Docker 20.10+ and Docker Compose 2.0+
  - Domain name with valid SSL certificate

#### Kubernetes (Optional)
- Kubernetes 1.20+
- Helm 3.0+
- Ingress controller (Nginx, Traefik, etc.)
- Persistent Volume provisioner

## 🚀 Quick Start

### Local Development Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/wronai/docker-platform.git
   cd docker-platform
   ```

2. **Set up environment**:
   ```bash
   # Copy and edit the environment file
   cp .env.example .env
   
   # Review and update configuration as needed
   # - Set ADMIN_EMAIL and ADMIN_PASSWORD
   # - Configure database credentials
   # - Update hostnames if needed
   ```

3. **Build and start services**:
   ```bash
   # Pull latest images and build containers
   make build
   
   # Start all services in detached mode
   make up
   
   # View logs (optional)
   make logs
   ```

4. **Verify the installation**:
   ```bash
   # Check running containers
   docker-compose ps
   
   # Check service health
   curl http://localhost:8080/health
   ```

5. **Access the application**:
   - Web UI: https://localhost
   - Admin Panel: https://admin.localhost
   - API Documentation: http://localhost:8080/api/docs
   - Keycloak Admin: https://auth.localhost

## 🔧 Configuration

### Key Configuration Files

| File | Description |
|------|-------------|
| `.env` | Main environment configuration |
| `docker-compose.yml` | Core services definition |
| `docker-compose.monitoring.yml` | Monitoring stack |
| `caddy/Caddyfile` | Reverse proxy configuration |
| `keycloak/import/realm-export.json` | Keycloak realm settings |

### Environment Variables

Key environment variables to configure:

```bash
# Application
APP_ENV=production
APP_SECRET=your-secret-key
APP_URL=https://your-domain.com

# Database
POSTGRES_USER=mediavault
POSTGRES_PASSWORD=secure-password
POSTGRES_DB=mediavault

# Keycloak
KEYCLOAK_ADMIN=admin
KEYCLOAK_ADMIN_PASSWORD=change-me
KEYCLOAK_FRONTEND_URL=https://auth.your-domain.com

# Storage
STORAGE_PATH=/data/media
BACKUP_PATH=/data/backups
```

## 🔌 Service Ports

| Service | Port | Protocol | Description |
|---------|------|----------|-------------|
| Web UI | 80/443 | HTTP/HTTPS | Main application interface |
| Admin | 8081 | HTTPS | Admin dashboard |
| Keycloak | 8080 | HTTPS | Authentication service |
| Grafana | 3000 | HTTP | Monitoring dashboards |
| Prometheus | 9090 | HTTP | Metrics collection |
| Alertmanager | 9093 | HTTP | Alert management |
| cAdvisor | 8081 | HTTP | Container metrics |
| Node Exporter | 9100 | HTTP | Host metrics |

## ⚖️ Scaling

### Horizontal Scaling

To scale backend services:
```bash
docker-compose up -d --scale backend=3
```

### Database Scaling
For production deployments, consider:
- Setting up PostgreSQL replication
- Using a managed database service
- Implementing connection pooling with PgBouncer

## 💾 Backup & Recovery

### Automated Backups

1. **Database Backups**:
   ```bash
   # Create a database backup
   make backup-db
   
   # Restore from backup
   make restore-db BACKUP_FILE=backup_20230530.sql
   ```

2. **Media Files**:
   ```bash
   # Backup media files
   make backup-media
   
   # Restore media
   make restore-media BACKUP_FILE=media_backup_20230530.tar.gz
   ```

### Disaster Recovery

1. **Full System Backup**:
   ```bash
   make full-backup
   ```

2. **Restore from Backup**:
   ```bash
   make restore-backup BACKUP_FILE=full_backup_20230530.tar.gz
   ```

## 🐛 Troubleshooting

### Common Issues

1. **Port Conflicts**
   - Check for running services: `sudo lsof -i -P -n | grep LISTEN`
   - Update `docker-compose.yml` with available ports

2. **Permission Issues**
   ```bash
   sudo chown -R $USER:$USER .
   chmod -R 755 data/
   ```

3. **Service Not Starting**
   ```bash
   # Check logs
   make logs service=backend
   
   # Restart services
   make restart
   ```

### Getting Help

- Check the [FAQ](docs/FAQ.md)
- Search [GitHub Issues](https://github.com/wronai/docker-platform/issues)
- Join our [community forum](#) (coming soon)

## 📝 Next Steps

- [Set up monitoring](MONITORING.md)
- [Configure authentication](SECURITY.md)
- [API Documentation](API.md)
