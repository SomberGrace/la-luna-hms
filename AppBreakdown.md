# 🌙 La Luna HMS — App Breakdown

> **A plain-language guide to what the app does, how its code fits together, and how information moves through it.**
>
> **Short version:** La Luna HMS is a hotel-management prototype with guest booking and separate work areas for reception, restaurant, bar, kitchen, and management. Its interface is built with SvelteKit and TypeScript, its API is written in Go, and its shared database is PostgreSQL.

## 1. The hotel as a software system

Imagine the app as a collection of hotel workstations connected to one central record book:

- The **frontend** is what guests and staff see and use: pages, buttons, forms, dashboards, and status messages.
- The **backend** receives requests from those screens, applies the operation, talks to the database, and returns a result.
- The **database** stores shared records such as users, rooms, reservations, folios, menu items, and stock.
- **Browser storage** keeps some information on one device. It is not the same as a shared database: it may be demo data, a local preference, or a queued action waiting to sync.

The main path is:

```text
Person → Svelte page → HTTP request (/api/v1/...) → Go route → handler/service
       → PostgreSQL query or transaction → JSON response → page updates
```

Not every screen follows this full path yet. Some screens use sample values or browser storage, so a convincing dashboard does not necessarily mean every number is live.

## 2. Backend walkthrough — the main engine

The backend is the best place to understand what an operation actually does. Its files live under `backend/`.

### How the Go server starts

The main, unified API entry point is [`backend/cmd/server/main.go`](./backend/cmd/server/main.go). Its `main()` function performs the startup sequence:

1. **Load settings** by calling `config.LoadConfig()`.
2. **Create the authentication service** with the configured JWT signing secret.
3. **Try to connect to PostgreSQL** using `db.InitDB()`. If that succeeds, the authentication service receives the database connection. If it fails, the server logs the problem and can still start in a limited offline/demo mode.
4. **Create handler objects** for authentication, front desk, restaurant, kitchen, bar, manager, and guest operations. The database connection is passed into the handlers that need it.
5. **Register URL paths** with Go's `http.HandleFunc`. Each path is connected to a handler method.
6. **Start listening** on the configured port (8080 by default).

For example, `/api/v1/frontdesk/checkin` is routed to `FrontDeskHandler.CheckIn`; `/api/v1/auth/login` is routed to `AuthHandler.HandleLogin`.

### What a handler does

A **handler** is a Go function that handles one HTTP request. It normally:

1. Sets response headers, including JSON and CORS headers.
2. Checks the HTTP method or handles browser preflight (`OPTIONS`) requests.
3. Decodes the request body from JSON into a Go struct.
4. Validates or maps the incoming values.
5. Reads or changes database records.
6. Encodes a JSON response and writes an HTTP status code.

The handlers live in [`backend/internal/handlers/`](./backend/internal/handlers/). The structs used to decode JSON and shape responses define the informal API contract between frontend and backend.

### Authentication and account creation

- [`backend/internal/handlers/auth_handler.go`](./backend/internal/handlers/auth_handler.go) accepts login and registration requests and returns HTTP responses.
- [`backend/internal/services/auth_service.go`](./backend/internal/services/auth_service.go) contains the core authentication work: it looks up an account, checks the password hash with bcrypt, registers guest accounts, and builds a signed JWT after successful authentication.
- A JWT contains identity and role claims and is set to expire after 24 hours. Guest endpoints that call `extractUserID()` use the bearer token to identify the guest.
- The service has demo/offline fallback accounts when its database connection is unavailable or a database account lookup does not find a user.

The frontend stores the returned token when it receives one. However, the current login screen also has a fallback that can navigate to a workstation after a failed request or network error. That fallback is useful for a demo, but it should not be treated as secure authentication.

### Database setup and shared records

- [`backend/internal/config/config.go`](./backend/internal/config/config.go) reads `PORT`, `JWT_SECRET`, and `DATABASE_URL` (or `SUPABASE_DB_URL`), with `.env` files as a fallback. Defaults include port `8080` and a JWT secret, so production deployment must provide safe environment values.
- [`backend/internal/db/db.go`](./backend/internal/db/db.go) opens the PostgreSQL connection, configures the connection pool, checks connectivity, runs small startup migrations, and seeds baseline data.
- [`database/schema.sql`](./database/schema.sql) defines the database types and tables.
- [`database/migrations/001_guest_accounts.sql`](./database/migrations/001_guest_accounts.sql) adds guest-registration support and room description/image fields.
- [`docker-compose.yaml`](./docker-compose.yaml) describes a local PostgreSQL container for development.

The schema groups records into hotel concepts:

| Hotel concept | Main database records |
|---|---|
| Staff and guest accounts | `user_accounts`, `user_role` |
| People and accommodation | `guests`, `rooms`, `bookings`, `folios` |
| Bar sales and stock | `bar_items`, `bar_ledgers`, `bar_reports` |
| Kitchen ingredients and recipes | `kitchen_stock`, `menu_item_ingredients`, `kitchen_reports` |
| Restaurant menu and orders | `menu_items`, `restaurant_dockets`, `docket_items`, `restaurant_reports` |
| Management reporting | `executive_summaries`, `consolidated_reports`, `report_documents` |

PostgreSQL constraints and enum types help keep values consistent. For example, room status and payment mode are limited to known values, while room rates and stock quantities have non-negative checks.

**Important setup warning:** `database/schema.sql` begins by dropping tables and custom types before recreating them. It is a reset-style schema script, not a safe update to run against a database with data you want to keep. The backend's startup migrations are separate from that full schema script.

### A concrete backend example: front-desk check-in

The check-in operation in [`backend/internal/handlers/frontdesk_handler.go`](./backend/internal/handlers/frontdesk_handler.go) makes a useful example of how backend code changes real records:

1. The frontend sends guest details, a room number, number of nights, and deposit/payment information as JSON.
2. The handler decodes the JSON into a Go request struct.
3. It begins a database **transaction**. A transaction makes several related changes act as one unit: either they all commit, or an unfinished transaction is rolled back.
4. It looks up and locks the room row, checks the room status, and gets the room rate and type.
5. It creates a guest record and a booking.
6. It creates a folio (the guest's running bill) with the room charges.
7. It marks the room occupied.
8. It commits the transaction and returns a success response with folio details. If an earlier step fails, the deferred rollback prevents a half-finished check-in.

This is the general pattern to look for in backend code: **decode → validate → query/update → commit → respond**.

### What each backend module handles

| Backend file | Responsibility |
|---|---|
| [`backend/internal/handlers/frontdesk_handler.go`](./backend/internal/handlers/frontdesk_handler.go) | Room status, check-in/walk-in, folios, payments, checkout, extensions, reservations, and daily reports. |
| [`backend/internal/handlers/guest_handler.go`](./backend/internal/handlers/guest_handler.go) | Available rooms, guest booking, guest bookings and folios, guest check-in, and resident login. It checks bearer tokens for guest-specific operations. |
| [`backend/internal/handlers/restaurant_handler.go`](./backend/internal/handlers/restaurant_handler.go) | Menu, ingredient-aware availability, restaurant dockets, summaries, reports, and meal reconciliation. |
| [`backend/internal/handlers/kitchen_handler.go`](./backend/internal/handlers/kitchen_handler.go) | Kitchen stock records, recipes, and kitchen reports. |
| [`backend/internal/handlers/bar_handler.go`](./backend/internal/handlers/bar_handler.go) | Bar item lists, sales, stock ledgers, and reports. Sales use a transaction when updating related records. |
| [`backend/internal/handlers/manager_handler.go`](./backend/internal/handlers/manager_handler.go) | Cross-department summaries, audit actions, and consolidated reports. |
| [`backend/internal/handlers/auth_handler.go`](./backend/internal/handlers/auth_handler.go) | HTTP endpoints for login and guest registration; delegates authentication logic to the service. |
| [`backend/internal/services/auth_service.go`](./backend/internal/services/auth_service.go) | Password checking/hashing, account registration, fallback demo users, and JWT creation. |

### API map

These are the main route groups registered by the unified API:

| URL group | What it serves |
|---|---|
| `/api/v1/health` | General service and database health |
| `/api/v1/auth/...` | Login and registration |
| `/api/v1/frontdesk/...` | Rooms, check-in, folios, payments, checkout, reservations, reports |
| `/api/v1/guest/...` | Available rooms, booking, guest bookings/folios, guest check-in |
| `/api/v1/restaurant/...` | Menu, dockets, summaries, reports, meal reconciliation |
| `/api/v1/kitchen/...` | Stock, recipes, reports |
| `/api/v1/bar/...` | Items, sales, stock ledger, reports |
| `/api/v1/manager/...` | Summary, audits, consolidated report |

For some paths, the server selects the operation by HTTP method: for example, `GET /restaurant/menu` reads the menu while `POST /restaurant/menu` adds an item.

## 3. Frontend walkthrough — the workstations

The frontend lives under `frontend/` and uses SvelteKit, Svelte, TypeScript, and Tailwind CSS.

- [`frontend/src/routes/`](./frontend/src/routes/) contains the pages. A `+page.svelte` file is a page, and a `+layout.svelte` file wraps pages with shared navigation or styling.
- [`frontend/src/routes/+layout.svelte`](./frontend/src/routes/+layout.svelte) loads shared styles, fonts, icons, and the favicon.
- [`frontend/src/routes/+page.svelte`](./frontend/src/routes/+page.svelte) is the staff login/role-selection portal.
- Department folders hold their own dashboards and task pages: `frontdesk`, `guest`, `book`, `restaurant`, `bar`, `kitchen`, and `manager`.
- [`frontend/src/lib/stores/authStore.ts`](./frontend/src/lib/stores/authStore.ts) holds the current browser session and saves it in local storage.
- [`frontend/src/lib/stores/roomStore.ts`](./frontend/src/lib/stores/roomStore.ts) manages front-desk room/folio state and the offline action queue.
- [`frontend/src/lib/stores/guestStore.ts`](./frontend/src/lib/stores/guestStore.ts) holds guest stay, dining-order, and service-request state in browser storage, including default sample records.
- [`frontend/src/lib/types/guest.ts`](./frontend/src/lib/types/guest.ts) and [`frontend/src/lib/types/room.ts`](./frontend/src/lib/types/room.ts) describe the shapes of frontend data.

Svelte pages use state to keep the screen responsive. A page may call `fetch()` in response to a user action or when it loads. It converts the response from JSON and updates the state used to render the page. A form generally gathers values, sends JSON to an API endpoint, then displays success or an error.

## 4. Data journeys from screen to database

### Guest makes an online booking

```text
Booking page → GET /api/v1/guest/rooms → available rooms displayed
Guest selects room and dates → POST /api/v1/guest/book + bearer token
GuestHandler → validates dates and room availability
             → finds/creates guest → inserts booking and folio in transaction
             → commits → returns booking number and total
Booking page → displays confirmation
```

The backend booking path requires a valid guest token. The guest account registration endpoint creates both an account and a guest profile.

### Front desk checks in a walk-in

```text
Front-desk form → POST /api/v1/frontdesk/checkin
FrontDeskHandler → locks/checks room → creates guest, booking, and folio
                 → marks room occupied → commits transaction
                 → returns folio and balance
Frontend room/folio state → refreshes the display
```

If offline, `roomStore.ts` can save selected front-desk operations locally and queue them for later synchronization. This is a frontend queue; it is not the same as a backend database transaction, and only the operations implemented in that store are covered.

### Restaurant order uses kitchen stock

```text
Restaurant menu page → GET /api/v1/restaurant/menu
RestaurantHandler → reads menu and ingredient links to kitchen stock
                 → marks an item unavailable when required stock is insufficient
Menu page → displays returned availability
```

Restaurant dockets and kitchen stock are separate records, connected through menu/ingredient relationships where configured. Do not assume every meal order automatically updates every ingredient count unless that behavior is implemented for that path.

### Bar sale and manager reporting

```text
Bar screen → POST /api/v1/bar/sales → BarHandler records sale and adjusts stock
Manager screen → GET /api/v1/manager/summary → ManagerHandler queries departments
             → returns a combined summary → dashboard displays totals
```

Some dashboard values in the frontend or manager reporting code are example/fallback values. Check the page and handler before describing a particular number as live.

## 5. Security and reliability — current reality

The code includes useful security building blocks: bcrypt password hashes, signed expiring JWTs, SQL queries with parameters in many operations, database constraints, and transactions for related changes.

There are also important gaps to understand:

- Role selection and role-specific navigation in the frontend are **not** by themselves authorization. The unified server registers many staff endpoints directly, and the inspected routes do not show a shared role-checking middleware protecting all of them.
- The login page can use an offline/demo fallback and navigate even when backend login fails. That is a demonstration convenience, not proof of identity.
- The auth service and database seeder include demo accounts/passwords. The config has a default JWT secret. These must not be relied on for a real deployment; set strong private values and remove or disable demo access.
- CORS headers allow requests from any origin in the handler code. Production access should be restricted to the intended frontend origin.
- Some handler methods assume the database is available. The server may start in offline mode, but that does not make every module work offline.
- The reset-style `schema.sql` drops existing tables. Never run it against valuable data without an intentional backup and reset plan.

**Plain-language takeaway:** authentication exists, but authorization must be consistently enforced by the backend before this is used with real guest or financial data.

## 6. Where the app is in development

The repository looks like a **feature-rich prototype that is partway through integration**:

**Already present**

- Separate interfaces for guest booking and major hotel departments.
- A unified Go API with route groups for authentication, rooms, bookings, restaurant, kitchen, bar, and management.
- A detailed PostgreSQL schema, startup seeding, and guest-account migration support.
- Several real database workflows, including guest booking and front-desk check-in.
- A front-desk local cache/offline queue for selected actions.

**Partly integrated or demo-oriented**

- Some pages call real API endpoints; some show preset dashboard figures or browser-stored sample records.
- Some pages fall back to sample data when the API cannot be reached.
- The front-desk queue exists, but the code does not establish that the entire application can operate offline.
- The existence of a route or handler does not prove its full end-to-end journey is complete or tested.

**Not yet demonstrated by the inspected code**

- Consistent role-based authorization for every staff API operation.
- Complete end-to-end tests for every module and workflow.
- Production-ready secret, CORS, deployment, and data-reset practices.

There is also a second Go entry point, [`backend/cmd/frontdesk/main.go`](./backend/cmd/frontdesk/main.go), which contains a standalone front-desk implementation and an in-memory fallback store. The primary module/API map above describes the unified server at `backend/cmd/server/main.go`; treat the standalone entry point as a separate implementation rather than assuming it is the one serving every unified route.

## 7. A short explanation to present

> “La Luna HMS is a hotel operations prototype. Guests can browse and book rooms, and hotel staff use different work areas for reception, restaurant, bar, kitchen, and management. The frontend sends requests to a Go API. The API validates the request, performs database work in PostgreSQL, and returns a response that the screen displays. Some workflows are connected to the database, while some dashboards and guest features still use sample or browser-stored data. The main work ahead is finishing integration and tests, then strengthening backend permissions and production configuration.”

## 8. A simple code-reading route

To follow one feature through the project, read these in order:

1. Find the page in `frontend/src/routes/<module>/.../+page.svelte`.
2. Search that page for `fetch(` to see which URL it calls and what JSON it sends.
3. Find that URL in `backend/cmd/server/main.go` to see which handler receives it.
4. Read the corresponding method in `backend/internal/handlers/`.
5. Follow the SQL table names into `database/schema.sql`.
6. Return to the page and see how it handles the JSON response and updates what is displayed.

That route—from **screen, to URL, to handler, to database, and back to screen**—is the clearest way to explain how any individual feature works.

## 9. Code Operations Breakdown and Understanding

This section follows the backend methods that carry out each department's work. The simplest way to read a Go handler is:

```text
HTTP request → decode JSON → check input → run SQL → encode JSON response
```

In the methods below, `r` is the incoming request and `w` is where the handler writes its response. `h.db` is the shared PostgreSQL connection. A **DTO** (data-transfer object) is a Go struct shaped for JSON; it translates database fields into values the frontend can display.

### A. Shared backend concepts to recognize

- **`QueryRow`** runs a query expected to return one row; `Scan(...)` copies its columns into Go variables.
- **`Query`** reads multiple rows; the handler loops over them, scans each row, and closes the result when done.
- **`Exec`** runs a statement that does not need to return a row, such as an `UPDATE`.
- **`$1`, `$2`, ...** are SQL parameters. The values are sent separately from the query text, which avoids building SQL by concatenating ordinary user input.
- **`Begin`, `Commit`, `Rollback`** mark a transaction. A transaction keeps a group of related changes together. `defer tx.Rollback()` is commonly used as a safety net; after a successful `Commit`, that rollback has no effect.
- **`json.NewDecoder(r.Body).Decode(&req)`** turns the incoming JSON into a typed Go value. The struct's `json:"..."` tags say which JSON names map to which Go fields.
- **`json.NewEncoder(w).Encode(...)`** turns a Go value into a JSON response.
- **HTTP status codes** communicate the outcome: for example, `400` means a bad request, `401` means authentication failed, `404` means a record was not found, and `500` means the server could not complete its work.

Handlers often use `setCORS(...)` to attach JSON and cross-origin headers and return early for an `OPTIONS` preflight request. This is browser communication setup, not user authentication.

### B. Login and guest registration

**Files:** [`auth_handler.go`](./backend/internal/handlers/auth_handler.go), [`auth_service.go`](./backend/internal/services/auth_service.go)

**Login path**

1. `HandleLogin` accepts `POST`, decodes username/password/role into `LoginRequest`, and rejects a malformed body.
2. It calls `AuthService.Login`.
3. If the service has a database connection, it looks up an active username in `user_accounts`, retrieves the stored password hash and role, and compares the submitted password with bcrypt.
4. If the database path does not find a matching row, the code tries its predefined demo-user map. A matching password results in a `LoginResponse`.
5. `buildResponse` signs a JWT with user ID, username, role, issuer, issue time, and 24-hour expiry. The handler returns the token and user details as JSON.

**Registration path**

1. `HandleRegister` decodes the guest's account and profile fields.
2. `AuthService.Register` requires a database, checks required fields, and hashes the password with bcrypt.
3. It starts a transaction, inserts a `GUEST` account into `user_accounts`, and inserts a corresponding guest profile into `guests`.
4. It commits both records together and returns a signed-in response with a JWT.

**Why separate handler and service?** The handler knows HTTP concerns—methods, JSON, and status codes. The service knows authentication concerns—passwords, account lookup, registration, and token creation.

### C. Front desk: the stay lifecycle

**File:** [`frontdesk_handler.go`](./backend/internal/handlers/frontdesk_handler.go)

| Method | What it does |
|---|---|
| `HealthCheck` | Pings the database when possible and reports whether the backend sees PostgreSQL. |
| `GetRooms` | Reads room number, type, bed type, rate, and status, then returns a JSON room list for the front desk. |
| `UpdateRoomStatus` | Converts screen labels such as `Cleaning` or `OutOfService` to the database's `Maintenance` status, then updates the room row. |
| `CheckIn` | Creates the guest, booking, and folio, and marks the room occupied in a transaction. `FOR UPDATE` locks the selected room while the operation checks it, reducing the chance of two simultaneous check-ins claiming it. |
| `GetFolio` | Finds the latest folio by room number or card number, joins guest and booking details, and adds room-charge restaurant dockets to the returned item list. |
| `RecordPayment` | Finds and locks the relevant booking/folio, adds the payment amount to `advance_paid`, marks the folio settled when the remaining balance is near zero, and commits the update. |
| `CheckOut` | Finds and locks the booking/folio, calculates balance due, and refuses checkout if money remains outstanding. If cleared, it closes the folio and sets the room to `Available`. |
| `ExtendStay` | Adds extra days to the departure date, calculates the additional room charge from the nightly rate, and updates booking and folio in one transaction. |
| `GetReservations` | Joins bookings to guests, rooms, and folios, then returns a screen-friendly reservation list. It marks an entry as online based on the clerk label. |
| `DailyReport` | Counts rooms and occupied rooms, sums open folio charges/balances, calculates occupancy percentage, and returns summary JSON. |

**Typical user journey:** a receptionist picks a room and guest, submits check-in, and receives a card number/folio. During the stay, restaurant room charges can appear in the folio. Payments reduce the balance. Checkout is blocked when a balance is still due; after settlement, the folio closes and the room is released.

**Detail worth explaining carefully:** checkout currently sets the room to `Available` even though its response says it is released for housekeeping inspection. The UI also maps cleaning to the database's `Maintenance` state. So the message and stored room status are not perfectly aligned yet.

### D. Guest booking and resident access

**File:** [`guest_handler.go`](./backend/internal/handlers/guest_handler.go)

| Method | What it does |
|---|---|
| `setCORS` | Adds response headers used by browser requests. |
| `extractUserID` | Reads `Authorization: Bearer ...`, parses the JWT with the configured secret, checks that it is valid, and extracts the `user_id` claim. |
| `GetAvailableRooms` | Selects rooms whose status is `Available` and returns their descriptions, image URLs, types, and rates. |
| `CreateBooking` | Requires a valid token, parses and validates arrival/departure dates, checks room availability, finds or creates the guest profile, inserts booking and folio rows in a transaction, then returns a confirmation. |
| `resolveGuestInfo` | Maps an authenticated user ID to a guest profile, first checking guest ID directly and then following account email to the guest record. |
| `GetMyBookings` | Requires a token, resolves that guest, and returns only booking rows matching the guest ID or email, including room and folio details. |
| `GetMyFolio` | Requires a token, finds the latest open folio belonging to that guest, and returns the charges and stay dates. |
| `CheckIn` | Always returns `403 Forbidden`: online self-check-in is intentionally disabled, and arrival must be processed by reception. |
| `ResidentLogin` | Accepts a room and optional card number or last name, finds an active booking, verifies that the room is occupied, then creates a guest JWT with a 72-hour expiry. |

**Booking versus check-in:** booking reserves a room for dates; it does not mark the room occupied in this handler. Occupancy is meant to follow front-desk identity verification and check-in. This distinction helps explain why a guest can have a booking but still be denied resident access until reception checks them in.

### E. Restaurant: menu, orders, and meal counts

**File:** [`restaurant_handler.go`](./backend/internal/handlers/restaurant_handler.go)

| Method | What it does |
|---|---|
| `GetMenu` | Reads menu items and joins their ingredient links to kitchen stock. It marks a dish unavailable when a linked ingredient's closing balance is too low. |
| `AddMenuItem` | Decodes a new item, maps its category to the database enum, and inserts it into the menu. |
| `GetDockets` | Reads recent dockets with their line items, using PostgreSQL JSON aggregation to package each docket's items as JSON. |
| `CreateDocket` | Calculates the order total, decides whether it is cash or a room charge, inserts a docket and its line items in one transaction, and deducts linked recipe ingredients from kitchen stock (the menu's Bill of Materials, or BOM). |
| `Summary` | Sums today's docket revenue and counts today's docket records. The displayed `active_tables` and `today_covers` values are currently fixed example values in the handler. |
| `Reports` | `POST` stores a submitted restaurant report; `GET` reads a report for a requested date. The returned category breakdown and top item are example values in the current code. |
| `GetBreakfastGuests` | Lists current/in-house rooms and estimates included covers from room type/bed type. It returns example guests when there is no database connection. |
| `ReconcileBreakfast`, `ReconcileDinner` | Thin wrappers that call the shared meal-reconciliation function with the chosen meal type. |
| `ReconcileMeal` | Reads room, meal type, covers, and notes, then delegates to the same reconciliation logic. |
| `executeReconcile` | Increases a breakfast or lunch/dinner cover count on today's restaurant report, if a database exists, and returns a reconciliation message. |

**Typical user journey:** an attendant chooses food items and quantities, then submits a docket. The backend calculates the total and saves the docket and item rows. When menu ingredients are linked to kitchen stock, it also reduces their closing balances. If the order is a room charge, the docket is connected to that room's active folio where one is found.

**BOM in plain language:** a Bill of Materials is a recipe map. If one dish needs a certain quantity of chicken and flour, a docket for several portions uses that recipe map to calculate ingredient usage.

### F. Kitchen: daily stock and recipes

**File:** [`kitchen_handler.go`](./backend/internal/handlers/kitchen_handler.go)

| Method | What it does |
|---|---|
| `GetStock` | Reads stock for a date (today by default) and returns opening stock, additions, usage, and closing balance. |
| `SaveStock` | Decodes stock items and, inside a transaction, updates an ingredient's record for that date or inserts a new one. |
| `GetRecipes` | Joins menu items to BOM ingredients and returns each dish with its ingredient list and approximate ingredient costs. |
| `SubmitReport` | `POST` stores the date, JSON stock summary, and submitted status in `kitchen_reports`; `GET` returns the report status and JSON summary for a date. |

**Typical user journey:** kitchen staff review the day's starting quantities, record deliveries/additions and usage, save closing balances, and submit a report. Restaurant menu availability can then use those recorded balances if the ingredient links exist.

### G. Bar: item sales, stock ledger, and report

**File:** [`bar_handler.go`](./backend/internal/handlers/bar_handler.go)

| Method | What it does |
|---|---|
| `GetItems` | Reads the bar catalog and current stock quantities. |
| `RecordSale` | Finds an item by a partial name match and, if found, decrements its stock inside a transaction. |
| `GetLedger` | Reads bar items and the requested day's ledger, returning opening, added, sold, closing, and revenue fields. If there is no ledger row, it presents current stock as opening/closing stock. |
| `SaveLedger` | For each submitted item, finds or creates the catalog item, limits closing stock to opening plus additions, then updates or inserts the daily ledger in a transaction. |
| `Reports` | `POST` saves total sales and cost of goods sold (COGS); the database calculates gross margin. `GET` returns those metrics for a date. |

**Important current behavior:** `RecordSale` decrements inventory when it finds a matching item, but the inspected method does not insert a sale transaction or update a ledger/report. Its success message says the sale was recorded, but the actual database work in this method is the stock decrement. That means the sale and reporting trail may not yet be complete.

### H. Manager: summaries, review, and archive

**File:** [`manager_handler.go`](./backend/internal/handlers/manager_handler.go)

| Method | What it does |
|---|---|
| `Summary` | Reads room payment totals, restaurant sales, bar sales, and room occupancy. The current handler substitutes example amounts when some queried totals are zero and uses a fixed petty-cash value. |
| `GetAudits` | Looks for submitted bar, kitchen, and restaurant reports for a date, then formats them as rows for manager review. Some clerk, timestamp, and detail text is currently hard-coded display information. |
| `AuditAction` | Reads a department, report ID, and requested status, chooses a report table, and updates the report to Approved or Disputed. |
| `ConsolidatedReport` | `POST` creates a report-document pointer and saves totals plus JSON department breakdown; `GET` reads a report for a date. When no saved record is found, the current code returns example totals. |

**Typical user journey:** department staff submit reports; the manager opens the audit list, reviews the department summaries, and approves or disputes them. A consolidated report stores the overall numbers and a pointer describing where a document is expected to live.

### I. Read one operation from request to response

Use front-desk check-in as a model for understanding unfamiliar backend code:

1. **Find the route** in `backend/cmd/server/main.go`. This tells you which method receives the URL.
2. **Find the method** `CheckIn` in `frontdesk_handler.go`.
3. **Read the request struct.** Its JSON tags show what fields the frontend must send.
4. **Follow each SQL call in order.** The sequence reveals which records are checked, inserted, or updated.
5. **Look for a transaction.** Here, guest, booking, folio, and room changes must agree.
6. **Read every early return.** These explain failures such as invalid JSON, unavailable room, database errors, and commit failures.
7. **Read the response map.** These JSON keys are what the frontend can display next.
8. **Compare with the caller page.** Confirm the frontend sends the expected names and uses the returned values.

This method works for every module. It also exposes unfinished integrations: a screen may collect a sale, but if its handler only changes stock and does not save a sale record, the full business process has not been implemented end to end.