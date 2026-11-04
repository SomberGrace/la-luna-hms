---
name: docker
description: "Use when creating or changing Dockerfiles, Docker Compose services, container builds, networking, volumes, or container troubleshooting for this Go and SvelteKit application."
---

# Docker Engineering

Make container changes reproducible, secure, and consistent with the repository. Keep investigation focused on the files and runtime behavior involved.

## Repository context

- `docker-compose.yaml` currently defines PostgreSQL only; do not assume the app is already containerized.
- The backend is Go and serves the API on port 8080 in the documented local workflow.
- The frontend is SvelteKit/Vite using the static adapter. Its API base URL is read at build time; the Vite `/api` proxy is a development-server feature, not a production reverse proxy.
- Preserve existing local workflows unless the requested change intentionally replaces them.

## Workflow

1. Inspect the relevant Docker/Compose files, application entry points, ports, environment variables, and build scripts before editing.
2. Establish the intended environment (local development, CI, or production). Keep environment-specific settings explicit and avoid turning development defaults into production configuration.
3. Implement the smallest complete change:
   - Use multi-stage builds where useful; pin base image versions and copy only required build artifacts.
   - Run application processes as non-root where practical.
   - Add a `.dockerignore` that excludes dependencies, generated output, VCS metadata, and secrets without excluding required build inputs.
   - Use Compose service DNS names for container-to-container connections; `localhost` inside a container refers to that container.
   - Add health checks and readiness-based dependencies when they reflect actual service readiness.
   - Persist database state with named volumes. Do not expose ports or grant privileges more broadly than needed.
4. Keep credentials out of Dockerfiles, committed Compose values, image layers, and logs. Use environment injection or a secret manager appropriate to the target environment; do not print resolved secret values.
5. Validate with the smallest relevant checks, typically `docker compose config` and targeted `docker compose build` or service smoke tests. Be mindful that rendered configuration may contain interpolated secrets.
6. Never run destructive volume, prune, or cleanup commands (including `docker compose down -v`) without explicit user authorization. Do not claim a production deployment was tested when only local checks ran.
7. Summarize changed files, validation performed, and any required manual secrets or deployment steps.
