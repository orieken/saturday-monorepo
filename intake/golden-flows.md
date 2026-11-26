---
intake: 1
id: golden-flows
type: chore
plan: deliver-feature
status: ready
priority: 2
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [golden-test-ids]
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/lib/**, apps/ye-olde-magic-shop/tests/**, docs/**]
---

# Register the demo app's login and checkout flows

## Request

> `MagicShop.initializeFlows()` registers nothing, so the Flow layer is never shown.

## Summary

Add the two flows the existing scenarios repeat — logging in and checking out — so the exemplar
shows Site → Flow → Page → Element.

## Acceptance Criteria

- [ ] Given `lib/site/`, then `login.flow.ts` and `checkout.flow.ts` exist, each extends `BaseFlow`,
      takes its inputs as typed parameters, and uses only pages from the site.
- [ ] Given `lib/site/flows/draw-line.flow.ts`, then the new flows follow its shape and sit beside it.
      It is the reagent map's flow from E1 and is already registered on the site.
- [ ] Given `MagicShop`, then both flows are registered in `initializeFlows()` and exposed through
      typed getters, like the pages.
- [ ] Given each flow, then a unit test drives it against stubbed pages and asserts the calls it
      makes, and fails when a call is removed.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/golden-flows/` excepted.

## Out of Scope

- Changing step definitions to use the flows; that is `golden-steps-through-site`.

## Domain Language

- **Flow**: a named user journey across pages, such as logging in.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — test credentials only, from the existing fixtures

## Base

Start from `main` with #35 and #36 merged: the site files are `name.type.ts` (`magic-shop.site.ts`,
`login.page.ts`, `cart.page.ts`, `home.page.ts`), and the pages use `data-testid` locators. An
earlier run of this order started before those merged and designed against the old names.

## Test Approach

Unit tests for the flows; existing suites unchanged.

Run the Playwright suite through `scripts/ui-suite-in-chrome.sh`, one runner image at a time, never
two at once, and wait for each run to finish before ending the stage. One run takes about two minutes
and a stage has 30. Run the Cucumber suite with the app's `test:bdd` script against the same running
app. Changing `.feature` files is out of scope for every golden order.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
