---
feature: "golden-test-ids"
bounded_context: "Saturday Site-Centric test layer (BasePage/BaseElement locators) for ye-olde-magic-shop"
domain_terms: ["BasePage", "BaseElement"]
files_touched: ["apps/ye-olde-magic-shop/lib/site/CartPage.ts", "apps/ye-olde-magic-shop/src/views/CartView.vue", "apps/ye-olde-magic-shop/lib/site/HomePage.ts"]
issue_refs: ["golden-test-ids"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: golden-test-ids

## Summary

CartPage.ts registers its complete-transaction button with a text-matching selector (button:has-text("COMPLETE TRANSACTION")) because no test id was ever added for it, and a second locator in HomePage.ts builds a CSS selector by nesting a bare tag under an attribute selector ([data-testid="item-card"] a). Both are fragile: a copy or markup change can silently break the test suite even though the feature itself is unchanged. This chore replaces both with flat data-testid locators — adding the missing complete-transaction-btn test id to CartView.vue and renaming its current, wrongly-named checkout-btn test id, and pointing HomePage.ts at the item-name test id ItemCard.vue already renders — with zero visual or behavioral change to the app and zero .feature/assertion edits.

### Acceptance Criteria

- Every registerElement call in lib/site/ locates its element with a [data-testid="..."] locator or a role-based locator, and the Vue component renders that test id.
  - Example: CartPage.ts's completeTransactionBtn becomes [data-testid="complete-transaction-btn"]; CartView.vue's button gets that same test id (replacing its current, incorrect checkout-btn test id).
- CartPage's complete-transaction button is located by data-testid="complete-transaction-btn".
- No locator in lib/site/ uses has-text, :text, a CSS nth- pseudo-selector, or a selector string built from element nesting.
  - Example: HomePage.ts's clickFirstProductDetails() currently uses [data-testid="item-card"] a, a single CSS string nesting a bare a tag under an attribute selector. It must become a flat testid locator ([data-testid="item-name"]), which ItemCard.vue already renders.
- Every scenario that passed before this change passes after it, with no .feature file or assertion changed.
- Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail, and no file outside scope.paths changes except this run's own delivery record under docs/features/golden-test-ids/.

### Non-Functional Requirements

None

## Proposed Fitness Functions

- **No locator string under apps/ye-olde-magic-shop/lib/site/ matches has-text(, :text(, a CSS :nth- pseudo-class, or a descendant-combinator (nested) selector.**
  - **Verification**: A qa-engineer-authored regression test that scans lib/site/**/*.ts source text for these patterns and fails the suite if any match (plain string/regex scan, not a CI pipeline change).
  - **Owner**: qa-engineer

## Out of Scope

- Editing any .feature file or step definition file.
- Any visual or behavioral change to the running app.
- The golden-site-file-names rename (CartPage.ts -> cart.page.ts, etc.) — a separate, not-yet-merged dependency feature; this run acts on the current file names.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Saturday Site-Centric test layer (BasePage/BaseElement locators) for ye-olde-magic-shop
- **Context Crossings**: ye-olde-magic-shop frontend components (Vue data-testid attributes) — the locator fix requires adding/renaming a test id in CartView.vue

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/lib/site/CartPage.ts` — Line 22's registerElement('completeTransactionBtn', ...) uses button:has-text("COMPLETE TRANSACTION"); the only has-text/:text violation in lib/site/. Must become [data-testid="complete-transaction-btn"] per AC1/AC2.
- `apps/ye-olde-magic-shop/src/views/CartView.vue` — Line 61 renders the complete-transaction button with the wrong test id (data-testid="checkout-btn"); grep confirms this is the only checkout-btn occurrence repo-wide, so renaming it to complete-transaction-btn is isolated and required by AC2.
- `apps/ye-olde-magic-shop/lib/site/HomePage.ts` — Line 39's clickFirstProductDetails() uses [data-testid="item-card"] a, a single CSS selector string nesting a bare tag under an attribute selector — a literal AC3 violation ('selector built from element nesting') found by direct grep beyond the context-manifest's single named case. Fix: swap to the existing [data-testid="item-name"] test id, which ItemCard.vue already renders on an anchor pointing at the same /item/:id route, so the fix is behavior-preserving.

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. In CartPage.ts, change registerElement('completeTransactionBtn', ButtonElement, 'button:has-text("COMPLETE TRANSACTION")') to use '[data-testid="complete-transaction-btn"]'; remove the now-obsolete inline comment about guessing the selector.
2. In CartView.vue (~line 61), rename the complete-transaction button's data-testid from 'checkout-btn' to 'complete-transaction-btn'.
3. In HomePage.ts, replace the nested selector string '[data-testid="item-card"] a' in clickFirstProductDetails() with '[data-testid="item-name"]' (already rendered by ItemCard.vue on the item-name link, which navigates to the same destination as the image link currently being matched).
4. Leave HomePage.ts's .nth(1)/.nth(randomIndex) Playwright locator-chaining calls and ProductDetailsPage.ts's page.locator('h2') untouched — neither is a has-text/:text/CSS-nth-/nested-selector-string match under AC3's literal wording, and fixing them is not requested by this spec.
5. After the edits, grep apps/ye-olde-magic-shop/lib/site/ for has-text(, :text(, a CSS :nth- pseudo-class, and any quoted selector string containing an internal space (nesting), and confirm zero matches.

### QA Tasks

1. Run the Playwright suite via scripts/ui-suite-in-chrome.sh, one runner image at a time, waiting for each ~2-minute run to finish before ending the stage (30-minute budget).
2. Run the Cucumber suite with the app's test:bdd script against the same running app.
3. Confirm no .feature file and no assertion changed, and that every scenario that passed before still passes.
4. Add the regression test proposed in proposedFitnessFunctions: scan lib/site/**/*.ts for has-text(, :text(, CSS :nth-, or nested selector strings, and fail if any are found.
5. Verify every changed line (CartPage.ts, CartView.vue, HomePage.ts) is covered by an existing test that would fail if the fix were reverted — the existing Cucumber checkout scenario and cart.spec.ts cover CartPage/CartView; confirm HomePage.ts's clickFirstProductDetails() path is exercised by an existing scenario, or note the gap.

### Tech Writer Tasks

None

### DevOps Tasks

None

## Edge Cases and Risks

- HomePage.ts's .nth(1) and .nth(randomIndex) are Playwright locator-chaining method calls, not CSS :nth-child/:nth-of-type pseudo-selector strings.: Leave untouched. AC3 targets selector strings (the literal substring 'nth-' as it would appear inside a CSS pseudo-class), not Playwright API methods operating on an already-flat data-testid locator.
- ProductDetailsPage.ts's page.locator('h2') is a bare tag selector outside any registerElement call.: Leave untouched. It is not has-text/:text/nth- and is not built from nesting (no descendant combinator) — not a violation under AC1 (scoped to registerElement) or AC3's literal wording. Flagged so it isn't 'fixed' as scope creep.
- HomePage.ts's clickFirstProductDetails() uses [data-testid="item-card"] a, a single CSS string combining an attribute selector with a nested bare tag via a descendant combinator.: This is a genuine AC3 violation even though it sits outside any registerElement call — AC3 says 'no locator in lib/site/', not 'no registered locator'. Fix by pointing at the existing [data-testid="item-name"] test id; confirmed behavior-preserving because both anchors in ItemCard.vue (the image link and the item-name link) navigate to the same /item/${item.id} route.
- The golden-site-file-names dependency (CartPage.ts -> cart.page.ts rename) has not merged into this branch.: Act on current file names now. If that rename lands first and this run is re-based, the same registerElement call is simply under a renamed file — the selector-string fix itself is unaffected.
- CartView.vue's button already carries a data-testid, just the wrong literal value.: Rename in place rather than adding a second attribute. Grep confirms checkout-btn has exactly one occurrence repo-wide, so the rename is isolated and safe.

## Definition of Done

- CartPage.ts's completeTransactionBtn registerElement call uses [data-testid="complete-transaction-btn"].
- CartView.vue's complete-transaction button carries data-testid="complete-transaction-btn" (renamed from checkout-btn).
- HomePage.ts's clickFirstProductDetails() uses a flat data-testid locator, not a nested selector string.
- Grep across apps/ye-olde-magic-shop/lib/site/ confirms zero remaining has-text, :text, CSS nth- pseudo-selector, or nested-selector-string matches.
- No .feature file changed; no step definition changed; no assertion changed.
- Full Playwright suite (scripts/ui-suite-in-chrome.sh) passes, one runner image at a time.
- Full Cucumber suite (test:bdd) passes against the same running app.
- Every new or changed function has cyclomatic complexity below 7.
- Every changed line is covered by a test that can fail.
- No file outside scope.paths (apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/lib/**, docs/**) changes, except this run's own delivery record under docs/features/golden-test-ids/.
