#!/bin/bash

# Exit on error
set -e

# Configuration
KEYCLOAK_URL="http://localhost:8082"
ADMIN_USER="admin"
ADMIN_PASSWORD="admin123"
REALM="photovault"
CLIENT_ID="photovault-client"
CLIENT_SECRET="your-client-secret-here"  # Change this in production
REDIRECT_URI="http://localhost:8000/*"

# Wait for Keycloak to be ready
echo "Waiting for Keycloak to be ready..."
until $(curl --output /dev/null --silent --head --fail ${KEYCLOAK_URL}/auth/realms/master); do
    printf '.'
    sleep 5
done
echo "\nKeycloak is ready!"

# Get admin token
function get_admin_token {
    curl -s \
        -d "client_id=admin-cli" \
        -d "username=${ADMIN_USER}" \
        -d "password=${ADMIN_PASSWORD}" \
        -d "grant_type=password" \
        "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" | jq -r '.access_token'
}

# Create realm
function create_realm {
    local token=$1
    
    echo "Creating realm ${REALM}..."
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${token}" \
        -d "{
            \"realm\": \"${REALM}\",
            \"enabled\": true,
            \"registrationAllowed\": true,
            \"registrationEmailAsUsername\": true,
            \"verifyEmail\": true,
            \"resetPasswordAllowed\": true,
            \"editUsernameAllowed\": false,
            \"loginWithEmailAllowed\": true,
            \"duplicateEmailsAllowed\": false,
            \"bruteForceProtected\": true
        }" \
        "${KEYCLOAK_URL}/admin/realms"
}

# Create client
function create_client {
    local token=$1
    
    echo "Creating client ${CLIENT_ID}..."
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${token}" \
        -d "{
            \"clientId\": \"${CLIENT_ID}\",
            \"enabled\": true,
            \"publicClient\": false,
            \"redirectUris\": [\"${REDIRECT_URI}\"],
            \"webOrigins\": [\"*\"],
            \"standardFlowEnabled\": true,
            \"implicitFlowEnabled\": false,
            \"directAccessGrantsEnabled\": true,
            \"serviceAccountsEnabled\": true,
            \"authorizationServicesEnabled\": true,
            \"secret\": \"${CLIENT_SECRET}\",
            \"bearerOnly\": false,
            \"consentRequired\": false,
            \"fullScopeAllowed\": true
        }" \
        "${KEYCLOAK_URL}/admin/realms/${REALM}/clients"
}

# Create roles
function create_roles {
    local token=$1
    
    local roles=("admin" "partner" "user")
    
    for role in "${roles[@]}"; do
        echo "Creating role ${role}..."
        curl -s -X POST \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer ${token}" \
            -d "{\"name\": \"${role}\"}" \
            "${KEYCLOAK_URL}/admin/realms/${REALM}/roles"
    done
}

# Create admin user
function create_admin_user {
    local token=$1
    local username=$2
    local password=$3
    
    echo "Creating admin user ${username}..."
    
    # Create user
    local user_id=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${token}" \
        -d "{
            \"username\": \"${username}\",
            \"email\": \"${username}@example.com\",
            \"enabled\": true,
            \"emailVerified\": true,
            \"credentials\": [{
                \"type\": \"password\",
                \"value\": \"${password}\",
                \"temporary\": false
            }]
        }" \
        "${KEYCLOAK_URL}/admin/realms/${REALM}/users" \
        -w "%{http_code}" -o /dev/null)
    
    # Get user ID
    local user_id=$(curl -s \
        -H "Authorization: Bearer ${token}" \
        "${KEYCLOAK_URL}/admin/realms/${REALM}/users?username=${username}" | jq -r '.[0].id')
    
    if [ "$user_id" == "null" ] || [ -z "$user_id" ]; then
        echo "Failed to get user ID for ${username}"
        return 1
    fi
    
    # Assign admin role
    curl -s -X POST \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${token}" \
        -d '[{"id":"admin","name":"admin"}]' \
        "${KEYCLOAK_URL}/admin/realms/${REALM}/users/${user_id}/role-mappings/realm"
    
    echo "Admin user created with ID: ${user_id}"
}

# Main execution
echo "Starting Keycloak setup..."

# Get admin token
echo "Getting admin token..."
TOKEN=$(get_admin_token)

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "Failed to get admin token. Check Keycloak credentials and URL."
    exit 1
fi

# Create realm
create_realm "$TOKEN"

# Create client
create_client "$TOKEN"

# Create roles
create_roles "$TOKEN"

# Create admin user
create_admin_user "$TOKEN" "admin@photovault.com" "admin123"

echo "Keycloak setup completed successfully!"
echo "Keycloak Admin URL: ${KEYCLOAK_URL}/admin"
echo "Realm: ${REALM}"
echo "Client ID: ${CLIENT_ID}"
echo "Client Secret: ${CLIENT_SECRET}"
