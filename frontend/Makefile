# Frontend Makefile

.PHONY: install dev build start lint test clean deps

# Variables
NODE_VERSION := $(shell node -v 2>/dev/null || echo "not installed")
NPM_VERSION := $(shell npm -v 2>/dev/null || echo "not installed")

# Check if Node.js is installed
check-node:
	@echo "Node.js version: $(NODE_VERSION)"
	@echo "npm version: $(NPM_VERSION)"
	@if [ "$(NODE_VERSION)" = "not installed" ]; then \
		echo "Error: Node.js is not installed"; \
		exit 1; \
	fi

# Install dependencies
install: check-node
	npm install

# Install dependencies (alias)
deps: install

# Start development server
dev: check-node
	npm run dev

# Build for production
build: check-node
	npm run build

# Start production server
start: check-node
	npm run start

# Lint code
lint: check-node
	npm run lint

# Fix linting issues
lint-fix: check-node
	npm run lint -- --fix

# Run tests
test: check-node
	npm run test

# Run tests in watch mode
test-watch: check-node
	npm run test -- --watch

# Type check
type-check: check-node
	npx tsc --noEmit

# Clean build artifacts and node_modules
clean:
	rm -rf .next
	rm -rf out
	rm -rf node_modules
	rm -rf .turbo

# Clean build artifacts only
clean-build:
	rm -rf .next
	rm -rf out

# Install and start development
setup: install dev

# Production build and start
prod: build start

# Show project info
info:
	@echo "Frontend Project Information:"
	@echo "============================="
	@echo "Node.js version: $(NODE_VERSION)"
	@echo "npm version: $(NPM_VERSION)"
	@echo "Project: PCAP Tracking Frontend"
	@echo "Framework: Next.js with TypeScript"
	@echo "Styling: Tailwind CSS"
	@echo ""
	@echo "Available commands:"
	@echo "  make install     - Install dependencies"
	@echo "  make dev         - Start development server"
	@echo "  make build       - Build for production"
	@echo "  make start       - Start production server"
	@echo "  make lint        - Run ESLint"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean all build artifacts"

# Docker commands
docker-build:
	docker build -t pcap-tracking-frontend .

docker-run:
	docker run -p 3000:3000 pcap-tracking-frontend

# Help command
help: info
