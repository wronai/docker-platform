# 🏗️ Media Vault - Infrastructure Administration

## 📊 Complete Infrastructure Stack

```ascii
                    🏗️ INFRASTRUCTURE OVERVIEW
┌─────────────────────────────────────────────────────────────────────┐
│                          MONITORING LAYER                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │  Prometheus │ │   Grafana   │ │AlertManager │ │   Jaeger    │   │
│  │ :9090       │ │ :3333       │ │ :9093       │ │ :16686      │   │
│  │ Metrics     │ │ Dashboards  │ │ Alerts      │ │ Tracing     │   │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────┐
│                           LOGGING LAYER                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │    Loki     │ │  Promtail   │ │    ELK      │ │  FileBeat   │   │
│  │ :3100       │ │ Log Ship    │ │ :5601       │ │ Log Ship    │   │
│  │ Log Aggreg  │ │             │ │ Advanced    │ │             │   │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────┐
│                         MANAGEMENT LAYER                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │  Portainer  │ │  pgAdmin    │ │Redis Cmd    │ │ Watchtower  │   │
│  │ :9000       │ │ :5050       │ │ :8081       │ │ Auto Update │   │
│  │ Docker Mgmt │ │ DB Admin    │ │ Redis Mgmt  │ │             │   │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────┐
│                         SECURITY & BACKUP                          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│  │   Traefik   │ │ Nginx Proxy │ │   Backup    │ │  Security   │   │
│  │ :8080       │ │ :8181       │ │ Service     │ │ Scanner     │   │
│  │ Load Bal    │ │ Proxy Mgmt  │ │ Automated   │ │ Trivy       │   │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start Commands

### **Infrastructure Deployment:**
```bash
# Full infrastructure stack
make infrastructure

# Basic monitoring only
make monitoring  

# Everything (Vault + Infrastructure)
make full-stack

# Specialized stacks
make infra-elk      # Elasticsearch + Kibana
make infra-traefik  # Traefik load balancer
make infra-nginx    # Nginx Proxy Manager
```

### **Management Commands:**
```bash
# Status and health
make infra-status   # Infrastructure status
make health-check   # Health check all services
make dev-logs       # Development logs

# Performance and security
make load-test      # K6 performance testing
make security-scan  # Trivy vulnerability scan
make performance-report  # Generate report

# Backup and maintenance
make backup-now     # Immediate backup
make cleanup-all    # Clean unused resources
```

## 📊 Monitoring & Observability

### **Prometheus Metrics:**
- ✅ **System metrics** (CPU, Memory, Disk, Network)
- ✅ **Container metrics** (Docker stats via cAdvisor)
- ✅ **Application metrics** (API response times, errors)
- ✅ **Keycloak metrics** (authentication events)
- ✅ **Custom alerts** (service down, high resource usage)

### **Grafana Dashboards:**
- 📊 **Infrastructure Overview** - System health
- 📈 **Application Performance** - API metrics
- 🔐 **Security Dashboard** - Auth events
- 💾 **Storage Analytics** - Vault usage
- 🐳 **Container Monitoring** - Docker stats

### **Alert Definitions:**
```yaml
Critical Alerts:
  - Service Down (> 1 min)
  - Disk Space < 10%
  - Memory Usage > 90%
  - Failed Login Spike

Warning Alerts:
  - High CPU (> 80%, 5 min)
  - Response Time > 2s
  - Memory Usage > 85%
```

## 🏗️ Infrastructure Components

### **Core Monitoring:**
| Service | Port | Purpose | Dashboard |
|---------|------|---------|-----------|
| **Prometheus** | 9090 | Metrics collection | http://localhost:9090 |
| **Grafana** | 3333 | Visualization | http://localhost:3333 |
| **AlertManager** | 9093 | Alert management | http://localhost:9093 |
| **Node Exporter** | 9100 | Host metrics | http://localhost:9100 |
| **cAdvisor** | 8888 | Container metrics | http://localhost:8888 |

### **Advanced Logging:**
| Service | Port | Purpose | Dashboard |
|---------|------|---------|-----------|
| **Loki** | 3100 | Log aggregation | via Grafana |
| **Promtail** | - | Log shipping | - |
| **Elasticsearch** | 9200 | Search & analytics | http://localhost:9200 |
| **Kibana** | 5601 | Log visualization | http://localhost:5601 |
| **Jaeger** | 16686 | Distributed tracing | http://localhost:16686 |

### **Management Tools:**
| Service | Port | Purpose | Dashboard |
|---------|------|---------|-----------|
| **Portainer** | 9000 | Docker management | http://localhost:9000 |
| **pgAdmin** | 5050 | PostgreSQL admin | http://localhost:5050 |
| **Redis Commander** | 8081 | Redis management | http://localhost:8081 |
| **Traefik** | 8080 | Load balancer | http://localhost:8080 |

## 🔒 Security & Compliance

### **Security Features:**
- 🛡️ **Container scanning** with Trivy
- 🔐 **Secret management** via Docker secrets
- 🌐 **Network policies** and segmentation
- 📝 **Audit logging** for all admin actions
- 🚨 **Intrusion detection** alerts

### **Backup Strategy:**
- 📅 **Automated daily backups** at 2 AM
- ☁️ **S3 upload** for off-site storage
- 🔄 **30-day retention** policy
- 🗜️ **Compressed archives** with encryption
- ✅ **Backup verification** checks

### **Production Security:**
```yaml
Security Headers:
  - Strict-Transport-Security
  - X-Content-Type-Options: nosniff
  - X-Frame-Options: DENY
  - Referrer-Policy: strict-origin

Rate Limiting:
  - 100 requests/minute per IP
  - Admin panel IP restriction
  - API throttling enabled
```

## 📈 Performance Optimization

### **Resource Monitoring:**
- 📊 **Real-time dashboards** 
- ⚡ **Performance alerts**
- 🎯 **SLA monitoring** (99.9% uptime)
- 📉 **Trend analysis**

### **Auto-scaling:**
- 🔄 **HPA** for API services
- 📈 **Load balancing** with Traefik
- 🎛️ **Resource limits** per container
- ⚖️ **Dynamic scaling** based on metrics

### **Performance Testing:**
```bash
# K6 load testing
make load-test

# Response time monitoring
curl -w "@curl-format.txt" http://localhost:8080/health

# Resource usage
docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

## 🚨 Incident Response

### **Alert Channels:**
- 📧 **Email notifications** to admin
- 💬 **Slack integration** for team alerts
- 📱 **PagerDuty** for critical issues
- 📊 **Dashboard notifications**

### **Runbooks:**
1. **Service Down:**
   ```bash
   make health-check     # Identify failed service
   docker-compose restart <service>
   make dev-logs        # Check logs
   ```

2. **High Resource Usage:**
   ```bash
   make performance-report  # Generate report
   docker stats            # Check current usage
   make cleanup-all        # Clean resources
   ```

3. **Security Incident:**
   ```bash
   make security-scan      # Vulnerability check
   docker-compose logs     # Review access logs
   # Block suspicious IPs in Caddy/Traefik
   ```

## 🎛️ Administration Workflows

### **Daily Operations:**
- ✅ Check health dashboard
- 📊 Review performance metrics
- 🔍 Scan security alerts
- 💾 Verify backup completion

### **Weekly Tasks:**
- 📈 Performance report review
- 🛡️ Security scan execution
- 🧹 Resource cleanup
- 📋 Capacity planning review

### **Monthly Maintenance:**
- 🔄 Update container images
- 📊 SLA report generation
- 🗄️ Backup retention cleanup
- 🔐 Security audit review

## 💡 Pro Tips

### **Optimization:**
- Use **profiles** for different environments
- Enable **automatic updates** with Watchtower
- Set up **custom metrics** for business logic
- Implement **circuit breakers** for resilience

### **Troubleshooting:**
- Always check **logs first**: `make dev-logs`
- Use **health endpoints** for service status
- Monitor **resource usage** trends
- Keep **runbooks updated** for common issues

This infrastructure setup provides **enterprise-grade** monitoring, security, and management capabilities for Media Vault! 🚀