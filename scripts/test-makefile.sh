#!/bin/bash

# Exit on any error
set -e

echo "🚀 Starting Makefile test in Docker environment..."

# Define colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to run a make target and check its status
run_make_target() {
    local target=$1
    local description=$2
    
    echo -e "\n${YELLOW}🔄 Testing: ${target} - ${description}${NC}"
    echo "--------------------------------------------------"
    
    # Run the make target and capture output and status
    if make -n $target >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Dry run successful for: ${target}${NC}"
        return 0
    else
        echo -e "${RED}❌ Dry run failed for: ${target}${NC}"
        return 1
    fi
}

# Main test function
test_makefile() {
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
    
    # Print summary
    echo -e "\n${YELLOW}📊 Test Summary:${NC}"
    echo "--------------------------------------------------"
    echo -e "${GREEN}✅ Passed: $((total_tests - failed_tests))${NC}"
    echo -e "${RED}❌ Failed: ${failed_tests}${NC}"
    echo -e "📊 Total:  ${total_tests}"
    
    # Exit with appropriate status
    if [ $failed_tests -gt 0 ]; then
        echo -e "\n❌ Some tests failed. Check the output above for details."
        exit 1
    else
        echo -e "\n🎉 All tests passed successfully!"
        exit 0
    fi
}

# Run the tests
test_makefile
