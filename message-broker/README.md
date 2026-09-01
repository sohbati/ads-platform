# Message Broker

Event message broker service for the ads platform. The current implementation embeds a NATS server — no external broker installation is required. The service is named generically so the backing broker can be swapped later.

## Quick Start

1. **Setup**:
   ```bash
   make setup
   cp config.example .env   # if .env does not exist
   ```

2. **Run**:
   ```bash
   make run
   # or
   ./run.sh
   ```

3. **Access**:
   - HTTP API: `http://localhost:8095`
   - Health: `http://localhost:8095/health`
   - Broker panel: `http://localhost:8095/broker/`
   - Swagger: `http://localhost:8095/swagger/index.html`

## Project Structure

```
message-broker/
├── cmd/server/main.go              # Application entry point
├── go.mod
├── config.example
├── Makefile
├── run.sh
└── internal/
    ├── core/
    │   ├── config/                 # Configuration
    │   ├── container/              # Dependency injection
    │   ├── exception/              # App errors
    │   ├── middleware/             # HTTP middleware
    │   ├── broker/                 # Broker client connection
    │   ├── brokerserver/           # Embedded broker server
    │   └── router/                 # HTTP routing
    └── business/
        └── event/                  # Event publish business logic
            ├── container/
            ├── errorcode/
            ├── handler/
            ├── model/
            └── service/
```

## Configuration

| Variable              | Default       | Description                              |
|-----------------------|---------------|------------------------------------------|
| `PORT`                | `8095`        | HTTP API server port                     |
| `BROKER_HOST`         | `127.0.0.1`   | Embedded broker host                     |
| `BROKER_PORT`         | `-1`          | Embedded broker port (`-1` = auto-assign)|
| `BROKER_MONITOR_PORT` | `-1`          | Broker monitoring HTTP port (internal)   |

## API Endpoints

| Method | Path                          | Description                |
|--------|-------------------------------|----------------------------|
| GET    | `/health`                     | Health + broker status     |
| GET    | `/broker/`                    | Broker monitoring panel    |
| POST   | `/api/v1/events/:subject`     | Publish event to subject   |

**Publish example:**
```bash
curl -X POST http://localhost:8095/api/v1/events/user.registered \
  -H "Content-Type: application/json" \
  -d '{"data": {"userId": 1, "mobile": "09123456789"}}'
```

## Available Commands

- `make help` - Show all commands
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make setup` - Install dependencies
