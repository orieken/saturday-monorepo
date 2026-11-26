---
feature: "golden-test-ids"
bounded_context: ""
domain_terms: ["BasePage", "BaseElement"]
files_touched: []
issue_refs: ["golden-test-ids"]
linked_adrs: []
linked_kis: []
---

# Code Review: golden-test-ids

## Overall Status

**APPROVED**

## Design Narrative

The developer replaced two fragile Playwright locators in the Site-Centric layer — CartPage's text-matching `button:has-text(...)` and HomePage's nested `[data-testid="item-card"] a` selector — with flat `data-testid` locators, including renaming CartView.vue's mislabeled `checkout-btn` to `complete-transaction-btn`. No public method signatures, control flow, or rendered markup changed; only selector strings and one attribute value moved to point at stable test ids. Behaviour change: no.

## Design Score

- **Clarity** 5/5
- **Cohesion** 5/5
- **Coupling** 5/5
- **Craft** 5/5

## Security Surface

None

## Performance Surface

None

## Test Design Review

- No tests were added or modified in this diff — the three changed files (CartPage.ts, HomePage.ts, CartView.vue) are production/test-infrastructure locator code, not test files themselves. Existing coverage (checkout.steps.ts calling completeTransaction(), inventory.steps.ts calling clickFirstProductDetails()) was verified by grep to exercise both changed call sites unmodified, which is the behavior-preservation evidence this stage required rather than a new assertion to judge.

## Verification of Developer Self-Review

- Verified 'reveals intention' claim: completeTransactionBtn's selector now matches its name and CartView's testid is no longer mislabeled checkout-btn — confirmed by reading both files.
- Verified 'no duplication' claim: grepped repo for item-name/item-card testids and confirmed HomePage.ts now reuses ItemCard.vue's existing item-name testid rather than adding a new one.
- Verified the dead-comment removal claim against the actual diff — the two-line has-text-guess comment is gone, consistent with the Named Refactoring Log entry.
- Verified the rename-in-place claim: grepped the full repo for 'checkout-btn' and found zero remaining occurrences, confirming the rename left no stale attribute.
- Could not independently verify the 'passes the tests' line (suite wasn't run, as the developer disclosed) — accepted as an honest limitation given no node_modules in this worktree, not a discrepancy, since the developer flagged rather than claimed it.

## Feedback for the Developer

None
