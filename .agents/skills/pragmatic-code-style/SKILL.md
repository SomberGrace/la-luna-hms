---
name: pragmatic-code-style
description: "Write and revise clear, idiomatic code with straightforward structure, contextual names, and concise comments. Use when implementing or refactoring code for readability and maintainability."
---

# Pragmatic Code Style

## Guidelines

- Follow the language and repository conventions already in use. Prefer familiar names that are clear in context; use concise names like `err`, `req`, or `idx` when the scope makes them obvious.
- Solve the immediate problem with the simplest design that handles the real requirements. Avoid extra layers, patterns, or abstractions that do not make the code easier to understand or change.
- Prefer early returns and straightforward conditionals over deeply nested ternaries or clever control flow.
- Keep formatting consistent with the surrounding files and the project's formatter. Do not introduce irregular spacing or inconsistent naming just to make code look hand-written.
- Add comments only when they explain a non-obvious reason, constraint, or workaround. Keep them brief and do not restate what the code already says.
- Preserve behavior unless the task calls for a behavior change. Make focused edits and avoid unrelated cleanup.
- Use existing types, helpers, tests, and error-handling conventions. Do not hide failures behind silent defaults or broad catches.
- After editing, review the diff for unnecessary complexity and run the relevant tests, lint, or type checks.
