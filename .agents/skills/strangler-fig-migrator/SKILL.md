---
name: strangler-fig-migrator
description: "Plan or implement a gradual strangler-fig migration from a legacy component to a replacement. Use when replacing a service, module, API, or system incrementally while keeping the application working."
---

# Strangler-Fig Migrator

## Migration workflow

1. Map the existing component's responsibilities, callers, data ownership, external contracts, and operational constraints. Identify the smallest useful seam for migration.
2. Agree on the migration boundary and success criteria before changing public behavior or persistent data. Prefer a narrow, reversible slice.
3. Keep the existing implementation as the working fallback while routing only the selected slice to the replacement. Make routing explicit and observable; do not add a flag or abstraction without a concrete rollout need.
4. Preserve API and data compatibility during coexistence. Define how reads and writes behave, how data stays consistent, and how failures fall back before routing real traffic.
5. Add focused tests for both paths, boundary behavior, and rollback or fallback behavior. Use existing tooling and avoid large rewrites.
6. Move additional slices only after the current one is verified. Track what remains on the legacy path and remove old code only after traffic, data, and dependencies have moved safely.

Keep migration steps small and independently reviewable. Surface unresolved compatibility or data-migration risks instead of assuming they are safe.
