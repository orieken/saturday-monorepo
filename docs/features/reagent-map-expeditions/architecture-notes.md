---
feature: "reagent-map-expeditions"
bounded_context: "Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7"
domain_terms: ["Expedition", "Adventurer", "Seeded Expedition"]
files_touched: ["apps/ye-olde-magic-shop/src/api/expeditions.ts", "apps/ye-olde-magic-shop/src/views/ReagentMapView.vue", "apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: reagent-map-expeditions

## Structural Decisions

### Add fetchExpeditions() in a new sibling module, apps/ye-olde-magic-shop/src/api/expeditions.ts — not inside ReagentMapView.vue and not appended to an existing client file.
- **Rationale**: Matches the established one-file-per-resource client pattern (reagentSites.ts, orders.ts, items.ts), each resolving API_BASE at runtime independently. Keeps fetch/timeout/header concerns out of the view, consistent with every other resource.
- **Trade-offs**:
  - Duplicates the runtime API_BASE/isLocal resolution block a fourth time (reagentSites.ts, orders.ts, items.ts, now expeditions.ts) — accepted because this order is additive, not a refactor; flagged below as a refactoring opportunity.
- **Fitness Function**: No fetch()/HTTP call exists outside apps/ye-olde-magic-shop/src/api/*.ts — verified by judgment-only — no eslint-plugin-boundaries config is wired for apps/ye-olde-magic-shop yet, though ARCHITECTURE_RULES.md section IV.2 names it as the intended mechanism; code-reviewer checks by grepping for fetch( outside src/api/ until that config exists.
### fetchExpeditions() attaches Authorization only when a token is present in localStorage('dnd_auth_token'), and omits Content-Type entirely (this is a bodyless GET) — a trimmed inline copy of orders.ts's getAuthHeaders(), not a shared helper.
- **Rationale**: GET /api/expeditions already tolerates an absent or invalid token per the NFR; sending Content-Type on a bodyless GET would be new, uncalled-for behavior no other client in this app exhibits.
- **Trade-offs**:
  - Introduces a second, slightly different header-building block instead of extracting a shared getAuthHeaders(token?, hasBody?) helper — accepted because that extraction is a cross-cutting refactor outside this order's additive scope; flagged below.
- **Fitness Function**: judgment-only — GET /api/expeditions already tolerates an absent or invalid token per the NFR; sending Content-Type on a bodyless GET would be new, uncalled-for behavior no other client in this app exhibits.
### Draw expeditions via a new addExpeditionEntities() helper, same shape and call site (onMounted, right after addSiteEntities) as the existing site-marker helper — no new base class, no typed entity wrapper.
- **Rationale**: Mirrors addSiteEntities() exactly: same parameter order (Cesium, viewer, ellipsoid, data[]), same ad-hoc `(entity as unknown as {...}).saturdayMetadata =` tagging cast already present in the file. Introducing a typed Entity-tagging wrapper now would be new pattern for a second call site, not justified by this order alone.
- **Trade-offs**:
  - The unknown-cast for saturdayMetadata is now duplicated across two helpers instead of extracted into one tagging utility — accepted, same reasoning as above.
- **Fitness Function**: judgment-only — Mirrors addSiteEntities() exactly: same parameter order (Cesium, viewer, ellipsoid, data[]), same ad-hoc `(entity as unknown as {...}).saturdayMetadata =` tagging cast already present in the file. Introducing a typed Entity-tagging wrapper now would be new pattern for a second call site, not justified by this order alone.
### Site-list items become <button type="button"> controls calling flyToSite(site); the Cesium module and ellipsoid returned by buildViewer() are promoted to component-scope variables (same lifetime as the existing `viewer` variable) so flyToSite can reach them. No composable, no Pinia store.
- **Rationale**: flyToSite is behavior local to this one view's own Cesium viewer instance; there is exactly one consumer. A composable or store would be premature generalization — nothing else in the app needs camera control.
- **Trade-offs**:
  - ReagentMapView.vue's component-scope state grows from one variable (viewer) to three (viewer, Cesium module, ellipsoid) — accepted; still one <script setup> module, no new file, no new public surface.
- **Fitness Function**: judgment-only — flyToSite is behavior local to this one view's own Cesium viewer instance; there is exactly one consumer. A composable or store would be premature generalization — nothing else in the app needs camera control.

## Component Placement

- `fetchExpeditions()` → Interface Adapters layer, `apps/ye-olde-magic-shop/src/api/expeditions.ts` — Converts the mock API's JSON response into the app's Expedition/GeoPoint shape — same role and layer as the existing fetchReagentSites() and fetchOrders().
- `addExpeditionEntities()` → Interface Adapters layer, `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Cesium entity construction is adapter-layer glue between fetched data and the Viewer, the same role the existing addSiteEntities() already plays in this file.
- `flyToSite()` → Interface Adapters layer, `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue` — Camera control is Viewer glue, not domain logic; stays beside buildViewer()/addSiteEntities() in the view rather than a new composable.

## Bounded Context

- **Owning Context**: Expedition Planning (Ye Olde Magic Shop demo app) — DOMAIN_DICTIONARY.md section 7
- **Context Crossings**: apps/mock-api GET /api/expeditions — consumed only; already shipped, no backend change in this order

## Stability Design

- **Explicit timeout (AbortController, 5s)** on fetchExpeditions() in apps/ye-olde-magic-shop/src/api/expeditions.ts — Mirrors fetchReagentSites()'s existing timeout, satisfying architecture-guardrails.md #5's explicit-timeout requirement. No retry loop is introduced anywhere — a rejected promise is handled by onMounted's existing Promise.allSettled failure path, not a custom retry/sleep loop.

## Observability Design

None

## Layer Boundary Checks

- [NOT_APPLICABLE] Domain/UseCases layers must not import Adapters or Frameworks (ARCHITECTURE_RULES.md section I) — This Vue demo app has no separate domain/use-case layer to protect — ReagentMapView.vue and api/*.ts are both adapter-level (UI + HTTP client). Confirmed by reading both files directly; no graph query was available/needed for a two-file app slice.
- [PASS] New apps/ye-olde-magic-shop/src/api/*.ts module has exactly one consumer (the view that needs it) — Read reagentSites.ts and orders.ts: each is a single-purpose client imported by exactly one view/store. expeditions.ts is expected to follow the same shape, imported only by ReagentMapView.vue.

## Anti-Pattern Check

- Premature generalization: not found — No new abstraction, composable, or store introduced; addExpeditionEntities() and flyToSite() copy the exact shape of the existing buildViewer()/addSiteEntities() functions rather than generalizing early.
- Shotgun surgery: not found — Additive only — one new file (expeditions.ts) plus one view file touched; reagentSites.ts, orders.ts, items.ts untouched.
- Leaky abstraction: not found — No new interface is exposed for other callers to depend on; expeditions.ts exports only fetchExpeditions() and its two plain data types, same surface shape as reagentSites.ts.
- God object: not found — ReagentMapView.vue grows by two helper functions, one ref, and three component-scope variables; it remains single-responsibility (render the Aurelis map and its controls), not accumulating unrelated concerns.

## Fitness Functions

- **Every new or changed function in this order keeps cyclomatic complexity below 7**
  - **Verification**: ESLint complexity rule capped at 6 (typescript-conventions.md / ARCHITECTURE_RULES.md section 140), run in CI lint
  - **Owner**: devops-engineer
- **No fetch()/HTTP call exists outside apps/ye-olde-magic-shop/src/api/*.ts**
  - **Verification**: judgment-only until eslint-plugin-boundaries is wired for this app; code-reviewer greps for fetch( outside src/api/
  - **Owner**: code-reviewer
- **The @visual globe baseline changes only for the reason this order states, and no other baseline changes**
  - **Verification**: scripts/ui-suite-in-chrome.sh --update-baselines and --ubi9 --update-baselines regenerate the PNGs; a human or agent views each diff before commit and states the reason in the delivery record
  - **Owner**: qa-engineer
- **No request from /reagent-map leaves the app's or mock API's own origin**
  - **Verification**: apps/ye-olde-magic-shop/tests/reagent-map-globe.spec.ts's existing off-origin NFR test continues to pass once fetchExpeditions() uses the same runtime API_BASE resolution as fetchReagentSites()
  - **Owner**: qa-engineer

## Refactoring Opportunities

- Extract Function on `The runtime API_BASE/isLocal resolution block, now duplicated across reagentSites.ts, orders.ts, items.ts, and expeditions.ts` — Four independent copies of the same 5-line block violate 'no duplication' (design-principles.md #1.3); a shared apiBase.ts would remove it. Out of scope for this additive order — flag for a future chore.
- Extract Function on `The optional-Authorization-header logic duplicated between orders.ts's getAuthHeaders() and expeditions.ts's trimmed inline equivalent` — A third API client needing an optional-Authorization pattern would make this duplication worth a shared getAuthHeaders(token?, hasBody?) helper; not required yet at two call sites.

## Developer Handoff Notes

- No new package, base class, or layer boundary is introduced by this order — follow the analyst's developerTasks list exactly; this architecture pass only confirms placement and names two deferred refactoring opportunities (API_BASE duplication, auth-header duplication) that are explicitly out of scope.
- expeditions.ts exports GeoPoint and Expedition types plus fetchExpeditions() only — keep the same minimal export surface as reagentSites.ts.
- Do not extract a shared apiBase.ts or getAuthHeaders() helper as part of this order even though the duplication is visible — that is a separate refactor, not this feature's scope.
- Keep the exact mapError text 'Failed to load reagent sites. Try refreshing the page.' on either sites or expeditions fetch failure — rewording breaks reagent-map-error.png's baseline.

## Open Architectural Questions

- Analyst's architecturalFlags already noted this order is additive enough that the architect stage was skippable (CLAUDE.md's 'Skip architect for purely additive features' clause); this pass concurs in hindsight — no blocking structural question was found. Noting it only so the pipeline can tune the architect-invocation heuristic for similar future orders.
