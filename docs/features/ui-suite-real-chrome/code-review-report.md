---
feature: "ui-suite-real-chrome"
bounded_context: ""
domain_terms: ["Runner image"]
files_touched: []
issue_refs: ["ui-suite-real-chrome"]
linked_adrs: []
linked_kis: []
---

# Code Review: ui-suite-real-chrome

## Overall Status

**APPROVED**

## Design Narrative

The diff wires the demo app's Playwright suite to the already-built, already-verified Saturday runner image: a new 'chrome' project consumes runner-check's literal ANGLE/SwiftShader launch-arg contract without re-deriving it, webgl.spec.ts asserts the resulting webgl2 context against the live page, and scripts/ui-suite-in-chrome.sh orchestrates build/run/health-check/runner-check/suite/report with a reserved, named exit-code namespace (2-5) that never remaps Playwright's own exit code. It is purely additive test/tooling surface consuming the runner image as a read-only dependency, exactly as architecture-notes.md specified. Behaviour change: no.

## Design Score

- **Clarity** 5/5
- **Cohesion** 5/5
- **Coupling** 4/5
- **Craft** 5/5

## Security Surface

- RUNNER_VIEWER is never set and never referenced by the script, matching the spec's trust-boundary requirement that the VNC viewer stay off and its password never printed.
- No auth, API, user input, or secret handling is introduced; no security-reviewer trigger per CLAUDE.md's criteria list.

## Performance Surface

- wait_for_app_healthy has an explicit, named 180s bound (APP_HEALTH_TIMEOUT_SECONDS) and the curl checks inside it use --max-time 5, matching the NFR.
- No N+1 queries, pagination, or synchronous fan-out concerns apply to this script/config/test surface.

## Test Design Review

- apps/ye-olde-magic-shop/tests/webgl.spec.ts — YES: the only assertion checks a real canvas.getContext('webgl2') result against the actual running page, not a mock or tautology (e.g. toBeDefined()); it fails if WebGL is genuinely unavailable. AC3's test-of-the-test (re-running with --disable-gpu --disable-software-rasterizer) is QA's task per analysis.md, not a gap in this test's own design.
- No other test files are added or modified in this diff; the existing suite's specs are untouched, consistent with the spec's Out of Scope clause.

## Verification of Developer Self-Review

- Developer's implementation-notes.md claims both prior CHANGES_REQUESTED findings were addressed (PRECONDITION/ENFORCED-BY comment at wait_for_app_healthy; explicit stderr warning replacing the silent cp swallow) — verified true against the current script, lines 65-72 and 100-105.
- Self-review checklist claims (no function >30 LOC, complexity <7, named constants, no layered-architecture concerns) verified true against the actual diff — no discrepancy found.

## Feedback for the Developer

None
