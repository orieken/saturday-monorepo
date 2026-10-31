---
feature: "ui-suite-real-chrome"
bounded_context: ""
domain_terms: ["Runner image"]
files_touched: []
issue_refs: ["ui-suite-real-chrome"]
linked_adrs: []
linked_kis: []
---

# Security Review: ui-suite-real-chrome

## Threat Model Summary

This change wires a local/CI shell script (scripts/ui-suite-in-chrome.sh) to an already-verified container image to run a Playwright suite in real Chrome under Xvfb. There is no authentication, authorization, API surface, user-supplied input, or secret handling introduced. The one named trust boundary — the VNC viewer must stay off and never print its password — is satisfied structurally: the script never sets or references RUNNER_VIEWER, and docker run only forwards RUNNER_SLOW_MO/RUNNER_VIDEO into the container, so the viewer stays unset regardless of the host environment. All docker-exec payloads use fixed, non-interpolated command strings; the only variables passed through (RUNNER_SLOW_MO, RUNNER_VIDEO, CONTAINER_NAME via $$) are quoted and never built from untrusted input. No new dependencies were added. Two residual observations are recorded below as Low/Info — neither blocks QA.

## Dependency Audit

- No new dependencies introduced by this change (implementation-notes.md confirms 'Dependencies Added: None') — pnpm audit was not run since there is nothing new to audit.

## STRIDE Analysis

None

### Spoofing

No authentication, session, or identity surface exists in this script or the test files it adds; nothing to spoof.

### Tampering

docker exec payloads are fixed command strings, not built from untrusted input. RUNNER_SLOW_MO/RUNNER_VIDEO pass through quoted docker -e flags and are consumed via Number()/strict-equality checks in playwright.config.ts, not shell-interpolated. pnpm install uses --frozen-lockfile, preventing silent lockfile drift.

### Repudiation

Local/CI dev tooling with no user-attributable actions to audit; not applicable at this layer.

### Information Disclosure

Checked whether the VNC viewer or its password could leak — RUNNER_VIEWER is never set or referenced by this script, satisfying the spec's trust boundary. Checked the app.log tail printed on health-check failure (wait_for_app_healthy) — apps/ye-olde-magic-shop/.env contains no secrets today, only non-sensitive dev flags, so the tail cannot currently disclose credentials; flagged as Info for future-proofing. Checked outbound destinations — none are introduced; the container's own network egress is a pre-existing, out-of-scope image property, not something this diff changes.

### Denial of Service

wait_for_app_healthy is bounded by the named APP_HEALTH_TIMEOUT_SECONDS (180s) constant with curl --max-time 5 per probe, matching the NFR, and exits a distinct non-zero code (EXIT_APP_UNHEALTHY) rather than hanging. No unbounded loops over user-controlled data exist in this script.

### Elevation of Privilege

Container runs as non-root USER runner (containers/saturday-runner/Containerfile, unmodified, out of scope). docker run uses no --privileged, no extra capabilities, and no host network/PID namespace sharing. The repo bind-mount (-v REPO_ROOT:/work) gives the container write access to the checkout, which is the intended, expected behavior of a local test runner the developer already chose to execute, not a privilege boundary this diff crosses.

## Findings

### [LOW] — Container network egress is not structurally restricted during the suite run
**Location**: `scripts/ui-suite-in-chrome.sh:39 (start_container)`
**Threat category**: Information Disclosure
**Description**: docker run uses the default bridge network rather than --network=none or an egress allowlist. The spec's NFR states the app and mock API make no request off the container after dependencies are installed, but that is an application-behavior assumption, not a network-policy guarantee — a compromised dependency pulled during pnpm install, or a regression in the app/mock-api, could reach the network undetected. This is a pre-existing property of how the runner image and this script are used together, not something this diff newly introduces, and restricting it would require reconfiguring container networking per-phase (install vs. test), which is an architectural change outside this spec's scope (containers/saturday-runner/ is explicitly out of scope).
**Fix applied**: Recommendation only
**Verification**: Not applicable — no fix applied; recorded as a defense-in-depth gap for a human to weigh, not blocking QA.
### [INFO] — Health-check failure path dumps the app's dev-server log verbatim to stderr
**Location**: `scripts/ui-suite-in-chrome.sh:81 (wait_for_app_healthy)`
**Threat category**: Information Disclosure
**Description**: On timeout, the script runs `tail -n 50 /tmp/app.log` and prints it to stderr. Today apps/ye-olde-magic-shop/.env carries only non-sensitive dev flags (ENABLE_OTEL, OTEL endpoint, HEADLESS, BROWSER, ENVIRONMENT), so there is nothing sensitive to leak. If a secret is ever added to that .env or logged by the dev server, this tail would surface it in CI console output. No action needed now; worth a one-line note for whoever next touches this script's error path.
**Fix applied**: Recommendation only
**Verification**: Not applicable — informational only.

## Files Modified

None

## Security Checklist

- [x] VNC viewer stays off and its password is never printed (spec Trust Boundaries) — RUNNER_VIEWER is never set or referenced by scripts/ui-suite-in-chrome.sh; docker run only forwards RUNNER_SLOW_MO and RUNNER_VIDEO.
- [x] No secrets, tokens, or credentials hardcoded in the script
- [x] No shell command built from untrusted/user-controlled input (injection check) — All docker exec payloads are fixed strings; the only pass-through variables (RUNNER_SLOW_MO, RUNNER_VIDEO) are consumed via typed checks in playwright.config.ts, not shell-interpolated.
- [x] Network calls have explicit timeouts — curl --max-time 5 inside wait_for_app_healthy; bounded by APP_HEALTH_TIMEOUT_SECONDS overall.
- [x] Container runs without elevated privileges (--privileged, extra capabilities, host namespaces)
- [x] New dependencies scanned for known CVEs — No new dependencies were introduced, so there was nothing to scan.
- [ ] Container network egress restricted for defense-in-depth — Default bridge network is used; not a regression from this diff and out of scope to fix here (see LOW finding).

## Notes for QA

- No Critical or High findings — nothing blocks QA from proceeding.
- Two Low/Info findings recorded for awareness only: (1) container network egress isn't structurally blocked during the run, it relies on the app/mock-api not making outbound calls, matching the NFR as a behavioral assertion QA can verify but not a technical guarantee; (2) the health-check failure path tails and prints /tmp/app.log to stderr, which is safe today because .env has no secrets, but would need redaction if that ever changes.
- AC3's test-of-the-test (re-run tests/webgl.spec.ts with --disable-gpu --disable-software-rasterizer and confirm it fails) is QA's responsibility per analysis.md — code-reviewer already confirmed the assertion itself is real (checks an actual webgl2 context, not a tautology).

## Notes for Tech Writer

None
