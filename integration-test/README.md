# integration-test

Integration tests for the ads platform using [Testcontainers for Go](https://golang.testcontainers.org/).

## Prerequisites

- Docker Desktop or Docker Engine running
- Go 1.25+

## What it spins up

| Container | Image / Build | Port |
|-----------|---------------|------|
| PostgreSQL | `postgres:16-alpine` | 5432 |
| ads-platform-back | Dockerfile | 8092 |
| ads-platform-cache-service | Dockerfile | 8093 |
| ads-platform-cdn | Dockerfile | 4000 |
| ads-platform-ui | Dockerfile | 8094 |
| nats-message-broker | Dockerfile | 8095 |
| ads-platform-notification | Dockerfile | 8096 |

All microservices are built from this repo and run on a shared Docker network.

## Run

```bash
cd integration-test
make tidy
make test-integration
```

Or from repo root:

```bash
./7-integration-test.sh
```

## Project layout

```
integration-test/
├── docker/                         # Dockerfiles for each microservice
├── fixtures/postgres/init.sql      # DB schema for tests
├── internal/testcontainers/
│   ├── postgres.go                 # PostgreSQL testcontainer
│   ├── services.go                 # Microservice testcontainers
│   └── stack.go                    # Full stack orchestration
└── tests/
    └── stack_test.go               # Integration tests
```

## Notes

- First run builds Docker images and can take several minutes.
- Use `go test -tags=integration -short` to skip integration tests.
- Set `-timeout=20m` for full stack tests.
