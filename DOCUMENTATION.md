# 📚 Media Vault - Comprehensive Documentation

## 📖 Table of Contents
1. [Overview](#-overview)
2. [Quick Start](#-quick-start)
3. [Architecture](#-architecture)
4. [API Reference](#-api-reference)
5. [Deployment](#-deployment)
6. [Development](#-development)
7. [Monitoring & Operations](#-monitoring--operations)
8. [Security](#-security)
9. [Contributing](#-contributing)
10. [Troubleshooting](#-troubleshooting)

## 🌟 Overview

Media Vault is a secure, scalable media management platform with AI-powered analysis and organization capabilities. It provides:

- Secure media storage and sharing
- AI-powered media analysis
- Role-based access control
- Comprehensive monitoring
- Scalable architecture

## 🚀 Quick Start

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ RAM recommended

### Local Development Setup

```bash
# Clone the repository
git clone https://github.com/wronai/docker-platform.git
cd docker-platform

# Copy example environment file
cp .env.example .env

# Start all services
make up

# Access the application
open http://localhost
```

## 🏗️ Architecture

For detailed architecture documentation, see [ARCHITECTURE.md](./docs/ARCHITECTURE.md).

### Core Components

1. **Frontend**: Flutter-based web interface
2. **Backend API**: Go-based REST API
3. **Authentication**: Keycloak
4. **Database**: PostgreSQL
5. **Object Storage**: MinIO
6. **AI Services**: Media analysis and processing
7. **Monitoring**: Prometheus, Grafana, Loki

## 📡 API Reference

For complete API documentation, see [API.md](./docs/API.md).

### Authentication

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "yourpassword"
}
```

## 🚀 Deployment

For production deployment instructions, see [DEPLOYMENT.md](./docs/DEPLOYMENT.md).

### Environment Variables

Key environment variables:

```env
# Database
POSTGRES_DB=mediavault
POSTGRES_USER=mediavault
POSTGRES_PASSWORD=changeme

# Keycloak
KEYCLOAK_ADMIN=admin
KEYCLOAK_ADMIN_PASSWORD=changeme

# MinIO
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
```

## 💻 Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- Flutter 3.0+

### Building

```bash
# Build all services
make build

# Run tests
make test

# Lint code
make lint
```

## 📊 Monitoring & Operations

For monitoring setup and operations, see [MONITORING.md](./docs/MONITORING.md).

### Accessing Monitoring Tools

- **Grafana**: http://localhost:3000
- **Prometheus**: http://localhost:9090
- **Loki**: http://localhost:3100

## 🔒 Security

For security best practices and guidelines, see [SECURITY.md](./docs/SECURITY.md).

### Reporting Security Issues

Please report security issues to security@wron.ai.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a pull request

## 🛠️ Troubleshooting

### Common Issues

1. **Port conflicts**
   - Check if ports 80, 443, 5432, 9000, 9090, 3000 are available

2. **Database connection issues**
   - Verify PostgreSQL is running
   - Check database credentials in .env

3. **Authentication problems**
   - Ensure Keycloak is running
   - Verify user credentials in Keycloak

### Getting Help

- [GitHub Issues](https://github.com/wronai/docker-platform/issues)
- support@wron.ai
