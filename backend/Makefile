# Backend Makefile

.PHONY: build run clean test deps install-deps dev

# Variables
BINARY_NAME=pcap-tracking
MAIN_PATH=./cmd/main.go

# Build the application
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH)

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install development dependencies
install-deps:
	go install github.com/air-verse/air@latest

# Run in development mode with hot reload
dev:
	air

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux $(MAIN_PATH)

# Build for macOS
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-macos $(MAIN_PATH)

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe $(MAIN_PATH)
