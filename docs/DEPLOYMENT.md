# 🚀 Deployment Guide

This guide provides detailed instructions for deploying the Media Vault platform in various environments.

## Table of Contents
- [Prerequisites](#-prerequisites)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [Environment Variables](#-environment-variables)
- [Service Ports](#-service-ports)
- [Scaling](#-scaling)
- [Backup & Recovery](#-backup--recovery)
- [Troubleshooting](#-troubleshooting)

## 📋 Prerequisites

### System Requirements
- **Minimum**:
  - 2 vCPUs
  - 4GB RAM
  - 20GB Storage
- **Recommended**:
  - 4+ vCPUs
  - 8GB+ RAM
  - 100GB+ Storage (SSD recommended)

### Required Software
- Docker 20.10+
- Docker Compose 2.0+
- Git
- Make (optional but recommended)

## 🚀 Quick Start

### Local Development Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/wronai/docker-platform.git
   cd docker-platform
   ```

2. **Set up environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the stack**:
   ```bash
   make up
   ```
   
   Or for development mode:
   ```bash
   make dev
   ```

4. **Verify the installation**:
   ```bash
   make status
   ```

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