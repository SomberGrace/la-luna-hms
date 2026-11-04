---
name: security-engineering
description: "Use for defensive application security reviews and hardening: secure coding, authorization checks, secret handling, dependency risk, and OWASP-aligned mitigations."
---

# Defensive Security Engineering

## Workflow

1. Define the assets, trust boundaries, user roles, and security-sensitive data in scope.
2. Review changed code and direct call paths for input validation, authentication, authorization, data exposure, secret handling, and safe error behavior.
3. Prioritize concrete risks by impact and likelihood. Cite the affected code and explain the unsafe condition without speculative claims.
4. Recommend minimal mitigations that fit the architecture; preserve existing behavior unless a security fix requires a deliberate change.
5. Check relevant dependency advisories, configuration, and tests using the repository's existing tools.
6. Verify fixes with focused tests and safe static checks. Clearly state what was and was not tested.

## Safety Boundaries

- Keep work authorized and defensive. Do not probe or exploit third-party systems or build autonomous attack workflows.
- Use synthetic or local test data; never expose credentials, personal data, or production secrets.
- Prefer secure defaults, least privilege, parameterized queries, validated input, and explicit authorization at server boundaries.
- Do not claim a system is secure based on a narrow review or a passing test suite alone.
