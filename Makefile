# Extract database configuration from YAML using yq
DB_HOST := $(shell yq '.db.pg.write.host' config.yaml 2>/dev/null || echo "localhost")
DB_PORT := $(shell yq '.db.pg.write.port' config.yaml 2>/dev/null || echo "5432")
DB_NAME := $(shell yq '.db.pg.write.name' config.yaml 2>/dev/null || echo "kelasgo")
DB_USER := $(shell yq '.db.pg.write.user' config.yaml 2>/dev/null || echo "postgres")
DB_PASSWORD := $(shell yq '.db.pg.write.password' config.yaml 2>/dev/null || echo "")
DB_SSLMODE := $(shell yq '.db.pg.write.sslmode' config.yaml 2>/dev/null || echo "disable")

# Default target
.PHONY: dev setup build wire-gen clean test run check-config

BINARY=kelasgo-api
MIGRATION_STEP=1
DB_CONN_POSTGRES=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Check configuration source and yq availability
check-config:
	@echo "🔧 Configuration Status:"
	@if [ -f config.yaml ]; then \
		echo "✅ config.yaml found"; \
		if command -v yq >/dev/null 2>&1; then \
			echo "✅ yq is installed and available"; \
			echo "📊 Database Configuration:"; \
			echo "   Host: $(DB_HOST)"; \
			echo "   Port: $(DB_PORT)"; \
			echo "   Database: $(DB_NAME)"; \
			echo "   User: $(DB_USER)"; \
			echo "   SSL Mode: $(DB_SSLMODE)"; \
		else \
			echo "❌ yq not found. Install with: brew install yq (macOS) or check https://github.com/mikefarah/yq"; \
		fi \
	else \
		echo "❌ config.yaml not found. Please create it from config.example.yaml"; \
		echo "   Copy: cp config.example.yaml config.yaml"; \
		echo "   Edit: Update database credentials and other settings"; \
	fi
	@echo ""

# Help target - shows available commands
help:
	@echo "📖 Available commands:"
	@echo ""
	@echo "🏗️  Development:"
	@echo "  dev            - Start development server (auto-detects OS)"
	@echo "  build          - Build the application binary"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
	@echo "  check-config   - Check configuration source and database settings"
	@echo ""
	@echo "⚡ Wire (Dependency Injection):"
	@echo "  wire-gen       - Generate wire dependency injection code"
	@echo "  wire-force     - Force regenerate wire files"
	@echo "  install-wire   - Install Google Wire tool"
	@echo "  check-wire     - Check if Wire is installed"
	@echo ""
	@echo "🗃️  Database Migration:"
	@echo "  migrate_create - Create new migration file"
	@echo "  migrate_up     - Run migrations up"
	@echo "  migrate_down   - Run migrations down"
	@echo "  migrate_force  - Force migration version"
	@echo "  migrate_version- Show current migration version"
	@echo ""
	@echo "🧹 Maintenance:"
	@echo "  clean          - Remove built binaries and generated files"
	@echo "  macos-setup    - Setup macOS dependencies (symlinks)"
	@echo ""

# Wire generation target - generates dependency injection code
wire-gen:
	@echo "⚡ Generating wire dependency injection..."
	@go mod download
	@go generate ./...
	@echo "✅ Wire generation complete"

# Main dev target
dev: wire-gen
	@echo "🌐 Detecting OS..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		echo "🖥️  macOS detected"; \
		$(MAKE) macos-setup; \
		echo "🚀 Starting backend API (macOS)..."; \
		bash macos-dev.sh; \
	elif [ "$$(uname)" = "Linux" ]; then \
		echo "🐧 Linux detected"; \
		echo "🚀 Starting backend API (Linux)..."; \
		$$(go env GOPATH)/bin/air; \
	else \
		echo "🪟 Windows detected"; \
		echo "🚀 Starting backend API (Windows)..."; \
		$$(go env GOPATH)/bin/air; \
	fi

# Build target - compiles the application
build: wire-gen
	@echo "🔨 Building application..."
	@go build -o bin/${BINARY} .
	@echo "✅ Build complete: bin/${BINARY}"

# macOS setup target (runs symlink script once)
macos-setup:
	@if [ ! -L /usr/local/include/leptonica ] || [ ! -L /usr/local/include/tesseract ]; then \
		echo "🔗 Running macOS symlink setup..."; \
		bash macos-setup.sh; \
	else \
		echo "✅ macOS symlinks already set up"; \
	fi

migrate_create:
	@read -p "migration name (do not use space): " NAME \
  	&& migrate create -ext sql -dir ./database/migrations/postgres $${NAME}

migrate_up:
	@migrate -path ./database/migrations/postgres -database "$(DB_CONN_POSTGRES)" up $(MIGRATION_STEP)

migrate_down:
	@migrate -path ./database/migrations/postgres -database "$(DB_CONN_POSTGRES)" down $(MIGRATION_STEP)

migrate_force:
	@read -p "please enter the migration version (the migration filename prefix): " VERSION \
  	&& migrate -path ./database/migrations/postgres -database "$(DB_CONN_POSTGRES)" force $${VERSION}

migrate_version:
	@migrate -path ./database/migrations/postgres -database "$(DB_CONN_POSTGRES)" version

migrate_drop:
	@migrate -path ./database/migrations/postgres -database "$(DB_CONN_POSTGRES)" drop

# Run target - builds and runs the application
run: build
	@echo "🚀 Running application..."
	@./bin/${BINARY}

# Test target - runs tests with wire generation
test: wire-gen
	@echo "🧪 Running tests..."
	@go test ./...

# Clean target - removes built binaries and generated files
clean:
	@echo "🧹 Cleaning up..."
	@rm -f bin/${BINARY}
	@rm -f wire_gen.go
	@echo "✅ Cleanup complete"

# Install wire if not present
install-wire:
	@echo "⚡ Installing Google Wire..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo "✅ Wire installation complete"

# Check wire installation
check-wire:
	@test -f $$(go env GOPATH)/bin/wire || (echo "❌ Wire not found. Run 'make install-wire' first" && exit 1)
	@echo "✅ Wire is installed"

# Force wire regeneration (useful when wire files is corrupted)
wire-force:
	@echo "⚡ Force regenerating wire dependency injection..."
	@rm -f wire_gen.go
	@go mod download
	@go generate ./...
	@echo "✅ Wire force regeneration complete"

# Handle numeric arguments for migrate commands
%:
	@:
