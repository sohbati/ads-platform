# OTP Verify Integration Tests

All OTP scenarios share **one** Docker stack via `TestMain` (started once, torn down once).

## Structure

```
tests/otp_verify/
├── main_test.go              # starts/stops shared OTP stack
├── valid_send_test.go
├── valid_verify_test.go
├── e2e_happy_path_test.go
├── wrong_otp_test.go
├── ...
└── verify_twice_test.go
```

## Run

```bash
cd integration-test
make test-otp
# or
go test -tags=integration -timeout=45m -v ./tests/otp_verify/
```

## Shared helpers

- `internal/otptest/helpers.go` — HTTP helpers
- `internal/testcontainers/otp_stack.go` — back, cache, NATS, notification, postgres
