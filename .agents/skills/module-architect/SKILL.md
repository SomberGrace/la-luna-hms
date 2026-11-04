---
name: module-architect
description: "Use when planning, creating, or restructuring a feature module or subsystem. Define clear responsibilities, boundaries, and contracts that fit the existing architecture before implementation."
---

# Module Architect

## Workflow

1. Clarify the module's purpose, users, inputs, outputs, and expected behavior from the request and existing code.
2. Trace the relevant entry points, callers, data models, persistence, and tests. Keep the investigation limited to the module and its direct dependencies.
3. Propose the smallest coherent boundary that fits the repository's existing architecture and conventions.
4. Define responsibilities and interfaces between the module and its callers. Keep business rules out of presentation and transport layers where the project separates them.
5. Consider validation, errors, authorization, persistence, concurrency, and lifecycle needs that apply to this module.
6. Implement the necessary pieces together, including directly related tests and documentation. Avoid speculative frameworks and unrelated refactors.
7. Verify the module through focused tests or checks, then summarize its responsibilities, integration points, and validation.

## Design Standards

- Prefer cohesive modules with explicit contracts and minimal coupling.
- Reuse existing types, helpers, and architectural patterns unless they prevent the requested behavior.
- Make dependencies and failure behavior explicit; do not hide errors behind success-shaped fallbacks.
- Keep public interfaces small and stable. Add abstractions only when they solve a demonstrated need.
- Preserve existing behavior outside the module's intended changes.
