---
name: code-reviewer
description: "Review code changes for concrete bugs, regressions, and maintainability risks. Use when asked to review a diff, pull request, or implementation; report findings before style preferences."
---

# Code Reviewer

## Review workflow

1. Read the relevant diff and enough surrounding code to understand intended behavior. Check related tests and callers where needed.
2. Look for defects introduced by the change: incorrect behavior, missed edge cases, unsafe assumptions, broken contracts, and missing error handling.
3. Verify each concern against the code. Do not report speculative issues, existing problems outside the change, or style preferences as bugs.
4. Report actionable findings first, ordered by severity. Include file and line, impact, and the condition needed to reproduce the issue. If there are no findings, say so and mention any validation limits.
5. Do not modify files during a review unless the user asks for fixes. When asked to fix findings, keep changes focused and run relevant checks.

Prefer clear, direct findings over exhaustive commentary. Treat security-sensitive changes carefully and use the project's dedicated security-review workflow when the user requests a security review.
