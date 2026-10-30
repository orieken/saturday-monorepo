---
intake: 1
id: golden-workspace-packages
type: chore
plan: deliver-feature
status: draft
priority: 1
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: []
budget:
  max-attempts: 2
  max-usd: 25
  max-hours: 3
scope:
  paths: [apps/ye-olde-magic-shop/**, pnpm-lock.yaml, docs/**]
---

# Move the demo app onto the workspace framework packages

## Request

> Agents left the demo app depending on vendored copies of the framework. Fix the golden tests
> before any factory copies them.

## Summary

The demo app depends on `file:./saturday-core` and `file:./saturday-cucumber`, copies inside the
app that have drifted from `packages/`. It should depend on the workspace packages, so its tests
exercise the framework as it is.

## Acceptance Criteria

- [ ] Given `apps/ye-olde-magic-shop/package.json`, then `@orieken/saturday-core`,
      `@orieken/saturday-cucumber` and `@orieken/saturday-playwright-heatmap` are `workspace:*`
      dependencies, and no dependency points at a path inside the app.
- [ ] Given the change is complete, then the vendored directories `saturday-core/`,
      `saturday-cucumber/` and `saturday-playwright-heatmap/` inside the app are gone. Any test in
      them that covers framework code and is not already in its package under `packages/` is moved
      there, not deleted.
- [ ] Given the app's existing feature files and Playwright specs, then every scenario that passed
      before the change passes after it, with no `.feature` or spec assertion changed. Record the
      before and after results in the delivery record.
- [ ] Given `apps/ye-olde-magic-shop/Dockerfile.runner`, which copies the vendored directories,
      then it builds from the workspace instead, or is removed with the reason recorded.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/golden-workspace-packages/` excepted.

## Out of Scope

- Renaming files, adding flows, changing step definitions. Those are later orders.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Cucumber and Saturday/Playwright: the existing suites are the characterization tests.

## Open Questions

- A vendored file has behaviour the package lacks: Assumption if unresolved: move the behaviour into
  the package with a test, in this order, and name it in the delivery record.

## Infrastructure / Deploy Notes

- None
