# E2E Tests

Playwright tests for Apptron. Tests run against a local `wrangler dev` server and use a real Hanko auth flow.

## Prerequisites

**1. Install dependencies**
```bash
npm install
```

**2. Install the Chromium browser**
```bash
npx playwright install chromium
```

**3. Set up environment variables**

Copy `.env.example` to `.env.local` and fill in the required values:

```bash
cp .env.example .env.local
```

| Variable | Where to get it |
|---|---|
| `AUTH_URL` | Hanko Cloud dashboard → your project URL (different for local and prod) |
| `MAILSLURP_API_KEY` | Create a free account at [mailslurp.com](https://mailslurp.com) |
| `HANKO_ADMIN_API_KEY` | Hanko Cloud dashboard → your project → API Keys → **secret** (different for local and prod) |

`AUTH_URL` and `HANKO_ADMIN_API_KEY` are environment-specific — your local Hanko project and prod Hanko project each have their own values. `.env.local` is for local development; set the prod equivalents in your deployment environment.

## Running the tests

```bash
npx playwright test
```

Or with the browser visible (useful for debugging):

```bash
npx playwright test --headed
```

## How it works

**Setup** (`setup.ts`) runs once before all tests. It:
1. Creates a real MailSlurp inbox to receive the Hanko verification email
2. Signs up a new test account via the Hanko auth flow
3. Uses a virtual WebAuthn authenticator to register a passkey — this is in-memory only and never saves to your browser's credential store
4. Saves the authenticated session to `tests/.auth/user.json`

**Tests** (`*.test.ts`) load the saved session automatically — no sign-in needed.

**Teardown** (`teardown.ts`) runs after all tests. It deletes the test account from Hanko via the admin API so nothing accumulates.

## Notes

- `tests/.auth/` is gitignored — it contains live session tokens and should never be committed
- The virtual authenticator means no real passkeys are created in your browser
- If setup fails due to a transient Hanko API error, Playwright will retry once automatically
