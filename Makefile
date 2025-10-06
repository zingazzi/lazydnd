# Makefile for LazyDnD

.PHONY: help build run clean docker-build docker-run docker-compose-up docker-compose-down docker-push test install

# Default target
help:
	@echo "LazyDnD - Makefile Commands"
	@echo ""
	@echo "Native Build:"
	@echo "  make build          - Build native binary"
	@echo "  make run            - Run native binary"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install        - Install to /usr/local/bin"
	@echo "  make test           - Run tests"
	@echo ""
	@echo "Cross-Platform Build:"
	@echo "  make build-all      - Build for all platforms"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"
	@echo "  make docker-compose-up   - Start with docker-compose"
	@echo "  make docker-compose-down - Stop docker-compose"
	@echo "  make docker-push    - Push to GHCR"
	@echo ""

# Native build
build:
	@echo "Building LazyDnD..."
	go build -ldflags "-s -w" -o lazydnd .
	@echo "✓ Build complete: ./lazydnd"

# Run native binary
run: build
	./lazydnd

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f lazydnd
	rm -rf build/
	@echo "✓ Clean complete"

# Install to system
install: build
	@echo "Installing to /usr/local/bin..."
	sudo mv lazydnd /usr/local/bin/
	@echo "✓ Installed: /usr/local/bin/lazydnd"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	./build.sh

# Docker build
docker-build:
	@echo "Building Docker image..."
	./docker-build.sh

# Docker run
docker-run:
	@echo "Running Docker container..."
	./docker-run.sh

# Docker compose up
docker-compose-up:
	@echo "Starting with docker-compose..."
	docker-compose up

# Docker compose down
docker-compose-down:
	@echo "Stopping docker-compose..."
	docker-compose down

# Push to GHCR
docker-push: docker-build
	@echo "Pushing to GitHub Container Registry..."
	@read -p "Enter version tag (e.g., v1.0.1): " VERSION; \
	docker tag lazydnd:latest ghcr.io/zingazzi/lazydnd:$$VERSION && \
	docker tag lazydnd:latest ghcr.io/zingazzi/lazydnd:latest && \
	docker push ghcr.io/zingazzi/lazydnd:$$VERSION && \
	docker push ghcr.io/zingazzi/lazydnd:latest
	@echo "✓ Pushed to GHCR"

# Development helpers
dev:
	@echo "Running in development mode..."
	go run .

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Formatted"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "✓ Linted"

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download
	@echo "✓ Dependencies updated"

