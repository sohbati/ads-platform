# ads-bff

Backend-for-Frontend service. The UI talks to `ads-bff`; `ads-bff` proxies `/api/v1/*` to `ads-platform-back`.

## Quick start

1. Start backend (and its dependencies):
   ```bash
   ../1-ads-platform-back.sh
   ```
2. Run BFF:
   ```bash
   ./run.sh
   ```
3. Health: `GET http://localhost:8097/health`

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8097` | BFF HTTP port |
| `BACKEND_API_BASE_URL` | `http://localhost:8092` | ads-platform-back base URL |

## Architecture

```
ads-platform-ui  →  ads-bff (:8097)  →  ads-platform-back (:8092)
```

All backend API routes are available under the same paths, e.g.:

- `GET /api/v1/users/mobile/:mobile`
- `POST /api/v1/otp/:mobile/send`
- `POST /api/v1/otp/:mobile/verify`
