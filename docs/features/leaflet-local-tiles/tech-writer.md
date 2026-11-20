# Documentation Report: leaflet-local-tiles

---
feature: "leaflet-local-tiles"
bounded_context: "ye-olde-magic-shop frontend — shop-locations map (Leaflet)"
domain_terms: []
files_touched: ["docs/features/README.md"]
issue_refs: ["leaflet-local-tiles"]
linked_adrs: []
linked_kis: []
---

## Summary

The locations map now loads its background image from the app's own origin instead of
`tile.openstreetmap.org`. `LocationsView.vue` calls `L.imageOverlay()` with a hand-authored,
in-repo SVG instead of `L.tileLayer()`. Leaflet keeps its default EPSG:3857 CRS, so every marker,
the Sword Coast territory polygon, and all existing `tests/map.spec.ts` assertions stay correct.
A new Playwright test proves the view loads with every off-origin request blocked. This report
updates the delivered-features index and records one known, accepted test-coverage gap for the
human record.

## Documentation Checklist

### Always Update

- [x] **CHANGELOG.md** — Not applicable. No `CHANGELOG.md` exists at the repository root or under
  `apps/ye-olde-magic-shop/`, and no changelog convention exists anywhere in this repo today
  (verified by search). Introducing one is a new project-wide convention this feature did not ask
  for, so none was created.
- [x] **README.md** — `apps/ye-olde-magic-shop/README.md` makes no mention of OpenStreetMap or
  map tiles; it only covers install and test-run commands, so it needs no change.
  `docs/features/README.md` (the delivered-features index) is updated: added a row for
  `leaflet-local-tiles` pointing at `./leaflet-local-tiles/`, matching the existing table format.

### Update if Applicable

- [ ] **API documentation** — N/A. No API surface changed.
- [ ] **Configuration docs** — N/A. No new env vars or config options.
- [x] **Architecture Decision Record** — No new ADR written. `architecture-notes.md` (already
  persisted to `docs/features/leaflet-local-tiles/architecture-notes.md`) already carries the two
  real decisions — replacing the tile layer with a same-origin `L.imageOverlay`, and keeping
  Leaflet's default EPSG:3857 CRS — each with rationale, trade-offs, and a fitness function. This
  repo has no `docs/adr/` directory and no ADR convention; architecture-notes.md is the durable
  record for this decision and a separate ADR would only duplicate it.
- [ ] **Getting started / setup guides** — N/A. No new setup steps.
- [ ] **Migration guides** — N/A. No breaking changes or DB migrations.

### Code-Level Docs

- [x] The inline comment already added in `LocationsView.vue`'s `onMounted` (above the
  `L.imageOverlay` call) states why the image needs no network access and how its bounds were
  sized — correct and sufficient; no module- or function-level docstring convention exists
  elsewhere in this Vue SFC codebase, so none was added.
- [ ] No new exported functions were added that need JSDoc.

## Licence Record Verified

- `apps/ye-olde-magic-shop/public/sword-coast-map.LICENSE.txt` sits directly beside
  `sword-coast-map.svg`, as the spec's acceptance criterion requires.
- States the SVG is original work by Oscar Rieken, copyright 2026, licensed CC BY 4.0 — the same
  licence already used elsewhere in this project's documentation. No third-party redistribution
  risk.

## Known Issues / Follow-ups

- `tests/map.spec.ts:97` (`'should not have broken tiles'`) only queries `.leaflet-tile` elements.
  Those no longer exist now that the OSM tile layer is gone, so this specific assertion is
  vacuously true on every run. This is a deliberate, reviewed trade-off, not an oversight: the
  spec requires every existing test in this file to keep passing with unchanged text, so neither
  the developer nor code-reviewer touched it. Real coverage for "did the local image actually
  load under a fully blocked network" now lives in the new test at `tests/map.spec.ts:140`
  (`'locations view loads with every off-origin request blocked (NFR)'`), which asserts the
  `.leaflet-image-layer` `<img>` element's `complete`/`naturalWidth` state directly. QA confirmed
  by breaking the image path that only the new test — not the old one — catches the regression.
- No other follow-ups. No bugs were found during QA; the security review raised no findings.

## Files Changed By This Report

- `docs/features/README.md` — added the `leaflet-local-tiles` row to the delivered-features index.