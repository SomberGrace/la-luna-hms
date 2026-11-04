---
name: postgres-supabase-data
description: "Use when changing or debugging PostgreSQL or Supabase SQL, database-backed app behavior, schema, persistence, or queries in this La Luna HMS project."
---

# PostgreSQL and Supabase Data

## Project Map

- Supabase is the hosted PostgreSQL database for this app. Application SQL uses PostgreSQL, not the Supabase JavaScript client.
- The Go backend connects through `github.com/lib/pq` and `database/sql`.
- Database configuration is read from `DATABASE_URL`, with `SUPABASE_DB_URL` as a fallback. Local `.env` loading is implemented in `backend/internal/config/config.go` and separately in `backend/cmd/frontdesk/main.go`. Do not read, repeat, log, or expose secret values from `.env`.
- `database/schema.sql` defines the PostgreSQL types, tables, foreign keys, generated columns, and indexes.
- The front-desk API and SQL handlers are in `backend/cmd/frontdesk/main.go`; login queries `user_accounts` through `backend/internal/services/auth_service.go` and is initialized in `backend/cmd/server/main.go`.
- The front-desk process reports its database state at `/api/v1/frontdesk/health`. A response reporting `In-Memory Cache` means front-desk operations may be using transient memory, not persisted Supabase data. The separate `/api/v1/health` endpoint only reports server health; it does not prove a database connection.

## Trace Before Changing

1. Follow the relevant frontend action to its API endpoint, Go handler, SQL statements, and schema constraints. Search callers and tests rather than assuming every screen is database-backed.
2. Verify whether the handler uses PostgreSQL, `memDB`, or both. Check all success, error, and fallback paths.
3. Use the actual schema definitions for table names, enum values, nullability, generated columns, foreign keys, and delete behavior. Do not infer the schema from frontend types or assume Supabase Auth is in use.
4. Preserve established API contracts unless the user requested a change. If persistence differs between database and memory paths, explain the difference and address it only within scope.
5. Add or update focused tests for the handler/service behavior and validate with the project's available Go checks.

## Query and Transaction Practices

- Use parameterized PostgreSQL placeholders (`$1`, `$2`, ...) for all external values. Never build SQL by concatenating request data.
- Scan query results into explicit Go types; close result rows and check `rows.Err()` after iteration. Handle `sql.ErrNoRows` separately from connection or query failures.
- Check errors from `Begin`, every query and write, `Commit`, and `Rollback`. Do not ignore database errors or return success if a write failed.
- Put related writes that must succeed together in one transaction. Roll back on intermediate failure and verify commit before reporting success.
- Use constraints as the final integrity boundary; still validate input and authorization in the backend.
- Keep SQL portable within PostgreSQL and consistent with schema enum spellings. Respect generated columns: do not explicitly insert or update columns PostgreSQL computes.
- For concurrent updates such as payments, protect read-modify-write operations with a transaction and an appropriate row lock or atomic SQL update; avoid lost updates.
- Prefer explicit columns over `SELECT *`. Keep result sets bounded when they can grow.
- Distinguish an intentional offline/demo mode from a successful database operation. Never claim that in-memory writes persisted to Supabase.

## Schema and Data Safety

- Treat `database/schema.sql` as destructive: it begins by dropping application tables with `CASCADE`. Never execute it against a populated, shared, or production Supabase database without an explicit user request, a verified backup, and a reviewed migration plan.
- Prefer forward-only, narrowly scoped migrations for existing databases. Do not rewrite or reset live data to make a local schema match.
- Before any write, migration, seed, delete, truncate, or other database-changing operation, explain the target and impact and get explicit user approval. Do not connect to or inspect the live database unless the user explicitly authorizes that access.
- Use synthetic/local data for testing. Never place database credentials, Supabase service-role keys, or personal guest data in source code, frontend bundles, logs, or prompts to external services.
- The browser/frontend must call the backend for privileged data access. Never expose a database URL or service-role key to the client.

## Debugging Checklist

1. Confirm which backend process and port are serving the requested endpoint.
2. Check the sanitized startup message and `/api/v1/frontdesk/health` for the front-desk connection state; never ask the user to paste a full connection URL.
3. For DNS, TLS, authentication, or schema errors, diagnose the exact reported failure without revealing credentials.
4. Confirm the query's tables and columns exist in the active schema and that enum values and constraints match.
5. Determine whether a result came from PostgreSQL or the in-memory fallback before treating it as durable data.

## MCP and Documentation

- Follow `mcp-usage-guard`: do not access an MCP server without first explaining why it is necessary and obtaining explicit approval for that specific access.
- Prefer the repository schema and code for app-specific behavior. Use current external documentation only when necessary and authorized.
