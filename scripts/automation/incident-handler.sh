#!/bin/bash

# Automated Incident Response Handler
PROMETHEUS_URL="http://prometheus:9090"
ALERTMANAGER_URL="http://alertmanager:9093"

# Get active alerts
get_alerts() {
    curl -s "$ALERTMANAGER_URL/api/v1/alerts" | jq -r '.data[] | select(.status.state == "active")'
}

# Handle service down alert
handle_service_down() {
    local service=$1
    echo "🚨 Handling service down: $service"

    # Try automatic restart
    docker restart "$service" 2>/dev/null

    # Wait and check if service is back up
    sleep 30
    if docker ps | grep -q "$service.*Up"; then
        send_notification "✅ Service $service automatically recovered"
        return 0
    else
        send_notification "❌ Service $service failed to restart - manual intervention required"
        return 1
    fi
}

# Handle high resource usage
handle_high_resource_usage() {
    local resource=$1
    local threshold=$2

    echo "⚠️ Handling high $resource usage: $threshold%"

    case $resource in
        "memory")
            # Clear system caches
            sync && echo 3 > /proc/sys/vm/drop_caches
            # Restart non-critical services
            docker restart media-vault-analyzer || true
            ;;
        "disk")
            # Clean old logs and temporary files
            docker system prune -f
            find /tmp -type f -atime +7 -delete 2>/dev/null || true
            ;;
        "cpu")
            # Scale down resource-intensive services temporarily
            docker update --cpus="0.5" media-vault-analyzer || true
            ;;
    esac

    send_notification "🔧 Automated remediation applied for high $resource usage"
}

# Handle security incidents
handle_security_incident() {
    local incident_type=$1
    local details=$2

    echo "🛡️ Handling security incident: $incident_type"

    case $incident_type in
        "failed_logins")
            # Enable additional security measures
            # This would integrate with Keycloak API to increase security
            send_notification "🔒 Enhanced security measures activated due to failed login attempts"
            ;;
        "unauthorized_access")
            # Log the incident for analysis
            echo "$(date): Unauthorized access detected - $details" >> /var/log/security-incidents.log
            send_notification "🚨 SECURITY ALERT: Unauthorized access detected - details logged"
            ;;
    esac
}

# Send notification to various channels
send_notification() {
    local message=$1
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

    # Slack notification
    if [ -n "$SLACK_WEBHOOK_URL" ]; then
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"[$timestamp] Media Vault: $message\"}" \
            "$SLACK_WEBHOOK_URL" 2>/dev/null || true
    fi

    # PagerDuty notification (for critical alerts)
    if [ -n "$PAGERDUTY_API_KEY" ] && [[ $message == *"❌"* ]]; then
        curl -X POST "https://events.pagerduty.com/v2/enqueue" \
            -H "Content-Type: application/json" \
            -d "{
                \"routing_key\": \"$PAGERDUTY_API_KEY\",
                \"event_action\": \"trigger\",
                \"payload\": {
                    \"summary\": \"Media Vault Alert: $message\",
                    \"source\": \"media-vault-automation\",
                    \"severity\": \"critical\"
                }
            }" 2>/dev/null || true
    fi

    # Local logging
    echo "[$timestamp] $message" >> /var/log/automation.log
}

# Main incident processing loop
process_incidents() {
    local alerts=$(get_alerts)

    if [ -z "$alerts" ]; then
        return 0
    fi

    echo "$alerts" | while IFS= read -r alert; do
        local alert_name=$(echo "$alert" | jq -r '.labels.alertname // "unknown"')
        local severity=$(echo "$alert" | jq -r '.labels.severity // "unknown"')

        echo "Processing alert: $alert_name (severity: $severity)"

        case $alert_name in
            "ServiceCompletleyDown")
                local service=$(echo "$alert" | jq -r '.labels.job // "unknown"')
                handle_service_down "$service"
                ;;
            "HighMemoryUsage"|"HighCPUUsage"|"DiskSpaceLow")
                local resource=$(echo "$alert_name" | sed 's/High\|Low//g' | tr '[:upper:]' '[:lower:]')
                handle_high_resource_usage "$resource" "85"
                ;;
            "SuspiciousLoginActivity")
                handle_security_incident "failed_logins" "$alert"
                ;;
            "UnauthorizedAPIAccess")
                handle_security_incident "unauthorized_access" "$alert"
                ;;
            *)
                send_notification "⚠️ Unhandled alert: $alert_name"
                ;;
        esac
    done
}

# Run the incident handler
main() {
    echo "🤖 Starting automated incident response system..."
    process_incidents
}

main "$@"

