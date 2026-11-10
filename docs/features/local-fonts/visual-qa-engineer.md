Now I have enough context. This feature is a pure font-source swap with no new UI surface — scope.paths explicitly excludes styling changes, and the only new test checks `document.fonts.load`, not a user interaction. I also confirmed `heatmap-data/` is a pre-existing, out-of-scope artifact from `tests/heatmap.spec.ts`, not something this feature's QA run generated, and `@orieken/saturday-ml-analyzer` isn't installed or built in this app. Producing the report now.

```markdown
---
feature: "local-fonts"
bounded_context: "ye-olde-magic-shop frontend — typography and font-serving"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/src/main.ts", "apps/ye-olde-magic-shop/src/styles.css", "apps/ye-olde-magic-shop/package.json", "apps/ye-olde-magic-shop/tests/local-fonts.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts"]
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
---

# Visual QA Report: local-fonts

## Summary

**UNCONFIGURED — NO VISUAL SURFACE FOR THIS CHANGE**

## Invocation Conditions Checked

- `heatmap-data/` **does exist** at `apps/ye-olde-magic-shop/heatmap-data/`, containing one file:
  `heatmap-spec-ts__heatmap-verification.json` (1.35 MB, last written 2026-11-10 01:29).
- Playwright snapshot baselines: **none found**. No `.playwright-baselines/` directory, no
  `*-snapshots/` directory, and no `*.png` under `test-results/` or `tests/` for this app (the only
  `.png` in the app tree is `public/logo.png`, a static asset, not a test baseline).

## Why This Report Stops Here

The `heatmap-data/` directory is **not output from this feature's QA run**. Three independent
checks confirm it predates and is unrelated to `local-fonts`:

1. **Provenance**: the single JSON file's name (`heatmap-spec-ts__heatmap-verification.json`) maps
   to `apps/ye-olde-magic-shop/tests/heatmap.spec.ts`, a pre-existing test unrelated to font
   loading. It is not produced by `tests/local-fonts.spec.ts`, the one test this feature added.
2. **Scope**: the work order (`.claude/factory/claims/local-fonts.md`) lists an explicit
   `scope.paths` for this feature. Neither `heatmap-data/**` nor `tests/heatmap.spec.ts` is in it,
   and the spec's own "Out of Scope" section rules out styling changes beyond the font swap.
3. **Shape of the change**: `tests/local-fonts.spec.ts` calls `document.fonts.load()` and compares
   returned `FontFace.weight` values — it performs no clicks, no navigation beyond one `page.goto`,
   and drives no UI interaction a heatmap could cluster. There is nothing for an interaction-coverage
   tool to measure here even in principle.

Given that, running `@orieken/saturday-ml-analyzer` against the directory would score an unrelated
feature's interaction data, not this one's — doing so would misattribute a stale artifact as this
feature's evidence, exactly the failure mode `test-repair-contract.md`'s evidence standard exists to
block applied to a different artifact type. I did not run it.

Separately, as a secondary check: `@orieken/saturday-ml-analyzer` is not installed as a dependency
of `apps/ye-olde-magic-shop` (`package.json` has no entry for it) and has no build output
(`packages/saturday-ml-analyzer/dist/` does not exist) — so the tool is not runnable here even if
the data were in scope.

## Coverage Score

Not applicable — no heatmap data exists for this feature's own test.

## Hotspots / Cold Spots

Not applicable.

## Visual Regression (Screenshot Diffs)

Not applicable — no baselines exist, and this feature introduces no new rendered markup. Per
`accessibility-engineer.md`'s review (already on file for this feature), font-family, weight, and
fallback stack are unchanged; only the font files' network origin moved from
`fonts.googleapis.com` to the app's own `@fontsource` bundle. There is no pixel-level change this
feature is expected to produce, so a screenshot baseline would have nothing to regress against even
if one existed.

## Recommendations

- No action needed for this feature. Visual/heatmap tooling has nothing to analyze here by design —
  the spec's "Out of Scope" section explicitly excludes any change to which fonts the app uses or
  how they render.
- The stale `heatmap-data/heatmap-spec-ts__heatmap-verification.json` and `heatmap-report.html`
  belong to `tests/heatmap.spec.ts` from a prior feature. They are outside `local-fonts`'
  `scope.paths` and untouched by this run; flagging their cleanup, if wanted, to whoever owns that
  test next — not actioned here.
- If a future feature on this app wants screenshot regression coverage, note for that work: no
  `.playwright-baselines/` convention exists yet in this app, so baselines would need to be
  established fresh rather than compared against anything current.

---
*Reviewed by visual-qa-engineer. This feature has no visual surface to analyze — the pipeline's
visual-QA gate is satisfied by confirming that fact against the actual artifacts (heatmap-data
provenance, tool installation state, and the accessibility-engineer's already-filed visual-impact
finding), not by running analysis against unrelated or absent data.*
```