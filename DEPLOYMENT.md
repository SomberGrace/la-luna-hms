# Render deployment

This deployment uses a Render Docker Web Service for the Go API and a Render Static Site for the SvelteKit frontend. A free Render web service is a long-running process while awake and may sleep when idle; it is not true serverless compute. Supabase remains the hosted PostgreSQL database.

## Go API Web Service

Create a Render Web Service from the repository with:

| Setting | Value |
|---|---|
| Runtime | Docker |
| Root directory | *(leave blank)* |
| Dockerfile path | `Dockerfile` |
| Docker context directory | `.` |
| Region | Frankfurt, if available for the account |
| Instance | Free |
| Health check path | `/api/v1/health` |

Use the repository root for the Docker context. The root `Dockerfile` copies and builds only `backend/`; the root `.dockerignore` excludes the frontend, prototype, Git metadata, and local environment files. This avoids ambiguity about whether Render is resolving Dockerfile paths from the repository root or a configured service root. The equivalent backend-only configuration is Root Directory `backend`, Dockerfile path `Dockerfile`, and Docker context directory `.`, but do not mix settings from these configurations.

Set these service environment variables. Do not commit their values:

| Variable | Value |
|---|---|
| `DATABASE_URL` | Supabase **session-mode** pooler URL on port `5432`, with `sslmode=require` |
| `JWT_SECRET` | Stable, randomly generated secret with at least 32 bytes |
| `ALLOW_DEMO_USERS` | `false` |
| `RUN_STARTUP_MIGRATIONS` | `false` |
| `SEED_DATABASE` | `false` |
| `DB_MAX_OPEN_CONNS` | `5` |
| `DB_MAX_IDLE_CONNS` | `2` |
| `DB_CONN_MAX_LIFETIME` | `5m` |
| `DB_CONN_MAX_IDLE_TIME` | `1m` |

Do not set `PORT`; Render injects it. Keep the Supabase password out of source control and logs. Verify that the existing Supabase schema contains `guest_accounts`, `bookings.checked_in`, and the docket workflow columns before deployment. Do **not** run `database/schema.sql` against the existing database; it drops application tables. Keep `JWT_SECRET` unchanged after deployment or existing sessions will become invalid.

The pool is capped at five open connections per running API instance. If Supabase's session pool is smaller, reduce `DB_MAX_OPEN_CONNS`. If scaling beyond one Render instance, calculate total connections as instances multiplied by five and stay within both Supabase limits.

Add a stable custom domain to the API service before building APKs. A domain change requires rebuilding installed APKs because the API origin is embedded in the app.

## SvelteKit Static Site

Create a Render Static Site from the same repository:

| Setting | Value |
|---|---|
| Root directory | `frontend` |
| Build command | `npm ci && npm run build` |
| Publish directory | `build` |
| `NODE_VERSION` | `24` |
| `PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD` | `1` |

Preferred routing is a Render rewrite from `/api/*` to `https://<API_CUSTOM_DOMAIN>/api/*`, followed by the SPA fallback `/*` to `/index.html`. Keep `VITE_API_BASE_URL` empty for this mode, so browser requests remain relative and same-origin. Verify that the rewrite preserves `Authorization` on login and authenticated requests.

If the rewrite does not forward authorization headers, set the static site's build variable `VITE_API_BASE_URL` to the API's HTTPS origin and redeploy. The frontend client will then prefix web API requests with that origin. The API handles CORS preflights globally, including `Authorization` for all operational routes; CORS is not a substitute for backend authentication and role checks. Do not put database credentials or JWT secrets in frontend build variables.

## Mobile APKs

On a machine with Node 24, JDK 17, and Android SDK/build tools 35, set `VITE_API_BASE_URL` to the stable API HTTPS origin and run from `frontend`:

```powershell
$env:VITE_API_BASE_URL = 'https://api.your-domain.example'
npm run mobile:build:all
```

The build checks that the unified API health endpoint reports a live Supabase connection and now allows up to 90 seconds for a sleeping free service to wake. APK outputs are written to `frontend/dist/mobile`. Test an APK on mobile data and verify login plus one workflow for each role you intend to deploy.

## Local cleanup and files to retain

The following are generated and can be recreated rather than committed or kept between builds:

- `frontend/node_modules/` with `npm ci`
- `frontend/.svelte-kit/` and `frontend/build/` with the frontend build
- `frontend/android/.gradle/` and `frontend/android/app/build/` with Gradle/Capacitor
- `frontend/android/app/src/main/assets/public/` with `npx cap sync android` (the mobile build script runs this before assembling)
- `frontend/dist/` with the mobile build script
- `backend/laluna-go.exe`, if it is only a local prebuilt binary; rebuild it from Go source when needed

Keep `frontend/android/` source files, Gradle wrapper, `backend/`, migrations, and `database/schema.sql` for APK builds, API builds, and fresh-database setup. The standalone `backend/cmd/frontdesk` program is not used by this Render API image, which builds only `cmd/server`; remove it only if no local workflow or external consumer still depends on it. Do not remove example tests or documentation solely to reduce deployment size; they are not included in the Docker runtime image or static publish directory.

This workspace is not currently a Git repository. Render requires a connected Git repository; initialize and push it only after confirming the ignore rules exclude local secrets and generated outputs. Rotate any credentials that may have been committed or shared previously.
