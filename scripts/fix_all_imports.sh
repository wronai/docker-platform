#!/bin/bash

# Fix imports in all Go files
find media-vault-backend -name "*.go" -type f -exec sed -i 's|media-vault/internal/|github.com/wronai/media-vault-backend/internal/|g' {} \;
find media-vault-backend -name "*.go" -type f -exec sed -i 's|media-vault-backend/internal/|github.com/wronai/media-vault-backend/internal/|g' {} \;

echo "Fixed all import paths"
