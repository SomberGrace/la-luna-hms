---
name: playwright
description: "Use when testing a web application in a real browser with Playwright, checking responsive layouts, or automating end-to-end user flows and regressions."
---

# Playwright Browser Testing

## Workflow

1. Inspect the project's existing test scripts, Playwright configuration, and browser-testing conventions before adding or running tests.
2. Start the application using its documented development or test command and wait for its ready signal.
3. Test the user's real workflow through visible UI controls. Prefer role, label, and text locators over brittle CSS or XPath selectors.
4. Verify observable outcomes, including navigation, validation, loading, empty, and error states that matter to the scenario.
5. Check responsive behavior at relevant viewport sizes and capture screenshots or traces when they help diagnose a failure.
6. Keep tests deterministic: isolate data, avoid arbitrary sleeps, and wait for meaningful page or network conditions.
7. Stop any server or browser process started for the test. Report the exact command and outcome; distinguish app failures from environment or browser-installation failures.

## Safety and Reliability

- Use the project's configured test environment and test data; do not submit real customer information or perform irreversible actions against production.
- Do not add a new browser-testing dependency if an existing test tool already covers the need.
- Make assertions against user-visible behavior rather than implementation details wherever practical.
