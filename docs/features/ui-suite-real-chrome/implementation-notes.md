---
feature: "ui-suite-real-chrome"
bounded_context: ""
domain_terms: []
files_touched: ["scripts/ui-suite-in-chrome.sh"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Implementation Notes: ui-suite-real-chrome

## Files Created

None

## Files Modified

- scripts/ui-suite-in-chrome.sh

## Interface Design

None

## Named Refactoring Log

None

## Self-Review Checklist

- [x] Every public method has an intention-revealing name
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new/changed functions
- [x] No primitive obsession — Not applicable — shell script, no new domain types introduced.
- [x] No feature envy
- [x] No magic numbers or strings — Timeout bound remains the existing named constant APP_HEALTH_TIMEOUT_SECONDS.
- [x] Dependency direction verified — Not applicable — single script, no layered architecture crossed.

## Simple Design Verification

- [x] Passes all tests — bash -n syntax check passed; no test suite exercises this script directly beyond the shell-level tests already planned by qa-engineer, which this change does not affect.
- [x] Reveals intention — PRECONDITION/ENFORCED-BY comment makes the judgment-only exception legible at the loop site, matching architecture-guardrails.md #7.
- [x] No duplication
- [x] Fewest elements — Both changes are minimal — a comment block and a one-line error message, no new functions or abstractions.

## Key Decisions

- Documented wait_for_app_healthy's until-loop-with-sleep as a judgment-only exception to architecture-guardrails.md #5 via a PRECONDITION/ENFORCED-BY comment at the loop site.: architecture-notes.md's Stability Design section already classified this as judgment-only (a bounded one-shot poll of a process the script itself started, matching containers/saturday-runner/entrypoint.sh's wait_for_display shape, not resilience against an unreliable remote dependency) and its Trade-offs note directed restating that as a one-line comment in the script; the comment was missing, which code-reviewer flagged as blocking.
- Replaced copy_report's silent `2>/dev/null || true` with an explicit stderr warning on cp failure, keeping the script's exit code unchanged.: A silently swallowed copy failure (permissions, disk full, missing report dir) was indistinguishable from success; AC2 still requires the script to exit with the Playwright suite's own exit code, so the warning only adds visibility, it doesn't change control flow.

## Deviations from Analysis

None

## Dependencies Added

None

## Notes for QA

- Both fixes are comment/stderr-only changes to scripts/ui-suite-in-chrome.sh — no behavioral change to exit codes or control flow, so prior test expectations for AC2/AC4 still hold.
- If there's a shell-level test asserting copy_report's failure path, it can now assert the new stderr warning text ("ui-suite-in-chrome: warning: failed to copy playwright report") instead of only checking exit code passthrough.

## Notes for DevOps

None
