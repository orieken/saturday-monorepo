---
feature: "golden-test-ids"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: golden-test-ids

## Every registered page element resolves through a stable data-testid or role locator

Verifies: Every registerElement call in lib/site/ locates its element with a [data-testid="..."] locator or a role-based locator, and the Vue component renders that test id.

```gherkin
Scenario: Every registered page element resolves through a stable data-testid or role locator
  Given A Saturday page object in lib/site/ registers an element via registerElement
  And The Vue component that element targets is rendered
  When Each registerElement call's locator string is inspected against the rendered component's markup
  Then The locator is either a [data-testid="..."] attribute selector or a role-based locator
  And The Vue component renders an element carrying that exact test id (or exposing that exact role)
```

## The complete-transaction button is found by a dedicated, literal test id

Verifies: CartPage's complete-transaction button is located by data-testid="complete-transaction-btn".

```gherkin
Scenario: The complete-transaction button is found by a dedicated, literal test id
  Given CartPage registers its completeTransactionBtn element
  And CartView.vue renders the complete-transaction button
  When CartPage locates the complete-transaction button
  Then The locator used is exactly [data-testid="complete-transaction-btn"]
  And CartView.vue's button element carries data-testid="complete-transaction-btn" (not checkout-btn, and not a text- or role-based guess)
```

## No locator in the site layer depends on visible text, CSS position, or DOM nesting

Verifies: No locator in lib/site/ uses has-text, :text, a CSS nth- pseudo-selector, or a selector string built from element nesting.

```gherkin
Scenario: No locator in the site layer depends on visible text, CSS position, or DOM nesting
  Given The full set of locator strings used anywhere under lib/site/
  When Each locator string is scanned for has-text(, :text(, a CSS :nth- pseudo-class, or a descendant-combinator (nested) selector such as [data-testid="item-card"] a
  Then Zero locator strings in lib/site/ contain has-text(, :text(, or a :nth- pseudo-class
  And Zero locator strings in lib/site/ combine a parent selector with a nested element via a descendant combinator
  And HomePage.ts's clickFirstProductDetails() locates its target through a flat [data-testid="item-name"] locator instead of [data-testid="item-card"] a
```

## The existing behavior-verification suites are unaffected by the locator cleanup

Verifies: Every scenario that passed before this change passes after it, with no .feature file or assertion changed.

```gherkin
Scenario: The existing behavior-verification suites are unaffected by the locator cleanup
  Given The full existing Playwright and Cucumber suites pass before this change
  When The same suites are run again after the locator changes land, with no .feature file and no assertion edited
  Then Every scenario that passed before still passes
  And No .feature file differs from before the change
  And No assertion in any test file differs from before the change
```

## The change stays small, verifiable, and inside its declared scope

Verifies: Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail, and no file outside scope.paths changes except this run's own delivery record under docs/features/golden-test-ids/.

```gherkin
Scenario: The change stays small, verifiable, and inside its declared scope
  Given The complete set of source and test edits made to satisfy this change
  When Each new or changed function is measured for cyclomatic complexity
  And Each changed line is checked for a covering test, and that test is checked for a path to failure
  And The full list of changed files is compared against scope.paths
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test that fails if the fix is reverted
  And No file outside apps/ye-olde-magic-shop/src/**, apps/ye-olde-magic-shop/lib/**, or docs/** changes, other than this run's delivery record under docs/features/golden-test-ids/
```
