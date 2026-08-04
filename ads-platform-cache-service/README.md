# ads platform Service

Business Logic Service for the ads platform Application.

## Quick Start

1. **Setup**:
   ```bash
   make setup
   ```

2. **Run**:
   ```bash
   make run
   ```

3. **Access**:
   - Service: `http://localhost:8093`
   - Health: `http://localhost:8093/health`
   - Swagger: `http://localhost:8093/swagger/index.html`

## Project Structure

```
ads-platform/
├── cmd/main.go                    # Alternative entry point
├── main.go                        # Main application
├── go.mod                         # Dependencies
├── config.example                 # Configuration template
├── Makefile                       # Build commands
├── docker-compose.yml             # Database setup
├── internal/
│   ├── config/config.go          # Configuration
│   ├── db/db.go                   # Database connection
│   ├── router/router.go           # HTTP routing
│   ├── programs/                  # Program management
│   ├── exercises/                 # Exercise management
│   ├── workouts/                  # Workout management
│   ├── sessions/                  # Session tracking
│   └── analytics/                 # Analytics
├── database/
│   ├── mysql/0_init.sql          # MySQL setup
│   └── postgres/0_init.sql        # PostgreSQL setup
└── docs/
    └── docs.go                    # Swagger docs
```

## Available Commands

- `make help` - Show all commands
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make swag` - Generate Swagger docs
- `make docker-up` - Start MySQL
- `make docker-down` - Stop services
- `make dev` - Start development environment
- `make setup` - Complete setup

## API Endpoints

All endpoints return "Not implemented yet" - ready for development:

- `/api/v1/programs` - Program management
- `/api/v1/exercises` - Exercise management  
- `/api/v1/workouts` - Workout management
- `/api/v1/sessions` - Session tracking
- `/api/v1/analytics` - Analytics and reporting

