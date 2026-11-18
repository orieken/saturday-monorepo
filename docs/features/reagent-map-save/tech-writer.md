Domain dictionary already has these terms defined (established by the prior reagent-map features), so no dictionary update is needed here.

Now producing the docs-report.md artifact.

```markdown
# Documentation Report: reagent-map-save

## Status: PASS

## Summary

Saving an expedition from the line tool is now documented at the point a reader needs it: inline
in the two files it touches, plus the acceptance criteria already carried by `analysis.md`. The
reagent map itself went from an opt-in preview to the app's default-on behavior
(`VITE_FEATURE_REAGENT_MAP=true`); no README, CHANGELOG, or ADR exists yet for the reagent map
feature family, so this report follows the same convention the three prior reagent-map orders
(`reagent-map-globe`, `reagent-map-line-tool`, `reagent-map-expeditions`) already set rather than
introducing a new doc file this feature's scope doesn't call for.

## What Was Built (plain English)

A signed-in adventurer can now press **Save** on a completed line (two points placed on the land)
and have it stored as their own expedition through the mock API's existing `POST /api/expeditions`.
It reappears, drawn, every time they reopen the reagent map. Saving is refused with a visible
message, and nothing is stored, when nobody is signed in or when the two points are the same place
within 1 km. Leaving the map with an unsaved line discards it — nothing is drawn or stored on
return. The reagent map, previously hidden behind a build-time flag, is now on by default; the flag
can still turn it off.

## Documentation Checklist

### Always Update
- [ ] **CHANGELOG.md** — Not applicable. No `CHANGELOG.md` exists at the repo root or under
  `apps/ye-olde-magic-shop/`, and none of the three prior reagent-map feature deliveries created
  one. Introducing a changelog file as a side effect of this order would be new scope outside
  `scope.paths` and outside this feature's request.
- [ ] **README.md** — Not updated. `apps/ye-olde-magic-shop/README.md` documents only the app's
  dev/test commands and has never mentioned the reagent map, Cesium globe, or expeditions feature
  across any of its three prior deliveries. Matching that established convention rather than
  introducing reagent-map documentation asymmetrically in this one order.

### Update if Applicable
- [ ] **API documentation** — Not applicable. `POST /api/expeditions` (`apps/mock-api/server.js`)
  is unchanged by this order (per `analysis.md`'s Out of Scope and confirmed in `security-report.md`
  and `qa-report.md`); it was already in service for `reagent-map-expeditions`.
- [ ] **Configuration docs** — No dedicated env-var reference doc exists for this app to update.
  `apps/ye-olde-magic-shop/.env` itself is the single source of truth and now reads
  `VITE_FEATURE_REAGENT_MAP=true`, the one line this order changes.
- [ ] **Architecture Decision Record** — Not applicable. `architecture-notes.md` records this
  order's structural decisions (the cross-file `expeditions.json` lock, the save-handler shape)
  as `judgment-only`/fitness-function entries, not an ADR-worthy constraint on future features.
- [ ] **Getting started / setup guides** — Not applicable. No new setup step, dependency, or
  environment variable was introduced; the one existing flag's default value flipped.
- [ ] **Migration guides** — Not applicable. No schema or data migration in this order.

### Code-Level Docs
- [x] **Module-level docs** — No new files were created (`implementation-notes.md` confirms
  `apps/ye-olde-magic-shop/.env`, `src/api/expeditions.ts`, and `src/views/ReagentMapView.vue` were
  all modified, none created), so no new module-level docstring is needed.
- [x] **Function docs** — `saveExpedition()`, `saveLine()`, and `promoteSavedExpedition()` carry
  intention-revealing names per `implementation-notes.md`'s Interface Design section and this
  codebase's existing convention of no comments beyond named helpers (`fetchExpeditions()`,
  `addExpeditionEntities()` have none either) — consistent, not a gap.
- [x] **Type hints** — `saveExpedition(from: GeoPoint, to: GeoPoint): Promise<Expedition>` reuses
  the existing `GeoPoint`/`Expedition` types; no new type was introduced that needs annotation.

## Feature Behavior (for the record)

- **Save**: a `data-testid="save-line-tool"` button (`aria-label="Save this expedition"`) appears
  once a line is drawn, calling `POST /api/expeditions` with `{ from, to }` only — no owner field;
  the server derives ownership from the session token.
- **Sign-in gate (AC12)**: signed out, saving shows "Sign in to save this expedition." and makes no
  network request.
- **Zero-length gate (AC13)**: two points under 1 km apart, saving shows the server's own message
  ("Zero-length line: from and to are the same place") verbatim and stores nothing.
- **Persistence (AC11)**: on success the server's returned expedition is drawn immediately and is
  present on every later visit, within 1 km of the placed points.
- **Discard on exit (AC15)**: navigating away from an unsaved line leaves nothing drawn or stored on
  return.
- **Feature flag**: `VITE_FEATURE_REAGENT_MAP` now defaults to `true` in
  `apps/ye-olde-magic-shop/.env`; setting it to `false` still hides the nav link and route.

## Test Isolation Note (for whoever next touches these specs)

Per `qa-report.md`, the three `reagent-map-*.spec.ts` files each keep their own local copy of
`acquireExpeditionsLock()`/`releaseExpeditionsLock()` rather than sharing one via import — importing
a sibling spec file re-executes its top-level `test.describe.configure({ mode: 'serial' })`, which
collides with the importing file's own call. `STALE_LOCK_MS` is `180_000`, not the `30_000`
originally proposed, because `30_000` measurably fired on a legitimate, still-running lock holder
under `--workers=3 --repeat-each=3`. Both deviations are recorded in the specs' own code comments
(code-level docs, not duplicated into a separate doc per this checklist's "don't repeat what code
already makes obvious" rule).

## Manual Steps Required

None. No secret, migration, or infrastructure change accompanies this order.

## Known Issues / Follow-ups

- `tests/reagent-map-pan-zoom.spec.ts` remains a known, pre-existing flake under load (drag misses by
  ~312px) — out of this order's scope per the claim's own Test Approach note, unchanged by this
  delivery.
- Server-side rate limiting / idempotency for `POST /api/expeditions` is still client-side-only
  (the `isSaving` guard); `security-report.md` tracks this as pre-existing [LOW] tech debt, not a
  regression, and recommends no fix within this order's scope.
```