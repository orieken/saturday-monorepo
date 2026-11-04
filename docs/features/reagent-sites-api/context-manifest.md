# Context Manifest: reagent-sites-api

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/mock-api/server.js` — The only file this slice adds an endpoint to. Shows the existing readJSON/writeJSON helper pattern, the __dirname-relative data path constants, and where app.get handlers and app.listen's endpoint log go — the new GET /api/reagent-sites must follow this file's existing shape, not introduce a new one.
- `apps/mock-api/package.json` — Current scripts/dependencies (express, cors only, no test runner). The analyst must plan adding a `test` script using node:test per the spec's AC — this file is where that script is added and is the evidence no test runner currently exists.
- `docs/specs/reagent-map-spec/acceptance-criteria.json` — Authoritative seed.reagentSites array (8 sites, name/lat/lng, no id yet) that apps/mock-api/data/reagent-sites.json must hold exactly. Digest-pinned by the intake order; verified sha256 match before pinning.
- `docs/specs/reagent-map-spec/acceptance-criteria.md` — Human-readable AC1 (Given/When/Then) and the settled gate answers (1 km tolerance, geodesic line, offLand, visibility, samePlace) the analyst needs to scope this API-only slice correctly against the larger spec.
- `features/reagent-map/reagent-sites.feature` — The @level-api scenario ('The shop publishes every reagent site') this order's own API test must satisfy — distance-0, this is the order's own acceptance scenario, digest-pinned.
- `docs/specs/reagent-map-spec/corrections.md` — Records why every spec gate is already settled (Q3/Q4/Q5/Q6/Q11), supporting the intake order's 'Open Questions: None' claim — the analyst should not re-litigate these.
- `docs/specs/reagent-map-spec/package.sha256` — The digest list the intake order pins against. Re-verified this session (shasum -c): all 9 files OK, so the spec package has not changed since the order was written.
- `DOMAIN_DICTIONARY.md` — Section 7 only (lines 107-118), 'Product Domain: Expedition Planning' — defines Adventurer, Reagent Site, Aurelis exactly as this feature must name them. Full file is 122 lines, no range constraint needed, but only this section applies.

## Token Budget

- **Status**: OK — ~10344 tokens, 8.6% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/mock-api/server.js` | 7338 | 1834 |
| `apps/mock-api/package.json` | 221 | 55 |
| `docs/specs/reagent-map-spec/acceptance-criteria.json` | 8116 | 2029 |
| `docs/specs/reagent-map-spec/acceptance-criteria.md` | 7969 | 1992 |
| `features/reagent-map/reagent-sites.feature` | 1335 | 333 |
| `docs/specs/reagent-map-spec/corrections.md` | 989 | 247 |
| `docs/specs/reagent-map-spec/package.sha256` | 1001 | 250 |
| `DOMAIN_DICTIONARY.md` | 14417 | 3604 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `apps/mock-api/data/items.json, orders.json, users.json — different bounded context (shop items/auth/orders), excluded by bounded-context rule; items.json's data-file shape is already visible via server.js's readJSON/readItems pattern`
- `docs/specs/reagent-map-spec/domain-model.json, domain-model.md, discovery.json, discovery.md, example-map.json, example-map.md — cover the full reagent-map feature (UI, expeditions, drawing); this order is API-only (8 reagent sites), out of scope per the intake order's own 'Out of Scope' section`
- `features/reagent-map/draw-line.feature, expeditions-visibility.feature, fly-to-site.feature, save-expedition.feature — scenarios for the other five slices of this spec (expeditions-api, reagent-map-globe, reagent-map-save, etc.), not this order's scope`
- `docs/specs/reagent-map-spec/orders/ — the other five intake orders split from the same spec package; this stage only needs reagent-sites-api's own order`
- `DOMAIN_DICTIONARY.md sections 1-6 (framework terms, core domains, entities, value objects, events, operations) — these describe the ai-assistant-dot-files framework itself, not the Expedition Planning product domain this feature touches`
