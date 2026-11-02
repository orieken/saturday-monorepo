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
| `confirm-scenarios` | scenarios | tester | **waiting** |
| `confirm-ship` | package | requester | not started |

In this pilot one person holds all three roles. Each gate's answer, and every correction it makes
to the draft, is recorded in `corrections.md` (the factory's quality metric).

## Stage outputs

- `discovery.json` / `discovery.md`: goal ladder, outcomes, actors, impact map, questions
- `domain-model.json` / `domain-model.md`: event storm draft with hotspots
- `example-map.json` / `example-map.md`: rules, examples (each rule has a negative), questions
- `acceptance-criteria.json` / `acceptance-criteria.md`: 15 criteria, seed data, traceability, critic results
- `../../../features/reagent-map/*.feature`: 21 scenarios (6 API, 15 UI), tagged `@AC-n`, `@level-*`, `@runner-playwright`
