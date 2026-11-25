---
intake: 1
id: golden-lint-and-exemplars
type: chore
plan: deliver-feature
status: hold
priority: 2
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [golden-steps-through-site]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/**, packages/**, docs/**, GOLDEN.md]
---

# Lint the golden rules and name the exemplar files

## Request

> Add a lint so the next agent cannot regress the golden tests, and a GOLDEN.md naming them.

## Summary

The rules the previous orders established become a check that fails, and the repaired files become
the named set a factory pins as its pattern.

## Acceptance Criteria

- [ ] Given a lint run over the demo app (an ESLint rule set or a small checker), then it fails on:
      `this.page`, `getByTestId`, `locator(` or `$(` in `step_definitions/`; a text or structural
      locator in `lib/site/`; `waitForTimeout` or `setTimeout` in `step_definitions/` or `lib/site/`; a file in `lib/site/`
      not named `name.type.ts`; a page or flow not registered on the site.
- [ ] Given the Playwright specs in `tests/`, then the fixed-wait rule does not apply to them: the
      pan and zoom specs hold the mouse button for a set time on purpose, which is gesture timing,
      not a wait for state.
- [ ] Given each rule, then a test shows the lint failing on a fixture that breaks it and passing on
      one that does not.
- [ ] Given the app's `lint` script, then it runs the rule set and the current code passes.
- [ ] Given `GOLDEN.md` at the repository root, then it lists, per kind (site, page, element, flow,
      step file, feature file, Playwright spec), the exemplar file and one sentence on what it shows.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/golden-lint-and-exemplars/` excepted.

## Out of Scope

- Applying the rules to other apps in the monorepo.

## Domain Language

- **Exemplar**: a file a factory's context-engineer pins as the pattern to follow.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Unit tests for each lint rule.

Run the Playwright suite through `scripts/ui-suite-in-chrome.sh`, one runner image at a time, never
two at once, and wait for each run to finish before ending the stage. One run takes about two minutes
and a stage has 30. Run the Cucumber suite with the app's `test:bdd` script against the same running
app. Changing `.feature` files is out of scope for every golden order.

## Open Questions

- ESLint rule or standalone checker: Assumption if unresolved: ESLint with `no-restricted-syntax`
  and a small custom rule where that cannot express it.

## Infrastructure / Deploy Notes

After merge, a human adds `{name: lint, command: pnpm --filter @orieken/ye-olde-magic-shop lint}` to
`.claude/delivery-policy.yaml` `preShipChecks`.
