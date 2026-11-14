# Context Manifest: reagent-map-expeditions

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `features/reagent-map/expeditions-visibility.feature` — The Gherkin for this order's AC3:ui and AC4:ui covers -- business-language traceability for 'seeded expeditions drawn where seeded' and 'a signed-in adventurer's own expedition drawn with the seeded ones, never another adventurer's'. draw-line.feature and save-expedition.feature in the same directory belong to the next two orders and are excluded.
- `features/reagent-map/fly-to-site.feature` — The Gherkin for this order's AC5 cover -- flying to Tidewrack Shoals from the far side brings it to the near side. Single scenario, directly maps to the claim's second acceptance criterion.
- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — The file this order extends. Already draws reagent-site point entities tagged with saturdayMetadata {id, name, kind:'site'} and sets window.__SATURDAY_CESIUM__ -- expeditions must follow the identical per-entity tagging pattern with kind:'expedition', and the camera fly-to must use viewer.camera.flyTo the same way buildViewer() already uses camera.setView. Note for the analyst: Cesium's PolylineGraphics defaults to a straight Cartesian3 chord, not a surface-following path -- rendering 'geodesic' per the spec's settled gate needs arcType: Cesium.ArcType.GEODESIC (or equivalent) against this view's custom 6,000km-sphere ellipsoid, not Earth's default.
- `apps/ye-olde-magic-shop/src/api/reagentSites.ts` — The exact client pattern (type, runtime API_BASE resolution, fetch-with-timeout, throw-on-!ok) a new src/api/expeditions.ts module should mirror for fetchExpeditions().
- `apps/ye-olde-magic-shop/src/api/orders.ts` — Shows the only existing pattern in this app for attaching the Authorization header from localStorage('dnd_auth_token') to a GET. GET /api/expeditions is unauthenticated-tolerant but visibility-sensitive (shop's seeded expeditions always; the caller's own only with a valid token) -- the new expeditions client needs this same optional-auth-header shape, not orders.ts's full getAuthHeaders Content-Type variant verbatim.
- `packages/saturday-playwright-cesium/src/fixture.ts` — Defines the seam contract already in production use: getEntities() reads entity.polyline.positions as a {line: GeoPoint[]} geometry when saturdayMetadata is present (no new fixture code needed for expedition lines), and flyTo(lat,lng,heightMeters) already drives viewer.camera.flyTo + waitForSettled -- this is the exact mechanism AC5's 'fly to a site' must trigger from the production view, and the exact mechanism the new test will assert against. Read before designing either the view change or its test.
- `packages/saturday-playwright-cesium/src/geometry.ts` — Pure distance()/lineLength() helpers, radius-parameterized (never Earth's) -- what the new test asserts 'both ends within 1 km of stored ends' with, matching reagent-map-globe's precedent test.
- `apps/mock-api/data/expeditions.json` — The seed fixture: 3 shop-owned expeditions (AC3) plus Elara's and Borin's own (AC4), each {id, name, owner, from:{lat,lng}, to:{lat,lng}, savedAt}. GET /api/expeditions (apps/mock-api/server.js, not in scope.paths, not pinned) filters this list server-side by owner==='shop' or the session's own id -- the frontend only needs to render whatever the endpoint returns and attach the Authorization header when signed in; no client-side filtering logic is needed or correct to add.
- `apps/ye-olde-magic-shop/lib/api/ShopApi.ts` — Already-established Expedition/GeoPoint types ({id, name?, owner, from, to}) and visibleExpeditions()/saveExpedition() methods, used by the existing (out-of-scope, already-shipped) API-level spec reagent-map-api.spec.ts. The new src/api/expeditions.ts client's Expedition type should match this shape so UI and API-test code agree on the same domain object.
- `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` — House style for the sibling spec this order adds: @orieken/saturday-playwright-cesium's test/cesium fixture import, data-testid='reagent-site-list' locator pattern, and the existing 'Frostcap Hollow and Tidewrack Shoals stay single entities after flying to each' test already exercises cesium.flyTo() against the production view's seam -- the new fly-to-site test is a close sibling of this one, not a new pattern.

## Token Budget

- **Status**: OK — ~8655 tokens, 7.2% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `features/reagent-map/expeditions-visibility.feature` | 1765 | 441 |
| `features/reagent-map/fly-to-site.feature` | 524 | 131 |
| `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` | 5962 | 1490 |
| `apps/ye-olde-magic-shop/src/api/reagentSites.ts` | 1101 | 275 |
| `apps/ye-olde-magic-shop/src/api/orders.ts` | 1865 | 466 |
| `packages/saturday-playwright-cesium/src/fixture.ts` | 10891 | 2722 |
| `packages/saturday-playwright-cesium/src/geometry.ts` | 2061 | 515 |
| `apps/mock-api/data/expeditions.json` | 1087 | 271 |
| `apps/ye-olde-magic-shop/lib/api/ShopApi.ts` | 1537 | 384 |
| `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` | 7841 | 1960 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `DOMAIN_DICTIONARY.md -- section 7's five terms (Adventurer, Reagent Site, Expedition, Seeded Expedition, Aurelis) are already named correctly throughout the pinned files and the claim's own Domain Language section; no separate read needed.`
- `apps/mock-api/server.js -- GET/POST /api/expeditions already shipped by expeditions-api; its visibility-filter behavior is summarized in apps/mock-api/data/expeditions.json's pinned reason rather than reading the full route handler, since no backend change is in this order's scope.paths.`
- `docs/specs/reagent-map-spec/acceptance-criteria.json, domain-model.json, discovery.json, example-map.json, README.md, corrections.md -- the full spec package; the claim document already pins the digest and restates AC3/AC4/AC5 plus every settled gate (1km tolerance, geodesic, same-place rule, visibility rule, Aurelis radius) in plain language, so re-reading the raw package is redundant for this order.`
- `docs/features/reagent-map-globe/analysis.md, implementation-notes.md -- the prior delivery in this same bounded context. Its established pattern (per-entity saturdayMetadata tagging, window.__SATURDAY_CESIUM__ lifecycle, no-ion/no-CDN Cesium config) is already carried forward via the pinned ReagentMapView.vue itself; the lesson is 'extend this file's existing pattern, do not re-derive it' -- folded into that file's pinned reason rather than re-reading the full prior artifact. No retrospective.md exists for reagent-map-globe, reagent-sites-api, expeditions-api, or reagent-map-visual-baselines to check for What-Went-Poorly/What-To-Improve notes.`
- `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts, reagent-map-pan-zoom.spec.ts, reagent-map-resilience.spec.ts, reagent-map-a11y-check.spec.ts, reagent-map-feature-flag.spec.ts -- existing UI test files with no expedition or fly-to coverage; unrelated to this order's change surface.`
- `apps/ye-olde-magic-shop/.playwright-baselines/**/*.png -- binary baseline images regenerated by scripts/ui-suite-in-chrome.sh per the claim's third acceptance criterion; not readable as text context.`
- `apps/ye-olde-magic-shop/src/router/index.ts, src/App.vue, src/featureFlags.ts -- the VITE_FEATURE_REAGENT_MAP route/nav-link plumbing shipped by reagent-map-globe; this order is additive within that existing flag and route, touching neither.`
- `apps/ye-olde-magic-shop/src/views/LocationsView.vue, src/components/ItemCard.vue, src/stores/cartStore.ts -- different bounded context (shop/cart, Leaflet locations view); excluded by bounded-context exclusion.`
- `mcp__loom__impact / mcp__loom__dependents_of / mcp__loom__find_files -- not available as registered tools in this session, so the blast radius below was not graph-derived; pinnedFiles were chosen by reading the claim, the spec package summaries, and Grep/manual exploration of imports and existing test patterns instead.`
