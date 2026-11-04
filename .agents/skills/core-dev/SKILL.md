---
name: core-dev
description: "Use for general software development work that benefits from disciplined planning, debugging, tests, linting, and safe version-control practices."
---

# Core Development Workflow

1. Identify the requested outcome, scope, and acceptance criteria. Clarify only consequential ambiguity.
2. Inspect relevant files, history or conventions, callers, and tests before editing; avoid unrelated exploration.
3. Make the smallest complete change that fits the existing architecture and preserves unrelated behavior.
4. Add or update focused tests for the requested behavior, including important failure and boundary cases.
5. Run the narrowest relevant test, lint, type-check, or build. Fix regressions introduced by the change.
6. Review the final diff for accidental edits, secrets, generated files, and unrelated changes.
7. Summarize the result and validation accurately. Never claim a check passed if it was not run successfully.

Do not commit, push, delete user data, or run destructive commands unless the user specifically asks.
