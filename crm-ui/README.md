# crm-ui

Ruab admin panel. Persian-first, RTL by default, with an English LTR mode.

## Run

From the repo root:

```bash
./10-crm-ui.sh
```

Dev server: `http://localhost:8100`

## Sign in

Staff sign in with their mobile number and password. A fresh database gets a default super admin (full access) from `script_8_seed-default-admin.sql`:

| Mobile | Password |
| --- | --- |
| `09120000000` | `Ruab@Admin1405` |

Change that password right after the first sign-in (Staff page → Set password). To create another super admin or reset a password from the server:

```bash
cd ads-platform-back && go run ./cmd/create-admin -name "Admin" -mobile 0912xxxxxxx -national-id 0012345678
```

It prompts for the password (or reads `STAFF_ADMIN_PASSWORD`). Other staff are added from the Staff page.

## Stack

Next.js App Router, TypeScript, Ant Design, TanStack Query, Zod, React Hook Form.

Feature code lives in `src/features`. Shared UI in `src/components`. API/auth/i18n in `src/lib`.

Fonts (Vazirmatn, Inter) are local files in `src/styles/fonts`. The panel does not load CSS or JS from the public internet.

## Config

| Variable               | Default                         | Meaning                                                                                                            |
| ---------------------- | ------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `PORT`                 | `8100`                          | HTTP port                                                                                                          |
| `BACKEND_API_BASE_URL` | `http://localhost:8092`         | ads-platform-back                                                                                                  |
| `SESSION_SECRET`       | `dev-only-change-me` (dev only) | Signs the session cookie. Required in production: a random value of 32+ characters, e.g. `openssl rand -base64 48` |

IAM screens (`/settings/roles`, `/actions`, `/groups`) proxy to `/api/v1/panel/*`. Other lists are wired and stay empty until those panel APIs exist.

## Checks

```bash
npm run check
```
