#!/bin/bash

# Find all Go files and update imports
find /home/tom/github/wronai/docker-platform/media-vault-backend -type f -name "*.go" -exec sed -i 's|media-vault/internal/|github.com/wronai/media-vault-backend/internal/|g' {} \;

echo "Imports updated to use github.com/wronai/media-vault-backend"
