---
name: golang-pro
description: "Use when implementing, debugging, or reviewing Go services and libraries. Apply idiomatic Go, clear error handling, concurrency discipline, and maintainable package design."
---

# Go Engineering

## Workflow

1. Read the module's `go.mod`, neighboring packages, and relevant tests to identify the supported Go version and local conventions.
2. Keep package responsibilities focused and APIs small; use standard-library facilities and existing dependencies when they meet the need.
3. Return errors with useful context and handle them at an appropriate boundary. Do not ignore errors or turn failures into success-shaped fallbacks.
4. Make ownership and concurrency behavior explicit. Protect shared state with synchronization or message passing and avoid copying synchronization primitives.
5. Respect `context.Context` for cancellable I/O and long-running operations; close resources and check iterator errors.
6. Add focused tests, including edge cases and concurrency tests where shared state or parallel execution is involved.
7. Run the narrowest relevant `go test` and project lint or static-analysis command. Use the repository's supported Go toolchain.

## Standards

- Prefer simple, idiomatic Go over unnecessary frameworks, reflection, or abstraction layers.
- Validate external input and use parameterized database queries.
- Avoid `panic` for ordinary error handling and avoid data races, goroutine leaks, and unbounded work.
