```markdown
# Documentation Report: reagent-map-expeditions

```yaml
feature: "reagent-map-expeditions"
bounded_context: "Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)"
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: ["reagent-map-expeditions", "reagent-map-globe"]
linked_adrs: []
linked_kis: []
```

## Status: PASS

## Summary

The reagent map, still behind the existing `VITE_FEATURE_REAGENT_MAP` flag (off by default), now
draws every expedition `GET /api/expeditions` returns as a geodesic line on the Aurelis globe, and
each site in the adjacent list becomes a button that flies the camera to that site — including
Tidewrack Shoals, which starts on the far side of the globe in the default view. Signed out, an
adventurer sees the shop's three seeded expeditions; signed in, they also see their own, never
another adventurer's. This slice adds no backend change — `GET /api/expeditions` was already
shipped by `expeditions-api` — and no line-drawing or saving tool; those are the next two orders in
the same spec package. The permanent record for this run persists to
`docs/features/reagent-map-expeditions/` once the pipeline completes.

## README / CHANGELOG

`apps/ye-olde-magic-shop/README.md` is unchanged — the feature stays behind the same build-time
flag `reagent-map-globe` already documented there (or, if that prior slice did not add an entry
either, the flag's behavior is unchanged by this order either way), and the repository has no root
`CHANGELOG.md`. If a changelog is adopted later, this feature's entry is: "The reagent map now
draws every visible expedition as a line on the globe, and clicking a site in the list flies the
camera to it."

## Feature Documentation

### What's new on the map

- Every expedition the signed-in (or signed-out) adventurer can see is drawn as a geodesic line
  between its two reagent sites, tagged `saturdayMetadata: { id, name, kind: 'expedition' }` — the
  same per-entity tagging convention `reagent-map-globe` established for sites (`kind: 'site'`).
- Signed out, the three seeded expeditions ("Spring to Crag," "Glade to Caves," "Cairn to Shoals")
  are drawn. Signed in, the adventurer's own expedition is drawn alongside the seeded ones; another
  adventurer's own expedition never is. This visibility rule is enforced server-side by
  `GET /api/expeditions` (`expeditions-api`), not by this view — the view renders whatever the
  endpoint returns.
- Each line both ends within 1 km of its stored coordinates and follows the sphere's surface
  (`arcType: Cesium.ArcType.GEODESIC`), not a straight chord through the globe.

### Flying to a site

- Every entry in the site list (`data-testid="reagent-site-list"`) is now a `<button type="button"
  data-testid="fly-to-site-{id}" aria-label="Fly to {name}">`, not plain text. Clicking, or
  pressing Enter/Space while it has keyboard focus, flies the camera so that site lands on screen
  on the near side of the globe — including a site that starts on the far side, such as Tidewrack
  Shoals from the default view.
- The button carries no custom keyboard handling; it relies entirely on the browser's native
  button semantics, so Tab order and activation keys are unchanged from any other button on the
  page.

### Accessible name for the globe

- The globe's `role="img"` accessible name now also describes the newly drawn expeditions ("…and
  every known expedition drawn as a line between its two sites; the full list of sites is in the
  adjacent list"), closing a WCAG 1.1.1 gap this order would otherwise have introduced: expeditions
  render only as Cesium scene geometry with no other text equivalent anywhere in the view.

### Visual baseline change

- `globe-default-view.png` changed in both runner images (`debian`, `ubi9`): a new line now appears
  between Moonwell Spring and Ember Crag — the "Spring to Crag" seeded expedition, the only one of
  the three whose both endpoints fall inside the default view's near-side hemisphere. The other two
  seeded expeditions each have at least one endpoint off-screen in this particular view, so they do
  not appear in this baseline even though they are drawn. No other baseline
  (`reagent-map-header.png`, `reagent-site-list.png`, `reagent-map-error.png`) changed, and
  `GLOBE_MAX_DIFF_PIXEL_RATIO` (1%) was not raised.

## Architecture Decision Records

No ADR was written. Every structural decision for this slice — a new sibling API client module
(`expeditions.ts`), a second entity-drawing helper mirroring the existing site one, promoting the
Cesium module and ellipsoid to component-scope variables so `flyToSite` can reach them, and plain
`<button>` controls with no composable or store — is additive within the pattern `reagent-map-globe`
already established, and each already carries its own fitness function or judgment-only
justification in `architecture-notes.md`. None constrains how a future feature must be built.

## Code-Level Docs

- `apps/ye-olde-magic-shop/src/api/expeditions.ts` exports `GeoPoint`, `Expedition`, and
  `fetchExpeditions()`, matching `reagentSites.ts`'s existing client shape — no further doc needed
  beyond the self-explanatory names and types.
- `tests/reagent-map-expeditions.spec.ts` and `tests/reagent-map-fly-to-site.spec.ts` carry `@issue`
  JSDoc annotations naming the covered acceptance criteria and feature file, per
  `shared/rules/testing-conventions.md`'s Test Annotation Convention.
- The selector-path repairs in `tests/reagent-map-globe.spec.ts` and
  `tests/reagent-map-pan-zoom.spec.ts` (filtering entity queries to `kind === 'site'` so the new
  expedition entities don't change their counts) are documented in `implementation-notes.md`'s
  Deviations from Analysis section, with the test-repair-contract.md evidence for why they are path
  repairs, not assertion weakening.

## Known Gaps / Follow-ups

- `security-report.md` recorded two Low/Info findings, neither blocking: `fetchExpeditions()` trusts
  the API response shape via a type assertion rather than runtime (Zod) validation — an exact mirror
  of `fetchReagentSites()`'s pre-existing pattern; and `Expedition.owner` travels to the client
  unused by the UI today, flagged for whoever first renders it to confirm it's treated as a
  per-user identifier, not a safe-to-display label.
- `code-review-report.md` / `implementation-notes.md` note the runtime `API_BASE` resolution block
  is now duplicated a fourth time across `reagentSites.ts`, `orders.ts`, `items.ts`, and
  `expeditions.ts`. Deferred by design — this order is additive, not a refactor.
- `qa-report.md` diagnosed two intermittent failures during the pre-ship full-suite run under
  3-worker parallelism (`tests/map.spec.ts`'s unrelated layer-visibility test, and
  `reagent-map-pan-zoom.spec.ts`'s drag/zoom tests) as resource contention from parallel
  software-rendered WebGL pages, not a regression from this order — confirmed by a clean 12/12
  isolated re-run. Worth flagging to whoever owns runner-image worker tuning.
- `apps/ye-olde-magic-shop` still has no statement/line coverage instrumentation; the claim's "every
  changed line is covered by a test that can fail" criterion was verified by reverting
  `ReagentMapView.vue` and confirming the new assertions fail against the pre-change code, per
  `qa-report.md`.

## Notes for DevOps

None.
```