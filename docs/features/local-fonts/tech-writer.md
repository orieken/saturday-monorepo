I have what I need. Here is the artifact.

```markdown
# Documentation Report: local-fonts

```yaml
feature: "local-fonts"
bounded_context: "ye-olde-magic-shop frontend — typography and font-serving"
domain_terms: []
files_touched: ["docs/testing-cesium.md"]
issue_refs: ["local-fonts"]
linked_adrs: []
linked_kis: []
```

## Status: PASS

## Summary

`apps/ye-olde-magic-shop` served Cinzel and Inter from `fonts.googleapis.com` on every route, so
every page reached an external origin and `tests/reagent-map-globe.spec.ts` carried a two-host
exception in its no-request check. This chore self-hosts both fonts through `@fontsource/cinzel`
and `@fontsource/inter`, pinned to exact version `5.3.0`, imported as side-effect CSS in
`src/main.ts` next to the existing `leaflet/dist/leaflet.css` import. `src/styles.css` now imports
only `tailwindcss`. The no-request test's host exception is gone, so it now fails on any off-origin
request with no fallback list. A new test, `tests/local-fonts.spec.ts`, proves the fonts render at
their pinned weight using `document.fonts.load`, not `document.fonts.check` — the latter returned a
false positive with the fonts removed during an earlier salvage run. `docs/testing-cesium.md`'s
Google Fonts entries are removed, since the gap they described no longer exists.

## README / CHANGELOG

`apps/ye-olde-magic-shop/README.md` does not mention fonts and needs no change. The repository has
no root `CHANGELOG.md` and the app has no app-level one; none was added, consistent with
`docs/features/reagent-map-globe/tech-writer.md`'s same finding for this repository. If a changelog
is adopted later, this feature's entry is: "Self-hosted Cinzel and Inter through `@fontsource`
instead of loading them from Google Fonts on every page."

## Feature Documentation

### Where the fonts come from now

| File | Change |
|---|---|
| `apps/ye-olde-magic-shop/src/main.ts` | Imports `@fontsource/cinzel/400.css`, `@fontsource/cinzel/700.css`, `@fontsource/inter/400.css`, and `@fontsource/inter/600.css` as side-effect CSS, before the app mounts |
| `apps/ye-olde-magic-shop/src/styles.css` | Its only `@import` is `tailwindcss`; the `fonts.googleapis.com` import is gone |
| `apps/ye-olde-magic-shop/package.json` | Depends on `@fontsource/cinzel` and `@fontsource/inter`, both pinned to the exact version `5.3.0` (no `^` or `~`), matching the existing exact pin on `cesium` in the same file |

No page requests a font, or anything else, from another origin. `--theme`'s `--font-magical:
'Cinzel', serif` and `--font-body: 'Inter', sans-serif` in `styles.css` are unchanged — only where
the two families are served from changed.

### Proving the fonts render

`tests/local-fonts.spec.ts` opens the app and calls `document.fonts.load('<weight> 1em
<family>')` for each of the four family/weight pairs (Cinzel 400, Cinzel 700, Inter 400, Inter
600), then checks the returned `FontFace.weight` matches the requested weight. The test does not
use `document.fonts.check`: that method returned `true` even with every font import removed during
an earlier run of this same work, so it could not have caught a regression.

### Proving the no-request test still catches a regression

`tests/reagent-map-globe.spec.ts`'s no-request test no longer carries
`PRE_EXISTING_OFF_ORIGIN_HOSTS` or the skip branch in `collectOffOriginRequests` that referenced it.
With the old Google Fonts `@import` temporarily restored in `src/styles.css`, the test failed with
exactly one off-origin URL —
`https://fonts.googleapis.com/css2?family=Cinzel:wght@400;700&family=Inter:wght@400;600&display=swap`
— confirming the test has no fallback left to rely on.

### Dependencies added

- `@fontsource/cinzel` `5.3.0` (exact version) — self-hosts Cinzel 400 and 700
- `@fontsource/inter` `5.3.0` (exact version) — self-hosts Inter 400 and 600

## Architecture Decision Records

No ADR was written. `architecture-notes.md` records three structural decisions for this change —
importing the font CSS directly in `main.ts` rather than adding a fonts-loading abstraction, exact-
pinning the two new dependencies to match the existing `cesium` precedent, and not introducing any
loading utility — and each is additive within an already-established pattern (the same side-effect
CSS import `leaflet.css` already uses) or reversible at low cost. None constrains how a future
feature must be built, so none rises to an ADR-worthy decision.

## Code-Level Docs

- `tests/local-fonts.spec.ts` carries a module comment stating the acceptance criterion it verifies
  and explaining, non-obviously, why it uses `document.fonts.load` instead of
  `document.fonts.check` — the comment names the specific prior false positive, which is exactly
  the kind of hidden constraint a reader would otherwise have no way to know about.
- `src/main.ts` and `src/styles.css` need no new comments: the four `@fontsource` imports sit at
  the same side-effect-import site as the pre-existing `leaflet/dist/leaflet.css` import, and
  `styles.css`'s single remaining `@import "tailwindcss"` is self-explanatory.

## `docs/testing-cesium.md` update

Removed the "Still open" paragraph (after the "To run with no network" table) and the matching
"Known gaps" bullet that described the sitewide Google Fonts `@import` as an unresolved exception
to the no-request check — that gap is closed by this feature. The two entries `main` added to this
file after the salvage commit this work was recovered from (the Viewer-ellipsoid table row, and the
wheel-at-centre-of-view "Known gaps" entry) are untouched. The file still follows
`.claude/rules/writing-style.md`: active voice, short sentences, and no new noun clusters were
introduced by the removal.

## Known Gaps / Follow-ups

- `qa-report.md` flags two earlier `scripts/ui-suite-in-chrome.sh --ubi9` runs that failed inside
  `tests/reagent-map-pan-zoom.spec.ts` on timing/precision, not on anything font-related. That file
  is outside this feature's `scope.paths` (added by the earlier `reagent-map-globe` feature) and
  must not be touched here. A clean standalone `--ubi9` run and two clean `debian` runs satisfy this
  feature's own acceptance criterion; the flake is left for whoever owns that spec next.
- `apps/ye-olde-magic-shop` has no instrumented coverage tool (pre-existing, already documented in
  `docs/testing-cesium.md`'s "Known gaps"). `qa-report.md`'s 100% coverage figure for this change is
  a manual, mutation-verified claim — every changed or added line was shown to fail a test when
  reverted — not a tool-measured statement percentage.

## Notes for DevOps

- No CI pipeline change is needed. `scripts/ui-suite-in-chrome.sh` and its `--ubi9` variant already
  run the full suite, including the no-request test, which is this feature's fitness function for
  "no page requests a font from another origin" — confirm both continue to exit `0`, now with no
  host exception backing that check.
```