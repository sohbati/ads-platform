# NATS Message Broker

Event message broker service for the ads platform with an **embedded NATS server** — no external NATS installation required.

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
   - NATS panel: `http://localhost:8095/nats/`
   - Swagger: `http://localhost:8095/swagger/index.html`

## Project Structure

```
nats-message-broker/
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
    │   ├── natsconn/               # NATS client connection
    │   ├── natsserver/             # Embedded NATS server
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

| Variable         | Default       | Description                             |
|------------------|---------------|-----------------------------------------|
| `PORT`           | `8095`        | HTTP API server port                    |
| `NATS_HOST`      | `127.0.0.1`   | Embedded NATS server host               |
| `NATS_PORT`      | `-1`          | Embedded NATS port (`-1` = auto-assign) |
| `NATS_HTTP_PORT` | `-1`          | NATS monitoring HTTP port (internal)    |

## API Endpoints

| Method | Path                          | Description              |
|--------|-------------------------------|--------------------------|
| GET    | `/health`                     | Health + NATS status     |
| GET    | `/nats/`                      | NATS monitoring panel    |
| POST   | `/api/v1/events/:subject`     | Publish event to subject |

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
