# Context Manifest: reagent-map-save

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Primary file to modify: holds draftPoints/draftEntities and the toolStarted/lineToolStatus state from reagent-map-line-tool; the save control, the sign-in gate (AC12), and the zero-length rejection message (AC13) all attach here, and AC15's 'leaves and reopens' behavior depends on its existing onMounted/onBeforeUnmount lifecycle holding no persistence of its own.
- `apps/ye-olde-magic-shop/src/api/expeditions.ts` — Only exports fetchExpeditions() today; this order adds a saveExpedition()-style POST client and must match its existing API_BASE/AbortController/REQUEST_TIMEOUT_MS/getAuthHeaders conventions exactly, including how a non-ok response's body is surfaced to the caller.
- `apps/ye-olde-magic-shop/src/stores/authStore.ts` — Defines isAuthenticated/token used to gate the save action before any network call (AC12) — the view must read this store, not re-derive auth state or call it after the fetch.
- `apps/ye-olde-magic-shop/tests/reagent-map-api.spec.ts` — In scope.paths (tests/reagent-map-*.spec.ts). Already writes to apps/mock-api/data/expeditions.json and restores it (beforeAll snapshot + afterEach rewrite) and already uses test.describe.configure({ mode: 'serial' }) for its own save/reject tests — this is the one existing precedent for the claim's new isolation requirement ('say how the writers are kept apart from the readers'), and it only serializes within its own file, not across files/workers, which is exactly the gap the new AC is pointing at.
- `apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts` — The claim names this file directly: 'reagent-map-expeditions.spec.ts counts exactly Elara's expeditions while the suite runs three workers.' It asserts exact entity counts (toHaveLength) against GET /api/expeditions and must stay correct if a save-writing spec runs concurrently in another worker against the same shared data file.
- `apps/ye-olde-magic-shop/playwright.config.ts` — Reference only, outside scope.paths (lives above src/**, not matched by the claim's globs). Defines fullyParallel: true, workers: undefined outside CI, and a single shared webServer — the analyst needs this to know the isolation fix for AC's worker requirement has to be achieved inside the in-scope spec files themselves (describe-level serial/ordering), since this config file cannot be edited under the claim's own scope.paths.
- `apps/mock-api/server.js` — Reference only, outside scope.paths. POST /api/expeditions (readExpeditions/writeExpeditions, lines ~29-163) does a synchronous full-file read on every GET and full-file overwrite on every POST, with no in-memory cache or locking, and the webServer in playwright.config.ts starts exactly one mock-api process shared by every worker — this is the concrete mechanism behind the new 'no test sees an expedition another test saved' requirement, and it must not be changed (already implements auth/validation/1km-rejection per prior orders).
- `apps/ye-olde-magic-shop/.env` — This order's own AC requires flipping VITE_FEATURE_REAGENT_MAP to true while confirming the flag can still turn the nav link/route off when set to false.
- `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` — In scope.paths. Defines the four existing @visual regions (header/nav link, site list, error state, globe) that must be reviewed per-PNG after regeneration in both runner images; the analyst needs its exact assertions to judge which baseline(s) the nav-link-now-visible-by-default change can plausibly touch.
- `docs/specs/reagent-map-spec/acceptance-criteria.md` — Authoritative Given/When/Then wording for AC11-AC13 and AC15 (rules R6-R8) that this order's acceptance criteria must trace back to exactly; the claim's test-isolation and visual-baseline requirements are this order's own additions and are not in this file, so no second source is needed for those.

## Token Budget

- **Status**: OK — ~13687 tokens, 11.4% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` | 11302 | 2825 |
| `apps/ye-olde-magic-shop/src/api/expeditions.ts` | 1605 | 401 |
| `apps/ye-olde-magic-shop/src/stores/authStore.ts` | 5898 | 1474 |
| `apps/ye-olde-magic-shop/tests/reagent-map-api.spec.ts` | 4226 | 1056 |
| `apps/ye-olde-magic-shop/tests/reagent-map-expeditions.spec.ts` | 5269 | 1317 |
| `apps/ye-olde-magic-shop/playwright.config.ts` | 4036 | 1009 |
| `apps/mock-api/server.js` | 11361 | 2840 |
| `apps/ye-olde-magic-shop/.env` | 164 | 41 |
| `apps/ye-olde-magic-shop/tests/reagent-map-visual.spec.ts` | 2930 | 732 |
| `docs/specs/reagent-map-spec/acceptance-criteria.md` | 7969 | 1992 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `No shared/knowledge, .claude/knowledge, or docs/adrs directories exist anywhere in this repository (verified by find) — the search-ki/Proactive RAG step had nothing to retrieve, not a pruning choice.`
- `docs/features/reagent-map-save/ — this is this same feature's own previous pipeline iteration (pre-dates the spec change that marked every stage STALE), not a different delivery in the bounded context, so it is excluded from 'Prior Deliveries' as a precedent and left for tech-writer/devops to reconcile, not the analyst.`
- `Prior deliveries in this bounded context (Expedition Planning): reagent-map-expeditions, reagent-map-globe, reagent-map-line-tool, reagent-map-visual-baselines, expeditions-api, reagent-sites-api — none of them has a retrospective.md (confirmed by find across docs/features/*), so there is no 'What Went Poorly'/'What To Improve' content to surface. Explicitly 'none found', not omitted.`
- `apps/ye-olde-magic-shop/tests/reagent-map-line-tool.spec.ts and reagent-map-resilience.spec.ts — distance-1 dependents of ReagentMapView.vue/expeditions.ts by Grep, but the draft-point/draft-line and GET-failure conventions they'd illustrate are already fully present in ReagentMapView.vue itself and in reagent-map-expeditions.spec.ts's getEntities()/metadata.kind pattern, which are already pinned.`
- `apps/ye-olde-magic-shop/tests/reagent-map-a11y-check.spec.ts — accessibility-engineer's own pass, not analyst-relevant for a save flow reusing existing control conventions.`
- `apps/ye-olde-magic-shop/src/router/index.ts and src/main.ts — only lazy-load ReagentMapView.vue behind the same flag App.vue already reads; flipping .env's default requires no change here.`
- `docs/specs/reagent-map-spec/domain-model.md and corrections.md — the gates they establish (1 km tolerance, geodesic lines, land/sea rule, sphere radius) are already restated verbatim in the feature claim's own Summary.`
- `Blast radius not graph-derived: ToolSearch found no mcp__loom__impact/dependents_of/find_files tools registered in this session, so dependents above were found by Grep over imports of ReagentMapView.vue and expeditions.ts instead of the repo graph; blastRadius is omitted rather than hand-filled.`
