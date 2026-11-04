---
name: token-efficient-coding
description: "Use when implementing, debugging, or reviewing code with limited context or token budget. Keep investigation and responses focused without sacrificing correctness, safety, or necessary validation."
---

# Token-Efficient Coding

## Principles

- Optimize for useful work per context token, not for the shortest answer at any cost.
- Preserve correctness, security, maintainability, and the user's requested behavior.
- Avoid repeating information already present in the conversation or repository.

## Workflow

1. Restate the requested outcome internally; ask only when a consequential requirement is ambiguous.
2. Locate relevant files and symbols with targeted searches. Read only the sections needed to understand the change and its callers.
3. Follow existing project patterns and reuse established helpers rather than exploring unrelated areas or duplicating logic.
4. Make a complete, focused change. Avoid unrelated cleanup and unnecessary abstractions.
5. Run the smallest relevant test, lint, or type-check. If validation cannot run, state why rather than implying it passed.
6. Respond with a concise summary of what changed and the validation result. Include caveats or follow-up steps only when useful.

## Context Discipline

- Prefer symbol or filename searches and bounded file reads over dumping directories or entire large files.
- Expand the investigation only when the initial evidence reveals a dependency, failure, or behavior that must be checked.
- Do not omit necessary requirements, edge cases, tests, or safety checks merely to save tokens.
- Do not silently guess when different interpretations would materially change the implementation.
