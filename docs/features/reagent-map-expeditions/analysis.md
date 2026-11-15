---
feature: "reagent-map-expeditions"
bounded_context: "Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7"
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/src/api/expeditions.ts", "apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-fly-to-site.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts", "apps/ye-olde-magic-shop/.playwright-baselines/debian/globe-default-view.png", "apps/ye-olde-magic-shop/.playwright-baselines/ubi9/globe-default-view.png", "docs/features/reagent-map-expeditions/"]
issue_refs: ["reagent-map-expeditions"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: reagent-map-expeditions

## Summary

Draws every expedition GET /api/expeditions returns (the shop's three seeded expeditions plus the signed-in adventurer's own) as geodesic polylines on the existing Aurelis globe in ReagentMapView.vue, tagged with saturdayMetadata like the site entities already are, and makes each entry in the site list a control that flies the camera to that site, bringing it onto the near side of the globe even when it starts on the far side. Purely additive within the pattern reagent-map-globe already established (per-entity saturdayMetadata tagging, the window.__SATURDAY_CESIUM__ seam, the fetch-with-timeout API client shape) — no new package, endpoint, or layer boundary.

### Acceptance Criteria

- The view MUST draw every expedition GET /api/expeditions returns for the current adventurer as a geodesic polyline, both ends within 1 km of its stored ends, tagged with saturdayMetadata { id, name, kind: 'expedition' }.
  - Example: Signed out: 'Spring to Crag', 'Glade to Caves', and 'Cairn to Shoals' are drawn (spec AC3).
  - Example: Signed in as Elara: the three seeded expeditions plus 'Moonwell Spring to Frostcap Hollow' are drawn; 'Ember Crag to Saltmarsh Cairn' is not drawn (spec AC4).
- Choosing a site from the list MUST fly the camera so that site lands on screen on the near side of the globe, including a site that starts on the far side.
  - Example: From the default view, Tidewrack Shoals is on the far side; choosing it from the list brings it on screen on the near side (spec AC5).
- The @visual globe baseline MUST change only because the newly drawn expeditions now appear in it. Regenerate it in both runner images, view each changed PNG, and record what changed and why. No other baseline changes, and no tolerance is raised.
  - Example: Run scripts/ui-suite-in-chrome.sh --update-baselines and --ubi9 --update-baselines; only globe-default-view.png differs in each image; reagent-map-header.png, reagent-site-list.png, and reagent-map-error.png stay byte-identical.
- Every new or changed function MUST keep cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except this order's own delivery record under docs/features/reagent-map-expeditions/.

### Non-Functional Requirements

- **security**: The reagent map MUST send no request outside the app's own origin or the mock API's own origin, including the new expeditions fetch.
- **security**: The expeditions client MAY attach an Authorization Bearer header from localStorage('dnd_auth_token'), mirroring orders.ts's pattern, but MUST NOT require one: GET /api/expeditions already tolerates an absent or invalid token.
- **accessibility**: Each site-list entry MUST stay operable by keyboard (Enter/Space activates it) and MUST expose an accessible name that announces it flies the camera to that site, once it becomes a clickable control.

## Proposed Fitness Functions

- **Every new or changed function in this order keeps cyclomatic complexity below 7**
  - **Verification**: ESLint complexity rule capped at 6 per typescript-conventions.md, run in CI lint
  - **Owner**: devops-engineer
- **Every changed line is covered by a Playwright spec that can fail**
  - **Verification**: No coverage tool is wired for this Playwright-only app (no vitest/istanbul layer exists under apps/ye-olde-magic-shop/src). qa-engineer proves each new assertion fails against the pre-change view, per test-repair-contract.md, and records this in qa-report.md.
  - **Owner**: qa-engineer
- **The @visual globe baseline changes only for the reason this order states, and no other baseline changes**
  - **Verification**: scripts/ui-suite-in-chrome.sh --update-baselines and --ubi9 --update-baselines regenerate the PNGs; a human or agent views each diff before commit and states the reason in the delivery record
  - **Owner**: qa-engineer
- **No request from /reagent-map leaves the app's or mock API's own origin**
  - **Verification**: apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts's existing off-origin NFR test already asserts this; it continues to pass once the expeditions fetch uses the same runtime API_BASE resolution as reagentSites.ts
  - **Owner**: qa-engineer

## Out of Scope

- The line tool and saving an expedition: the next two orders.
- The full scenario suites E1 and F1: run by hand from the spec package, not by this pipeline.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7
- **Context Crossings**: apps/mock-api GET /api/expeditions — already shipped by expeditions-api; this order only consumes it, no backend change

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: Expedition (existing aggregate; read-only in this order)
- **Read Models / Projections**: Visible expeditions for the current adventurer (GET /api/expeditions)

### Affected Components

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Extend onMounted to also fetch expeditions and draw each as a geodesic polyline tagged saturdayMetadata kind:'expedition'; make each site-list entry a control that flies the camera to that site.
- `apps/ye-olde-magic-shop/src/api/expeditions.ts` — New client module: fetchExpeditions(), mirroring reagentSites.ts's runtime API_BASE resolution and fetch-with-timeout, and orders.ts's optional-Authorization-header pattern for a GET that tolerates no token.
- `apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts` — New Playwright UI spec covering expeditions-visibility.feature's two @level-ui scenarios (AC3, AC4).
- `apps/ye-olde-magic-shop/tests/reagent-map-fly-to-site.spec.ts` — New Playwright UI spec covering fly-to-site.feature (AC5).
- `apps/ye-olde-magic-shop/tests/reagent-map-resilience.spec.ts` — Extend with a third failure-path case: fetchExpeditions() rejecting after buildViewer() and fetchReagentSites() both resolve, mirroring the two existing cases.
- `apps/ye-olde-magic-shop/.playwright-baselines/debian/globe-default-view.png` — Regenerate only this baseline: the default view now also shows the three seeded expeditions.
- `apps/ye-olde-magic-shop/.playwright-baselines/ubi9/globe-default-view.png` — Same regeneration, ubi9 runner image.
- `docs/features/reagent-map-expeditions/` — Delivery record stating what changed in the globe baseline and why, per this order's third acceptance criterion.

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Create apps/ye-olde-magic-shop/src/api/expeditions.ts: types GeoPoint {lat,lng} and Expedition {id, name?, owner, from: GeoPoint, to: GeoPoint} matching lib/api/ShopApi.ts's shape; fetchExpeditions(): Promise<Expedition[]> against `${API_BASE}/api/expeditions`, same runtime API_BASE resolution and AbortController 5s timeout as reagentSites.ts.
2. In fetchExpeditions(), attach `Authorization: Bearer <token>` from localStorage('dnd_auth_token') only when a token is present, mirroring orders.ts's getAuthHeaders but omitting Content-Type (this is a GET).
3. In ReagentMapView.vue, import fetchExpeditions and add an `expeditions` ref alongside the existing `sites` ref.
4. Extend onMounted's Promise.allSettled to [buildViewer(), fetchReagentSites(), fetchExpeditions()]; if either the sites or the expeditions promise rejects, destroy the viewer and set mapError to the existing exact text 'Failed to load reagent sites. Try refreshing the page.' — do not reword it, since reagent-map-error.png's baseline must not change.
5. Add addExpeditionEntities(Cesium, viewer, ellipsoid, expeditions): for each expedition add a polyline entity from Cesium.Cartesian3.fromDegreesArray([from.lng, from.lat, to.lng, to.lat], ellipsoid) with arcType: Cesium.ArcType.GEODESIC, tagged saturdayMetadata { id, name: expedition.name ?? '', kind: 'expedition' }.
6. Call addExpeditionEntities after addSiteEntities in onMounted's happy path, once both sites and expeditions have resolved.
7. Wrap each site-list item's text in a `<button type="button">` that calls flyToSite(site) on click; reset its appearance (no border/background/margin change from the current plain text) so reagent-site-list.png's baseline does not change; give it an accessible name such as aria-label="Fly to {{ site.name }}" and a stable data-testid per site.
8. Add flyToSite(site: ReagentSite): keep the Cesium module and ellipsoid returned by buildViewer() in component-scope variables (same lifetime as `viewer`) and call viewer.camera.flyTo({ destination: Cesium.Cartesian3.fromDegrees(site.lng, site.lat, FLY_TO_HEIGHT_METERS, ellipsoid) }); do not special-case any one site.
9. Add a third failure-path case to reagent-map-resilience.spec.ts: fetchExpeditions() rejecting after buildViewer() and fetchReagentSites() both resolve still destroys the viewer and leaves no canvas, no __SATURDAY_CESIUM__, and no site-list items, across repeated mounts — same assertion shape as the two existing cases.

### QA Tasks

1. Write apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts for expeditions-visibility.feature's two @level-ui scenarios: signed out, exactly 'Spring to Crag', 'Glade to Caves', and 'Cairn to Shoals' are drawn, each end within 1 km of its seeded site; Elara signed in (log in via the mock API, inject the returned token into localStorage('dnd_auth_token') before navigating) also sees 'Moonwell Spring to Frostcap Hollow' and not 'Ember Crag to Saltmarsh Cairn'.
2. Write apps/ye-olde-magic-shop/tests/reagent-map-fly-to-site.spec.ts for fly-to-site.feature: from the default view confirm cesium.toScreen(Tidewrack Shoals) is null (far side), click its entry in the site list, then confirm cesium.toScreen(Tidewrack Shoals) returns a point within the canvas bounds (near side, on screen).
3. For every new assertion, confirm it fails against the pre-change view (or cite why the behavior is new) per test-repair-contract.md, and record this in qa-report.md.
4. Run the full reagent-map-*.spec.ts set under --project=chrome locally while developing; before ship, run scripts/ui-suite-in-chrome.sh and scripts/ui-suite-in-chrome.sh --ubi9 (no flags) to check the baselines still match.
5. Regenerate only globe-default-view.png in each image with scripts/ui-suite-in-chrome.sh --update-baselines and --ubi9 --update-baselines; view each PNG; confirm reagent-map-header.png, reagent-site-list.png, and reagent-map-error.png are unchanged.
6. Record in the delivery record exactly what changed in globe-default-view.png and why.

### Tech Writer Tasks

1. Write the docs/features/reagent-map-expeditions/ delivery record per CLAUDE.md's Delivery Summary Format, stating what changed in globe-default-view.png (both runner images) and why, and confirming no other baseline changed.

### DevOps Tasks

None

## Edge Cases and Risks

- fetchExpeditions() rejects after buildViewer() and fetchReagentSites() both resolve: Treat identically to a reagent-sites fetch failure: destroy the viewer and set mapError to the existing message, unchanged, so reagent-map-error.png's baseline does not change.
- An expedition has no name (Expedition.name is optional per lib/api/ShopApi.ts): saturdayMetadata.name falls back to '' rather than undefined, so every consumer of cesium.getEntities() always receives a string.
- A site is clicked again while a previous camera flight is still animating: Cesium's camera.flyTo cancels the in-flight animation and starts the new one; no extra guard is added.
- Choosing a site already on the near side of the globe: flyToSite behaves the same for any site; AC5's far-side scenario (Tidewrack Shoals) is the one the test exercises, not a special case in the implementation.
- A signed-in adventurer who is neither Elara nor Borin, with no expeditions of their own: GET /api/expeditions returns only the three seeded expeditions, same as signed out; no client-side filtering is added — the server already enforces visibility.

## Definition of Done

- Every expedition GET /api/expeditions returns is drawn as a geodesic polyline tagged saturdayMetadata, verified by reagent-map-expeditions.spec.ts for both the signed-out and Elara-signed-in cases.
- Choosing any site from the list flies the camera so that site lands on screen on the near side, verified by reagent-map-fly-to-site.spec.ts using Tidewrack Shoals from the default view.
- globe-default-view.png is regenerated in both runner images, each diff is reviewed by eye, and the delivery record states what changed and why; no other baseline changed.
- No file outside scope.paths changed, except this order's own delivery record under docs/features/reagent-map-expeditions/.
- Every new or changed function keeps cyclomatic complexity below 7.
- code-reviewer has approved the implementation.
- qa-engineer has proven every new assertion can fail.
