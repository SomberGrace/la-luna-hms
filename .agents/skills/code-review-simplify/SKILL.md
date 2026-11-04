---
name: code-review-simplify
description: "Use when reviewing a change for correctness or simplifying code before review. Prioritize concrete bugs, regressions, security, and unnecessary complexity; keep findings evidence-based."
---

# Code Review and Simplification

## Review Workflow

1. Read the requested behavior, the complete relevant diff, and nearby call sites or tests needed to understand it.
2. Trace changed behavior through its callers and data boundaries; check error paths, authorization, input validation, and compatibility where relevant.
3. Report only actionable issues introduced by the change. Rank by impact and cite precise file and line locations.
4. Explain the triggering condition and consequence for each finding. Do not report speculative concerns as confirmed bugs.
5. If no actionable issues are found, say so plainly and note any validation limits.

## Simplification Workflow

1. Preserve externally visible behavior and error semantics.
2. Remove duplication, dead branches, and indirection only when their redundancy is clear.
3. Prefer direct control flow and existing project helpers; avoid clever compression and speculative abstractions.
4. Re-run the narrowest relevant tests, lint, and type checks after edits.

Do not silently reformat unrelated code or make behavior-changing edits during a review-only request.
