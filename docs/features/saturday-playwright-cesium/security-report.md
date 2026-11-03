---
feature: "saturday-playwright-cesium"
bounded_context: ""
domain_terms: []
files_touched: ["packages/saturday-playwright-cesium/package.json", "pnpm-lock.yaml"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Security Review: saturday-playwright-cesium

## Threat Model Summary

This is local-only Playwright test fixture infrastructure for E2E-testing CesiumJS globes — a dev dependency, never shipped to production, with no auth, no user accounts, no persisted data, and no network exposure beyond a server bound explicitly to 127.0.0.1. The only component that parses input shaped like it crosses a trust boundary is fixtures/serve.mjs, a static file server invoked by Playwright's webServer during test runs; it serves Cesium's bundled assets and the fixture's own planet texture. STRIDE was applied with that server, the pure geometry module (src/geometry.ts), and the Cesium page-evaluate seam (src/fixture.ts) as the surfaces of interest. The one finding worth recording was a dependency-version drift (vitest pinned to an old major pulling in a vite release with a known high-severity advisory) rather than a code-level flaw; it has been fixed directly. The static file server's path-resolution guard was verified by hand with both raw and URL-encoded `../` traversal attempts and correctly returns 404 for both while still serving legitimate assets.

## Dependency Audit

- pnpm audit --audit-level=high (full workspace) flagged packages__saturday-playwright-cesium>vitest>vite as High: 'launch-editor fs.existsSync() bypass on Windows' (GHSA-587h-7wmj-h5pm), caused by this package pinning vitest ^2.1.9 while every other package in the monorepo (saturday-core, saturday-cucumber, saturday-cucumber-indexer, saturday-ml-analyzer, saturday-playwright, saturday-playwright-heatmap) pins ^4.0.16.
- Fixed: bumped packages/saturday-playwright-cesium/package.json's vitest devDependency from ^2.1.9 to ^4.0.16 to match the rest of the monorepo, ran pnpm install, and confirmed via pnpm audit that this package no longer appears in any High or Critical finding. Re-ran `pnpm vitest run` under the new version — all 7 geometry.spec.ts tests pass unchanged.
- No other new dependency (cesium@1.146.0, @playwright/test) produced a High or Critical finding in the scoped audit output.

## STRIDE Analysis

None

### Spoofing

No authentication, sessions, or identity of any kind exists in this package — it is a local test-fixture server and a Playwright page-evaluate seam. Not applicable.

### Tampering

The only externally-reachable input is the HTTP request path to fixtures/serve.mjs. resolveWithinRoot() resolves the path against each route's root directory and rejects any result that does not start with that root (with a trailing separator, closing the classic prefix-bypass gap). Verified by hand: both `GET /assets/../../../../etc/passwd` and its %2e%2e-encoded form return 404; legitimate assets under /assets/ and /cesium/ still return 200. decodeURIComponent() runs before path resolution, so an encoded traversal attempt is normalized and still caught by the boundary check rather than slipping past it. src/geometry.ts and src/fixture.ts take only numeric lat/lng/height arguments from the test author's own code, not from any external or remote source.

### Repudiation

Not applicable — there are no user actions to attribute or deny; this is test tooling run by a developer against their own machine.

### Information Disclosure

(egress) The fixture page's CESIUM_BASE_URL is hardcoded to the local '/cesium/' route and the imagery provider loads only the local fixtures/assets/aurelis.jpg — both enforced by a code comment stating 'never a CDN.' No outbound network call exists anywhere in src/ or fixtures/. Nothing in this package can be induced to send data anywhere.

### Denial of Service

serve.mjs has no request size limits or rate limiting, but it is a local, single-purpose static file server with no write operations, no user uploads, and no internet exposure (127.0.0.1 bind) — a DoS surface here has no attacker to exploit it. waitForSettled()/pollUntilSettled() in fixture.ts take an explicit, caller-supplied timeoutMs with a bounded default (15000ms) rather than an unbounded wait.

### Elevation of Privilege

No roles, permissions, or privilege levels exist in this package. Not applicable.

## Findings

### [MEDIUM] — vitest pinned to an outdated major version, pulling in a vite release with a known High-severity advisory
**Location**: `packages/saturday-playwright-cesium/package.json:18`
**Threat category**: Information Disclosure
**Description**: vitest ^2.1.9 resolved to vite <=6.4.3, which carries GHSA-587h-7wmj-h5pm (launch-editor fs.existsSync() bypass on Windows, exploitable only through vite's dev-mode launch-editor middleware). Every other TypeScript package in the monorepo already pins vitest ^4.0.16, so this was version drift introduced by the new package rather than an intentional choice. Low real-world exploitability (dev-only, Windows-specific, requires opening a stack trace in-editor) but it is an unforced, easily-avoided inconsistency, so rated Medium rather than Low.
**Fix applied**: Bumped the vitest devDependency to ^4.0.16 in packages/saturday-playwright-cesium/package.json and ran pnpm install, which also removed the now-unused vite 5 / esbuild 0.21.5 chain from pnpm-lock.yaml.
**Verification**: `pnpm audit --audit-level=high` no longer lists this package in any finding's Paths column. `pnpm vitest run` still passes all 7 tests in tests/geometry.spec.ts under vitest 4.0.16.
### [INFO] — Static file server's path-traversal guard confirmed correct by direct test
**Location**: `packages/saturday-playwright-cesium/fixtures/serve.mjs:39`
**Threat category**: Tampering
**Description**: resolveWithinRoot() uses resolve() plus a trailing-separator boundary check, which closes both the raw `../` traversal path and the classic sibling-directory prefix-bypass bug (e.g. root '/foo/bar' matching requests for '/foo/barbaz'). Manually verified against both a raw-encoded and percent-encoded traversal attempt; both returned 404 while legitimate asset requests returned 200. No fix needed — documenting the verification so QA does not need to re-derive it.
**Fix applied**: Recommendation only
**Verification**: 

## Files Modified

- packages/saturday-playwright-cesium/package.json
- pnpm-lock.yaml

## Security Checklist

- [x] No hardcoded secrets, tokens, or credentials anywhere in the package
- [x] No authentication, session, or authorization surface exists (verified not applicable)
- [x] HTTP server binds to 127.0.0.1 only, not 0.0.0.0 or all interfaces
- [x] Static file server rejects path traversal (raw and percent-encoded), verified by direct request
- [x] No outbound/egress network calls; imagery and Cesium assets are served locally only, never a CDN
- [x] Dependency audit run (pnpm audit --audit-level=high) against the full workspace and scoped to this package's dependency paths — One Medium finding (vitest version drift) found and fixed; re-audited clean for this package afterward.
- [x] No new dependency introduces a Critical-severity advisory

## Notes for QA

- No security-driven test additions are required — this package has no auth, input validation, or trust-boundary logic that QA needs to add coverage for beyond what fixture.spec.ts and geometry.spec.ts already exercise.
- vitest was bumped from ^2.1.9 to ^4.0.16 to match the rest of the monorepo; re-run the full test command (`pnpm test` in packages/saturday-playwright-cesium) to confirm both the vitest suite and the Playwright fixture suite still pass under the new version — the vitest-only geometry suite was already re-verified in this review.
- The fixture server's traversal guard was exercised manually with curl during this review (raw `../` and percent-encoded `%2e%2e%2f`); no automated regression test exists for it. Consider adding one if this server's routing logic is ever extended beyond its current two fixed routes.

## Notes for Tech Writer

None
