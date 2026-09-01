# ads-platform-notification

Notification and SMS service for the ads platform. Consumes OTP events from the message broker and sends SMS (log-only provider for now).

## Quick Start

1. **Start message broker** (`message-broker` on port 8095)
2. **Setup**:
   ```bash
   make setup
   cp config.example .env   # if .env does not exist
   ```
3. **Configure** `BROKER_HTTP_URL` (default `http://localhost:8095`) — broker URL is resolved automatically from broker `/health`
4. **Run**:
   ```bash
   make run
   # or
   ./run.sh
   ```

## Configuration

| Variable          | Default                  | Description                              |
|-------------------|--------------------------|------------------------------------------|
| `PORT`            | `8096`                   | HTTP server port                         |
| `BROKER_HTTP_URL` | `http://localhost:8095`  | Broker health URL for connection discovery |
| `BROKER_URL`      | *(auto from broker)*     | Optional broker URL override             |
| `OTP_SUBJECT`     | `notifications.otp.send` | Subject for OTP events                   |

## OTP Event Flow

1. `ads-platform-back` generates OTP and publishes to `notifications.otp.send`
2. This service receives the event and logs it
3. SMS is sent via `sms.Provider` (currently a log-only implementation)

**Event payload:**
```json
{ "mobile": "09123456789", "otp": "123456" }
```

## Project Structure

```
ads-platform-notification/
├── cmd/server/main.go
└── internal/
    ├── core/
    │   ├── config/
    │   ├── container/
    │   ├── broker/
    │   ├── middleware/
    │   └── router/
    └── business/otp/
        ├── container/
        ├── listener/      # OTP subscriber
        ├── model/
        ├── service/
        └── sms/           # SMS provider interface
```
