# Documentation Report: golden-test-ids

## Status: PASS

## Summary

This is an internal test-hygiene chore: it replaces two fragile Saturday Site-Centric locators
(`CartPage.ts`'s `has-text()` guess and `HomePage.ts`'s nested `[data-testid="item-card"] a`
selector) with flat `data-testid` locators, and renames `CartView.vue`'s misnamed `checkout-btn`
test id to `complete-transaction-btn`. No user-facing behavior, UI, or API changed — only test
infrastructure did. No documentation describes the old locator strings or the `checkout-btn` test
id value, so nothing in the docs corpus went stale.

## Documentation Checklist

### Always Update
- **CHANGELOG.md** — Not applicable. No `CHANGELOG.md` exists at the repo root or under
  `apps/ye-olde-magic-shop/`. Nothing to add an entry to, and this chore has no user-facing change
  to log even if one existed.
- **README.md** — Not applicable. `apps/ye-olde-magic-shop/README.md` documents how to run the
  app and its test suites (`test:e2e`, `test:bdd`); those commands and their behavior are unchanged.
  No new capability for a user or developer to learn about.

### Update if Applicable
- **API documentation** — Not applicable. No endpoint changed.
- **Configuration docs** — Not applicable. No new env var or config option.
- **Architecture Decision Record** — Not applicable. Swapping a brittle locator for a
  `data-testid` locator follows the existing, already-documented Site-Centric convention
  (`BasePage`/`BaseElement`/`registerElement`); it applies a convention already in force rather
  than deciding a new one. No ADR-worthy decision was made.
- **Getting started / setup guides** — Not applicable. No new setup step.
- **Migration guides** — Not applicable. No breaking change, no database migration.

### Code-Level Docs
- No new source files. `CartPage.ts` and `HomePage.ts` keep their existing method and property
  names; only locator-string literals and one Vue template attribute value changed, so no
  docstring/JSDoc was stale or needed updating.
- `apps/ye-olde-magic-shop/tests/locator-conventions.spec.ts` (new, written by `qa-engineer`)
  already carries its own header comment with the `@issue golden-test-ids` annotation and an
  explanation of what the static scan checks and why it's a source scan rather than a browser
  test — this meets the test-annotation convention and needed no addition from this stage.

## Files Updated

None. No documentation file was created or modified in this pass — see the checklist above for
why each category was judged not applicable.

## Persistence Note

Per the Documentation Persistence Convention, the orchestrator will persist this run's artifacts
(`analysis.md`, `architecture-notes.md`, `implementation-notes.md`, `code-review-report.md`,
`security-report.md`, `qa-report.md`, this `docs-report.md`, and `devops-report.md` once written)
to `docs/features/golden-test-ids/` as the permanent, append-only record. `docs/features/README.md`
is the feature archive index; the orchestrator (or a human, at delivery) should add a one-line row
there for `golden-test-ids`, consistent with the existing entries for `leaflet-local-tiles` and
`backfill-workflow-tool-tests` — that index update is outside this stage's own documentation
checklist, since it is a delivery-record action, not a product-doc update.

## Notes for DevOps / Follow-up

- **QA flagged an unresolved scope question, not a documentation gap**: the new fitness-function
  test `apps/ye-olde-magic-shop/tests/locator-conventions.spec.ts` sits outside the claim's literal
  `scope.paths` (`apps/ye-olde-magic-shop/src/**`, `apps/ye-olde-magic-shop/lib/**`, `docs/**`) —
  `apps/ye-olde-magic-shop/tests/` is not under any of those three globs, and Playwright's
  `testDir` is fixed to `./tests`, so there is no in-scope path where the test would actually run.
  This needs a human or `devops-engineer` scope decision, not a documentation fix — noting it here
  so it isn't lost between stages.
- QA's suggestion to add "a one-line mention in any Saturday testing-conventions doc for this app"
  was evaluated and not acted on: no app-specific testing-conventions doc exists for
  `ye-olde-magic-shop`, and the only candidate, `.claude/rules/testing-conventions.md`, is outside
  this claim's `scope.paths` and outside this agent's writable tool set — editing it here would be
  an out-of-boundary write (Approval Gate #6). Flagging as a possible future `create-ki` candidate
  for a human to act on, not doing it unilaterally.

## Acceptance Criteria — Documentation Coverage

- [x] No AC in this spec requires a documentation artifact — all five ACs are locator-convention,
      regression-safety, and code-quality criteria verified by `qa-engineer`'s tests, not by docs.
- [x] Confirmed no doc anywhere in the repo references the old `checkout-btn` test id or the
      `has-text("COMPLETE TRANSACTION")` / `[data-testid="item-card"] a` selector strings, so no
      stale reference was left behind by the implementation change.