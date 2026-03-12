#!/bin/bash

# Flashduty Provider Test Script
# Usage: ./scripts/test.sh [unit|acc|manual|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_header() {
    echo ""
    echo "============================================"
    echo -e "${GREEN}$1${NC}"
    echo "============================================"
}

print_error() {
    echo -e "${RED}ERROR: $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}WARNING: $1${NC}"
}

check_api_key() {
    if [ -z "$FLASHDUTY_API_KEY" ]; then
        print_error "FLASHDUTY_API_KEY environment variable is not set"
        echo "Please set it: export FLASHDUTY_API_KEY='your-api-key'"
        exit 1
    fi
    echo "✓ FLASHDUTY_API_KEY is set"
}

run_unit_tests() {
    print_header "Running Unit Tests"
    cd "$PROJECT_DIR"
    go test -v -cover -timeout=120s -parallel=10 ./...
}

run_acceptance_tests() {
    print_header "Running Acceptance Tests"
    check_api_key
    cd "$PROJECT_DIR"
    TF_ACC=1 go test -v -cover -timeout 120m ./internal/provider/...
}

run_single_acceptance_test() {
    print_header "Running Single Acceptance Test: $1"
    check_api_key
    cd "$PROJECT_DIR"
    TF_ACC=1 go test -v -timeout 30m -run "$1" ./internal/provider/...
}

install_provider() {
    print_header "Installing Provider Locally"
    
    cd "$PROJECT_DIR"
    
    # Build the provider
    echo "Building provider..."
    go build -o terraform-provider-flashduty .
    
    # Copy to Go bin
    GO_BIN="${GOPATH:-$HOME/go}/bin"
    mkdir -p "$GO_BIN"
    cp terraform-provider-flashduty "$GO_BIN/"
    echo "✓ Provider installed to $GO_BIN"
}

build_provider() {
    print_header "Building Provider"
    cd "$PROJECT_DIR"
    go build -v -o terraform-provider-flashduty .
    echo "✓ Build successful"
}

show_help() {
    echo "Flashduty Provider Test Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build    Build the provider"
    echo "  unit     Run unit tests (no API key required)"
    echo "  acc      Run all acceptance tests (requires API key)"
    echo "  test     Run a single acceptance test (requires test name)"
    echo "  install  Build and install provider to GOPATH/bin"
    echo "  all      Run build + unit + acceptance tests"
    echo "  help     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 build"
    echo "  $0 unit"
    echo "  $0 acc"
    echo "  $0 test TestAccTeamResource"
    echo "  $0 install"
}

case "${1:-help}" in
    build)
        build_provider
        ;;
    unit)
        run_unit_tests
        ;;
    acc)
        run_acceptance_tests
        ;;
    test)
        if [ -z "$2" ]; then
            print_error "Please specify test name"
            echo "Example: $0 test TestAccTeamResource"
            exit 1
        fi
        run_single_acceptance_test "$2"
        ;;
    install)
        install_provider
        ;;
    all)
        build_provider
        run_unit_tests
        run_acceptance_tests
        ;;
    help|*)
        show_help
        ;;
esac
