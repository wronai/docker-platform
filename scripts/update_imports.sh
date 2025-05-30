#!/bin/bash

# Update imports in all Go files
find media-vault-backend -name "*.go" -type f -exec sed -i 's|media-vault-backend|github.com/wronai/media-vault-backend|g' {} \;

# Update the module path in go.mod
cat > media-vault-backend/go.mod << 'EOL'
module github.com/wronai/media-vault-backend

go 1.21

require (
    github.com/gofiber/fiber/v2 v2.50.0
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/google/uuid v1.3.1
    github.com/mattn/go-sqlite3 v1.14.17
)

require (
    github.com/andybalholm/brotli v1.0.5 // indirect
    github.com/klauspost/compress v1.16.7 // indirect
    github.com/mattn/go-colorable v0.1.13 // indirect
    github.com/mattn/go-isatty v0.0.19 // indirect
    github.com/mattn/go-runewidth v0.0.15 // indirect
    github.com/rivo/uniseg v0.2.0 // indirect
    github.com/valyala/bytebufferpool v1.0.0 // indirect
    github.com/valyala/fasthttp v1.50.0 // indirect
    github.com/valyala/tcplisten v1.0.0 // indirect
)
EOL

echo "Updated imports and module path"
