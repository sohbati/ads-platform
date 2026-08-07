# OTP Verify Integration Tests

Each folder is one scenario. All tests use testcontainers and the `integration` build tag.

## Structure

```
tests/otp_verify/
├── valid_send/           # POST send returns 200, OTP stored in cache
├── valid_verify/         # POST verify with correct OTP returns 200
├── e2e_happy_path/       # send → cache → verify → wrong otp rejected
├── wrong_otp/            # verify with incorrect code → 401
├── expired_otp/          # verify without send → 404 (skipped: 500s TTL)
├── missing_body/         # verify with no body → 400
├── missing_otp_field/    # verify with {} → 400
├── otp_too_short/        # 5 digits → 400
├── otp_too_long/         # 7 digits → 400
├── non_numeric_otp/      # letters in otp → 400
├── empty_mobile_send/    # send with empty mobile → 400
├── empty_mobile_verify/  # verify with empty mobile → 400
├── wrong_mobile/         # send to A, verify on B → 404
├── resend_overwrites/    # second send replaces cached OTP
└── verify_twice/         # same OTP verifies twice (not consumed)
```

## Run all OTP tests

```bash
cd integration-test
go test -tags=integration -timeout=20m -v ./tests/otp_verify/...
```

## Run one scenario

```bash
go test -tags=integration -timeout=20m -v ./tests/otp_verify/valid_send/...
```

## Shared helpers

- `internal/otptest/helpers.go` — HTTP client, stack setup
- `internal/testcontainers/otp_stack.go` — back, cache, NATS, notification containers
