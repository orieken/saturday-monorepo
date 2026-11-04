---
feature: "reagent-map-globe"
bounded_context: "Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)"
domain_terms: ["Reagent Site", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/featureFlags.ts", "apps/ye-olde-magic-shop/src/router/index.ts", "apps/ye-olde-magic-shop/src/App.vue", "apps/ye-olde-magic-shop/src/api/reagentSites.ts", "apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/src/main.ts", "apps/ye-olde-magic-shop/public/aurelis.jpg", "apps/ye-olde-magic-shop/package.json", "apps/ye-olde-magic-shop/vite.config.ts", "apps/ye-olde-magic-shop/playwright.config.ts", "apps/ye-olde-magic-shop/.env", "scripts/ui-suite-in-chrome.sh", "apps/ye-olde-magic-shop/tests/reagent-map-feature-flag.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts"]
issue_refs: ["reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: reagent-map-globe

## Summary

Add a feature-flagged `/reagent-map` route to the ye-olde-magic-shop Vue app showing Aurelis as an offline 3D Cesium globe, with the eight reagent sites from the existing GET /api/reagent-sites endpoint plotted as entities and listed by name beside the globe. The view implements the saturday-playwright-cesium seam contract (`window.__SATURDAY_CESIUM__`, `saturdayMetadata` on each site entity) so the sibling test package can drive it. The route and its nav link exist only when the new build-time flag VITE_FEATURE_REAGENT_MAP is `true`; `.env` keeps it `false` so the app is otherwise unchanged. This is slice 3 of 6 in the reagent-map-spec D1 split — no expeditions, fly-to, line tool, or save behavior, and no backend changes (reagent-sites-api already shipped the endpoint this view consumes).

### Acceptance Criteria

- The /reagent-map route and its nav link exist only when VITE_FEATURE_REAGENT_MAP is exactly "true"; with it unset, empty, or any other value, there is no route and no link, and the app is otherwise unchanged. .env sets the flag to false.
  - Example: VITE_FEATURE_REAGENT_MAP=true -> route exists, nav link renders
  - Example: VITE_FEATURE_REAGENT_MAP=false -> no route, no link
  - Example: VITE_FEATURE_REAGENT_MAP unset -> no route, no link
- The view runs Cesium with no network dependency: no ion token or ion imagery, baseLayerPicker and geocoder off, one equirectangular image of Aurelis served from the app's own files, Cesium's own static assets served same-origin via CESIUM_BASE_URL, scene mode SCENE3D, and an ellipsoid that is a sphere of radius 6,000 km.
  - Example: viewer.scene.globe.ellipsoid has radii (6000000, 6000000, 6000000)
  - Example: viewer.scene.mode === Cesium.SceneMode.SCENE3D
  - Example: viewer.baseLayerPicker and viewer.geocoder are both false/absent
- The view sets window.__SATURDAY_CESIUM__ = { viewer, Cesium } and attaches saturdayMetadata { id, name, kind: 'site' } to every reagent site entity, per packages/saturday-playwright-cesium/README.md's seam contract.
- Every one of the eight reagent sites is shown within 1 km of its published position and listed by name beside the globe, including Frostcap Hollow near the north pole and Tidewrack Shoals on the 180th meridian, each shown exactly once with no clipping or duplication.
  - Example: Frostcap Hollow (lat 85, lng 0) appears once, within 1 km of that position
  - Example: Tidewrack Shoals (lat 0, lng 180) appears once, within 1 km of that position, not duplicated at the antimeridian wrap
- No request made while the reagent map view is open leaves the app's own origin or the mock API's own origin.
- A small pure function decides whether the /reagent-map route exists, tested for "true", "false", and unset without a browser. scripts/ui-suite-in-chrome.sh starts the app with VITE_FEATURE_REAGENT_MAP=true so the pre-ship check exercises the view, while .env keeps the flag false. New specs named tests/reagent-map-*.spec.ts run only in the chrome Playwright project; npm run test:e2e (bundled Chromium) still passes without them.
- Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail, and no file outside scope.paths changes except this run's own delivery record under docs/features/reagent-map-globe/.

### Non-Functional Requirements

- **security**: GET requests, static asset loads, and Cesium's own asset/worker loads triggered by the /reagent-map view MUST NOT leave the app's own origin or the mock API's own origin — no ion imagery, no CDN, no external geocoder.

## Proposed Fitness Functions

- **Zero off-origin network requests during /reagent-map page load and settle**
  - **Verification**: Playwright test on page.on('request') collecting every request hostname and asserting each is the dev server's own origin or the proxied mock API origin — mirrors packages/saturday-playwright-cesium/tests/fixture.spec.ts's existing 'makes no request off its own origin' test.
  - **Owner**: qa-engineer
- **isReagentMapEnabled(value) returns true only for the literal string "true", false for "false" and for undefined**
  - **Verification**: tests/reagent-map-feature-flag.spec.ts using @playwright/test's test/expect with no page fixture, same shape as the existing tests/playwright-config.spec.ts pattern — runs with no browser launched.
  - **Owner**: qa-engineer
- **Files named tests/reagent-map-*.spec.ts run only under the Playwright 'chrome' project; npm run test:e2e (chromium project) passes unaffected by their presence**
  - **Verification**: playwright.config.ts's chromium project carries a testIgnore for the reagent-map-*.spec.ts glob; CI runs both npm run test:e2e and npm run test:chrome (or scripts/ui-suite-in-chrome.sh) and checks both exit codes.
  - **Owner**: devops-engineer
- **Every reagent site entity renders within 1 km of its seeded position, including the pole and antimeridian cases, with no duplicate entity**
  - **Verification**: tests/reagent-map-globe.spec.ts using @orieken/saturday-playwright-cesium's cesium.getEntities() and distance() geometry helper against apps/mock-api/data/reagent-sites.json's seeded lat/lng.
  - **Owner**: qa-engineer
- **No file outside scope.paths changes, excepting this run's own delivery record under docs/features/reagent-map-globe/**
  - **Verification**: git diff --name-only against the base branch, checked by loom's own run gate.
  - **Owner**: devops-engineer
- **Cyclomatic complexity stays below 7 for every function added or changed in apps/ye-olde-magic-shop/src/** and vite.config.ts**
  - **Verification**: Manual review by code-reviewer — apps/ye-olde-magic-shop has no ESLint complexity rule configured yet, so this cannot run as an automated CI check for this app.
  - **Owner**: code-reviewer

## Out of Scope

- Expeditions, fly-to, the line tool, and saving — the next three orders in the six-way split (reagent-map-expeditions, reagent-map-line-tool, reagent-map-save).
- The existing Leaflet-based /locations view — unchanged.
- Backend/API changes — GET /api/reagent-sites already exists from reagent-sites-api; this order only consumes it.
- Running the full E1 (UI) and F1 (API) scenario suites — the spec package says these are run by hand, not by this order's automated tests.
- Any change to apps/ye-olde-magic-shop/index.html, apps/ye-olde-magic-shop/README.md, or any .gitignore file — none are in scope.paths.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Expedition Planning — Reagent Map UI (apps/ye-olde-magic-shop frontend)
- **Context Crossings**: apps/mock-api — read-only consumer of the existing GET /api/reagent-sites endpoint (reagent-sites-api); no backend change in this order, packages/saturday-playwright-cesium — Test Automation bounded context; this order implements its published seam contract (window.__SATURDAY_CESIUM__, saturdayMetadata) and reuses its already-reviewed fixture-globe.js Cesium configuration and fixtures/assets/aurelis.jpg texture as reference/source material, read-only (packages/saturday-playwright-cesium/** is not in this order's scope.paths)

### Domain Events (Event Storming Lite)

- **Commands**: OpenReagentMap
- **Events Produced**: ReagentMapOpened
- **Owning Aggregates**: ReagentSite
- **Read Models / Projections**: ReagentSiteCatalog

### Affected Components

- `apps/ye-olde-magic-shop/src/featureFlags.ts` — New module exporting the pure function isReagentMapEnabled(value: string | undefined): boolean, the small pure predicate the route and nav link both gate on, tested without a browser.
- `apps/ye-olde-magic-shop/src/router/index.ts` — Conditionally append the { path: '/reagent-map', component: ReagentMapView } route only when isReagentMapEnabled(import.meta.env.VITE_FEATURE_REAGENT_MAP) is true, so the route genuinely does not exist when the flag is off.
- `apps/ye-olde-magic-shop/src/App.vue` — Add a nav <router-link to="/reagent-map"> (data-testid="reagent-map-link") gated by the same isReagentMapEnabled(...) check, alongside the existing Shop/Locations/Cart links.
- `apps/ye-olde-magic-shop/src/api/reagentSites.ts` — New API client module: ReagentSite type ({id, name, lat, lng}) and fetchReagentSites(): Promise<ReagentSite[]>, mirroring src/api/items.ts's existing API_BASE/fetch convention against GET /api/reagent-sites.
- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — New view: builds the no-network Cesium Viewer (ellipsoid/globe/imagery config lifted verbatim from packages/saturday-playwright-cesium/fixtures/assets/fixture-globe.js), fetches reagent sites, adds one tagged point entity per site, renders the site-name list beside the globe, exposes window.__SATURDAY_CESIUM__ on mount and tears it down on unmount (mirroring LocationsView.vue's __SATURDAY_LEAFLET_MAP__ lifecycle pattern).
- `apps/ye-olde-magic-shop/src/main.ts` — Set (window as any).CESIUM_BASE_URL = '/cesium/' before mount — inert unless the reagent map route is actually visited, so the app stays unchanged for every other route when the flag is off; keeps Cesium itself out of main.ts's own static imports.
- `apps/ye-olde-magic-shop/public/aurelis.jpg` — New static asset: a copy of the already-reviewed, same-origin equirectangular Aurelis texture at packages/saturday-playwright-cesium/fixtures/assets/aurelis.jpg (1024x512px). Do not regenerate it — copy the existing, reviewed file.
- `apps/ye-olde-magic-shop/package.json` — Add "cesium": "1.146.0" (exact version, no caret range) to dependencies, matching packages/saturday-playwright-cesium/package.json and the version already resolved in pnpm-lock.yaml — avoids a second resolved Cesium version in the lockfile.
- `apps/ye-olde-magic-shop/vite.config.ts` — Add a small inline dev-server/build plugin that serves node_modules/cesium/Build/Cesium under /cesium/ in dev (configureServer) and copies the same directory into the build's outDir/cesium on build (closeBundle/writeBundle) — no new dependency, and no write into public/ or need for a new .gitignore entry (dist/ is already ignored; .gitignore is not in scope.paths).
- `apps/ye-olde-magic-shop/playwright.config.ts` — Add a testIgnore for the reagent-map-*.spec.ts glob on the 'chromium' project only, so npm run test:e2e/test:pw stay unaffected; the 'chrome' project picks the new specs up with no change needed on its side.
- `apps/ye-olde-magic-shop/.env` — Add VITE_FEATURE_REAGENT_MAP=false, keeping the flag off for every normal local run.
- `scripts/ui-suite-in-chrome.sh` — start_app_and_mock_api() must export VITE_FEATURE_REAGENT_MAP=true to the containerized npm run dev process, so the pre-ship chrome-project run exercises the view while every other invocation keeps the flag at .env's false.
- `apps/ye-olde-magic-shop/tests/reagent-map-feature-flag.spec.ts` — New test: isReagentMapEnabled('true')===true, isReagentMapEnabled('false')===false, isReagentMapEnabled(undefined)===false — no page fixture, no browser launched.
- `apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts` — New chrome-project test using @orieken/saturday-playwright-cesium's cesium fixture: asserts the seam contract, all 8 sites within 1 km (including the pole/antimeridian cases), SCENE3D, the 6,000 km ellipsoid, and zero off-origin requests.

### Data Model Changes

None

### API Changes

None

### New Dependencies

- cesium@1.146.0 (apps/ye-olde-magic-shop runtime dependency — exact version, matching packages/saturday-playwright-cesium and the version already present in pnpm-lock.yaml; no new transitive version introduced)

## Task List

None

### Developer Tasks

1. Add "cesium": "1.146.0" to apps/ye-olde-magic-shop/package.json dependencies (exact version, matching pnpm-lock.yaml's existing resolution from packages/saturday-playwright-cesium); run pnpm install and confirm no second Cesium version appears in pnpm-lock.yaml.
2. Copy packages/saturday-playwright-cesium/fixtures/assets/aurelis.jpg into apps/ye-olde-magic-shop/public/aurelis.jpg verbatim — reuse the already-reviewed 1024x512px equirectangular texture, do not regenerate it.
3. Create apps/ye-olde-magic-shop/src/featureFlags.ts exporting `export function isReagentMapEnabled(value: string | undefined): boolean { return value === 'true'; }` — single equality check, complexity 1.
4. In apps/ye-olde-magic-shop/vite.config.ts, add a small inline plugin object (no new dependency) with a configureServer(server) hook serving node_modules/cesium/Build/Cesium under the /cesium/ URL prefix in dev, and a closeBundle (or writeBundle) hook copying the same directory tree into <outDir>/cesium on build. Keep generated output out of public/ and out of any committed path — dist/ is already in the root .gitignore, and .gitignore itself is outside this order's scope.paths.
5. In apps/ye-olde-magic-shop/src/main.ts, set `(window as any).CESIUM_BASE_URL = '/cesium/';` before `app.mount('#app')` — inert for every route except /reagent-map, which imports cesium dynamically.
6. Modify apps/ye-olde-magic-shop/src/router/index.ts: import featureFlags and a lazy/dynamic ReagentMapView, and push the { path: '/reagent-map', component: ReagentMapView } route into the routes array only when isReagentMapEnabled(import.meta.env.VITE_FEATURE_REAGENT_MAP) is true, so the route object itself does not exist when the flag is off.
7. Modify apps/ye-olde-magic-shop/src/App.vue: add a nav <router-link to="/reagent-map" data-testid="reagent-map-link"> next to the existing Locations link, rendered only when isReagentMapEnabled(import.meta.env.VITE_FEATURE_REAGENT_MAP) is true.
8. Create apps/ye-olde-magic-shop/src/api/reagentSites.ts: `export type ReagentSite = { id: string; name: string; lat: number; lng: number }` and `export async function fetchReagentSites(): Promise<ReagentSite[]>` fetching `${API_BASE}/api/reagent-sites`, mirroring src/api/items.ts's existing API_BASE resolution and error-handling shape.
9. Create apps/ye-olde-magic-shop/src/views/ReagentMapView.vue. On mount: dynamically `const Cesium = await import('cesium')` and dynamically import 'cesium/Build/Cesium/Widgets/widgets.css'; build an Ellipsoid/Globe sized to AURELIS_RADIUS_METERS = 6_000_000 (all three radii equal); build a SingleTileImageryProvider from '/aurelis.jpg' with Cesium.Rectangle.fromDegrees(-180,-90,180,90); construct the Viewer with that globe and `baseLayer: new Cesium.ImageryLayer(imageryProvider)`, and baseLayerPicker/geocoder/animation/timeline/fullscreenButton/homeButton/sceneModePicker/navigationHelpButton/infoBox/selectionIndicator/skyBox/skyAtmosphere/shouldAnimate all false — the exact option set already proven in packages/saturday-playwright-cesium/fixtures/assets/fixture-globe.js. Explicitly set viewer.scene.mode = Cesium.SceneMode.SCENE3D.
10. In the same onMounted: call fetchReagentSites(), add one point entity per site via viewer.entities.add({ position: Cesium.Cartesian3.fromDegrees(site.lng, site.lat, 0, ellipsoid), point: { pixelSize: 10 } }), and set entity.saturdayMetadata = { id: site.id, name: site.name, kind: 'site' }. Render the fetched sites as a plain name list (data-testid="reagent-site-list") in the template beside the globe container. Set window.__SATURDAY_CESIUM__ = { viewer, Cesium } once construction is complete.
11. In onBeforeUnmount: call viewer.destroy() and set window.__SATURDAY_CESIUM__ = undefined, mirroring LocationsView.vue's onBeforeUnmount handling of __SATURDAY_LEAFLET_MAP__.
12. In apps/ye-olde-magic-shop/playwright.config.ts, add `testIgnore: [/reagent-map-.*\.spec\.ts$/]` to the 'chromium' project's config only, so npm run test:e2e / test:pw skip the new reagent-map-*.spec.ts files; leave the 'chrome' project as-is so it collects them.
13. Add `VITE_FEATURE_REAGENT_MAP=false` to apps/ye-olde-magic-shop/.env.

### QA Tasks

1. Write apps/ye-olde-magic-shop/tests/reagent-map-feature-flag.spec.ts using @playwright/test's test/expect with no page fixture (same shape as the existing tests/playwright-config.spec.ts), asserting isReagentMapEnabled('true') === true, isReagentMapEnabled('false') === false, and isReagentMapEnabled(undefined) === false.
2. Write apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts (chrome project), importing { test } from '@orieken/saturday-playwright-cesium': navigate to /reagent-map, await cesium.waitForSettled(), then cesium.getEntities() and assert exactly 8 entities with kind: 'site', each within 1000m (cesium.globeRadiusMeters() should read 6,000,000) of its seeded lat/lng from apps/mock-api/data/reagent-sites.json, via the package's exported distance() helper.
3. In the same or a sibling chrome-project spec, cover AC2 explicitly: cesium.flyTo Frostcap Hollow's (85, 0) and Tidewrack Shoals' (0, 180) positions in turn and assert each still appears exactly once, within 1 km, with no duplicate entity at the antimeridian wrap.
4. Assert window.__SATURDAY_CESIUM__ is { viewer, Cesium } (via page.evaluate) and that every site entity carries saturdayMetadata { id, name, kind: 'site' }.
5. Assert viewer.scene.mode === Cesium.SceneMode.SCENE3D and viewer.scene.globe.ellipsoid's three radii all equal 6,000,000, via page.evaluate against window.__SATURDAY_CESIUM__.
6. Add a zero-off-origin-request assertion for the /reagent-map page load: collect every request via page.on('request'), assert each hostname matches the dev server's own origin or the proxied mock API origin — mirror packages/saturday-playwright-cesium/tests/fixture.spec.ts's existing pattern for this exact check.
7. Assert the site-name list (data-testid="reagent-site-list") renders all 8 reagent site names from GET /api/reagent-sites.
8. Run npm run test:e2e and confirm it passes unaffected by the new reagent-map-*.spec.ts files; run npm run test:chrome (or scripts/ui-suite-in-chrome.sh) and confirm the new specs run and pass under the chrome project.
9. Confirm via git diff --name-only against the base branch that no file outside scope.paths changed, excepting docs/features/reagent-map-globe/**.

### Tech Writer Tasks

None

### DevOps Tasks

1. Modify scripts/ui-suite-in-chrome.sh's start_app_and_mock_api() to export VITE_FEATURE_REAGENT_MAP=true to the containerized `npm run dev` process (e.g. `docker exec -d -e VITE_FEATURE_REAGENT_MAP=true "$CONTAINER_NAME" bash -lc '...'`), so the pre-ship chrome-project run exercises the reagent map view while apps/ye-olde-magic-shop/.env keeps the flag false for every other invocation.
2. Confirm CI (or the equivalent pre-ship check) runs both npm run test:e2e (chromium) and npm run test:chrome / scripts/ui-suite-in-chrome.sh (chrome), and that both exit codes gate the build.

## Edge Cases and Risks

- VITE_FEATURE_REAGENT_MAP is set to a near-true value such as "TRUE", "1", or " true" (whitespace/case variants).: isReagentMapEnabled uses strict equality against the literal string 'true' — any other value, including case or whitespace variants, is treated as off. This is deliberate and exactly what the pure-function unit test must cover for 'false' and 'unset'; no normalization is added.
- A reagent site's seeded longitude is exactly 180 (Tidewrack Shoals) or latitude is near 90 (Frostcap Hollow, lat 85).: Cesium.Cartesian3.fromDegrees handles the antimeridian and near-pole cases natively when given the correct ellipsoid; the fixture-globe.js pattern already proves this renders once, not duplicated — no special-casing is added in the view beyond passing the ellipsoid through.
- GET /api/reagent-sites is slow, fails, or the mock API is not running.: fetchReagentSites() follows items.ts's existing convention: throw on a non-OK response. The view should let the error surface (e.g. an unhandled rejection becomes visible in dev) rather than silently showing an empty globe — no new retry or fallback logic is introduced, consistent with architecture-guardrails.md's prohibition on ad-hoc retry loops; a real CircuitBreaker is out of scope for a single same-origin GET with no stated latency SLA.
- The flag is true but the Cesium static assets at /cesium/ are missing or stale (e.g. a build that forgot the closeBundle copy step).: Cesium's worker bootstrap fails fast with its own console error when CESIUM_BASE_URL assets are absent; this is a build-configuration defect to catch via the QA chrome-project run (cesium.waitForSettled() will reject naming the stuck condition), not a runtime fallback the app needs to handle gracefully.
- A developer navigates directly to /reagent-map via URL while the flag is off.: Because the route object itself is never added to the router when the flag is off (not merely hidden via a guard), vue-router falls through to its normal unmatched-route behavior — identical to typing any other nonexistent path today. No new 404 view is introduced.

## Definition of Done

- /reagent-map route and its nav link exist if and only if VITE_FEATURE_REAGENT_MAP is the literal string 'true'; apps/ye-olde-magic-shop/.env sets it to false.
- The view's Cesium Viewer uses a 6,000,000m-radius ellipsoid, scene mode SCENE3D, baseLayerPicker/geocoder off, a same-origin equirectangular Aurelis image, and same-origin Cesium static assets served via CESIUM_BASE_URL — with zero network requests leaving the app's or mock API's own origin.
- window.__SATURDAY_CESIUM__ = { viewer, Cesium } is set, and every one of the 8 reagent site entities carries saturdayMetadata { id, name, kind: 'site' }.
- All 8 reagent sites, including Frostcap Hollow and Tidewrack Shoals, are shown exactly once, each within 1 km of its published position, and listed by name beside the globe.
- isReagentMapEnabled is covered by a no-browser unit test for 'true', 'false', and unset.
- tests/reagent-map-*.spec.ts run only under the Playwright 'chrome' project; npm run test:e2e passes unaffected by their presence; scripts/ui-suite-in-chrome.sh starts the container with the flag on and its chrome-project run passes.
- Every new/changed function is reviewed for cyclomatic complexity below 7 (judgment-only — no ESLint complexity rule exists yet in apps/ye-olde-magic-shop).
- git diff shows no file changed outside scope.paths, excepting docs/features/reagent-map-globe/**.
