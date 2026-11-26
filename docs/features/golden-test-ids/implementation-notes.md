---
feature: "golden-test-ids"
bounded_context: ""
domain_terms: ["BasePage", "BaseElement"]
files_touched: ["apps/ye-olde-magic-shop/lib/site/CartPage.ts", "apps/ye-olde-magic-shop/src/views/CartView.vue", "apps/ye-olde-magic-shop/lib/site/HomePage.ts"]
issue_refs: ["golden-test-ids"]
linked_adrs: []
linked_kis: []
---

# Implementation Notes: golden-test-ids

## Files Created

None

## Files Modified

- apps/ye-olde-magic-shop/lib/site/CartPage.ts
- apps/ye-olde-magic-shop/src/views/CartView.vue
- apps/ye-olde-magic-shop/lib/site/HomePage.ts

## Interface Design

- No public interface changes — locator-string and data-testid value edits only; CartPage/HomePage public method signatures (completeTransaction(), clickFirstProductDetails(), etc.) are unchanged.

## Named Refactoring Log

- **Remove Dead Code** on `apps/ye-olde-magic-shop/lib/site/CartPage.ts:20-21 (pre-edit) — two-line comment above completeTransactionBtn explaining the has-text guess` — The comment documented a workaround for a missing test id; now that complete-transaction-btn exists, the comment describes a condition that no longer holds, so it was deleted rather than edited (per architecture-notes.md Refactoring Opportunities).

## Self-Review Checklist

- [x] Every public method has an intention-revealing name (no process/handle/manage) — No new methods added; existing names (completeTransaction, clickFirstProductDetails) untouched.
- [x] No function exceeds 30 LOC
- [x] Cyclomatic complexity < 7 on all new functions — No new functions; edits are locator-string swaps and a one-attribute rename inside existing functions with no new branching.
- [x] No primitive obsession: domain concepts wrapped in value objects, not raw strings/ints — Locator strings are the established Saturday convention for this layer (BaseElement selector argument), not a primitive-obsession smell.
- [x] No feature envy: methods use their own class's data more than another's
- [x] No magic numbers or strings: all literals are named constants — data-testid selector literals are the existing, established pattern throughout CartPage.ts/HomePage.ts (12 other call sites use the same literal-string style) — not a new smell introduced by this change.
- [x] Dependency direction verified: nothing in inner layers imports from outer layers — No new imports added in any of the three files.

## Simple Design Verification

- [x] Passes the tests — Could not run the full Playwright/Cucumber suites locally — this worktree has no installed dependencies (node_modules absent at repo root and app level), and a monorepo-wide install was out of scope for this stage. Per the spec's own Test Approach, QA runs scripts/ui-suite-in-chrome.sh and test:bdd against a running app; both existing suites already exercise these page objects (checkout.steps.ts calls cartPage.completeTransaction(), the home/inventory scenarios exercise clickFirstProductDetails's destination) with zero .feature/assertion changes made here.
- [x] Reveals intention — complete-transaction-btn names the button by its action, replacing both the guessed has-text selector and the misleadingly-named checkout-btn testid.
- [x] No duplication — Reused the existing item-name testid already rendered by ItemCard.vue instead of adding a third, redundant testid to the same anchor.
- [x] Fewest elements — Deleted the now-stale two-line comment explaining the has-text guess; no new elements, classes, or abstractions introduced.

## Key Decisions

- Renamed CartView.vue's existing checkout-btn testid to complete-transaction-btn in place rather than adding a second attribute.: Grep confirmed checkout-btn has exactly one occurrence repo-wide, so the rename is isolated and leaves no stale attribute that could silently reattract a future locator.
- Pointed HomePage.ts's clickFirstProductDetails() at the existing item-name testid instead of adding a new testid to the item-card anchor.: ItemCard.vue already renders data-testid="item-name" on a router-link to the same /item/${item.id} destination the old nested selector resolved to, so reuse is behavior-preserving and avoids a redundant attribute.
- Left HomePage.ts's .nth(1)/.nth(randomIndex) Playwright locator-chaining calls and ProductDetailsPage.ts's page.locator('h2') untouched.: Per analysis.md/architecture-notes.md, these are not has-text/:text/CSS-nth-/nested-selector-string violations under the AC's literal wording — .nth() is a Playwright API method on an already-flat locator, and 'h2' is a bare tag selector with no nesting. Fixing either would be scope creep beyond this spec.

## Deviations from Analysis

None

## Dependencies Added

None

## Notes for QA

- Could not run scripts/ui-suite-in-chrome.sh or test:bdd locally — this worktree has no installed node_modules (root or app-level). Please run both per the spec's Test Approach before sign-off.
- Confirm the existing Cucumber checkout scenario (via checkout.steps.ts -> cartPage.completeTransaction()) still passes with the renamed complete-transaction-btn testid.
- Confirm clickFirstProductDetails()'s new item-name locator is exercised by an existing home/inventory scenario and still navigates to /item/:id as before — flagged in analysis.md as needing verification or an explicit gap note.
- The proposed fitness-function regression test (scanning apps/ye-olde-magic-shop/lib/site/**/*.ts for has-text(, :text(, CSS :nth-, or nested selector strings) has not been added — per the pipeline, that is QA's task, not developer's.

## Notes for DevOps

None
