---
name: dead-code-finder
description: "Find and safely remove unused code, imports, variables, exports, and unreachable branches. Use when asked to find, review, or clean up dead code."
---

# Dead Code Finder

## Workflow

1. Check the project's language tools and conventions. Prefer compiler, linter, and framework diagnostics over guesses.
2. Trace each candidate's references, including routes, registries, reflection, generated code, build tags, and dynamic loading where relevant.
3. Separate confirmed dead code from uncertain candidates. For each candidate, give its location and why it appears unused.
4. Remove code only when asked and when the evidence is strong. A missing text search result alone is not enough; preserve public APIs, persisted data, and generated files unless their removal is explicitly in scope.
5. Run the smallest relevant tests or checks after removal and report what was verified.

Prefer a small cleanup over broad restructuring. Keep names and code style consistent with the surrounding project.
