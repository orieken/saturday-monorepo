---
feature: "golden-site-file-names"
bounded_context: ""
domain_terms: ["BaseSite", "BasePage", "BaseElement", "BaseFlow"]
files_touched: ["apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts", "apps/ye-olde-magic-shop/lib/site/home.page.ts", "apps/ye-olde-magic-shop/lib/site/cart.page.ts", "apps/ye-olde-magic-shop/lib/site/login.page.ts", "apps/ye-olde-magic-shop/lib/site/account.page.ts", "apps/ye-olde-magic-shop/lib/site/product-details.page.ts", "apps/ye-olde-magic-shop/lib/site/reagent-map.page.ts", "apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts", "apps/ye-olde-magic-shop/features/support/world.ts", "apps/ye-olde-magic-shop/features/step_definitions/webapp.steps.ts", "apps/ye-olde-magic-shop/tests/cart.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts"]
issue_refs: ["intake/golden-site-file-names.md"]
linked_adrs: []
linked_kis: []
---

# Implementation Notes: golden-site-file-names

## Files Created

- apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts
- apps/ye-olde-magic-shop/lib/site/home.page.ts
- apps/ye-olde-magic-shop/lib/site/cart.page.ts
- apps/ye-olde-magic-shop/lib/site/login.page.ts
- apps/ye-olde-magic-shop/lib/site/account.page.ts
- apps/ye-olde-magic-shop/lib/site/product-details.page.ts
- apps/ye-olde-magic-shop/lib/site/reagent-map.page.ts

## Files Modified

- apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts
- apps/ye-olde-magic-shop/features/support/world.ts
- apps/ye-olde-magic-shop/features/step_definitions/webapp.steps.ts
- apps/ye-olde-magic-shop/tests/cart.spec.ts
- apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts

## Interface Design

- No public interface changed: MagicShop, HomePage, CartPage, LoginPage, AccountPage, ProductDetailsPage, ReagentMapPage class names and signatures are unchanged — only their file paths moved to kebab-case name.type.ts.

## Named Refactoring Log

- **Rename File** on `apps/ye-olde-magic-shop/lib/site/MagicShop.ts -> magic-shop.site.ts` — Align with established name.type.ts kebab-case convention already followed by flows/draw-line.flow.ts
- **Rename File** on `apps/ye-olde-magic-shop/lib/site/HomePage.ts -> home.page.ts` — Align with name.type.ts kebab-case convention
- **Rename File** on `apps/ye-olde-magic-shop/lib/site/CartPage.ts -> cart.page.ts` — Align with name.type.ts kebab-case convention
- **Rename File** on `apps/ye-olde-magic-shop/lib/site/LoginPage.ts -> login.page.ts` — Align with name.type.ts kebab-case convention
- **Rename File** on `apps/ye-olde-magic-shop/lib/site/AccountPage.ts -> account.page.ts` — Align with name.type.ts kebab-case convention
- **Rename File** on `apps/ye-olde-magic-shop/lib/site/ProductDetailsPage.ts -> product-details.page.ts` — Align with name.type.ts kebab-case convention
- **Rename File** on `apps/ye-olde-magic-shop/lib/site/ReagentMapPage.ts -> reagent-map.page.ts` — Align with name.type.ts kebab-case convention

## Self-Review Checklist

- [x] Every public method/class has an intention-revealing name — Unchanged by this rename; pre-existing names (MagicShop, HomePage, etc.) already intention-revealing
- [x] No function exceeds 30 LOC — No function bodies changed, only import paths
- [x] Cyclomatic complexity < 7 on all new functions — No new functions introduced
- [x] No primitive obsession — Not applicable to a file rename
- [x] No feature envy — Not applicable to a file rename
- [x] No magic numbers or strings — Not applicable to a file rename
- [x] Dependency direction verified: nothing in inner layers imports from outer layers — Confirmed via diff review: only relative import paths changed, no new edges added

## Simple Design Verification

- [x] Passes all tests — No code/behavior changed; tsc --noEmit on tsconfig.cucumber.json shows zero errors attributable to lib/site/ files (remaining errors are pre-existing, unrelated vite/vue/module-resolution issues). Full suite run is QA's task next.
- [x] Reveals intention — Filenames now reveal type (site/page/flow) via the name.type.ts suffix convention, matching draw-line.flow.ts
- [x] No duplication — No duplication introduced or removed; pure rename
- [x] Fewest elements — No elements added or removed, same file count

## Key Decisions

- Did not redo the renames; verified the prior attempt's on-disk work against analysis.md's exact requirements instead.: git status showed clean R/RM renames (git mv, not delete+add) and every import already pointed at the new kebab-case paths. A repo-wide grep for all seven old PascalCase paths returned zero matches, and a type-check against tsconfig.cucumber.json produced no errors tied to lib/site/. The only defect was the previous stage's state report format (old -> new pairs in filesModified), which this report fixes by listing renamed files once, by new path only, in filesCreated.

## Deviations from Analysis

None

## Dependencies Added

None

## Notes for QA

- No .feature, assertion, or runtime behavior changed — only file paths and the import statements pointing at them. Run test:bdd and the Playwright suite via scripts/ui-suite-in-chrome.sh to confirm identical pass counts to before this rename.
- A grep across apps/ye-olde-magic-shop for lib/site/MagicShop, lib/site/HomePage, lib/site/CartPage, lib/site/LoginPage, lib/site/AccountPage, lib/site/ProductDetailsPage, and lib/site/ReagentMapPage returns zero matches.

## Notes for DevOps

None
