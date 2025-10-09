# Makefile
.PHONY: build test test-verbose test-coverage clean run help

# Build the application
build:
	@echo "Building LazyDnD..."
	@go build -o lazydnd

# Run all tests
test:
	@echo "Running tests..."
	@go test ./tests/...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./tests/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./tests/... -cover
	@echo ""
	@echo "Detailed coverage report:"
	@go test ./tests/... -coverprofile=coverage.out -coverpkg=./...
	@go tool cover -func=coverage.out

# Run tests and generate HTML coverage report
test-coverage-html:
	@echo "Generating HTML coverage report..."
	@go test ./tests/... -coverprofile=coverage.out -coverpkg=./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f lazydnd
	@rm -f coverage.out
	@rm -f coverage.html
	@rm -rf build/

# Run the application
run: build
	@./lazydnd

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@go vet ./...

# Run all quality checks
check: fmt lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  test               - Run all tests"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  test-coverage-html - Generate HTML coverage report"
	@echo "  clean              - Remove build artifacts"
	@echo "  run                - Build and run the application"
	@echo "  fmt                - Format code with gofmt"
	@echo "  lint               - Lint code with go vet"
	@echo "  check              - Run fmt, lint, and test"
	@echo "  help               - Show this help message"
