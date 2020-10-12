# PCAP Tracking System Makefile

.PHONY: help install dev build start stop clean test docker-build docker-up docker-down backend frontend

# Variables
BACKEND_DIR := backend
FRONTEND_DIR := frontend

# Help command
help:
	@echo "PCAP Tracking System - Available Commands"
	@echo "========================================"
	@echo ""
	@echo "Setup Commands:"
	@echo "  make install       - Install all dependencies (backend + frontend)"
	@echo "  make setup         - Complete setup (install + swagger docs)"
	@echo ""
	@echo "Development Commands:"
	@echo "  make dev           - Start both backend and frontend in development mode"
	@echo "  make backend       - Start only backend server"
	@echo "  make frontend      - Start only frontend server"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build         - Build both backend and frontend"
	@echo "  make build-backend - Build only backend"
	@echo "  make build-frontend- Build only frontend"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-up     - Start services with Docker Compose"
	@echo "  make docker-down   - Stop Docker services"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make clean         - Clean all build artifacts"
	@echo "  make test          - Run all tests"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo "  make lint          - Run linters"
	@echo "  make traffic       - Generate test network traffic"
	@echo ""
	@echo "API Documentation:"
	@echo "  Backend API: http://localhost:8080/swagger/index.html"
	@echo "  Frontend:    http://localhost:3000"

# Install all dependencies
install:
	@echo "Installing backend dependencies..."
	cd $(BACKEND_DIR) && go mod tidy
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install
	@echo "✅ All dependencies installed"

# Complete setup
setup: install swagger
	@echo "✅ Setup complete!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	cd $(BACKEND_DIR) && swag init -g cmd/main.go
	@echo "✅ Swagger documentation generated"

# Development mode - start both services
dev:
	@echo "Starting development servers..."
	@echo "Backend will start on http://localhost:8080"
	@echo "Frontend will start on http://localhost:3000"
	@echo "API Documentation: http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "Starting backend in background..."
	cd $(BACKEND_DIR) && go run cmd/main.go &
	@sleep 3
	@echo "Starting frontend..."
	cd $(FRONTEND_DIR) && npm run dev

# Start only backend
backend:
	@echo "Starting backend server on http://localhost:8080"
	@echo "API Documentation: http://localhost:8080/swagger/index.html"
	@echo "Auto-detecting active network interface..."
	cd $(BACKEND_DIR) && go run cmd/main.go

# Start only frontend
frontend:
	@echo "Starting frontend server on http://localhost:3000"
	cd $(FRONTEND_DIR) && npm run dev

# Build all
build: build-backend build-frontend

# Build backend
build-backend:
	@echo "Building backend..."
	cd $(BACKEND_DIR) && make build
	@echo "✅ Backend built successfully"

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && make build
	@echo "✅ Frontend built successfully"

# Clean all build artifacts
clean:
	@echo "Cleaning backend..."
	cd $(BACKEND_DIR) && make clean
	@echo "Cleaning frontend..."
	cd $(FRONTEND_DIR) && make clean
	@echo "✅ All build artifacts cleaned"

# Run tests
test:
	@echo "Running backend tests..."
	cd $(BACKEND_DIR) && make test
	@echo "Running frontend tests..."
	cd $(FRONTEND_DIR) && make test
	@echo "✅ All tests completed"

# Run linters
lint:
	@echo "Running backend linter..."
	cd $(BACKEND_DIR) && make fmt
	@echo "Running frontend linter..."
	cd $(FRONTEND_DIR) && make lint
	@echo "✅ Linting completed"

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build
	@echo "✅ Docker images built"

docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d
	@echo "✅ Services started"
	@echo "Frontend: http://localhost:3000"
	@echo "Backend API: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/swagger/index.html"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "✅ Services stopped"

# Stop all running processes
stop:
	@echo "Stopping all processes..."
	@pkill -f "go run cmd/main.go" || true
	@pkill -f "npm run dev" || true
	@echo "✅ All processes stopped"

# Quick start for development
quick-start: install swagger
	@echo "🚀 Quick starting PCAP Tracking System..."
	@echo "This will start both backend and frontend servers"
	@make dev

# Production mode
prod: build
	@echo "Starting in production mode..."
	cd $(BACKEND_DIR) && ./bin/pcap-tracking &
	cd $(FRONTEND_DIR) && npm start

# Generate test network traffic
traffic:
	@echo "🚀 Generating test network traffic..."
	./scripts/generate_traffic.sh
