# Delivery Route: local-fonts

## Stages

| Stage | Runs | Why |
|---|---|---|
| `architect` | yes | structural work: the feature crosses a bounded context |
| `performance-engineer` | **no** | no performance requirement carries a threshold with a number and a unit |
| `data-engineer` | **no** | no data-model change to sequence |
| `acceptance-scenarios` | yes | always runs; not skippable by routing |
| `developer` | yes | always runs; not skippable by routing |
| `code-reviewer` | yes | always runs; not skippable by routing |
| `accessibility-engineer` | yes | the analysis declares a UI surface |
| `security-reviewer` | yes | always runs; not skippable by routing |
| `qa-engineer` | yes | always runs; not skippable by routing |
| `visual-qa-engineer` | yes | the analysis declares a UI surface |
| `sre-engineer` | **no** | the analysis declares no served runtime surface, so no availability or latency SLI applies |
| `tech-writer` | yes | always runs; not skippable by routing |
| `pre-ship-checks` | yes | always runs; not skippable by routing |
| `devops-engineer` | yes | the analysis lists DevOps tasks |

## Skipped

- `performance-engineer` — no performance requirement carries a threshold with a number and a unit
- `data-engineer` — no data-model change to sequence
- `sre-engineer` — the analysis declares no served runtime surface, so no availability or latency SLI applies

## How this was decided

Computed from `analysis.md` after the analyst stage, by predicates in `internal/state/` —
not by a model re-reading the analysis (roadmap L3.0).

Editing this file resets the design gate's approval: the run will halt until a human
approves the route as it then stands.
