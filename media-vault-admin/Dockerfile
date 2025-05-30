# Use nginx to serve static files
FROM nginx:alpine

# Create config directory
RUN mkdir -p /usr/share/nginx/html/config

# Copy the configuration file
COPY config/admin-config.js /usr/share/nginx/html/config/

# Copy the static files to the nginx web root
COPY src/ /usr/share/nginx/html/

# Copy the nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Set environment variables for the admin interface
ENV API_URL=http://media-vault-api:8080 \
    KEYCLOAK_URL=http://keycloak:8080 \
    KEYCLOAK_REALM=media-vault \
    KEYCLOAK_CLIENT_ID=media-vault-admin \
    ADMIN_ROLE=vault-admin \
    VAULT_TITLE="Media Vault Admin" \
    MONITORING_GRAFANA_URL=http://grafana:3000 \
    MONITORING_PROMETHEUS_URL=http://prometheus:9090

# Replace environment variables in the config file
RUN envsubst < /usr/share/nginx/html/config/admin-config.js > /usr/share/nginx/html/config/admin-config.tmp.js && \
    mv /usr/share/nginx/html/config/admin-config.tmp.js /usr/share/nginx/html/config/admin-config.js

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]