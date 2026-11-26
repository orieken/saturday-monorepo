# Context Manifest: golden-test-ids

## Scope

- **Tier**: analyst
- **Built for**: analyst

## Pinned Files

- `apps/ye-olde-magic-shop/lib/site/CartPage.ts` — The only locator violation in lib/site/ (verified by grepping registerElement + has-text|:text|nth- across the whole directory): line 22 registers completeTransactionBtn with `button:has-text("COMPLETE TRANSACTION")`. This is the file AC2 names directly — must become `[data-testid="complete-transaction-btn"]`.
- `apps/ye-olde-magic-shop/src/views/CartView.vue` — Renders the complete-transaction button. It already carries a data-testid, but the wrong one (`data-testid="checkout-btn"` at line 61) — must be renamed to `complete-transaction-btn` to satisfy AC2. Grepped repo-wide for `checkout-btn`: this is the only occurrence, so the rename is isolated.
- `apps/ye-olde-magic-shop/features/step_definitions/checkout.steps.ts` — Confirms the only caller of CartPage's complete-transaction button is `this.site.cartPage.completeTransaction()` — a page-object method, not a raw selector. Proves the locator/testid rename is safe and requires no step-definition edit, consistent with the spec's 'steps are out of scope.'
- `apps/ye-olde-magic-shop/lib/site/HomePage.ts` — Boundary case for AC3: line 20 uses the Playwright locator method `.nth(1)` (not a CSS `nth-` pseudo-selector) against an existing data-testid. Flag so the developer doesn't misread AC3's literal 'nth-' wording as requiring a change here — it isn't a selector string and isn't in scope.
- `apps/ye-olde-magic-shop/lib/site/ProductDetailsPage.ts` — Boundary case for AC1/AC3: line 18 does `this.page.locator('h2')` — a bare tag selector outside any registerElement call. Not has-text/:text/nth-/nesting, so it is not a violation under the AC's literal wording; flagged so it is explicitly left alone rather than 'fixed' as scope creep.
- `apps/ye-olde-magic-shop/lib/site/MagicShop.ts` — Site registry importing CartPage.ts etc. by their current names. Confirms this branch (factory/golden-test-ids, HEAD d096fae) does NOT yet have the golden-site-file-names rename (CartPage.ts -> cart.page.ts) merged in — that work is on a sibling branch (commit 7a4a304, not an ancestor of HEAD). The analyst/developer must act on the file names as they exist now, not the renamed ones, unless/until that dependency lands first.

## Token Budget

- **Status**: OK — ~3328 tokens, 2.8% of the analyst tier budget (120000 tokens)
- Measured by loom from the files above (bytes ÷ 4), not estimated by an agent.

| File | Bytes | ~Tokens |
|---|---:|---:|
| `apps/ye-olde-magic-shop/lib/site/CartPage.ts` | 3160 | 790 |
| `apps/ye-olde-magic-shop/src/views/CartView.vue` | 5496 | 1374 |
| `apps/ye-olde-magic-shop/features/step_definitions/checkout.steps.ts` | 1243 | 310 |
| `apps/ye-olde-magic-shop/lib/site/HomePage.ts` | 1206 | 301 |
| `apps/ye-olde-magic-shop/lib/site/ProductDetailsPage.ts` | 548 | 137 |
| `apps/ye-olde-magic-shop/lib/site/MagicShop.ts` | 1665 | 416 |

## Knowledge Items

None

## ADRs

None

## Pruned

- `apps/ye-olde-magic-shop/lib/site/AccountPage.ts — no registerElement calls and no locator()/getByRole/getByText calls at all; irrelevant to this task`
- `apps/ye-olde-magic-shop/lib/site/LoginPage.ts — same: zero registerElement or locator calls`
- `apps/ye-olde-magic-shop/lib/site/ReagentMapPage.ts — same: zero registerElement or locator calls`
- `apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts — same: zero registerElement or locator calls`
- `apps/ye-olde-magic-shop/tests/cart.spec.ts — exercises CartPage only through page-object methods (getItemQuantity, changeItemQuantity, getCartTotal); never touches completeTransactionBtn or any raw selector, so unaffected by this change`
- `docs/features/golden-site-file-names/** — prior delivery's artifacts for the dependency feature; not yet merged into this branch, read only if the dependency's merge status needs re-checking, not needed to make this change`
