---
intake: 1
id: golden-site-file-names
type: chore
plan: deliver-feature
status: draft
priority: 2
created: 2026-10-29
requested-by: oscar
source: manual
trust: internal
base-branch: main
finish: pr
autonomy: inherit
depends-on: [golden-workspace-packages]
budget:
  max-attempts: 2
  max-usd: 20
  max-hours: 2
scope:
  paths: [apps/ye-olde-magic-shop/lib/**, apps/ye-olde-magic-shop/features/**, apps/ye-olde-magic-shop/tests/**, docs/**]
---

# Name the demo app's site files name.type.ts

## Request

> The site model's files are PascalCase (`CartPage.ts`, `MagicShop.ts`); the rule is `cart.page.ts`.

## Summary

A pure rename, so the files a factory copies follow the naming rule.

## Acceptance Criteria

- [ ] Given `apps/ye-olde-magic-shop/lib/site/`, then every file is `name.type.ts` in kebab-case:
      pages as `*.page.ts`, the site as `magic-shop.site.ts`, elements as `*.element.ts`.
- [ ] Given the change is complete, then class names are unchanged and every import is updated.
- [ ] Given the existing suites, then every scenario that passed before passes after, with no
      `.feature` file or assertion changed.
- [ ] Given the change is complete, then every new or changed function has cyclomatic complexity
      below 7, every changed line is covered by a test that can fail, and no file outside
      `scope.paths` changes, the run's own delivery record under `docs/features/golden-site-file-names/` excepted.

## Out of Scope

- Any behaviour change.

## Domain Language

- No new terms introduced.

## Non-Functional Requirements

- None

## Trust Boundaries

- None — no auth, no user input, no external calls

## Test Approach

Saturday/Cucumber and Saturday/Playwright, existing suites.

## Open Questions

- None

## Infrastructure / Deploy Notes

- None
