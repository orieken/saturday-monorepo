Now producing the final artifact.

```markdown
---
feature: "reagent-map-globe"
bounded_context: "Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)"
domain_terms: ["Aurelis", "Reagent Site"]
files_touched: ["apps/ye-olde-magic-shop/heatmap-data/heatmap-spec-ts__heatmap-verification.json", "apps/ye-olde-magic-shop/tests/heatmap.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-a11y-check.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-feature-flag.spec.ts"]
issue_refs: ["reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Visual QA Report: reagent-map-globe

## Summary: ⚠️ UNCONFIGURED (no violation, not blocking)

Both trigger conditions are technically present at the repo root (`heatmap-data/` exists;
Playwright produced a `.webm` under `test-results/`), but neither contains data that actually
describes this feature. There is no interaction-heatmap or visual-baseline coverage of the
reagent map / globe view to analyze. This is reported as UNCONFIGURED, not FAIL, per the
invocation contract — the project hasn't adopted either tool for this feature yet, which is
expected, not a defect.

## Heatmap Analysis

**Probe**: `apps/ye-olde-magic-shop/heatmap-data/` exists and contains exactly one capture file:
`heatmap-spec-ts__heatmap-verification.json`.

**Finding — this capture is not about the reagent map feature.** It was produced by
`apps/ye-olde-magic-shop/tests/heatmap.spec.ts`, a pre-existing, generic smoke test that predates
this feature and exists only to verify the `@orieken/saturday-playwright-heatmap` fixture wiring
itself. That test loads `/` (the Shop homepage), clicks the first `button` with text "Buy" and
the first `a` link — it never navigates to `/reagent-map` and never touches Cesium, the globe
canvas, or the reagent-site list. Its 65 captured interactables are all homepage elements (`Shop`,
`Locations`, `Cart`, `🔮 Login` nav links, item cards, `Add to Cart` buttons).

None of the four reagent-map test files added or modified for this feature
(`reagent-map-globe.spec.ts`, `reagent-map-resilience.spec.ts`, `reagent-map-a11y-check.spec.ts`,
`reagent-map-feature-flag.spec.ts`) import the heatmap fixture, so no interaction data for this
feature's own UI exists anywhere in the repo.

I built and ran `@orieken/saturday-ml-analyzer` (`packages/saturday-ml-analyzer`, not yet
published/built in this worktree — built locally via `tsc` to run the probe) against the one file
present, as a diagnostic only:

```
Test: heatmap verification (passed)
  Coverage Score: 1.5%
  Interactables: 65
  Hotspots (Clusters): 1
  Cold Spots (Untested): 64
  Top 3 Cold Spots:
    - a "Shop" ([data-testid="shop-link"])
    - a "Locations" ([data-testid="locations-link"])
    - a "Cart" ([data-testid="cart-link"])
```

This 1.5% score and these cold spots describe the Shop homepage's generic smoke test, **not** the
reagent map feature — reporting them as this feature's coverage would misattribute a pre-existing,
unrelated test's shallow click-path to a feature it never exercises. No score is reported for
reagent-map-globe itself because no heatmap data for it exists.

### Coverage Score
Not applicable — no heatmap capture exists for the reagent map view.

### Hotspots
Not applicable.

### Cold Spots
Not applicable to this feature. (See Finding above for why the one existing capture's cold spots —
Shop/Locations/Cart nav links — are out of scope.)

## Visual Regression (Screenshot Baselines)

**Probe**: searched `test-results/`, `tests/`, and for a `.playwright-baselines/` directory.

- No `.playwright-baselines/` directory exists anywhere in `apps/ye-olde-magic-shop`.
- `test-results/` contains one leftover artifact,
  `zz-tmp-flag-off-check-flag-.../video.webm` — a Playwright trace video from an unrelated
  flag-off run, not a screenshot baseline.
- No `*.png` snapshot files exist under `tests/` or `test-results/`.
- None of the four reagent-map spec files use `toHaveScreenshot()`, a `screenshot()` call, or an
  `@visual` tag — grepped directly, zero matches.

No visual regression tests ran because none exist for this feature. This is INFO, not a finding:
the project has not adopted Playwright visual snapshots for this view.

## Screenshot Diffs

None — no visual tests exist to produce a diff.

## Recommendations

- **If visual regression is wanted for this view going forward**, be aware the primary surface —
  the Cesium WebGL globe canvas (`#cesiumContainer`) — is a poor target for pixel-diff screenshot
  testing: camera framing, GPU/driver-dependent antialiasing, and asynchronous tile/texture load
  timing make byte-identical renders unreliable even on a true no-op change. A snapshot of the
  DOM-level chrome around it (the `reagent-site-list` panel, nav link state, loading/error states)
  would be a more stable target than the canvas itself.
- **If interaction-coverage analysis is wanted for this view**, the reagent-map spec files would
  need to import `@orieken/saturday-playwright-heatmap`'s `test` fixture the way
  `tests/heatmap.spec.ts` already does, and actually drive the feature's own interactables (nav
  link, site list items, globe) rather than the homepage. None of that wiring exists today — this
  is a gap to flag for a future QA pass, not something this pass should add silently (qa-engineer
  already closed its own loop in `qa-report.md`; adding fixture wiring here would be scope creep
  into QA's territory without a corresponding scenario review).
- `packages/saturday-ml-analyzer` has no committed `dist/` in this worktree — it was built locally
  (`tsc`) only to run the diagnostic above. If this analyzer is meant to run as part of CI or a
  packaged pipeline step, it needs a build step wired in (devops-engineer's concern, not this
  report's).

## Known Issues / Follow-ups

- No heatmap or visual-snapshot coverage exists for the reagent-map-globe feature itself. Not
  blocking — the feature's functional and accessibility coverage was already verified by
  `qa-engineer` (16/16 tests passing, 7/7 acceptance criteria) and `accessibility-engineer`. This
  report only confirms there is nothing further for visual QA to analyze yet.
```