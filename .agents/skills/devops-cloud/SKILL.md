---
name: devops-cloud
description: "Use when changing Docker, CI/CD, deployment, cloud infrastructure, observability, or managed database configuration. Prefer least privilege and safe, repeatable operations."
---

# DevOps and Cloud Engineering

1. Inspect existing deployment files, environments, secret handling, and operational conventions before changing infrastructure.
2. Keep configuration explicit and environment-specific. Store credentials in a secret manager or deployment environment settings, never in source control.
3. Prefer least-privilege identities, encrypted connections, private networking where appropriate, and bounded service permissions.
4. Make infrastructure changes reproducible and reviewable; avoid unrequested destructive changes to data or live resources.
5. Include health checks, useful logs, and observable failure states without exposing credentials or sensitive data.
6. Validate configuration with available formatters, linters, dry runs, or local container checks before deployment.
7. Report required manual credentials, approvals, and deployment actions. Do not claim a production deployment was verified if it was not.
