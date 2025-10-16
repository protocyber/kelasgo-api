# Simple, clean Makefile for kelasgo-api
# Complex logic moved to scripts/ for better maintainability

# Default target
.PHONY: dev dev-env build clean test run check-config help migrate-config migrate-create migrate-up migrate-down migrate-force migrate-version migrate-drop

BINARY=kelasgo-api

# Check configuration and dependencies
check-config:
	@./scripts/dev-server.sh check

# Help target - shows available commands
help:
	@echo "📖 Available commands:"
	@echo ""
	@echo "🏗️  Development:"
	@echo "  dev            - Start development server (auto-detects OS)"
	@echo "  dev-env        - Show development environment info"
	@echo "  build          - Build the application binary"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
	@echo "  check-config   - Check configuration and dependencies"
	@echo ""
	@echo "🗃️  Database Migration:"
	@echo "  migrate-config - Show database configuration" 
	@echo "  migrate-create - Create new migration file"
	@echo "  migrate-up     - Run migrations up"
	@echo "  migrate-down   - Run migrations down"
	@echo "  migrate-force  - Force migration version"
	@echo "  migrate-version- Show current migration version"
	@echo "  migrate-drop   - Drop database (WARNING: destructive)"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "  clean          - Remove built binaries and generated files"
	@echo ""

# Start development server (OS-aware)
dev:
	@./scripts/dev-server.sh start

# Build target - compiles the application
build:
	@echo "🔨 Building application..."
	@go build -o bin/${BINARY} ./cmd/kelasgo-api
	@echo "✅ Build complete: bin/${BINARY}"

# Show development environment info
dev-env:
	@./scripts/dev-server.sh env

# Database migration targets (simplified)
migrate-config:
	@./scripts/db-migrate.sh config

migrate-create:
	@read -p "Migration name (no spaces): " NAME && ./scripts/db-migrate.sh create $$NAME

migrate-up:
	@./scripts/db-migrate.sh up

migrate-down:
	@./scripts/db-migrate.sh down

migrate-force:
	@read -p "Migration version: " VERSION && ./scripts/db-migrate.sh force $$VERSION

migrate-version:
	@./scripts/db-migrate.sh version

migrate-drop:
	@./scripts/db-migrate.sh drop

# Run target - builds and runs the application
run: build
	@echo "🚀 Running application..."
	@./bin/${BINARY}

# Test target - runs tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Clean target - removes built binaries and generated files
clean:
	@echo "🧹 Cleaning up..."
	@rm -f bin/${BINARY}
	@rm -rf tmp/
	@echo "✅ Cleanup complete"

# Handle numeric arguments for migrate commands
%:
	@:
