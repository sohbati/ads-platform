# ads-platform-notification

Notification and SMS service for the ads platform. Consumes OTP events from the message broker and sends SMS through a pluggable provider (`log` for development, `smsir` for sms.ir).

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
| `APP_ENV`         | `local`                  | `local`, `uat` or `production`           |
| `PORT`            | `8096`                   | HTTP server port                         |
| `BROKER_HTTP_URL` | `http://localhost:8095`  | Broker health URL for connection discovery |
| `BROKER_URL`      | *(auto from broker)*     | Optional broker URL override             |
| `OTP_SUBJECT`     | `notifications.otp.send` | Subject for OTP events                   |
| `SMS_PROVIDER`    | `log`                    | `log` (only logs) or `smsir`             |
| `SMS_ALLOWLIST`   |                          | Comma-separated numbers allowed to receive SMS outside production. Empty = nobody |
| `SMSIR_API_KEY`   |                          | sms.ir API key (Sandbox or production)   |
| `SMSIR_BASE_URL`  | `https://api.sms.ir`     | sms.ir API host (sandbox uses the same host) |
| `SMSIR_OTP_TEMPLATE_ID` | `123456`           | Verify template for OTP (123456 is the sandbox default) |
| `SMSIR_OTP_CODE_PARAM`  | `Code`             | Template parameter that holds the code   |
| `SMSIR_LINE_NUMBER`     |                    | Sender line, only for free-text messages |

## SMS Providers

Business code depends only on `internal/core/sms.Sender`. It sends a **logical template**
(`sms.TemplateOTP` with param `sms.ParamCode`) and never a vendor template ID. Each provider
maps logical templates to its own IDs and parameter names through config, so switching
vendors only changes `SMS_PROVIDER` and that vendor's env vars.

### Environments and the recipient allowlist

| | local | uat | production |
|---|---|---|---|
| `APP_ENV` | `local` (default) | `uat` | `production` |
| Provider | `log`, or `smsir` + sandbox key | `smsir` + live key (real SMS, same as production) | `smsir` + live key |
| `SMS_ALLOWLIST` | your own number | testers' numbers | ignored |

Unless `APP_ENV=production` is set explicitly, every message goes through
`internal/core/sms/allowlist`. Only numbers in `SMS_ALLOWLIST` reach the provider. Other
numbers are blocked and logged as `[SMS:allowlist] BLOCKED`. Numbers are normalized before
the check, so `0912…`, `+98912…`, `98912…`, `912…` and Persian digits all match. The guard
fails closed: a missing `APP_ENV` or an empty list blocks every message, and an invalid entry
stops the service at startup. The startup log shows the active setup, for example
`env=uat sms=smsir+allowlist allowlist=3 numbers`.

**Testing with the sms.ir sandbox:** in the sms.ir panel, go to Developers → API keys → New key
and choose type *Sandbox*. Set `SMS_PROVIDER=smsir` and `SMSIR_API_KEY=<sandbox key>`. Requests
go to the real API but no SMS is delivered and no credit is used. Only template `123456` / `Code` works there.

**Adding a provider** (for example Kavenegar):
1. Create `internal/core/sms/provider/<name>/` that implements `sms.Sender`
   (`SendTemplate` for OTP/lookup templates, `SendText` for free text).
2. Return `*sms.ProviderError` for vendor failures, and set `Retryable` for transient ones.
3. Add its env vars to `config.SMSConfig` and a `case` in `internal/core/sms/factory`.

## OTP Event Flow

1. `ads-platform-back` generates OTP and publishes to `notifications.otp.send`
2. This service receives the event and logs it
3. The OTP is sent as `sms.TemplateOTP` through the configured `sms.Sender`

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
    │   ├── router/
    │   └── sms/           # provider-neutral Sender interface
    │       ├── factory/   # picks the provider from SMS_PROVIDER
    │       └── provider/  # logsms, smsir, ...
    └── business/otp/
        ├── container/
        ├── listener/      # OTP subscriber
        ├── model/
        └── service/
```
