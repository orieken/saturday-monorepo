---
feature: "golden-workspace-packages"
bounded_context: ""
domain_terms: ["SaturdayConfig"]
files_touched: []
issue_refs: ["intake/golden-workspace-packages.md", "#7 (fix/cart-spec-magicshop-import)", "salvage/golden-workspace-packages @ f372a30"]
linked_adrs: []
linked_kis: []
---

# Security Review: golden-workspace-packages

## Threat Model Summary

This feature has no security surface. It switches apps/ye-olde-magic-shop's three framework dependencies (@orieken/saturday-core, @orieken/saturday-cucumber, @orieken/saturday-playwright-heatmap) from vendored file: copies to workspace:* links, removes the now-redundant vendored directories and the app's two Dockerfiles, and fixes a one-line default-value drift (defaultConfig.browser.headless now reads process.env.HEADLESS === 'true' instead of a hardcoded false) in packages/saturday-cucumber/src/config/index.ts. There is no authentication, authorization, user-supplied input, API endpoint, secret/token handling, or trust-boundary crossing anywhere in the diff, matching the spec's own stated Trust Boundaries ('None — no auth, no user input, no external calls'). The HEADLESS env var is a test-runner-local launch flag read at process start, not data crossing into or out of the application, and process.env access in Node is the ordinary mechanism for configuring a local CLI/test tool. pnpm-lock.yaml's diff (6 insertions, 117 deletions) is purely the removal of duplicate vendored file: resolutions; no new external package or version was introduced. No fixes were required or applied.

## Dependency Audit

- pnpm-lock.yaml diff is net-negative (117 lines removed, 6 added) and contains no new external package or version — it only removes resolutions for the deleted vendored copies. No new dependency surface was introduced by this feature.
- Ran `pnpm audit --audit-level=high` against the full monorepo: 130 pre-existing vulnerabilities (8 critical, 59 high, 54 moderate, 9 low) exist repo-wide, none introduced by this change. Spot-checked the two largest findings: vite (via packages/saturday-core's vitest devDependency) and body-parser (via apps/mock-api's express dependency) — both are pre-existing transitive devDependencies of packages/apps this feature does not touch or newly depend on, and neither sits in this feature's scope.paths. No Critical or High CVE was introduced by this feature; none blocks QA.

## STRIDE Analysis

None

### Spoofing

No authentication, session, or identity mechanism exists anywhere in this diff. Not applicable.

### Tampering

No data crosses a trust boundary. The only runtime input considered is process.env.HEADLESS, a local test-runner launch flag read at process start by the machine running the test suite, not externally supplied or attacker-controlled in any deployed sense. Comparison is a strict `=== 'true'` string check with a safe (headed) default on any other value, so there is no injection or coercion path.

### Repudiation

No user-facing or audit-relevant action is performed by this change. Not applicable.

### Information Disclosure

Checked both what can be read and what can carry it out: no secrets, credentials, or PII are introduced, read, or logged by the config default change or the dependency-resolution switch. Removing the vendored directories and both Dockerfiles reduces surface (one less place secrets/config could drift) rather than adding any egress path. No outbound network call is introduced.

### Denial of Service

No new loop, external call, or unbounded resource use is introduced. The dependency switch and config default change affect test-runner configuration only, not a running service.

### Elevation of Privilege

No permission, role, or authorization check exists in this diff. Not applicable.

## Findings

None

## Files Modified

None

## Security Checklist

- [x] No hardcoded secrets, API keys, tokens, or passwords introduced
- [x] No new dependency with a Critical/High CVE introduced (pnpm audit run against full repo; pre-existing findings unrelated to this feature's changed files)
- [x] No auth/session/authorization logic touched — Feature has no auth surface
- [x] No user-supplied input handling introduced — process.env.HEADLESS is a local launch flag, not user input
- [x] No OTel/logging change that could leak PII or secrets — No telemetry code touched
- [x] Dependency direction (app -> framework packages, never reverse) verified clean, consistent with code-reviewer's independent check

## Notes for QA

- No security findings to verify. Proceed directly with the functional/regression verification already planned (AC3 before/after golden-suite parity, AC5 scripts/ui-suite-in-chrome.sh exit code).

## Notes for Tech Writer

None
