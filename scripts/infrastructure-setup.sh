#!/bin/bash

echo "🏗️  Setting up Infrastructure Administration..."

# Create monitoring directories
mkdir -p monitoring/{prometheus/rules,alertmanager,grafana/{provisioning,dashboards},loki,promtail,filebeat}
mkdir -p traefik backups scripts

# Set permissions
chmod +x scripts/*.sh

# Create network if not exists
docker network create media-vault-network 2>/dev/null || true

echo "✅ Infrastructure setup completed"
echo ""
echo "📋 Available profiles:"
echo "  - default: Basic monitoring (Prometheus, Grafana, AlertManager)"
echo "  - elk: Elasticsearch, Logstash, Kibana"
echo "  - traefik: Traefik load balancer"
echo "  - redis-tools: Redis management tools"
echo "  - auto-update: Watchtower auto-updates"
echo "  - nginx-proxy: Nginx Proxy Manager"
echo ""
echo "🚀 Start with: docker-compose -f docker-compose.infrastructure.yml up -d"