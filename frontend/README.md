# La Luna HMS Frontend

The SvelteKit frontend is also packaged as six separate Android apps from one shared codebase. Each app has its own role-specific UI, application ID, display name, and launcher color:

| App        | Android application ID     | Permitted account role |
| ---------- | -------------------------- | ---------------------- |
| Guest      | `com.lalunahms.guest`      | `GUEST`                |
| Front Desk | `com.lalunahms.frontdesk`  | `FRONT_DESK`           |
| Kitchen    | `com.lalunahms.kitchen`    | `CHEF`                 |
| Manager    | `com.lalunahms.manager`    | `MANAGER`, `ADMIN`     |
| Restaurant | `com.lalunahms.restaurant` | `WAITER`               |
| Bar        | `com.lalunahms.bar`        | `BAR_ATTENDANT`        |

## Requirements

- Node.js compatible with this project (Node 24 is supported)
- JDK 17
- Android SDK with Android platform/build tools 35
- A deployed backend reachable over HTTPS for mobile app builds

Set the API origin before building. It must point to the deployed **unified HMS API** from `backend/cmd/server`, not the standalone front-desk service in `backend/cmd/frontdesk`. The build checks that `/api/v1/health` identifies the unified API and confirms a live PostgreSQL connection. Do not put database credentials or other server secrets in this value:

```powershell
$env:VITE_API_BASE_URL = 'https://api.your-domain.example'
```

The API origin is embedded in the mobile client. The build script refuses HTTP URLs, local URLs, and unset URLs, and allows 90 seconds for the health check so a sleeping Render free service can wake. For the web build, leave `VITE_API_BASE_URL` empty when using a same-origin `/api/*` rewrite. If the static-site rewrite cannot forward authorization headers, set it to the API's HTTPS origin; the client then sends web API requests directly to that origin. Web development with an empty value continues to use the Vite `/api` proxy to `http://localhost:8080`.

## Build and sync Android apps

Create debug APKs for all role apps:

```powershell
npm run mobile:build:all
```

Create just one role app:

```powershell
npm run mobile:build -- guest
npm run mobile:build -- frontdesk
npm run mobile:build -- kitchen
npm run mobile:build -- manager
npm run mobile:build -- restaurant
npm run mobile:build -- bar
```

Debug APKs are copied to `dist/mobile/`. The command rebuilds the web assets for each role, syncs them to Capacitor, and assembles the matching Android flavor so apps do not accidentally contain another role's frontend.

To build signed release APKs and Play Store AABs, set the Android signing variables in the current PowerShell session and use `--release`:

```powershell
$env:ANDROID_KEYSTORE_PATH = 'C:\secure\laluna-release.jks'
$env:ANDROID_KEYSTORE_PASSWORD = '...'
$env:ANDROID_KEY_ALIAS = 'laluna'
$env:ANDROID_KEY_PASSWORD = '...'
npm run mobile:build:all -- --release
```

Keep the keystore and passwords out of source control. Release outputs are also copied to `dist/mobile/`.

To build web assets and sync without invoking Gradle:

```powershell
npm run mobile:sync -- frontdesk
```

Installed APKs do not update automatically; each APK contains the API address and frontend code from its build. After deploying backend changes, rebuild and deploy the unified API. After frontend changes, rebuild the affected APK with `npm run mobile:build -- <role>` (or rebuild all apps with `npm run mobile:build:all`), then reinstall the updated APK from `dist/mobile/`.

## Backend authorization

Route filtering and app flavors improve the user experience but are not a security boundary. The Go backend validates JWTs and requires an appropriate role on each protected API route. The client attaches the login token to API calls. Health checks, authentication/registration, guest room availability, and resident login are public; operational APIs require a role-appropriate token.

The role-specific apps do not grant manager accounts access to other department APIs. A manager or admin account must use the Manager app.

Set `JWT_SECRET` on every backend instance to the same cryptographically random secret containing at least 32 bytes. The backend refuses to start without it. Demo-account authentication is disabled by default; enable `ALLOW_DEMO_USERS=true` only in a local development environment.

## Web build

The Vite/SvelteKit plugin uses `adapter-static` with an `index.html` fallback for Capacitor's local WebView. The Go API remains a separate online service; no database credentials or privileged backend code are bundled in the app.

```powershell
npm run check
npm run build
```

Mobile use requires a deployed HTTPS API. An Android phone's `localhost` is the phone itself, not the development computer. Offline front-desk data is not an authorization substitute; queued changes are still validated by the backend when synced.
