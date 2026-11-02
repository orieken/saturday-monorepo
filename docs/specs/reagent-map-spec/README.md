# Spec package: reagent-map-spec (pilot wave 1, order A1)

Run by hand, stage by stage, per ai-assistant-dot-files `docs/prompts/factories/02-spec-factory.md`
and `05-pilot-cesium-reagent-map.md`. The `deliver-spec` plan does not exist yet, so this order is
not in `intake/` (`loom intake check` rejects a plan the project cannot run); it lives here.

| Field | Value |
|---|---|
| id | `reagent-map-spec` |
| plan | `deliver-spec` (by hand) |
| requested-by | oscar |
| target | `saturday-monorepo` (`apps/ye-olde-magic-shop`, `apps/mock-api`) |
| publishes | spec package: `features/reagent-map/*.feature`, this directory, draft orders D1, E1, F1 |

## Gates and owners

| Gate | Stage | Owner | Status |
|---|---|---|---|
| `confirm-problem` | discovery | requester (Oscar) | accepted 2026-11-02, no corrections |
| `confirm-domain` | domain-model | analyst (Oscar) | accepted 2026-11-02, no corrections |
| `answer-questions` | example-map | analyst (Oscar) | answered 2026-11-02: every assumption accepted |
| `confirm-scenarios` | scenarios | tester (Oscar) | accepted 2026-11-02, no corrections |
| `confirm-ship` | package | requester | **waiting** (this PR) |

In this pilot one person holds all three roles. Each gate's answer, and every correction it makes
to the draft, is recorded in `corrections.md` (the factory's quality metric).

## Stage outputs

- `discovery.json` / `discovery.md`: goal ladder, outcomes, actors, impact map, questions
- `domain-model.json` / `domain-model.md`: event storm draft with hotspots
- `example-map.json` / `example-map.md`: rules, examples (each rule has a negative), questions
- `acceptance-criteria.json` / `acceptance-criteria.md`: 15 criteria, seed data, traceability, critic results
- `../../../features/reagent-map/*.feature`: 21 scenarios (6 API, 15 UI), tagged `@AC-n`, `@level-*`, `@runner-playwright`

## The package

- Digest: `sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413`, over the files in
  `package.sha256` (the four stage JSON files and the five `.feature` files).
- Draft orders, each pinning that digest:
  - D1 `intake/reagent-map-view.md` (`deliver-feature`, status draft; promote to `ready` after C1 ships)
  - E1 `orders/reagent-map-ui-tests.md` (by hand)
  - F1 `orders/reagent-map-api-tests.md` (by hand)
- Check the pin: `shasum -a 256 docs/specs/reagent-map-spec/package.sha256` and
  `shasum -a 256 -c docs/specs/reagent-map-spec/package.sha256`.
