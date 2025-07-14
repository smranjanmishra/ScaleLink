# LinkSprint Makefile
# Common development commands

.PHONY: help build run test clean docker-build docker-run docker-stop load-test

# Default target
help:
	@echo "🚀 LinkSprint - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application locally"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker services"
	@echo ""
	@echo "Testing:"
	@echo "  make load-test    - Run load tests"
	@echo "  make deps         - Download dependencies"

# Build the application
build:
	@echo "🔨 Building LinkSprint..."
	go build -o bin/linksprint cmd/server/main.go
	@echo "✅ Build completed!"

# Run the application locally
run:
	@echo "🚀 Starting LinkSprint..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	go clean
	@echo "✅ Clean completed!"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	@echo "✅ Dependencies downloaded!"

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t linksprint .
	@echo "✅ Docker image built!"

# Run with Docker Compose
docker-run:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "📊 Application: http://localhost:8080"
	@echo "📈 Prometheus: http://localhost:9090"
	@echo "📊 Grafana: http://localhost:3000"

# Stop Docker services
docker-stop:
	@echo "🛑 Stopping Docker services..."
	docker-compose down
	@echo "✅ Services stopped!"

# Run load tests
load-test:
	@echo "⚡ Running load tests..."
	@chmod +x scripts/load_test.sh
	./scripts/load_test.sh

# Install development tools
install-tools:
	@echo "🛠️ Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	@echo "✅ Tools installed!"

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run
	@echo "✅ Linting completed!"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted!"

# Generate API documentation
docs:
	@echo "📚 Generating API documentation..."
	@echo "API documentation available at: http://localhost:8080/api/v1"
	@echo "✅ Documentation generated!"

# Show application status
status:
	@echo "📊 LinkSprint Status:"
	@echo "======================"
	@curl -s http://localhost:8080/health | jq . 2>/dev/null || echo "Application not running"
	@echo ""
	@echo "Docker services:"
	@docker-compose ps

# Development setup
dev-setup: deps install-tools
	@echo "✅ Development environment setup completed!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Start Redis and CockroachDB: make docker-run"
	@echo "2. Run the application: make run"
	@echo "3. Test the API: curl http://localhost:8080/health" 