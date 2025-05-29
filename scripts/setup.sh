#!/bin/bash

# setup.sh - Complete Media Vault Setup Script
set -e

echo "🔐 Media Vault Complete Setup"
echo "============================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_status() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

# Check dependencies
check_dependencies() {
    echo "🔍 Checking dependencies..."

    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi

    print_status "All dependencies are available"
}

# Create directory structure
create_directories() {
    echo "📁 Creating directory structure..."

    directories=(
        "data"
        "uploads/originals"
        "uploads/thumbnails"
        "uploads/processed"
        "processing/incoming"
        "processing/temp"
        "backups"
        "logs"
        "models"
    )

    for dir in "${directories[@]}"; do
        mkdir -p "$dir"