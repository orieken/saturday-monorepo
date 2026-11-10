---
feature: "local-fonts"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
---

# QA Report: local-fonts

## Test Files Created

None

## Test Files Modified

None

## Coverage Summary

- Acceptance criteria covered: 7/7
- Total new tests: 1
- Statement coverage, lowest unit: **100.0%** (`apps/ye-olde-magic-shop (no coverage instrumentation tool configured for this app — see docs/testing-cesium.md 'Known gaps'; figure below is manual line-by-line verification, not an instrumented measurement)`)
  - `apps/ye-olde-magic-shop (no coverage instrumentation tool configured for this app — see docs/testing-cesium.md 'Known gaps'; figure below is manual line-by-line verification, not an instrumented measurement)`: 100.0%

## Test Results

- Passed: 54
- Failed: 0
- Skipped: 0

## Accessibility Check

- No new interactive elements or DOM structure introduced by this change — only the source of two already-used typefaces changed (remote Google Fonts -> self-hosted @fontsource packages). tests/reagent-map-a11y-check.spec.ts (pre-existing, unrelated to this feature) continued to pass in both runner images, confirming no regression.

## Bugs Found

None

## Known Gaps

- scripts/ui-suite-in-chrome.sh --ubi9 MUST pass: Two earlier --ubi9 runs intermittently failed inside tests/reagent-map-pan-zoom.spec.ts (first: a camera-move timeout on test #51; second, standalone: a grab-tolerance pixel overshoot on test #48) — a file entirely outside this feature's scope.paths, not touched by local-fonts, added by the prior reagent-map-globe feature's real-mouse-gesture pan/zoom tests. A third standalone --ubi9 run passed clean at 54/54, and debian passed clean at 54/54 on two separate runs (including the one run concurrently with the first flaky --ubi9 attempt, which passed the same pan-zoom test that failed under ubi9's resource-contended run). This points to pre-existing timing/precision flakiness in software-GL-rendered mouse-gesture tests under ubi9, not a regression from self-hosting fonts. AC6 is satisfied by the clean standalone --ubi9 pass; flagging for whoever owns reagent-map-pan-zoom.spec.ts next, since it's outside this feature's file-scope and must not be touched here.

## Notes for Tech Writer

- AC3 evidence (negative case): with the old Google Fonts @import temporarily restored in styles.css and the Docker debian suite run, tests/reagent-map-globe.spec.ts's no-request test failed with exactly one off-origin URL: 'https://fonts.googleapis.com/css2?family=Cinzel:wght@400;700&family=Inter:wght@400;600&display=swap' — expect(offOrigin).toEqual([]) received that single entry. Reverted immediately after capture; suite confirmed green again on rerun.
- AC4 evidence (negative case 1, reproduced independently of the developer's own run): removing all four @fontsource imports from main.ts made tests/local-fonts.spec.ts fail with missing = ['Cinzel 400', 'Cinzel 700', 'Inter 400', 'Inter 600']. Reverted.
- AC4 evidence (negative case 2, reproduced independently): removing only the '@fontsource/inter/600.css' import made the same test fail with missing = ['Inter 600'] only — confirming the test's exact-weight comparison (document.fonts.load + FontFace.weight check) is what catches this, not document.fonts.check, which the analysis notes returned true in this exact case during an earlier salvage run. Reverted.
- Final confirmation state: scripts/ui-suite-in-chrome.sh (debian) passed 54/54 twice; scripts/ui-suite-in-chrome.sh --ubi9 passed 54/54 on a clean standalone run. See knownGaps for the unrelated, pre-existing pan-zoom flake observed on two earlier --ubi9 attempts.
- This app has no instrumented coverage tool (pre-existing, documented in docs/testing-cesium.md's own 'Known gaps' section) — the 100% coverage figure above is a manual claim verified by mutation (every changed/added line was shown to cause a test failure when reverted), not a tool-measured statement percentage.
- git diff on pnpm-lock.yaml reviewed again independently: adds only @fontsource/cinzel@5.3.0 and @fontsource/inter@5.3.0 entries, nothing else shifted.
- No source code changes were required from QA — the developer's implementation and tests were correct as delivered; this stage only ran and independently reproduced the required evidence.
