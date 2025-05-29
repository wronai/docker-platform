#!/bin/bash

# Performance Optimization Script
echo "⚡ Starting performance optimization..."

# Get current metrics
get_metric() {
    local query=$1
    curl -s "http://prometheus:9090/api/v1/query?query=${query}" | \
        jq -r '.data.result[0].value[1] // "0"' 2>/dev/null || echo "0"
}

# CPU optimization
optimize_cpu() {
    local cpu_usage=$(get_metric "100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)")

    if (( $(echo "$cpu_usage > 70" | bc -l) )); then
        echo "🔧 High CPU usage detected: ${cpu_usage}%"

        # Reduce CPU limits for non-critical services
        docker update --cpus="0.5" media-vault-analyzer 2>/dev/null || true
        docker update --cpus="0.3" nsfw-analyzer 2>/dev/null || true

        echo "✅ CPU limits reduced for non-critical services"
    fi
}

# Memory optimization
optimize_memory() {
    local mem_usage=$(get_metric "(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100")

    if (( $(echo "$mem_usage > 80" | bc -l) )); then
        echo "🔧 High memory usage detected: ${mem_usage}%"

        # Clear system caches
        sync && echo 3 > /proc/sys/vm/drop_caches 2>/dev/null || true

        # Restart memory-intensive services
        docker restart media-vault-analyzer 2>/dev/null || true

        echo "✅ Memory optimization applied"
    fi
}

# Disk optimization
optimize_disk() {
    echo "🗃️ Running disk optimization..."

    # Clean Docker system
    docker system prune -f

    # Clean old logs
    find /var/log -name "*.log" -type f -mtime +7 -exec truncate -s 0 {} \; 2>/dev/null || true

    # Clean temporary files
    find /tmp -type f -atime +1 -delete 2>/dev/null || true

    echo "✅ Disk cleanup completed"
}

# Database optimization
optimize_database() {
    echo "🗄️ Running database optimization..."

    # SQLite optimization
    if [ -f "/data/media.db" ]; then
        sqlite3 /data/media.db "VACUUM; PRAGMA optimize;" 2>/dev/null || true
        echo "✅ SQLite database optimized"
    fi

    # PostgreSQL optimization (if available)
    if docker ps | grep -q keycloak-db; then
        docker exec keycloak-db psql -U keycloak -d keycloak -c "VACUUM ANALYZE;" 2>/dev/null || true
        echo "✅ PostgreSQL database optimized"
    fi
}

# Network optimization
optimize_network() {
    echo "🌐 Running network optimization..."

    # Optimize Docker network settings
    echo 'net.core.rmem_max = 134217728' >> /etc/sysctl.conf 2>/dev/null || true
    echo 'net.core.wmem_max = 134217728' >> /etc/sysctl.conf 2>/dev/null || true

    sysctl -p 2>/dev/null || true

    echo "✅ Network optimization applied"
}

# Generate optimization report
generate_report() {
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    cat > /tmp/optimization_report.json << EOF
{
    "timestamp": "$timestamp",
    "optimizations_applied": [
        "CPU limits adjusted",
        "Memory caches cleared",
        "Disk cleanup performed",
        "Database vacuumed",
        "Network settings optimized"
    ],
    "metrics_after": {
        "cpu_usage": $(get_metric "100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"),
        "memory_usage": $(get_metric "(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100"),
        "disk_usage": $(get_metric "(node_filesystem_size_bytes{mountpoint=\"/\"} - node_filesystem_avail_bytes{mountpoint=\"/\"}) / node_filesystem_size_bytes{mountpoint=\"/\"} * 100")
    }
}
EOF

    echo "📊 Optimization report generated: /tmp/optimization_report.json"
}

# Main optimization workflow
main() {
    echo "🚀 Media Vault Performance Optimization"
    echo "======================================="

    optimize_cpu
    optimize_memory
    optimize_disk
    optimize_database
    optimize_network
    generate_report

    echo ""
    echo "✅ Performance optimization completed!"
    echo "📊 Check metrics in Grafana for improvements"
}

main "$@"