include .env

# Default target
.PHONY: dev setup build wire-gen clean test run

BINARY=kelasgo-api
MIGRATION_STEP=1
DB_CONN_POSTGRES=postgres://$(DB.PG.WRITE.USER:"%"=%):$(DB.PG.WRITE.PASSWORD:"%"=%)@$(DB.PG.WRITE.HOST:"%"=%):$(DB.PG.WRITE.PORT:"%"=%)/$(DB.PG.WRITE.NAME:"%"=%)?sslmode=$(DB.PG.WRITE.SSLMODE:"%"=%)

# Help target - shows available commands
help:
	@echo "📖 Available commands:"
	@echo ""
	@echo "🏗️  Development:"
	@echo "  dev            - Start development server (auto-detects OS)"
	@echo "  build          - Build the application binary"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run tests"
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
wire-gen: check-wire
	@echo "⚡ Generating wire dependency injection..."
	@wire
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

# migrate_drop:
# 	@migrate -path ./database/migrations/postgres -database "$(DB_CONN_POSTGRES)" drop

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
	@which wire > /dev/null || (echo "❌ Wire not found. Run 'make install-wire' first" && exit 1)
	@echo "✅ Wire is installed"

# Force wire regeneration (useful when wire files are corrupted)
wire-force: check-wire
	@echo "⚡ Force regenerating wire dependency injection..."
	@rm -f wire_gen.go
	@wire
	@echo "✅ Wire force regeneration complete"

# Handle numeric arguments for migrate commands
%:
	@:
