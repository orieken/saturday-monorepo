# Context Manifest: expeditions-api

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `.claude/factory/claims/expeditions-api.md` — The feature spec itself — primary task input naming the two endpoints, the five acceptance criteria, the test-isolation requirement, and the scope.paths boundary (apps/mock-api/**, docs/features/expeditions-api/**).
- `apps/mock-api/server.js` — The file both new endpoints land in. Established patterns to follow exactly: readJSON/writeJSON helpers, requireAuth middleware reading req.userId from the session map, the unauthenticated GET /api/reagent-sites block to mirror for GET /api/expeditions, and the authenticated POST /api/orders block to mirror for POST /api/expeditions. Distance-1 by direct Grep — loom's repo-graph MCP tools (mcp__loom__impact etc.) are not registered in this session, so blast radius was not graph-derived; this is the only file any new route touches.
- `apps/mock-api/server.test.js` — The node:test pattern this feature's tests must follow exactly: startServer/stopServer helpers, the t.after() snapshot-and-restore trick for data-file isolation (test-repair-contract's data-is-unchanged requirement), and the @issue/@ac annotation convention already used for the sibling reagent-sites-api tests.
- `apps/mock-api/data/users.json` — Current adventurer records. Elara and Borin (named in the spec and in acceptance-criteria.json's seed.adventurers) do not exist here yet and must be added or fixtured, matching the existing user shape (id, email, password, name, createdAt, address).
- `apps/mock-api/data/reagent-sites.json` — Seeded expeditions reference site names (Moonwell Spring, Ember Crag, etc.) by from/to — this file gives the site ids and lat/lng needed to resolve those references into concrete expedition endpoints.
- `docs/specs/reagent-map-spec/acceptance-criteria.md` — Authoritative text for AC3, AC4, AC11, AC12, AC13 (the five this slice must satisfy) plus the settled seed table (3 seeded expeditions, Elara's and Borin's saved expeditions) and the settled tolerance/samePlace rules the 1 km zero-length check depends on.
- `docs/specs/reagent-map-spec/acceptance-criteria.json` — Machine-readable seed.seededExpeditions and seed.adventurers the spec text references directly by name — server.test.js's existing pattern reads this file at test time (SEED_PATH) rather than hand-copying seed data, and this feature's tests will need the same.
- `docs/specs/reagent-map-spec/domain-model.md` — Defines the Expedition aggregate's exact fields (id, owner, shape as a two-endpoint line, savedAt) and its two invariants (both endpoints on the surface; not the same place) that POST /api/expeditions must enforce.
- `docs/features/reagent-sites-api/architecture-notes.md` — Prior delivery in the same Owning Context ('Expedition Planning — Reagent Site catalog'), found by grepping analysis.md's Owning Context field across docs/features/*/analysis.md. Its structural decision directly carries over: stay inside server.js's flat style, no new abstraction/package/repository interface (premature-generalization check already made and reasoned through), server.js already exports app and guards app.listen behind require.main === module so node:test can import it. No retrospective.md exists for reagent-sites-api — checked, none found, so no 'What Went Poorly' lessons to carry forward beyond this architecture note.
- `DOMAIN_DICTIONARY.md` — Section 7 (Product Domain: Expedition Planning) is the Ubiquitous Language source for Adventurer, Reagent Site, Expedition, Seeded Expedition, and Aurelis — the analyst's task breakdown and the developer's identifiers must match these terms exactly (design-principles.md rule 6), not synonyms like 'route' or 'trip'.

## Token Budget

- **Status**: OK — ~15088 tokens, 12.6% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `.claude/factory/claims/expeditions-api.md` | 4362 | 1090 |
| `apps/mock-api/server.js` | 7756 | 1939 |
| `apps/mock-api/server.test.js` | 3678 | 919 |
| `apps/mock-api/data/users.json` | 1037 | 259 |
| `apps/mock-api/data/reagent-sites.json` | 781 | 195 |
| `docs/specs/reagent-map-spec/acceptance-criteria.md` | 7969 | 1992 |
| `docs/specs/reagent-map-spec/acceptance-criteria.json` | 8116 | 2029 |
| `docs/specs/reagent-map-spec/domain-model.md` | 3700 | 925 |
| `docs/features/reagent-sites-api/architecture-notes.md` | 8546 | 2136 |
| `DOMAIN_DICTIONARY.md` | 14417 | 3604 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `features/reagent-map/*.feature (draw-line, expeditions-visibility, fly-to-site, reagent-sites, save-expedition) — explicitly Out of Scope per the spec: the E1/F1 scenario suites are run by hand, not by this slice`
- `docs/specs/reagent-map-spec/example-map.md and example-map.json — the rules/examples they contain for this slice's ACs are already distilled into acceptance-criteria.md/.json, which are pinned instead`
- `docs/specs/reagent-map-spec/discovery.md and discovery.json, README.md, corrections.md — earlier-stage spec artifacts superseded by the settled acceptance-criteria.md for this order's purposes`
- `docs/specs/reagent-map-spec/orders/reagent-map-api-tests.md and orders/reagent-map-ui-tests.md — draft order for the separate, not-yet-started deliver-tests-sunday slice (depends on this one being done first); out of scope for this analyst pass`
- `apps/mock-api/data/items.json and orders.json — different bounded context slice (shop catalog/orders), untouched by this feature`
- `apps/mock-api/Dockerfile, .dockerignore, public/images/** — infra/static assets with no API-logic surface for this change`
- `apps/mock-api/package.json — trivial 14-line file; already confirmed directly that the test script is `node --test`, matching the spec's node:test requirement, so pinning it adds no new information`
- `docs/features/reagent-sites-api/implementation-notes.md, code-review-report.md, security-report.md, qa-report.md, tech-writer.md, sre-engineer.md, pre-ship-checks.md, route.md — sibling slice's downstream artifacts; architecture-notes.md alone carries the reusable structural precedent, the rest is execution detail for a different endpoint`
