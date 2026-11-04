---
name: typescript-pro
description: "Use when writing, debugging, reviewing, or restructuring TypeScript. Apply strict, expressive types, modern language features, and maintainable application architecture."
---

# TypeScript Pro

## Workflow

1. Check the project's TypeScript version, compiler options, framework conventions, and nearby types before designing a change.
2. Model the domain with precise interfaces, discriminated unions, generics, and literal types where they improve correctness.
3. Treat untrusted data as `unknown`; validate it at boundaries before use. Avoid `any`, unsafe assertions, and non-null assertions when a guard or better type can express the invariant.
4. Keep modules cohesive and APIs explicit. Reuse shared types and helpers rather than creating duplicate or overly generic abstractions.
5. Use async control flow, error handling, and resource cleanup that make failure behavior explicit.
6. Add or update tests for relevant edge cases and run the narrowest applicable test, lint, and type-check commands.

## Standards

- Respect existing public contracts and runtime behavior unless a change is requested.
- Prefer inference for obvious local values and explicit types at module boundaries.
- Do not suppress compiler errors to force a build through; fix the underlying type or document a justified exception.
