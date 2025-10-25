#!/bin/bash

# Subspace Backend - Automated Quick Start Script
# This script will get your backend running with minimal effort

set -e  # Exit on error

echo "🚀 Subspace Backend Quick Start"
echo "================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_success() { echo -e "${GREEN}✓${NC} $1"; }
print_error() { echo -e "${RED}✗${NC} $1"; }
print_info() { echo -e "${YELLOW}ℹ${NC} $1"; }

# Check if .env exists
if [ ! -f .env ]; then
    print_info "Creating .env file from .env.example..."
    cp .env.example .env
    print_success ".env file created"
else
    print_success ".env file already exists"
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed"
    echo "Please install Docker Desktop from: https://www.docker.com/products/docker-desktop"
    exit 1
fi
print_success "Docker is installed"

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    print_info "Docker daemon is not running. Starting Docker Desktop..."

    # Try to start Docker Desktop (macOS)
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open -a Docker
        print_info "Waiting for Docker to start (this may take 20-30 seconds)..."

        # Wait for Docker to be ready (max 60 seconds)
        counter=0
        while ! docker info &> /dev/null && [ $counter -lt 30 ]; do
            sleep 2
            counter=$((counter + 1))
            echo -n "."
        done
        echo ""

        if ! docker info &> /dev/null; then
            print_error "Docker failed to start. Please start Docker Desktop manually."
            exit 1
        fi
    else
        print_error "Please start Docker manually and run this script again"
        exit 1
    fi
fi
print_success "Docker daemon is running"

# Check for docker compose
HAS_COMPOSE=false
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
    HAS_COMPOSE=true
    print_success "Docker Compose is available"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
    HAS_COMPOSE=true
    print_success "docker-compose is available"
else
    print_error "Docker Compose is not available"
    print_info "Installing Docker Compose..."
    brew install docker-compose
    COMPOSE_CMD="docker-compose"
    HAS_COMPOSE=true
fi

# Start services
echo ""
print_info "Starting PostgreSQL and API server..."
$COMPOSE_CMD up -d

# Wait for services to be healthy
print_info "Waiting for services to be ready..."
sleep 10

# Check if services are running
if docker ps | grep -q subspace-db && docker ps | grep -q subspace-api; then
    print_success "All services are running!"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✨ Backend is ready at: http://localhost:8080"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    # Test health endpoint
    print_info "Testing health endpoint..."
    if curl -s http://localhost:8080/health | grep -q "healthy"; then
        print_success "Health check passed!"
        echo ""
        echo "Quick test commands:"
        echo "  • Health:   curl http://localhost:8080/health"
        echo "  • Register: curl -X POST http://localhost:8080/api/v1/auth/register \\"
        echo "                -H 'Content-Type: application/json' \\"
        echo "                -d '{\"name\":\"Test\",\"email\":\"test@example.com\",\"password\":\"password123\"}'"
        echo ""
        echo "View logs:    $COMPOSE_CMD logs -f"
        echo "Stop:         $COMPOSE_CMD down"
        echo ""
        echo "Test credentials (from init.sql):"
        echo "  • admin@subspace.dev / admin123"
        echo "  • test@subspace.dev / admin123"
    else
        print_error "Health check failed. Check logs with: $COMPOSE_CMD logs"
    fi
else
    print_error "Some services failed to start. Check logs with: $COMPOSE_CMD logs"
    exit 1
fi
