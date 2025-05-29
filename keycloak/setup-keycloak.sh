#!/bin/bash

# setup-keycloak.sh - Automatyczne skonfigurowanie Keycloak dla Media Vault

set -e

echo "🔐 Konfigurowanie Keycloak dla Media Vault..."

# Kolory dla output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Konfiguracja
KEYCLOAK_URL="https://localhost:8445"
ADMIN_USER="admin"
ADMIN_PASS="admin123"
REALM_NAME="media-vault"

# Funkcje pomocnicze
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Sprawdź czy Keycloak jest dostępny
wait_for_keycloak() {
    echo "⏳ Oczekiwanie na uruchomienie Keycloak..."
    for i in {1..60}; do
        if curl -f -s "$KEYCLOAK_URL/health/ready" > /dev/null 2>&1; then
            print_status "Keycloak jest gotowy"
            return 0
        fi
        echo "Próba $i/60 - Keycloak nie jest jeszcze gotowy..."
        sleep 5
    done
    print_error "Keycloak nie uruchomił się w czasie 5 minut"
    return 1
}

# Uzyskaj access token
get_admin_token() {
    echo "🔑 Uzyskiwanie tokenu administratora..."

    RESPONSE=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid_connect/token" \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "username=$ADMIN_USER" \
        -d "password=$ADMIN_PASS" \
        -d "grant_type=password" \
        -d "client_id=admin-cli")

    TOKEN=$(echo $RESPONSE | jq -r '.access_token')

    if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
        print_error "Nie udało się uzyskać tokenu administratora"
        echo "Response: $RESPONSE"
        return 1
    fi

    print_status "Token administratora uzyskany"
    export ADMIN_TOKEN="$TOKEN"
}

# Stwórz realm Media Vault
create_realm() {
    echo "🏰 Tworzenie realm 'media-vault'..."

    # Sprawdź czy realm już istnieje
    REALM_CHECK=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" \
        "$KEYCLOAK_URL/admin/realms/$REALM_NAME" \
        -w "%{http_code}" -o /dev/null)

    if [ "$REALM_CHECK" == "200" ]; then
        print_warning "Realm '$REALM_NAME' już istnieje"
        return 0
    fi

    # Stwórz realm
    curl -s -X POST "$KEYCLOAK_URL/admin/realms" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "realm": "media-vault",
            "enabled": true,
            "displayName": "Media Vault",
            "registrationAllowed": true,
            "registrationEmailAsUsername": true,
            "rememberMe": true,
            "verifyEmail": false,
            "loginWithEmailAllowed": true,
            "duplicateEmailsAllowed": false,
            "resetPasswordAllowed": true,
            "editUsernameAllowed": false,
            "bruteForceProtected": true
        }'

    print_status "Realm '$REALM_NAME' utworzony"
}

# Stwórz role
create_roles() {
    echo "👥 Tworzenie ról..."

    # Role do utworzenia
    ROLES=("vault-user" "vault-admin" "vault-analyzer")

    for role in "${ROLES[@]}"; do
        curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/roles" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{
                \"name\": \"$role\",
                \"description\": \"$role role for Media Vault\"
            }"
        print_status "Rola '$role' utworzona"
    done
}

# Stwórz klientów
create_clients() {
    echo "📱 Tworzenie klientów..."

    # Frontend Client (public)
    curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "clientId": "media-vault-frontend",
            "name": "Media Vault Frontend",
            "enabled": true,
            "publicClient": true,
            "standardFlowEnabled": true,
            "implicitFlowEnabled": false,
            "directAccessGrantsEnabled": false,
            "serviceAccountsEnabled": false,
            "protocol": "openid-connect",
            "fullScopeAllowed": true,
            "redirectUris": [
                "http://localhost:3000/*",
                "http://localhost/*"
            ],
            "webOrigins": [
                "http://localhost:3000",
                "http://localhost"
            ],
            "defaultClientScopes": [
                "web-origins",
                "profile",
                "roles",
                "email"
            ]
        }'

    print_status "Klient 'media-vault-frontend' utworzony"

    # API Client (confidential)
    curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "clientId": "media-vault-api",
            "name": "Media Vault API",
            "enabled": true,
            "publicClient": false,
            "standardFlowEnabled": true,
            "implicitFlowEnabled": false,
            "directAccessGrantsEnabled": true,
            "serviceAccountsEnabled": true,
            "protocol": "openid-connect",
            "secret": "vault-api-secret-123",
            "redirectUris": ["*"],
            "webOrigins": ["*"],
            "defaultClientScopes": [
                "web-origins",
                "profile",
                "roles",
                "email"
            ]
        }'

    print_status "Klient 'media-vault-api' utworzony"

    # Admin Client (public)
    curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/clients" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "clientId": "media-vault-admin",
            "name": "Media Vault Admin Panel",
            "enabled": true,
            "publicClient": true,
            "standardFlowEnabled": true,
            "implicitFlowEnabled": false,
            "directAccessGrantsEnabled": false,
            "serviceAccountsEnabled": false,
            "protocol": "openid-connect",
            "redirectUris": [
                "http://localhost:3001/*",
                "http://localhost/admin/*"
            ],
            "webOrigins": [
                "http://localhost:3001",
                "http://localhost"
            ],
            "defaultClientScopes": [
                "web-origins",
                "profile",
                "roles",
                "email"
            ]
        }'

    print_status "Klient 'media-vault-admin' utworzony"
}

# Stwórz użytkowników testowych
create_users() {
    echo "👤 Tworzenie użytkowników testowych..."

    # Admin user
    USER_ID=$(curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "vaultadmin",
            "enabled": true,
            "emailVerified": true,
            "firstName": "Vault",
            "lastName": "Administrator",
            "email": "admin@mediavault.local",
            "credentials": [
                {
                    "type": "password",
                    "value": "admin123",
                    "temporary": false
                }
            ]
        }' -D - | grep -i location | cut -d'/' -f8 | tr -d '\r')

    # Przypisz role admin
    curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users/$USER_ID/role-mappings/realm" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '[
            {"name": "vault-admin"},
            {"name": "vault-user"}
        ]'

    print_status "Użytkownik 'vaultadmin' utworzony z rolą admin"

    # Regular user
    USER_ID=$(curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "vaultuser",
            "enabled": true,
            "emailVerified": true,
            "firstName": "Vault",
            "lastName": "User",
            "email": "user@mediavault.local",
            "credentials": [
                {
                    "type": "password",
                    "value": "user123",
                    "temporary": false
                }
            ]
        }' -D - | grep -i location | cut -d'/' -f8 | tr -d '\r')

    # Przypisz rolę user
    curl -s -X POST "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users/$USER_ID/role-mappings/realm" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '[
            {"name": "vault-user"}
        ]'

    print_status "Użytkownik 'vaultuser' utworzony z rolą user"
}

# Konfiguruj ustawienia realm
configure_realm_settings() {
    echo "⚙️ Konfigurowanie ustawień realm..."

    # Update realm settings
    curl -s -X PUT "$KEYCLOAK_URL/admin/realms/$REALM_NAME" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "realm": "media-vault",
            "enabled": true,
            "displayName": "Media Vault",
            "registrationAllowed": true,
            "registrationEmailAsUsername": true,
            "rememberMe": true,
            "verifyEmail": false,
            "loginWithEmailAllowed": true,
            "duplicateEmailsAllowed": false,
            "resetPasswordAllowed": true,
            "editUsernameAllowed": false,
            "bruteForceProtected": true,
            "passwordPolicy": "length(8) and digits(1) and lowerCase(1) and upperCase(1)",
            "sslRequired": "none",
            "loginTheme": "keycloak",
            "adminTheme": "keycloak",
            "emailTheme": "keycloak",
            "accountTheme": "keycloak",
            "accessTokenLifespan": 300,
            "accessTokenLifespanForImplicitFlow": 900,
            "ssoSessionIdleTimeout": 1800,
            "ssoSessionMaxLifespan": 36000,
            "offlineSessionIdleTimeout": 2592000,
            "offlineSessionMaxLifespanEnabled": false,
            "accessCodeLifespan": 60,
            "accessCodeLifespanUserAction": 300,
            "accessCodeLifespanLogin": 1800
        }'

    print_status "Ustawienia realm zaktualizowane"
}

# Sprawdź czy wymagane narzędzia są zainstalowane
check_dependencies() {
    echo "🔍 Sprawdzanie zależności..."

    if ! command -v curl &> /dev/null; then
        print_error "curl nie jest zainstalowany"
        exit 1
    fi

    if ! command -v jq &> /dev/null; then
        print_warning "jq nie jest zainstalowany - instalowanie..."
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y jq
        elif command -v yum &> /dev/null; then
            sudo yum install -y jq
        elif command -v brew &> /dev/null; then
            brew install jq
        else
            print_error "Nie można zainstalować jq automatycznie. Zainstaluj ręcznie."
            exit 1
        fi
    fi

    print_status "Wszystkie zależności są dostępne"
}

# Wyexportuj konfigurację realm
export_realm_config() {
    echo "📤 Eksportowanie konfiguracji realm..."

    curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        > "media-vault-realm-export.json"

    print_status "Konfiguracja realm wyeksportowana do 'media-vault-realm-export.json'"
}

# Pokaż podsumowanie
show_summary() {
    echo ""
    echo "🎉 Konfiguracja Keycloak dla Media Vault zakończona!"
    echo ""
    echo "📋 Podsumowanie:"
    echo "├── Keycloak URL: $KEYCLOAK_URL"
    echo "├── Realm: $REALM_NAME"
    echo "├── Admin Console: $KEYCLOAK_URL/admin/"
    echo "└── User Account: $KEYCLOAK_URL/realms/$REALM_NAME/account/"
    echo ""
    echo "👥 Konta testowe:"
    echo "├── Admin: vaultadmin / admin123 (role: vault-admin, vault-user)"
    echo "└── User:  vaultuser / user123 (role: vault-user)"
    echo ""
    echo "📱 Klienci:"
    echo "├── media-vault-frontend (public)"
    echo "├── media-vault-api (confidential, secret: vault-api-secret-123)"
    echo "└── media-vault-admin (public)"
    echo ""
    echo "🔧 Następne kroki:"
    echo "1. Zaktualizuj .env file z odpowiednimi secretami"
    echo "2. Uruchom Media Vault: docker-compose up -d"
    echo "3. Otwórz http://localhost i zaloguj się"
    echo ""
    print_status "Keycloak jest gotowy do użycia!"
}

# Główna funkcja
main() {
    echo "🔐 Media Vault + Keycloak Setup"
    echo "================================"

    check_dependencies
    wait_for_keycloak
    get_admin_token
    create_realm
    create_roles
    create_clients
    create_users
    configure_realm_settings
    export_realm_config
    show_summary
}

# Sprawdź argumenty
if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    echo "Użycie: $0 [opcje]"
    echo ""
    echo "Opcje:"
    echo "  --help, -h     Pokaż tę pomoc"
    echo "  --clean        Wyczyść istniejący realm przed konfiguracją"
    echo "  --export-only  Tylko wyeksportuj konfigurację"
    echo ""
    echo "Zmienne środowiskowe:"
    echo "  KEYCLOAK_URL   URL do Keycloak (domyślnie: http://localhost:8443)"
    echo "  ADMIN_USER     Admin username (domyślnie: admin)"
    echo "  ADMIN_PASS     Admin password (domyślnie: admin123)"
    echo "  REALM_NAME     Nazwa realm (domyślnie: media-vault)"
    exit 0
fi

if [ "$1" == "--clean" ]; then
    echo "🧹 Czyszczenie istniejącego realm..."
    get_admin_token
    curl -s -X DELETE "$KEYCLOAK_URL/admin/realms/$REALM_NAME" \
        -H "Authorization: Bearer $ADMIN_TOKEN"
    print_status "Realm '$REALM_NAME' usunięty"
fi

if [ "$1" == "--export-only" ]; then
    get_admin_token
    export_realm_config
    exit 0
fi

# Uruchom główną funkcję
main

---

# update-env.sh - Automatyczne aktualizowanie .env z Keycloak
#!/bin/bash

echo "🔧 Aktualizowanie .env file dla Keycloak..."

ENV_FILE=".env"
BACKUP_FILE=".env.backup.$(date +%Y%m%d_%H%M%S)"

# Backup istniejącego .env
if [ -f "$ENV_FILE" ]; then
    cp "$ENV_FILE" "$BACKUP_FILE"
    echo "✅ Backup utworzony: $BACKUP_FILE"
fi

# Aktualizuj lub dodaj zmienne Keycloak
cat >> "$ENV_FILE" << 'EOF'

# Keycloak Configuration
KEYCLOAK_URL=http://keycloak:8080
KEYCLOAK_REALM=media-vault
KEYCLOAK_CLIENT_ID=media-vault-api
KEYCLOAK_CLIENT_SECRET=vault-api-secret-123
JWT_ISSUER=http://localhost:8443/realms/media-vault
JWT_AUDIENCE=media-vault-api
OAUTH2_ENABLED=true

# Frontend Keycloak Config
KEYCLOAK_FRONTEND_URL=http://localhost:8443
KEYCLOAK_FRONTEND_REALM=media-vault
KEYCLOAK_FRONTEND_CLIENT_ID=media-vault-frontend

# Admin Panel Keycloak Config
KEYCLOAK_ADMIN_CLIENT_ID=media-vault-admin
ADMIN_ROLE=vault-admin

# Service Account (dla analyzer)
SERVICE_ACCOUNT_CLIENT_ID=analyzer-service
SERVICE_ACCOUNT_SECRET=analyzer-secret-123
EOF

echo "✅ .env file zaktualizowany z konfiguracją Keycloak"

---

