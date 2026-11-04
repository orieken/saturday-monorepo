# Documentation Report: reagent-map-globe

```yaml
feature: "reagent-map-globe"
bounded_context: "Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)"
domain_terms: ["Reagent Site", "Aurelis", "Adventurer"]
files_touched: []
issue_refs: ["reagent-map-globe"]
linked_adrs: []
linked_kis: []
```

## Status: PASS

## Summary

`apps/ye-olde-magic-shop` gains a second map view. At `/reagent-map`, behind the new build-time
flag `VITE_FEATURE_REAGENT_MAP` (`.env` sets it `false`), the app shows Aurelis — an invented,
6,000 km-radius spherical planet — as an offline 3D Cesium globe. The eight seeded reagent sites
from the existing `GET /api/reagent-sites` endpoint are plotted on the globe and listed by name
beside it. With the flag off or unset, the route and its nav link do not exist and the app is
unchanged. This slice adds no backend change and no expeditions, fly-to, or line-tool behavior —
those are later slices in the same spec package. The permanent record for this run persists to
`docs/features/reagent-map-globe/` once the pipeline completes.

## README / CHANGELOG

`apps/ye-olde-magic-shop/README.md` is explicitly out of this order's `scope.paths`
(`analysis.md`'s Out of Scope), and the repository has no root `CHANGELOG.md`. Neither was
touched. If a changelog is adopted later, this feature's entry is: "Added a feature-flagged 3D
globe view of Aurelis showing every reagent site, off by default behind
`VITE_FEATURE_REAGENT_MAP`."

## Feature Documentation

### Enabling the view

| Setting | Value | Effect |
|---|---|---|
| `VITE_FEATURE_REAGENT_MAP` unset, empty, or anything other than the literal string `true` | default (`.env` sets `false`) | No `/reagent-map` route, no nav link. The app behaves exactly as it does without this feature. |
| `VITE_FEATURE_REAGENT_MAP=true` | opt-in | The `/reagent-map` route and its "Reagent Map" nav link both exist. |

`isReagentMapEnabled(value)` in `src/featureFlags.ts` is the single pure function both the router
and the nav link read. It returns `true` only for the exact string `"true"`.
`scripts/ui-suite-in-chrome.sh` sets `VITE_FEATURE_REAGENT_MAP=true` for the dev server it starts,
so the pre-ship Chrome check exercises the view while every other run keeps the flag at `.env`'s
`false`.

### What the view shows

- A 3D Cesium globe sized to Aurelis: a sphere of radius 6,000 km, rendered with one equirectangular
  texture (`public/aurelis.jpg`) and no imagery, geocoder, or base-layer picker.
- One point entity per reagent site, plotted at its seeded latitude/longitude, including Frostcap
  Hollow near the north pole and Tidewrack Shoals on the 180th meridian.
- A list of all site names beside the globe (`data-testid="reagent-site-list"`), so a screen-reader
  or keyboard user has the full data set without needing to operate the globe.

### Running fully offline

Cesium normally calls Cesium ion for imagery, terrain, and geocoding. This view disables all of
that: no ion token, `baseLayerPicker` and `geocoder` are both `false`, and the only imagery is the
app's own `public/aurelis.jpg`. Cesium's own static assets (widgets CSS, workers, `Cesium.js`) are
served same-origin under `/cesium/` by a small Vite plugin in `vite.config.ts` that copies
`node_modules/cesium/Build/Cesium` at dev-server start and at build time — nothing in that tree is
committed to source control. No request made while this view is open leaves the app's own origin
or the mock API's own origin.

### Test automation seam

The view sets `window.__SATURDAY_CESIUM__ = { viewer, Cesium }` on mount and clears it on unmount,
and tags every site entity with `saturdayMetadata: { id, name, kind: 'site' }`, per
`packages/saturday-playwright-cesium/README.md`'s published seam contract. This lets that package's
Playwright fixture drive and query the view without `apps/ye-olde-magic-shop` taking a runtime
dependency on it — the app only reuses the package's already-reviewed `aurelis.jpg` texture and
Cesium `Viewer` option set as reference material.

### Dependency added

- `cesium` `1.146.0` (exact version, matching `packages/saturday-playwright-cesium`'s already-resolved
  version — avoids a second resolved Cesium version in `pnpm-lock.yaml`).

## Architecture Decision Records

No ADR was written. Every structural decision for this slice (view placement, the exact Cesium
version pin, the inline Vite static-asset plugin, the single-predicate feature flag, the
no-runtime-dependency seam with `saturday-playwright-cesium`) is additive within an established
pattern or reversible at low cost, and each already carries its own fitness function in
`architecture-notes.md` — none constrains how a future feature must be built, so none rises to an
ADR-worthy decision.

## Code-Level Docs

- `src/featureFlags.ts`'s `isReagentMapEnabled` carries a one-line comment naming what it gates and
  who reads it (router, `App.vue`) — no further doc needed; the name and signature are
  self-explanatory.
- `src/views/ReagentMapView.vue`'s `loadCesiumGlobal` carries a comment explaining why Cesium is
  loaded as a classic global script rather than bundled (an upstream Rollup/`@cesium/engine`
  packaging mismatch) — this is non-obvious and would otherwise look like an unnecessary pattern
  deviation from the rest of the app's ES-module imports.
- `tests/reagent-map-feature-flag.spec.ts`, `tests/reagent-map-globe.spec.ts`,
  `tests/reagent-map-resilience.spec.ts`, and `tests/reagent-map-a11y-check.spec.ts` carry `@issue`
  / `@ac` JSDoc annotations per `shared/rules/testing-conventions.md`'s Test Annotation Convention.

## Known Gaps / Follow-ups

- `security-report.md` recorded one Low-severity finding: `vite.config.ts`'s dev-only Cesium
  static-asset middleware can throw on malformed percent-encoding in a request path (e.g.
  `GET /cesium/%`). Not fixed — tracked as tech debt, dev-server-only, no production exposure.
- `qa-report.md` notes the view's own code makes zero off-origin requests, but the page as a whole
  still loads two pre-existing, unrelated Google Fonts origins via a sitewide `@import` in
  `src/styles.css` that predates this feature and is outside this order's `scope.paths`. The NFR
  test carves out exactly those two named hosts with a comment citing this evidence, so it still
  fails on any new off-origin dependency the reagent map itself might introduce.
- `architecture-notes.md` flags `ReagentMapView.vue`'s sibling, `LocationsView.vue`, as a
  Boy-Scout-Rule opportunity (one long inline `onMounted`) for whoever next touches that file — it
  is unchanged by this order and out of `scope.paths`.
- `code-review-report.md` left one non-blocking DRY finding: `src/api/reagentSites.ts` duplicates
  `src/api/items.ts`'s API-base-resolution logic verbatim. Deferred per `architecture-notes.md`;
  extract a shared `resolveApiBase()` helper the next time a third API client module needs it.
- Cesium's camera accepts keyboard focus for panning but has no verified, documented
  keyboard-operable equivalent beyond its own default bindings (`accessibility-engineer.md`).
  Flagged for human/product judgment, not fixed here — the adjacent site-name list already gives
  keyboard and screen-reader users the full data set without needing to drive the globe.

## Notes for DevOps

- `architecture-notes.md` names two fitness functions still needing a CI check: (1) exactly one
  resolved `cesium` version across the workspace after `pnpm install`, and (2) `dist/cesium`
  contains the Cesium `Build` tree after `npm run build` (e.g. `ls dist/cesium/Widgets`).
- `playwright.config.ts`'s `chromium` project now carries a `testIgnore` for the
  `reagent-map-*.spec.ts` glob; confirm CI runs both `npm run test:e2e` (chromium, unaffected) and
  `npm run test:chrome` / `scripts/ui-suite-in-chrome.sh` (chrome project, picks up the new specs)
  and checks both exit codes.