# Draft order F1: reagent-map-api-tests (by hand)

| Field | Value |
|---|---|
| plan | `deliver-tests-sunday` (by hand; the plan does not exist yet, so not in `intake/`, D-27) |
| status | draft |
| depends on | `reagent-map-view` (D1) |
| spec package | `sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413` (`package.sha256`) |
| gates | testability → developer; model, ship → tester |

## Scope

The 6 `@level-api` scenarios in `features/reagent-map/`, as Playwright Test API specs against
`/api/reagent-sites` and `/api/expeditions`, plus the API client E1's cross-check (AC11) uses.

## Rules

- Visibility (AC3, AC4) is the trust boundary: assert both that the caller's own expeditions are
  returned and that another adventurer's are not.
- Distances use a sphere of radius 6,000 km (Aurelis), never Earth's.

## Judge

Sunday's mutations (prompt 04) on both routes: every scenario fails on at least one, including an
API that returns every adventurer's expeditions and one that accepts a zero-length line.
