# kelasgo-api

## Description

This is a Go-based REST API service for educational management system (KelasGo).
The service provides multi-tenant architecture with comprehensive user, student, and class management capabilities.

## Technical Information

### Architecture

This is a Go-based REST API service built with the following architecture:

- **Language**: Go 1.21+
- **Database**: PostgreSQL with read/write splitting support
- **Dependency Injection**: Google Wire for compile-time DI
- **Configuration**: YAML-based configuration with environment variable overrides
- **Hot Reloading**: Air for development
- **Migration**: golang-migrate for database schema management

### Project Structure

```text
├── internal/
│   ├── config/          # Configuration management
│   ├── handler/         # HTTP handlers/controllers
│   ├── service/         # Business logic layer
│   ├── repository/      # Data access layer
│   ├── model/           # Data models/entities
│   ├── dto/             # Data transfer objects
│   ├── middleware/      # HTTP middleware
│   └── util/            # Utility functions
├── database/
│   └── migrations/      # Database migration files
├── bin/                 # Compiled binaries
├── config.yaml          # Configuration file
├── wire.go             # Wire dependency injection setup
├── main.go             # Application entry point
└── Makefile            # Build and development commands
```

### API Endpoints

The service provides REST API endpoints for:

- **Authentication**: User login/logout, JWT token management
- **User Management**: CRUD operations for users
- **Student Management**: Student-related operations
- **Multi-tenancy**: Tenant-based data isolation

### Dependencies

Key Go modules used in this project:

- **Web Framework**: Standard `net/http` with custom routing
- **Database**: PostgreSQL driver
- **Configuration**: `spf13/viper` for YAML config management
- **Migration**: `golang-migrate/migrate` for database migrations
- **Dependency Injection**: `google/wire` for compile-time DI

### Development Guidelines

- Follow Go naming conventions and best practices
- Use dependency injection for loose coupling
- Implement proper error handling and logging
- Write unit tests for business logic
- Use migrations for all database schema changes
- Keep configuration in YAML files, not hardcoded values

## Getting Started

### Setup

1. Clone repo
1. Copy and paste `config.example.yaml` to `config.yaml`. Set your database and other configurations
1. Make sure you install the required tools:
   - [go-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) - for database migrations
   - [yq](https://github.com/mikefarah/yq#install) - for YAML processing in Makefile
   - [Google Wire](https://github.com/google/wire) - for dependency injection (run `make install-wire`)
   - [Air](https://github.com/air-verse/air) - for hot reloading during development (`go install github.com/air-verse/air@latest`)
1. Verify your configuration: `make check-config`
1. Generate dependency injection code: `make wire-gen`
1. Run the migration: `make migrate_up`
1. Run the service: `make dev`

### Configuration

This project uses YAML-based configuration. The configuration file structure follows this format:

```yaml
# Server settings
server:
  port: "8080"
  env: "development"
  log_level: "debug"

# Database configuration
db:
  pg:
    read:
      host: "localhost"
      port: "5432"
      name: "kelasgo"
      user: "postgres"
      password: "your_password_here"
      sslmode: "disable"
    write:
      host: "localhost" 
      port: "5432"
      name: "kelasgo"
      user: "postgres"
      password: "your_password_here"
      sslmode: "disable"

# Other configurations (CORS, Redis, S3, etc.)
app:
  cors:
    enable: true
    allowed_origins: "http://localhost:8080"
```

**Configuration Commands:**

| Command | Description |
| --- | --- |
| `make check-config` | Verify your configuration and show current database settings |
| `cp config.example.yaml config.yaml` | Create your configuration file from the example |

**Environment Variable Overrides:**

You can still override any configuration using environment variables with dot notation converted to underscores:

- `DB_PG_WRITE_HOST=newhost` overrides `db.pg.write.host`
- `SERVER_PORT=9000` overrides `server.port`

### Wire Dependency Injection

This project uses [Google Wire](https://github.com/google/wire) for compile-time dependency injection. Wire generates the dependency injection code automatically based on your provider functions.

**Wire Commands:**

| Command | Description |
| --- | --- |
| `make install-wire` | Install Google Wire tool |
| `make check-wire` | Verify Wire is installed |
| `make wire-gen` | Generate dependency injection code (auto-runs before build/dev) |
| `make wire-force` | Force regenerate Wire files (useful when corrupted) |

**Important Notes:**

- Wire code generation is automatically triggered when you run `make dev`, `make build`, or `make test`
- The generated `wire_gen.go` file should not be manually edited
- If you modify dependency providers in `wire.go`, run `make wire-gen` to regenerate the code
- Wire files are cleaned up with `make clean`

More info: [Google Wire Documentation](https://github.com/google/wire)

### Development Tools

This project uses several development tools to enhance the development experience:

#### Air (Hot Reloading)

[Air](https://github.com/air-verse/air) provides live reloading for Go applications during development.

**Installation:**

```bash
go install github.com/air-verse/air@latest
```

**Usage:**

- Air is automatically used when you run `make dev` on Linux/Windows
- Watches for file changes and automatically rebuilds/restarts the application
- Configuration can be customized in `.air.toml` (if present)

#### Development Commands

| Command | Description |
| --- | --- |
| `make dev` | Start development server with hot reloading (auto-detects OS) |
| `make build` | Build the application binary |
| `make run` | Build and run the application (no hot reloading) |
| `make test` | Run tests with Wire generation |
| `make clean` | Remove built binaries and generated files |

**Platform-specific behavior:**

- **Linux/Windows**: Uses Air for hot reloading
- **macOS**: Runs additional setup scripts and uses Air

### Database migration

As you know above, we are using [go-migrate](https://github.com/golang-migrate/migrate) for the migration tool. We have simplify the frequently executed command into the Makefile.

**Note: Make sure you configure your database settings in `config.yaml`:**

```yaml
db:
  pg:
    write:
      host: "localhost"
      port: "5432"  
      name: "kelasgo"
      user: "postgres"
      password: "your_password_here"
      sslmode: "disable"
```

The migration commands will automatically read these settings from your `config.yaml` file using `yq`.

| Command | Description |
| --- | --- |
| `make check-config` | Verify configuration and database settings |
| `make migrate_create` | Create migration file |
| `make migrate_up` | Apply all or N up migrations. If you want to specify the step, use: `make migrate_up MIGRATION_STEP=<number>` |
| `make migrate_down` | Apply all or N down migrations. If you want to specify the step, use: `make migrate_down MIGRATION_STEP=<number>` |
| `make migrate_force` | Set version V but don't run migration (fix the dirty state) |
| `make migrate_version` | Print current migration version |
| `make migrate_drop` | Drop everything inside database |

More info: [go-migrate documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

## Management

### Communication Channels

For technical questions or contributions, please contact the project maintainer directly.

### Product Owner & Engineering Manager

For product and engineering management inquiries, please contact the project maintainer.

## Maintainers

Please contact the person below for technical inquiries, feature requests, or contributions.

### Project Maintainer & Fullstack Developer

#### M. Fitrah Muttaqin

- Email: [fitrah.pro@gmail.com](mailto:fitrah.pro@gmail.com)
- Role: Fullstack Developer & Project Maintainer
- Responsibilities:
  - Backend API development (Go)
  - Frontend development
  - Database design and migrations
  - DevOps and deployment
  - Code review and project management
