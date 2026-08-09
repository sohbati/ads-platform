# ads-bff

Backend-for-Frontend service. The UI talks to `ads-bff`; `ads-bff` proxies backend APIs and owns login sessions.

## Quick start

1. Start backend dependencies (cache, NATS, notification, back).
2. Run BFF:
   ```bash
   ./run.sh
   ```
3. Health: `GET http://localhost:8097/health`

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8097` | BFF HTTP port |
| `BACKEND_API_BASE_URL` | `http://localhost:8092` | ads-platform-back |
| `CACHE_SERVICE_URL` | `http://localhost:8093` | Session storage |
| `SESSION_COOKIE_NAME` | `ads_session` | HttpOnly session cookie |
| `SESSION_TTL` | `24h` | Session lifetime |
| `COOKIE_SECURE` | `false` | Set `true` in production (HTTPS) |

## Auth flow (OTP login)

```
UI  →  POST /api/v1/auth/otp/:mobile/send   →  back (send OTP)
UI  →  POST /api/v1/auth/otp/:mobile/verify →  back (verify) + cache (session) + Set-Cookie
UI  →  GET  /api/v1/auth/me                 →  cache lookup via cookie
UI  →  POST /api/v1/auth/logout             →  delete session + clear cookie
```

Session data is stored in cache-service under key `session:{uuid}`.

Other backend routes remain available via proxy:

- `GET /api/v1/users/mobile/:mobile`
- etc.

## Architecture

```
ads-platform-ui  →  ads-bff (:8097)  →  ads-platform-back (:8092)
                         ↓
                 ads-platform-cache-service (:8093)  [sessions]
```
