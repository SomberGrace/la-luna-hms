---
name: rename-symbol
description: "Rename identifiers such as functions, types, fields, and variables across a codebase. Use for symbol-renaming requests; distinguish source identifiers from public names and stored values."
---

# Rename Symbol

## Workflow

1. Identify the symbol's declaration, references, and project conventions. Use language-aware rename tools when available; otherwise trace references carefully.
2. Check whether similarly named strings are part of an API, serialized format, database schema, environment variable, configuration key, or persisted value. Do not change those as a side effect of renaming a source identifier.
3. If the requested rename affects a public contract or stored data, explain the impact and confirm the intended scope before changing it.
4. Apply the smallest complete rename, updating directly related tests and documentation. Keep the new name idiomatic for the language and consistent with nearby code.
5. Review the diff for missed or unintended changes, then run focused tests and type or build checks for affected code.
