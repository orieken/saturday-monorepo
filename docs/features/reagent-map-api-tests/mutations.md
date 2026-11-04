# F1 judge: reagent-map-api-tests

Order `docs/specs/reagent-map-spec/orders/reagent-map-api-tests.md`, run by hand. Spec package
`sha256:18f9820d26a2013205ab0bec91e32bc6fe3a49166f2eb079342ee71b32604413`.

Suite: `apps/ye-olde-magic-shop/tests/reagent-map-api.spec.ts`, 6 scenarios, all passing in
`saturday-runner:debian` under the `chrome` project. No browser page is opened.

## Faults injected into `apps/mock-api`

Each fault was applied alone, the suite run, and the fault reverted.

| Fault | Scenario that failed |
|---|---|
| A published site is missing | AC1 |
| A site is published at the wrong place | AC1 |
| A seeded expedition is missing | AC3 |
| Signed out sees every expedition | AC3 |
| Signed in sees every adventurer's expeditions | AC4 |
| Own expeditions are not returned | AC4 |
| Saved for the shop instead of the adventurer | AC11 |
| An endpoint moved about 5 km when stored | AC11 |
| Save needs no sign-in | AC11, and AC12 when run alone (expected 401, received 200) |
| Zero-length lines accepted | AC13 |
| Zero-length check uses twice the planet's radius | AC13 |

None survived. The suite runs in order and stops at its first failure, so the "no sign-in" fault
stopped at AC11 in the full run. AC12 was then run on its own with `--grep @AC12` to show it fails.

These faults were applied by a throwaway script, not committed. Sunday's mutation tooling (prompt
04) does not exist for this repository yet.
