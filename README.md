# 🔐 Media Vault - Complete Solution

**Enterprise-grade secure media storage with AI analysis, role-based access, and comprehensive monitoring.**

[![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)](https://github.com/wronai/docker-platform)
[![License](https://img.shields.io/badge/license-Apache%202.0-green.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/docker-ready-brightgreen.svg)](docker-compose.yml)
[![Documentation](https://img.shields.io/badge/docs-📘-blueviolet)](docs/README.md)
[![Project Status](https://img.shields.io/badge/status-active%20development-yellowgreen)](#project-status)

## 📋 Table of Contents
- [🚀 Quick Start](#-quick-start)
- [✨ Key Features](#-key-features)
- [🏗️ Architecture Overview](#-architecture-overview)
- [📚 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## 🚀 Quick Start

1. **Prerequisites**:
   - Docker 20.10+ and Docker Compose
   - 4GB RAM minimum (8GB recommended)
   - Ports 80, 443, 8080, 3000 available

2. **Deploy the stack**:
   ```bash
   git clone https://github.com/wronai/docker-platform.git
   cd docker-platform
   cp .env.example .env
   make up
   ```

3. **Access the applications**:
   - Web UI: https://localhost
   - Admin Panel: https://admin.localhost
   - Keycloak: https://auth.localhost
   - Grafana: http://localhost:3000

For detailed setup and configuration, see the [Deployment Guide](docs/DEPLOYMENT.md).

## ✨ Key Features

- **Secure File Storage**: End-to-end encryption for all media
- **Role-Based Access**: Fine-grained permission controls
- **AI Analysis**: Automatic content moderation and tagging
- **High Availability**: Designed for 99.9% uptime
- **Comprehensive Monitoring**: Built-in observability stack

## 🏗️ Architecture Overview

Media Vault is built on a modern microservices architecture:

- **Frontend**: Flutter-based responsive web interface
- **Backend**: High-performance Go services
- **Authentication**: Keycloak for identity management
- **Database**: PostgreSQL for data persistence
- **Monitoring**: Prometheus, Grafana, and more

For a complete architecture deep dive, see the [Architecture Documentation](docs/ARCHITECTURE.md).

## 📚 Documentation

Comprehensive documentation is available in the `docs/` directory:

- [📘 User Guide](docs/USER_GUIDE.md) - End-user documentation
- [🔧 Deployment Guide](docs/DEPLOYMENT.md) - Setup and configuration
- [🏗️ Architecture](docs/ARCHITECTURE.md) - System design and components
- [🔐 Security](docs/SECURITY.md) - Security best practices
- [📊 Monitoring](docs/MONITORING.md) - Observability and alerting
- [📝 API Reference](docs/API.md) - API documentation

## 🤝 Contributing

We welcome contributions from the community! Here's how you can help:

1. **Report Bugs**: File an issue on our [issue tracker](https://github.com/wronai/docker-platform/issues).
2. **Submit Fixes**: Fork the repository and submit a pull request.
3. **Improve Docs**: Help us enhance our documentation.

Please read our [Contributing Guide](CONTRIBUTING.md) for development setup and contribution guidelines.

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <p>Made with ❤️ by the Media Vault Team</p>
  <p>For support, please open an issue or contact support@wron.ai</p>
</div>
# Run linters
make lint

# Format code
make format

# Update dependencies
make deps
```

## 🧪 Testing

### Running Tests

```bash
# Run unit tests
make test-unit

# Run integration tests
make test-integration

# Run end-to-end tests
make test-e2e
```

### Test Coverage

```bash
# Generate coverage report
make coverage

# View HTML coverage report
make coverage-html
```

## 📊 Monitoring

### Access Monitoring Tools

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093

### Key Metrics

- API response times
- Error rates
- Resource usage
- User activity
- Storage utilization

## 🔐 Authentication

### Keycloak Setup

1. Access Keycloak admin console: https://auth.localhost/admin
2. Log in with admin credentials
3. Import realm configuration from `keycloak/import/realm-export.json`
4. Configure identity providers and clients as needed

### User Management

- Create users in Keycloak admin console
- Assign roles and permissions
- Set up password policies
- Configure multi-factor authentication

## 📚 Documentation

### API Documentation

Access the interactive API documentation at:
- Swagger UI: https://localhost/api/docs
- OpenAPI Spec: https://localhost/api/docs.json

### Additional Resources

- [Developer Guide](docs/DEVELOPER_GUIDE.md)
- [API Reference](docs/API_REFERENCE.md)
- [Deployment Guide](docs/DEPLOYMENT.md)
- [Troubleshooting](docs/TROUBLESHOOTING.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Keycloak](https://www.keycloak.org/) for authentication
- [Fiber](https://gofiber.io/) for the Go web framework
- [Flutter](https://flutter.dev/) for the frontend
- [Prometheus](https://prometheus.io/) and [Grafana](https://grafana.com/) for monitoring
