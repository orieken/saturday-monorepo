---
feature: "golden-test-ids"
bounded_context: "Saturday Site-Centric test layer (BasePage/BaseElement locators) for ye-olde-magic-shop"
domain_terms: []
files_touched: ["apps/ye-olde-magic-shop/lib/site/CartPage.ts", "apps/ye-olde-magic-shop/src/views/CartView.vue", "apps/ye-olde-magic-shop/lib/site/HomePage.ts"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: golden-test-ids

## Structural Decisions

### Fix both fragile locators in place using the existing flat `[data-testid="..."]` pattern already dominant in CartPage.ts and HomePage.ts — no new selector strategy, base class, or package.
- **Rationale**: Confirmed by reading both files: every other registerElement/locator call in CartPage.ts and HomePage.ts already uses a flat data-testid attribute selector. The two violations (CartPage.ts:22 has-text, HomePage.ts:39 nested selector) are exceptions to an established pattern, not evidence a new pattern is needed. This is additive-within-pattern per the architect invocation criteria — no layer boundary, package, or abstraction decision is actually in play here.
- **Trade-offs**:
  - None material — confirmed by direct file inspection rather than assumed.
- **Fitness Function**: No locator string under apps/ye-olde-magic-shop/lib/site/ matches has-text(, :text(, a CSS :nth- pseudo-class, or a descendant-combinator (nested) selector string. — verified by qa-engineer-authored regression test scanning lib/site/**/*.ts source text for these patterns; fails the suite on any match.
### Rename CartView.vue's complete-transaction button test id from `checkout-btn` to `complete-transaction-btn` in place, rather than adding a second, parallel test id.
- **Rationale**: grep confirms `checkout-btn` has exactly one occurrence repo-wide (CartView.vue:61). A rename is safe and leaves no stale attribute that could silently reattract a future locator.
- **Trade-offs**:
  - None — single call site verified by grep before deciding.
- **Fitness Function**: judgment-only — grep confirms `checkout-btn` has exactly one occurrence repo-wide (CartView.vue:61). A rename is safe and leaves no stale attribute that could silently reattract a future locator.
### Point HomePage.ts's clickFirstProductDetails() at the existing `item-name` test id (already rendered by ItemCard.vue) instead of inventing a new, purpose-built test id for the details link.
- **Rationale**: ItemCard.vue:26 already renders data-testid="item-name" on a router-link to the same /item/:id destination the current nested selector `[data-testid="item-card"] a` resolves to. Reusing it is behavior-preserving and avoids a redundant attribute on markup that already has an addressable id.
- **Trade-offs**:
  - The locator's semantic intent (navigate to product details) is now carried by an attribute named after the item's name text rather than its own purpose — acceptable: a third testid on the same anchor would duplicate markup for no behavioral gain.
- **Fitness Function**: judgment-only — ItemCard.vue:26 already renders data-testid="item-name" on a router-link to the same /item/:id destination the current nested selector `[data-testid="item-card"] a` resolves to. Reusing it is behavior-preserving and avoids a redundant attribute on markup that already has an addressable id.

## Component Placement

- `CartPage.completeTransactionBtn locator` → adapter layer, `apps/ye-olde-magic-shop/lib/site/CartPage.ts` — Saturday BasePage element registration (interface-adapter layer); selector string changes only, no new component.
- `CartView complete-transaction-btn data-testid` → framework layer, `apps/ye-olde-magic-shop/src/views/CartView.vue` — Vue template attribute in the app under test; renames an existing attribute, adds no new UI.
- `HomePage.clickFirstProductDetails locator` → adapter layer, `apps/ye-olde-magic-shop/lib/site/HomePage.ts` — Saturday BasePage method; swaps a nested CSS selector for the existing item-name testid already rendered by ItemCard.vue.

## Bounded Context

- **Owning Context**: Saturday Site-Centric test layer (BasePage/BaseElement locators) for ye-olde-magic-shop
- **Context Crossings**: ye-olde-magic-shop Vue components (data-testid attributes) — CartView.vue's test id must be renamed in lockstep with CartPage.ts's locator

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [PASS] lib/site/ (interface-adapter layer) must not acquire new dependencies on application/domain internals or new base classes. — Verified by reading CartPage.ts and HomePage.ts directly: both edits are locator-string-only, no new imports, no new classes. Checked by direct file read, not the loom graph (not invoked for a change this small).
- [PASS] Vue component test-id rename must not alter component logic, props, or markup structure. — CartView.vue edit is limited to one attribute value; ItemCard.vue needs no edit since its item-name testid already exists (confirmed by grep).

## Anti-Pattern Check

- Leaky Abstraction: not found — Flat data-testid locators reduce coupling to DOM structure versus the nested-selector and has-text locators being removed; this fix closes leakage rather than adding it.
- Premature Generalization: not found — No new abstraction, package, or base class introduced; reuses the existing testid convention already dominant in both files.

## Fitness Functions

- **No locator string under apps/ye-olde-magic-shop/lib/site/ matches has-text(, :text(, a CSS :nth- pseudo-class, or a descendant-combinator (nested) selector string.**
  - **Verification**: qa-engineer-authored regression test scanning lib/site/**/*.ts source text for these patterns; fails the suite on any match.
  - **Owner**: qa-engineer

## Refactoring Opportunities

- Remove Dead Code on `CartPage.ts — inline comment above completeTransactionBtn explaining the has-text guess ('Using a broad selector for now...')` — The comment documents a workaround for a missing test id; once complete-transaction-btn exists the comment is stale and should be deleted, not edited.

## Developer Handoff Notes

- CartPage.ts:22 — change registerElement('completeTransactionBtn', ButtonElement, 'button:has-text("COMPLETE TRANSACTION")') to '[data-testid="complete-transaction-btn"]'; delete the two-line comment above it.
- CartView.vue:61 — rename data-testid="checkout-btn" to data-testid="complete-transaction-btn".
- HomePage.ts:39 — replace '[data-testid="item-card"] a' with '[data-testid="item-name"]' in clickFirstProductDetails().
- Do not touch HomePage.ts's .nth(1)/.nth(randomIndex) calls or ProductDetailsPage.ts's page.locator('h2') — out of scope per analyst.
- After edits, grep lib/site/ for has-text(, :text(, :nth-, and nested selector strings and confirm zero matches before handing off to code-reviewer.

## Open Architectural Questions

- Analyst's scope question (whether AC3 covers HomePage.ts:39 beyond the single context-manifest-named case) is resolved here: AC3's text applies to every locator in lib/site/, so HomePage.ts:39 is in scope — confirmed by direct read of both files.
- Whether the has-text/:text/nth-/nesting regression scan (this run's fitness function) should become a permanent, standing test beyond this cleanup is still open for qa-engineer or a human to decide; this run's fitness function covers only the current change set.
