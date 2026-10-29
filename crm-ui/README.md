# crm-ui

Ruab admin panel. Persian-first, RTL by default, with an English LTR mode.

## Run

From the repo root:

```bash
./10-crm-ui.sh
```

Dev server: `http://localhost:8100`

Temporary login: `09120000000` / `admin`

## Stack

Next.js App Router, TypeScript, Ant Design, TanStack Query, Zod, React Hook Form.

Feature code lives in `src/features`. Shared UI in `src/components`. API/auth/i18n in `src/lib`.

Fonts (Vazirmatn, Inter) are local files in `src/styles/fonts`. The panel does not load CSS or JS from the public internet.

## Config

| Variable | Default | Meaning |
| --- | --- | --- |
| `PORT` | `8100` | HTTP port |
| `BACKEND_API_BASE_URL` | `http://localhost:8092` | ads-platform-back |
| `SESSION_SECRET` | `dev-only-change-me` | Signs the admin session cookie |

IAM screens (`/settings/roles`, `/actions`, `/groups`) proxy to `/api/v1/panel/*`. Other lists are wired and stay empty until those panel APIs exist.

## Checks

```bash
npm run check
```
