# ads-platform-ui

UI service for the ads platform, structured around the site map in `ads-platform-doc/document.drawio`.

## Site map

```
/  →  /query-ads
├── /query-ads              Browse / search ads
├── /my-info                Account hub
│   ├── /my-info/user-details
│   ├── /my-info/user-ads
│   ├── /my-info/marked-ads
│   └── /my-info/setting
├── /new-ad                 Post a new ad
├── /category               Categories
└── /location               City / location
```

## Quick Start

1. Start CDN (cities JSON):
   ```bash
   cd ../ads-platform-cdn && ./run.sh
   ```
2. Setup / run UI:
   ```bash
   make setup
   cp config.example .env   # optional
   ./run.sh
   ```
3. Open `http://localhost:8093/` (redirects to `/query-ads`)

## Project Structure

```
ads-platform-ui/
├── cmd/server/
├── internal/
│   ├── business/
│   │   ├── queryads/     # Browse ads (main page)
│   │   ├── myinfo/       # Account + sub-pages
│   │   ├── newad/        # Create ad
│   │   ├── category/
│   │   └── location/
│   ├── domain/
│   ├── web/page/         # Shared HTML page helpers
│   └── core/
│       ├── cities/       # CDN cities client
│       ├── i18n/bundles/ # fa.json, en.json
│       ├── config/
│       ├── container/
│       ├── router/
│       └── view/
├── templates/
│   ├── partials/
│   └── pages/
└── static/
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8093` | HTTP port |
| `APP_NAME` | `Ruab` | Brand name |
| `DEFAULT_CITY` | `tehran` | Default city slug |
| `CDN_BASE_URL` | `http://localhost:4000` | ads-platform-cdn |
| `BFF_BASE_URL` | `http://localhost:8097` | ads-bff (proxies to ads-platform-back) |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/categories` | Proxies CDN `GET /api/categories` |
| GET | `/api/v1/cities` | Proxies CDN `GET /api/cities` |
| GET | `/api/v1/users/mobile/:mobile` | Via ads-bff → backend user lookup |
| POST | `/api/v1/otp/:mobile/send` | Via ads-bff → backend send OTP |
| POST | `/api/v1/otp/:mobile/verify` | Via ads-bff → backend verify OTP |

CDN base URL is configured via `CDN_BASE_URL` (default `http://localhost:4000`).

## Translations

One JSON file per language under `internal/core/i18n/bundles/`. Locale follows the selected city (RTL/LTR from the bundle `dir` field).
