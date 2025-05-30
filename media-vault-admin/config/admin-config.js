// Media Vault Admin Configuration
window.APP_CONFIG = {
  // API Configuration
  API_BASE_URL: process.env.API_URL || 'http://localhost:8080',
  
  // Authentication Configuration
  AUTH: {
    KEYCLOAK_URL: process.env.KEYCLOAK_URL || 'http://localhost:8080',
    REALM: process.env.KEYCLOAK_REALM || 'media-vault',
    CLIENT_ID: process.env.KEYCLOAK_CLIENT_ID || 'media-vault-admin',
    ADMIN_ROLE: process.env.ADMIN_ROLE || 'vault-admin'
  },
  
  // UI Configuration
  UI: {
    TITLE: process.env.VAULT_TITLE || 'Media Vault Admin',
    THEME: {
      PRIMARY_COLOR: '#4f46e5',
      SECONDARY_COLOR: '#6366f1',
      DANGER_COLOR: '#ef4444',
      SUCCESS_COLOR: '#10b981',
      WARNING_COLOR: '#f59e0b',
      INFO_COLOR: '#3b82f6'
    }
  },
  
  // Feature Flags
  FEATURES: {
    USER_MANAGEMENT: true,
    SYSTEM_MONITORING: true,
    VAULT_MANAGEMENT: true,
    SECURITY_SETTINGS: true
  },
  
  // Monitoring URLs
  MONITORING: {
    GRAFANA: process.env.MONITORING_GRAFANA_URL || 'http://localhost:3000',
    PROMETHEUS: process.env.MONITORING_PROMETHEUS_URL || 'http://localhost:9090'
  }
};
