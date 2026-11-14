---
feature: "reagent-map-expeditions"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: ["apps/ye-olde-magic-shop/src/views/ReagentMapView.vue"]
issue_refs: ["reagent-map-expeditions", "reagent-map-globe"]
linked_adrs: []
linked_kis: []
---

# Code Review: reagent-map-expeditions

## Overall Status

**APPROVED**

## Design Narrative

fetchExpeditions() in a new sibling api/expeditions.ts client mirrors fetchReagentSites()'s runtime API_BASE resolution and timeout pattern exactly, and ReagentMapView.vue draws the returned expeditions as geodesic polyline entities via a new addExpeditionEntities() helper that is a structural twin of the existing addSiteEntities(). Site-list items became accessible <button> controls calling a new flyToSite(), which required promoting the Cesium module and ellipsoid to component-scope variables alongside the pre-existing viewer variable so the click handler can reach them after onMounted's scope has exited. The onMounted failure path was extended from a 2-promise to a 3-promise Promise.allSettled, destroying the viewer and reusing the exact unchanged error string on either a sites or an expeditions rejection, keeping the error baseline stable. Behaviour change: yes.

## Design Score

- **Clarity** 4/5
- **Cohesion** 4/5
- **Coupling** 4/5
- **Craft** 4/5

## Security Surface

- New GET fetch to /api/expeditions, same-origin-resolved via the same runtime API_BASE logic as the existing reagentSites client
- Optional Authorization: Bearer <token> header read from localStorage('dnd_auth_token'), attached only when present, never required — mirrors orders.ts's pattern on a bodyless GET
- Site/expedition names from the API render into DOM text and an aria-label via Vue's template interpolation, which auto-escapes — no innerHTML or v-html introduced

## Performance Surface

- fetchExpeditions() runs in parallel with buildViewer() and fetchReagentSites() inside the existing Promise.allSettled — no added sequential round-trip
- addExpeditionEntities() is a single flat forEach over the fetched array with no per-item network or lookup call — no N+1 introduced
- fetchExpeditions() carries the same explicit 5s AbortController timeout as fetchReagentSites(), satisfying architecture-guardrails.md #5

## Test Design Review

- reagent-map-globe.spec.ts (2 assertions repaired): YES — filtering to kind==='site' before the existing toHaveLength(SEEDED_SITES.length) and per-entity checks still fails if site-tagging regresses or count drifts; the filter is a selector-path repair (test-repair-contract.md), not a weakened assertion, since addExpeditionEntities is what newly populates the shared entity list this filter now excludes.
- reagent-map-pan-zoom.spec.ts (1 assertion repaired): YES — same reasoning as the globe.spec.ts repair, same filter added before the same unmodified length/content assertions.
- reagent-map-resilience.spec.ts new test 'destroys the viewer ... when fetchExpeditions() rejects': YES — traced the failure mode: without the expeditionsResult-rejected branch in onMounted's guard, the viewer would never be destroyed and its canvas would remain, which assertNoViewerArtifactsLeft's canvasCount-equals-0 assertion would catch.

## Verification of Developer Self-Review

- Developer marked 'No duplication' unchecked and disclosed the API_BASE/isLocal block now duplicated a 4th time plus a trimmed copy of orders.ts's auth-header logic — matches architecture-notes.md's explicitly deferred refactoring opportunities; this is an honest disclosure, not a discrepancy.
- Complexity claims (flyToSite: 2, addExpeditionEntities: 2, extended onMounted: ~4) match a visual read of the diff — no function here approaches the <7 budget.
- All other checked boxes (intention-revealing names, <30 LOC, no primitive obsession, no feature envy, no magic numbers, dependency direction) hold up against the actual diff.

## Feedback for the Developer

### 1. Rename Variable (non-blocking)
- **File**: `apps/ye-olde-magic-shop/src/views/ReagentMapView.vue`
- **Smell**: addExpeditionEntities()'s parameters are named cesiumEllipsoid and items, while its structural twin addSiteEntities() right above it uses ellipsoid and reagentSites for the same roles — the naming drifts between two functions meant to read as the same pattern.
- **Instruction**: Rename to ellipsoid and expeditions to match addSiteEntities()'s parameter names exactly, so the two helpers read as the same pattern at a glance.
