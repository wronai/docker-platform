
# Media Vault - Makefile

# Color definitions
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: help \
        init \
        up down restart \
        build rebuild \
        test test-unit test-integration test-e2e \
        lint format \
        logs \
        clean \
        setup-keycloak \
        deploy \
        monitor

## Help
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

## Environment
init: ## Initialize development environment
	@echo "${GREEN}🚀 Initializing development environment...${RESET}"
	cp .env.example .env
	@echo "✅ Created .env file"
	@echo "${YELLOW}ℹ️  Please edit .env with your configuration${RESET}

## Docker Compose
up: ## Start all services
	@echo "${GREEN}🚀 Starting all services...${RESET}"
	docker-compose up -d

up-build: ## Rebuild and start all services
	@echo "${GREEN}🚀 Rebuilding and starting all services...${RESET}"
	docker-compose up -d --build

down: ## Stop all services
	@echo "${YELLOW}🛑 Stopping all services...${RESET}"
	docker-compose down

restart: down up ## Restart all services

logs: ## View logs from all services
	docker-compose logs -f

## Development
dev: up ## Start development environment

watch-backend: ## Watch backend for changes and rebuild
	@echo "${GREEN}👀 Watching backend for changes...${RESET}"
	docker-compose -f docker-compose.dev.yml up media-vault-api

watch-frontend: ## Watch frontend for changes and rebuild
	@echo "${GREEN}👀 Watching frontend for changes...${RESET}"
	docker-compose -f docker-compose.dev.yml up frontend

## Testing
test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	@echo "🧪 Running unit tests..."
	docker-compose run --rm media-vault-api go test -v ./... -short

test-integration: ## Run integration tests
	@echo "${GREEN}🧪 Running integration tests...${RESET}"
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit

test-e2e: ## Run end-to-end tests
	@echo "${GREEN}🧪 Running E2E tests...${RESET}"
	# TODO: Add E2E test command

coverage: ## Generate test coverage report
	@echo "${GREEN}📊 Generating test coverage...${RESET}"
	docker-compose run --rm media-vault-api go test -coverprofile=coverage.out ./...
	docker-compose run --rm media-vault-api go tool cover -html=coverage.out -o coverage.html

## Code Quality
lint: ## Run linters
	@echo "${GREEN}🔍 Running linters...${RESET}"
	docker-compose run --rm media-vault-api golangci-lint run

format: ## Format code
	@echo "${GREEN}🎨 Formatting code...${RESET}"
	docker-compose run --rm media-vault-api gofmt -w .

## Monitoring
monitor: ## Open monitoring dashboard
	@echo "${GREEN}📊 Opening monitoring dashboard...${RESET}"
	@echo "${YELLOW}Grafana: http://localhost:3000${RESET} (admin/grafana)"
	@echo "${YELLOW}Prometheus: http://localhost:9090${RESET}"

## Keycloak
setup-keycloak: ## Set up Keycloak with initial configuration
	@echo "${GREEN}🔐 Setting up Keycloak...${RESET}"
	./keycloak/setup-keycloak.sh

## Deployment
deploy: ## Deploy to production
	@echo "${GREEN}🚀 Deploying to production...${RESET}"
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

## Cleanup
clean: ## Remove all containers and volumes
	@echo "${YELLOW}🧹 Cleaning up...${RESET}"
	docker-compose down -v

dev-clean: clean ## Clean development environment
	@echo "${YELLOW}🧹 Cleaning development environment...${RESET}
	rm -rf node_modules

dist-clean: clean ## Remove all build artifacts and dependencies
	@echo "${YELLOW}🧹 Deep cleaning...${RESET}
	docker system prune -a --volumes

## Documentation
docs: ## Generate documentation
	@echo "${GREEN}📚 Generating documentation...${RESET}"
	# TODO: Add documentation generation command

## Complete Setup
full-stack: setup-keycloak up monitor ## Start complete stack with monitoring
	@echo "${GREEN}🎉 Media Vault is ready!${RESET}"
	@echo ""
	@echo "${YELLOW}📊 Service URLs:${RESET}"
	@echo "  🌐 Main App: http://localhost"
	@echo "  📊 Grafana: http://localhost:3000 (admin/grafana)"
	@echo "  📈 Prometheus: http://localhost:9090"
	@echo "  🔐 Keycloak: http://localhost:8080/admin (admin/admin123)"
	@echo "  ⚠️  AlertManager: http://localhost:9093"

## Infrastructure
infrastructure: ## Start infrastructure services
	@echo "${GREEN}🏗️  Starting infrastructure stack...${RESET}"
	docker-compose -f docker-compose.infrastructure.yml up -d

monitoring: ## Start monitoring stack
	@echo "${GREEN}📊 Starting monitoring stack...${RESET}"
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

# Podstawowe komendy
build:
	docker-compose build

up:
	docker-compose up -d

up-full:
	docker-compose --profile full up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	docker system prune -f

# Keycloak setup
setup-keycloak:
	@echo "🔐 Konfigurowanie Keycloak..."
	docker-compose up -d keycloak keycloak-db
	sleep 30
	chmod +x keycloak/setup-keycloak.sh
	./keycloak/setup-keycloak.sh
	./update-env.sh

# Keycloak management
keycloak-clean:
	./setup-keycloak.sh --clean

keycloak-export:
	./setup-keycloak.sh --export-only

# Media Vault with Keycloak
up-with-auth: setup-keycloak
	docker-compose up -d
	@echo ""
	@echo "🎉 Media Vault z Keycloak uruchomiony!"
	@echo "🌐 Aplikacja: http://localhost"
	@echo "🔐 Keycloak Admin: http://localhost:8443/admin"
	@echo "👤 Test login: vaultadmin / admin123"

# Restart specific services
restart-api:
	docker-compose restart media-vault-api

restart-keycloak:
	docker-compose restart keycloak

# Status and debugging
status:
	docker-compose ps

shell-api:
	docker-compose exec media-vault-api sh

shell-keycloak:
	docker-compose exec keycloak bash

# Quick start dla development
dev-start: clean up-with-auth
	@echo "🚀 Development environment ready!"

# Production setup
prod-setup:
	@echo "🏭 Konfiguracja produkcyjna..."
	@echo "⚠️  Pamiętaj o:"
	@echo "   1. Zmianie haseł w production"
	@echo "   2. Konfiguracji SSL"
	@echo "   3. Backup strategii"
	@echo "   4. Monitoringu"

# Backup
backup:
	@echo "💾 Tworzenie backup..."
	docker-compose exec media-vault-api sqlite3 /data/media.db ".backup /data/backup_$(shell date +%Y%m%d_%H%M%S).db"
	./setup-keycloak.sh --export-only
	tar -czf media-vault-backup-$(shell date +%Y%m%d).tar.gz data/ uploads/ media-vault-realm-export.json


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
	@echo "Media Vault + Keycloak - Dostępne komendy:"
	@echo ""
	@echo "Podstawowe:"
	@echo "  make up-with-auth    Uruchom z pełną konfiguracją Keycloak"
	@echo "  make setup-keycloak  Tylko konfiguracja Keycloak"
	@echo "  make dev-start       Szybki start dla development"
	@echo ""
	@echo "Zarządzanie:"
	@echo "  make up              Uruchom podstawowe serwisy"
	@echo "  make down            Zatrzymaj wszystko"
	@echo "  make clean           Wyczyść wszystko"
	@echo "  make logs            Pokaż logi"
	@echo ""
	@echo "Keycloak:"
	@echo "  make keycloak-clean  Wyczyść konfigurację Keycloak"
	@echo "  make keycloak-export Eksportuj konfigurację"
	@echo "  make restart-keycloak Restart Keycloak"
	@echo ""
	@echo "Debugging:"
	@echo "  make status          Status kontenerów"
	@echo "  make shell-api       Wejdź do API container"
	@echo "  make backup          Stwórz backup"


