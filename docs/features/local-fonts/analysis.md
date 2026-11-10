---
feature: "local-fonts"
bounded_context: "ye-olde-magic-shop frontend — typography and font-serving"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/src/styles.css", "apps/ye-olde-magic-shop/src/main.ts", "apps/ye-olde-magic-shop/package.json", "pnpm-lock.yaml", "apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts", "apps/ye-olde-magic-shop/tests/local-fonts.spec.ts", "docs/testing-cesium.md"]
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: local-fonts

## Summary

apps/ye-olde-magic-shop/src/styles.css imports Cinzel and Inter from fonts.googleapis.com on every route, so the app reaches an external origin on every page load and tests/reagent-map-globe.spec.ts has to name and skip two hosts in its no-request check. This chore self-hosts both fonts via @fontsource/cinzel and @fontsource/inter pinned to exact versions, removes the cross-origin @import and the matching test exception, adds a new document.fonts.load-based test proving the fonts actually render, and removes the now-resolved Google Fonts gap from docs/testing-cesium.md. A verified salvage commit (3629f44, branch salvage/local-fonts) already contains this work; the job is to re-apply it, hand-port the docs/testing-cesium.md hunk (that file changed on main since the salvage commit), and verify every acceptance criterion rather than rewrite the solution.

### Acceptance Criteria

- Cinzel (400, 700) and Inter (400, 600) MUST load from the app's own origin, through @fontsource/cinzel and @fontsource/inter pinned to exact versions, and src/styles.css MUST have no @import from another origin.
  - Example: package.json dependencies carry "@fontsource/cinzel": "5.3.0" and "@fontsource/inter": "5.3.0" (no ^ or ~)
  - Example: src/styles.css's only @import is the local "tailwindcss" import
- tests/reagent-map-globe.spec.ts MUST drop PRE_EXISTING_OFF_ORIGIN_HOSTS and the code that skips those hosts, and its no-request test MUST pass with no exception.
- The no-request test MUST fail when the old cross-origin @import is restored. Record the failing output in the delivery record.
  - Example: Temporarily re-add the Google Fonts @import line to styles.css, re-run the no-request test, capture the failure, then revert
- A page opened in the runner image MUST still render headings in Cinzel and body text in Inter, checked in tests/local-fonts.spec.ts with document.fonts.load at the exact weight, never document.fonts.check.
  - Example: document.fonts.check returned true even with the font imports removed in the salvage run — this is why the AC forbids it
  - Example: Record the test failing with all four @fontsource imports removed
  - Example: Record the test failing with only the Inter 600 import removed
- docs/testing-cesium.md's Google Fonts entry under "Known gaps" MUST be gone, and the file MUST still follow .claude/rules/writing-style.md.
- scripts/ui-suite-in-chrome.sh and scripts/ui-suite-in-chrome.sh --ubi9 MUST both pass.
- Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change (the run's own delivery record under docs/features/local-fonts/ excepted).

### Non-Functional Requirements

- **performance**: No page in the app requests a font, or any other resource formerly covered by the removed test exception, from another origin.

## Proposed Fitness Functions

- **No page makes a request off the app's own origin (fonts included — no host exception remains)**
  - **Verification**: Playwright: tests/reagent-map-globe.spec.ts's existing no-request test, run via scripts/ui-suite-in-chrome.sh in both the debian and --ubi9 images
  - **Owner**: devops-engineer
- **Headings render in self-hosted Cinzel and body text in self-hosted Inter at the exact pinned weights**
  - **Verification**: Playwright: tests/local-fonts.spec.ts using document.fonts.load per family/weight pair, comparing the matched FontFace.weight to the requested weight
  - **Owner**: qa-engineer
- **@fontsource/cinzel and @fontsource/inter stay pinned to exact versions (no ^ or ~ range)**
  - **Verification**: code-reviewer inspects apps/ye-olde-magic-shop/package.json and the pnpm-lock.yaml diff for exact version strings
  - **Owner**: code-reviewer

## Out of Scope

- Changing which fonts the shop uses, or any other styling decision.
- The map tiles the Leaflet locations view loads from OpenStreetMap — tracked separately as order leaflet-local-tiles.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: ye-olde-magic-shop frontend — typography and font-serving
- **Context Crossings**: reagent-map-globe's E2E no-request test (tests/reagent-map-globe.spec.ts) — this chore removes its PRE_EXISTING_OFF_ORIGIN_HOSTS exception

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/src/styles.css` — Remove the cross-origin @import of Google Fonts (Cinzel/Inter); only the local tailwindcss import remains
- `apps/ye-olde-magic-shop/src/main.ts` — Single entry point loaded by every route; add the four self-hosted @fontsource CSS imports (cinzel 400/700, inter 400/600) here
- `apps/ye-olde-magic-shop/package.json` — Add @fontsource/cinzel and @fontsource/inter as exact-pinned (no ^/~) dependencies
- `pnpm-lock.yaml` — Regenerate via `pnpm install` after the package.json change so the lockfile reflects the two new pinned dependencies; do not hand-edit
- `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` — Remove PRE_EXISTING_OFF_ORIGIN_HOSTS and the skip branch in collectOffOriginRequests so the no-request test runs with no exception
- `apps/ye-olde-magic-shop/tests/local-fonts.spec.ts` — New test file: asserts Cinzel/Inter load at the exact weight via document.fonts.load, never document.fonts.check
- `docs/testing-cesium.md` — Remove the 'Still open' Google Fonts paragraph after the 'To run with no network' table, and the matching 'Known gaps' bullet; main added two unrelated entries to this file since the salvage commit (a Viewer-ellipsoid table row and a wheel-at-centre gap) that must stay untouched

### Data Model Changes

None

### API Changes

None

### New Dependencies

- @fontsource/cinzel@5.3.0
- @fontsource/inter@5.3.0

## Task List

None

### Developer Tasks

1. Reuse the verified salvage work instead of rewriting it: `git diff --binary 83beb6e 3629f44 -- apps/ye-olde-magic-shop/package.json apps/ye-olde-magic-shop/src/main.ts apps/ye-olde-magic-shop/src/styles.css apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts apps/ye-olde-magic-shop/tests/local-fonts.spec.ts | git apply --index`
2. Hand-port the docs/testing-cesium.md change: remove the two-line 'Still open' paragraph right after the 'To run with no network' table, and remove only the 'Google Fonts' bullet under 'Known gaps'; leave the Viewer-ellipsoid table row and the wheel-at-centre-of-view gap (both added by main after the salvage commit) exactly as they are
3. Run `pnpm install` from the repo root; confirm the pnpm-lock.yaml diff adds only @fontsource/cinzel and @fontsource/inter entries at version 5.3.0
4. Confirm apps/ye-olde-magic-shop/src/styles.css has no @import pointing at any origin other than the app's own build output
5. Confirm apps/ye-olde-magic-shop/src/main.ts imports @fontsource/cinzel/400.css, @fontsource/cinzel/700.css, @fontsource/inter/400.css and @fontsource/inter/600.css before the app mounts

### QA Tasks

1. Run `scripts/ui-suite-in-chrome.sh` (debian image); confirm tests/reagent-map-globe.spec.ts's no-request test and tests/local-fonts.spec.ts both pass
2. Run `scripts/ui-suite-in-chrome.sh --ubi9`; confirm the same two tests pass
3. Temporarily restore the removed Google Fonts @import in styles.css, run the no-request test alone, confirm it fails, record the failing output in the delivery record, then revert
4. Temporarily remove all four @fontsource imports from main.ts, run tests/local-fonts.spec.ts, confirm it fails, record the failing output, then revert
5. Temporarily remove only the @fontsource/inter/600.css import from main.ts, run tests/local-fonts.spec.ts, confirm it fails because document.fonts.load does not resolve a 600-weight face, record the failing output, then revert
6. Confirm cyclomatic complexity stays below 7 for every changed function and that no file outside scope.paths changed

### Tech Writer Tasks

1. Confirm docs/testing-cesium.md still follows .claude/rules/writing-style.md (active voice, sentence length, RFC 2119 modals where applicable) after the Google Fonts passages are removed

### DevOps Tasks

1. Confirm scripts/ui-suite-in-chrome.sh and scripts/ui-suite-in-chrome.sh --ubi9 both exit 0; no CI pipeline changes are needed since both scripts already run this suite

## Edge Cases and Risks

- The old cross-origin @import is put back in styles.css: The no-request test must fail (it no longer has a host exception to fall back on); capture and record the failing output per AC3
- All four self-hosted font imports are removed from main.ts: tests/local-fonts.spec.ts must fail — document.fonts.load resolves to no matching face for each family/weight pair; record the failing output per AC4
- Only the @fontsource/inter/600.css import is removed: tests/local-fonts.spec.ts must still fail: document.fonts.load for Inter 600 either matches nothing or matches a neighbouring weight, and the test compares FontFace.weight explicitly rather than using document.fonts.check (which returned true in this exact case during the salvage run); record the failing output per AC4
- pnpm install touches more of pnpm-lock.yaml than the two new dependencies: Review the lockfile diff before committing; it must only add @fontsource/cinzel and @fontsource/inter entries, not shift unrelated resolved versions
- docs/testing-cesium.md diverged from the salvage commit's version: Apply the rest of the salvage diff with `git apply --index --exclude=docs/testing-cesium.md`, then hand-remove only the Google Fonts passages from the current file, preserving the two unrelated entries main added since 83beb6e

## Definition of Done

- All seven acceptance criteria verified, with failing-test evidence for the three negative cases recorded in the delivery record under docs/features/local-fonts/
- src/styles.css has no @import from another origin; @fontsource/cinzel and @fontsource/inter are exact-pinned dependencies
- tests/reagent-map-globe.spec.ts's no-request test runs with no host exception and passes
- tests/local-fonts.spec.ts exists, uses document.fonts.load (never document.fonts.check), and passes
- docs/testing-cesium.md's Google Fonts passages are removed and the file still follows .claude/rules/writing-style.md
- scripts/ui-suite-in-chrome.sh and scripts/ui-suite-in-chrome.sh --ubi9 both pass
- Every new or changed function has cyclomatic complexity below 7
- No file outside scope.paths changed (the delivery record under docs/features/local-fonts/ excepted)
