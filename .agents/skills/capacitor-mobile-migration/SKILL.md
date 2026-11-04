---
name: capacitor-mobile-migration
description: "Use Capawesome-informed guidance to migrate this existing SvelteKit web app into native Android or iOS apps with Capacitor. Activate only when the user explicitly asks to build, migrate, or package the app as a mobile application; do not use for ordinary web changes or assume mobile conversion is authorized."
license: MIT
metadata:
  author: Capawesome Team; workspace adaptation
  source: https://github.com/capawesome-team/skills
---

# Capacitor Mobile Migration

## Activation Gate

- Use this skill only when the user explicitly asks to make this project a native mobile app, migrate it to Capacitor, or build/install an Android or iOS application.
- A request to install, discuss, or plan this skill does not authorize implementation of the migration.
- Do not add Capacitor dependencies, edit project configuration, create native projects, or run platform builds until the user explicitly requests the mobile implementation.
- Do not use this workflow for routine SvelteKit web development.

## Migration Workflow

1. **Confirm scope.** Identify the requested platforms (Android, iOS, or both), intended application name and reverse-domain app ID, available SDKs, and whether the user wants a migration plan or implementation. Ask only for decisions that block safe progress.
2. **Inspect before changing.** Review `package.json`, SvelteKit configuration, Vite configuration, routes, server-only code, authentication/API calls, and existing build output. Check for `capacitor.config.*`, `android/`, and `ios/` before assuming Capacitor is absent.
3. **Check web build compatibility.** Capacitor packages built web assets into a native WebView. Determine whether the current SvelteKit app can be exported statically with its actual routes and data requirements. Identify server-side rendering, server routes, form actions, dynamic routes, or runtime dependencies that require an online backend or architectural changes. Do not silently replace server behavior with static output.
4. **Agree on necessary changes.** Explain any incompatible SvelteKit features and propose the smallest viable migration. Preserve the existing web deployment and behavior unless the user approves changes that affect them.
5. **Add Capacitor deliberately.** Choose mutually compatible `@capacitor/core`, CLI, and platform package versions based on the project and supported Node.js version. Configure the correct built asset directory as `webDir`, then create only the platform projects the user requested. Build the web app before syncing its assets to native projects.
6. **Adapt native behavior.** Review navigation, external links, status bars and safe areas, app lifecycle, file handling, permissions, and any required native plugins. Request permission before adding paid, account-backed, or externally hosted integrations.
7. **Protect secrets and data.** Never bundle database connection strings, service-role keys, or other server credentials in the mobile client. Keep privileged operations behind the backend and validate/authenticate them server-side.
8. **Validate the requested target.** Run the web checks and build first, then the narrowest available Capacitor sync/build or device check for the requested platform. Clearly report unavailable SDKs, signing requirements, device-only steps, and any unverified behavior.

## MCP and External Documentation

- Follow the workspace `mcp-usage-guard`: do not access an MCP server without first explaining why it is necessary and getting explicit user approval for that access.
- Prefer local project configuration and bundled knowledge. Do not install or configure Capawesome or Capacitor MCP servers as part of this migration skill.
- If current external documentation is essential, ask before fetching it and request only the relevant Capacitor/SvelteKit topic.

## References

This focused workflow is informed by Capawesome's [`capacitor-app-creation`](https://github.com/capawesome-team/skills/tree/main/skills/capacitor-app-creation) and [`capacitor-app-development`](https://github.com/capawesome-team/skills/tree/main/skills/capacitor-app-development) guidance. The creation skill targets new apps and is not copied as a migration procedure.
