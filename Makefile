
# Makefile (updated with monitoring commands)
.PHONY: help infrastructure monitoring full-stack

# Complete setup commands
full-stack-with-monitoring: setup-keycloak
	@echo "🚀 Starting complete Media Vault stack with monitoring..."
	docker-compose up -d
	docker-compose -f docker-compose.monitoring.yml up -d
	@echo ""
	@echo "🎉 Media Vault with monitoring is ready!"
	@echo ""
	@echo "📊 Service URLs:"
	@echo "├── 🌐 Main App: http://localhost"
	@echo "├── 📊 Grafana: http://localhost:3333 (admin/grafana123)"
	@echo "├── 📈 Prometheus: http://localhost:9090"
	@echo "├── 🔐 Keycloak: http://localhost:8443/admin (admin/admin123)"
	@echo "├── ⚠️  AlertManager: http://localhost:9093"
	@echo "└── 🐳 Portainer: http://localhost:9000"
	@echo ""
	@echo "👤 Test Accounts:"
	@echo "├── Admin: vaultadmin / admin123"
	@echo "└── User:  vaultuser / user123"

infrastructure:
	@echo "🏗️  Starting infrastructure stack..."
	docker-compose -f docker-compose.infrastructure.yml up -d
	@echo "✅ Infrastructure started"

monitoring:
	@echo "📊 Starting monitoring stack..."
	docker-compose -f docker-compose.monitoring.yml up -d
	@echo "✅ Monitoring started"
	@echo "📊 Grafana: http://localhost:3333"
	@echo "📈 Prometheus: http://localhost:9090"

monitoring-down:
	@echo "📊 Stopping monitoring stack..."
	docker-compose -f docker-compose.monitoring.yml down

# Quick monitoring access
grafana:
	@echo "🌐 Opening Grafana..."
	@echo "URL: http://localhost:3333"
	@echo "Login: admin / grafana123"

prometheus:
	@echo "🌐 Opening Prometheus..."
	@echo "URL: http://localhost:9090"

alerts:
	@echo "🌐 Opening AlertManager..."
	@echo "URL: http://localhost:9093"

# Monitoring management
monitoring-restart:
	docker-compose -f docker-compose.monitoring.yml restart

monitoring-logs:
	docker-compose -f docker-compose.monitoring.yml logs -f

monitoring-status:
	@echo "📊 Monitoring Status:"
	@echo "===================="
	@docker-compose -f docker-compose.monitoring.yml ps

# Dashboard management
dashboard-import:
	@echo "📊 Importing custom dashboards..."
	# Custom script to import dashboards via Grafana API

dashboard-backup:
	@echo "💾 Backing up Grafana dashboards..."
	# Custom script to export dashboards

# Alerting management
test-alerts:
	@echo "🚨 Testing alert system..."
	curl -XPOST http://localhost:9093/api/v1/alerts \
		-H "Content-Type: application/json" \
		-d '[{"labels":{"alertname":"TestAlert","severity":"warning","instance":"test"},"annotations":{"summary":"Test alert"}}]'

silence-alerts:
	@echo "🔇 Creating alert silence..."
	# Script to create alert silence via AlertManager API

# Performance testing with monitoring
load-test-monitored:
	@echo "⚡ Running monitored load test..."
	docker run --rm -i --network media-vault-network \
		grafana/k6 run - <<EOF
		import http from 'k6/http';
		import { check } from 'k6';
		export let options = {
			stages: [
				{ duration: '2m', target: 50 },
				{ duration: '5m', target: 50 },
				{ duration: '2m', target: 100 },
				{ duration: '5m', target: 100 },
				{ duration: '2m', target: 0 },
			],
			thresholds: {
				http_req_duration: ['p(95)<500'],
				http_req_failed: ['rate<0.1'],
			},
		};
		export default function() {
			let response = http.get('http://media-vault-api:8080/health');
			check(response, {
				'status is 200': (r) => r.status === 200,
				'response time < 500ms': (r) => r.timings.duration < 500,
			});
		}
	EOF

# System health check
health-check-full:
	@echo "🏥 Comprehensive Health Check:"
	@echo "=============================="
	@echo "📊 Checking main services..."
	@curl -s http://localhost:8080/health | jq '.' || echo "❌ API: DOWN"
	@curl -s http://localhost:8443/health/ready | jq '.' || echo "❌ Keycloak: DOWN"
	@echo ""
	@echo "📊 Checking monitoring..."
	@curl -s http://localhost:9090/-/healthy || echo "❌ Prometheus: DOWN"
	@curl -s http://localhost:3333/api/health || echo "❌ Grafana: DOWN"
	@curl -s http://localhost:9093/-/healthy || echo "❌ AlertManager: DOWN"
	@echo ""
	@echo "📊 Checking containers..."
	@docker-compose ps
	@echo ""
	@echo "📊 Resource usage:"
	@docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# Backup monitoring data
backup-monitoring:
	@echo "💾 Backing up monitoring data..."
	docker run --rm -v grafana_data:/source -v $(PWD)/backups:/backup alpine tar czf /backup/grafana-backup-$(shell date +%Y%m%d).tar.gz -C /source .
	docker run --rm -v prometheus_data:/source -v $(PWD)/backups:/backup alpine tar czf /backup/prometheus-backup-$(shell date +%Y%m%d).tar.gz -C /source .

# Help
help:
	@echo "Media Vault - Available Commands:"
	@echo "================================="
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make full-stack-with-monitoring  Complete setup with monitoring"
	@echo "  make up                          Basic application stack"
	@echo "  make monitoring                  Monitoring stack only"
	@echo ""
	@echo "📊 Monitoring:"
	@echo "  make grafana                     Open Grafana (localhost:3333)"
	@echo "  make prometheus                  Open Prometheus (localhost:9090)"
	@echo "  make alerts                      Open AlertManager (localhost:9093)"
	@echo "  make monitoring-status           Show monitoring status"
	@echo "  make monitoring-logs             Show monitoring logs"
	@echo ""
	@echo "🔧 Management:"
	@echo "  make health-check-full           Complete system health check"
	@echo "  make load-test-monitored         Performance test with monitoring"
	@echo "  make backup-monitoring           Backup monitoring data"
	@echo "  make test-alerts                 Test alert system"
	@echo ""
	@echo "🔐 Authentication:"
	@echo "  make setup-keycloak              Configure Keycloak"
	@echo "  make keycloak-clean              Reset Keycloak config"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  make down                        Stop all services"
	@echo "  make clean                       Clean up everything"
	@echo "  make logs                        Show application logs"