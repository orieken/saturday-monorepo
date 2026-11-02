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
| `confirm-domain` | domain-model | analyst | **waiting** |
| `answer-questions` | example-map | analyst | not started |
| `confirm-scenarios` | scenarios | tester | not started |
| `confirm-ship` | package | requester | not started |

In this pilot one person holds all three roles. Each gate's answer, and every correction it makes
to the draft, is recorded in `corrections.md` (the factory's quality metric).

## Stage outputs

- `discovery.json` / `discovery.md`: goal ladder, outcomes, actors, impact map, questions
- `domain-model.json` / `domain-model.md`: event storm draft with hotspots
