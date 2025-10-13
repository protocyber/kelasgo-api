include .env

# Default target
.PHONY: dev setup build

BINARY=kelasgo-api
MIGRATION_STEP=1
DB_CONN_POSTGRES=postgres://$(DB.PG.WRITE.USER:"%"=%):$(DB.PG.WRITE.PASSWORD:"%"=%)@$(DB.PG.WRITE.HOST:"%"=%):$(DB.PG.WRITE.PORT:"%"=%)/$(DB.PG.WRITE.NAME:"%"=%)?sslmode=$(DB.PG.WRITE.SSLMODE:"%"=%)

# Main dev target
dev:
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
build:
	@echo "🔨 Building application..."
	@go build -o bin/${BINARY} main.go
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

# Handle numeric arguments for migrate commands
%:
	@:
