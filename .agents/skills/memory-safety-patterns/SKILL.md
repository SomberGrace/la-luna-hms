---
name: memory-safety-patterns
description: "Use when working with unsafe code, native interfaces, pointers, shared mutable state, resource lifetimes, or security-sensitive memory handling in any language."
---

# Memory Safety Patterns

## Workflow

1. Identify memory and resource ownership: who creates, borrows, mutates, and releases each object or handle?
2. Prefer safe language constructs and established libraries; isolate unavoidable unsafe operations behind small, documented boundaries.
3. Validate lengths, indexes, integer conversions, nullability, and input encodings before accessing memory or allocating buffers.
4. Ensure resource lifetimes are explicit and cleanup happens on every exit path. Prevent use-after-free, double-free, leaks, and stale handles.
5. For shared mutable data, define synchronization and ownership rules; avoid data races and unsafe aliasing.
6. Treat FFI and deserialization boundaries as untrusted. Validate invariants before exposing data to the rest of the program.
7. Add tests for boundary conditions and use available sanitizers, race detectors, fuzzing, or static analysis where they fit the repository.

## Standards

- Do not introduce unsafe code as a convenience or bypass a type-system guarantee without a demonstrated need.
- Keep externally supplied sizes and counts bounded to prevent overflow and uncontrolled allocation.
- State any invariants required by a low-level helper and enforce them at its boundary.
