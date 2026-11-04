---
name: full-stack-development
description: "Use when a feature spans frontend and backend, including API contracts, validation, persistence, authentication, and end-to-end behavior."
---

# Full-Stack Feature Development

1. Trace the user workflow from the UI through its API, service, and persistence layers; identify the authoritative source of data.
2. Define or confirm the request and response contract, including validation, errors, authorization, and loading or empty states.
3. Implement each layer using its existing responsibilities and conventions. Keep business rules out of presentation and transport code where the architecture separates them.
4. Validate and authorize at the server boundary; never rely solely on client-side checks. Keep secrets server-side.
5. Preserve data consistency and handle failures explicitly. Consider transactions and retries only where the data operations require them.
6. Add focused tests at the affected layers and an end-to-end test when the project supports one.
7. Run the narrowest checks for the changed frontend and backend surfaces, then report any untested integration step.
