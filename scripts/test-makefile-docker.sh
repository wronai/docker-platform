#!/bin/bash

# Exit on any error
set -e

# Define colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Create a temporary directory for testing
TEMP_DIR=$(mktemp -d)
echo -e "${YELLOW}📁 Created temporary directory: ${TEMP_DIR}${NC}"

# Function to clean up
trap_handler() {
    echo -e "\n${YELLOW}🧹 Cleaning up...${NC}"
    docker-compose -f ${TEMP_DIR}/docker-compose.yml down --volumes --remove-orphans >/dev/null 2>&1 || true
    # Don't remove the temp dir for debugging
    # rm -rf "${TEMP_DIR}"
    echo -e "${GREEN}✅ Cleanup complete${NC}
"
}
trap trap_handler EXIT

# Copy necessary files to temp directory
echo -e "${YELLOW}📦 Copying project files...${NC}"
cp Makefile "${TEMP_DIR}/"

# Create a minimal .env file if .env.example doesn't exist
if [ ! -f .env.example ]; then
    echo -e "${YELLOW}ℹ️  Creating minimal .env file...${NC}"
    echo "# Minimal .env file for testing" > "${TEMP_DIR}/.env"
    echo "COMPOSE_PROJECT_NAME=makefile-test" >> "${TEMP_DIR}/.env"
else
    cp .env.example "${TEMP_DIR}/.env"
fi

# Copy docker-compose files
cp -r docker-compose* "${TEMP_DIR}/" 2>/dev/null || true

# Copy Ansible files if they exist
if [ -d "ansible" ]; then
    echo -e "${YELLOW}📦 Copying Ansible files...${NC}"
    mkdir -p "${TEMP_DIR}/ansible"
    cp -r ansible/* "${TEMP_DIR}/ansible/"
    
    # Create a test inventory file
    cat > "${TEMP_DIR}/ansible/inventory.ini" << 'EOL'
[local]
localhost ansible_connection=local

[all:vars]
ansible_python_interpreter=/usr/bin/python3
EOL
    
    # Create a test playbook if none exists
    if [ ! -f "${TEMP_DIR}/ansible/playbook.yml" ]; then
        cat > "${TEMP_DIR}/ansible/playbook.yml" << 'EOL'
---
- name: Test Playbook
  hosts: localhost
  connection: local
  gather_facts: yes
  become: no

  tasks:
    - name: Check Docker is running
      command: docker info
      register: docker_info
      changed_when: false

    - name: Print Docker info
      debug:
        var: docker_info.stdout_lines
EOL
    fi
fi

# Create a test Dockerfile
cat > "${TEMP_DIR}/Dockerfile.test" << 'EOL'
FROM docker:24.0.7-cli

# Install dependencies
RUN apk add --no-cache \
    bash \
    curl \
    git \
    make \
    go \
    nodejs \
    npm \
    python3 \
    py3-pip \
    docker-compose \
    openssh-client \
    sshpass

# Install Go linter compatible with Go 1.21.10
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

# Create and activate virtual environment for Python packages
RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install Ansible and dependencies in the virtual environment
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir ansible ansible-lint docker && \
    ansible --version && \
    ansible-lint --version

# Set working directory
WORKDIR /app

# Copy project files
COPY . .

EOL

# Build test image
echo -e "${YELLOW}🐳 Building test Docker image...${NC}"
docker build -t makefile-test -f "${TEMP_DIR}/Dockerfile.test" "${TEMP_DIR}"

# Function to test Ansible playbook
test_ansible_playbook() {
    if [ ! -d "${TEMP_DIR}/ansible" ]; then
        echo -e "${YELLOW}ℹ️  No Ansible directory found, skipping Ansible tests${NC}"
        return 0
    fi
    
    local failed=0
    
    echo -e "\n${YELLOW}🔍 Testing Ansible playbook...${NC}"
    echo "--------------------------------------------------"
    
    # Test playbook syntax
    echo -e "${YELLOW}🔧 Checking playbook syntax...${NC}"
    if ! docker run --rm \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v "${TEMP_DIR}/ansible:/ansible" \
        -w /ansible \
        makefile-test \
        ansible-playbook -i inventory.ini --syntax-check playbook.yml; then
        echo -e "${RED}❌ Ansible playbook syntax check failed${NC}"
        ((failed++))
    else
        echo -e "${GREEN}✅ Ansible playbook syntax check passed${NC}"
    fi
    
    # Lint the playbook
    echo -e "\n${YELLOW}🔍 Linting playbook...${NC}"
    if ! docker run --rm \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v "${TEMP_DIR}/ansible:/ansible" \
        -w /ansible \
        makefile-test \
        ansible-lint playbook.yml; then
        echo -e "${YELLOW}⚠️  Ansible linting found issues (non-fatal)${NC}"
    else
        echo -e "${GREEN}✅ Ansible linting passed${NC}"
    fi
    
    # Run the playbook in check mode
    echo -e "\n${YELLOW}🔍 Running playbook in check mode...${NC}"
    if ! docker run --rm \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v "${TEMP_DIR}/ansible:/ansible" \
        -w /ansible \
        makefile-test \
        ansible-playbook -i inventory.ini --check playbook.yml; then
        echo -e "${RED}❌ Ansible playbook check mode failed${NC}"
        ((failed++))
    else
        echo -e "${GREEN}✅ Ansible playbook check mode passed${NC}"
    fi
    
    return $failed
}

# Function to run a make target and check its status
run_make_target() {
    local target=$1
    local description=$2
    
    echo -e "\n${YELLOW}🔄 Testing: ${target} - ${description}${NC}"
    echo "--------------------------------------------------"
    
    # Run the make target in the container
    if docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
        -v "${TEMP_DIR}:/app" \
        -w /app \
        makefile-test \
        make -n "${target}" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Dry run successful for: ${target}${NC}"
        return 0
    else
        echo -e "${RED}❌ Dry run failed for: ${target}${NC}"
        return 1
    fi
}

# Main test function
run_tests() {
    local failed_tests=0
    local total_tests=0
    
    # Define the targets to test with their descriptions
    declare -A targets=(
        ["init"]="Initialize development environment"
        ["up"]="Start all services"
        ["up-build"]="Rebuild and start all services"
        ["down"]="Stop all services"
        ["restart"]="Restart all services"
        ["logs"]="View logs from all services"
        ["test-unit"]="Run unit tests"
        ["test-integration"]="Run integration tests"
        ["lint"]="Run linters"
        ["format"]="Format code"
        ["monitoring"]="Start monitoring stack"
        ["monitoring-down"]="Stop monitoring stack"
        ["clean"]="Remove all containers and volumes"
    )
    
    # Run tests for each target
    for target in "${!targets[@]}"; do
        if ! run_make_target "$target" "${targets[$target]}"; then
            ((failed_tests++))
        fi
        ((total_tests++))
    done
    
    # Test Ansible playbook if it exists
    if [ -d "${TEMP_DIR}/ansible" ]; then
        echo -e "\n${YELLOW}🔍 Running Ansible tests...${NC}"
        ((total_tests++))
        if ! test_ansible_playbook; then
            ((failed_tests++))
        fi
    fi
    
    # Print summary
    echo -e "\n${YELLOW}📊 Test Summary:${NC}"
    echo "--------------------------------------------------"
    echo -e "${GREEN}✅ Passed: $((total_tests - failed_tests))${NC}"
    echo -e "${RED}❌ Failed: ${failed_tests}${NC}"
    echo -e "📊 Total:  ${total_tests}"
    
    # Exit with appropriate status
    if [ $failed_tests -gt 0 ]; then
        echo -e "\n❌ Some tests failed. Check the output above for details."
        echo -e "Temporary directory preserved at: ${TEMP_DIR} for debugging"
        return 1
    else
        echo -e "\n🎉 All tests passed successfully!"
        return 0
    fi
}

# Run the tests
run_tests
