---
feature: "golden-site-file-names"
bounded_context: "Test Automation"
domain_terms: ["BaseSite", "BasePage", "BaseElement", "BaseFlow"]
files_touched: ["apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts", "apps/ye-olde-magic-shop/lib/site/home.page.ts", "apps/ye-olde-magic-shop/lib/site/cart.page.ts", "apps/ye-olde-magic-shop/lib/site/login.page.ts", "apps/ye-olde-magic-shop/lib/site/account.page.ts", "apps/ye-olde-magic-shop/lib/site/product-details.page.ts", "apps/ye-olde-magic-shop/lib/site/reagent-map.page.ts", "apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts", "apps/ye-olde-magic-shop/features/support/world.ts", "apps/ye-olde-magic-shop/features/step_definitions/webapp.steps.ts", "apps/ye-olde-magic-shop/tests/cart.spec.ts", "apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts", "docs/features/golden-site-file-names/"]
issue_refs: ["intake/golden-site-file-names.md"]
linked_adrs: []
linked_kis: []
---

# Feature Analysis: golden-site-file-names

## Summary

The demo app's site model under apps/ye-olde-magic-shop/lib/site/ used PascalCase filenames (MagicShop.ts, CartPage.ts, HomePage.ts, LoginPage.ts, AccountPage.ts, ProductDetailsPage.ts, ReagentMapPage.ts), violating the project's name.type.ts kebab-case file-naming rule that flows/draw-line.flow.ts already followed. This is a pure rename: the site file becomes magic-shop.site.ts and the six page files become *.page.ts, every class name stays exactly as it is today, and every import of those files is rewritten to the new path with no other change to its line. A prior developer attempt already performed every rename and import update correctly on disk, but the pipeline rejected that stage because it reported the renames as `old -> new` pairs inside `filesModified`, a format loom's state validator cannot resolve (it checks each listed path exists and has no field for a removed file). This analysis update exists to make the developer stage report the same, already-correct work in the format loom accepts: each renamed file listed once by its new path in `filesCreated`, and only files whose import lines changed (but which were not themselves renamed) listed as plain paths in `filesModified`. No .feature file, assertion, or runtime behavior changes; the existing Cucumber and Playwright suites must pass identically before and after.

### Acceptance Criteria

- Every file under apps/ye-olde-magic-shop/lib/site/ MUST use kebab-case name.type.ts naming: pages as *.page.ts, the site as magic-shop.site.ts, elements as *.element.ts.
  - Example: MagicShop.ts -> magic-shop.site.ts
  - Example: HomePage.ts -> home.page.ts
  - Example: CartPage.ts -> cart.page.ts
  - Example: LoginPage.ts -> login.page.ts
  - Example: AccountPage.ts -> account.page.ts
  - Example: ProductDetailsPage.ts -> product-details.page.ts
  - Example: ReagentMapPage.ts -> reagent-map.page.ts
  - Example: flows/draw-line.flow.ts already matches name.type.ts and does not move
  - Example: No standalone *.element.ts file exists today — CartPage.ts's inline, non-exported TextElement class is not a separate file, so this rule has no current target to rename
  - Example: Verified on disk: apps/ye-olde-magic-shop/lib/site/ currently contains exactly account.page.ts, cart.page.ts, home.page.ts, login.page.ts, magic-shop.site.ts, product-details.page.ts, reagent-map.page.ts, and flows/draw-line.flow.ts — no PascalCase .ts file remains
- Class names MUST stay unchanged and every import of a renamed file MUST be updated to its new path.
  - Example: magic-shop.site.ts's `export class MagicShop extends BaseSite` keeps the name MagicShop
  - Example: magic-shop.site.ts's imports now read from ./home.page, ./cart.page, ./product-details.page, ./login.page, ./account.page, ./reagent-map.page — verified on disk
  - Example: flows/draw-line.flow.ts's `import type { ReagentMapPage } from '../reagent-map.page'` — verified on disk
  - Example: features/support/world.ts imports `MagicShop` from '../../lib/site/magic-shop.site' — verified on disk
  - Example: features/step_definitions/webapp.steps.ts imports `MagicShop` from '../../lib/site/magic-shop.site.ts', preserving its existing explicit .ts extension style — verified on disk
  - Example: tests/cart.spec.ts and tests/reagent-map-ui.spec.ts import `MagicShop` from '../lib/site/magic-shop.site' — verified on disk
  - Example: A grep for lib/site/MagicShop, lib/site/HomePage, lib/site/CartPage, lib/site/LoginPage, lib/site/AccountPage, lib/site/ProductDetailsPage, or lib/site/ReagentMapPage across apps/ye-olde-magic-shop currently returns zero matches
- Every scenario that passed before the change MUST still pass after it, with no .feature file or assertion changed.
  - Example: test:bdd (TS_NODE_PROJECT=tsconfig.cucumber.json cucumber-js) produces the same pass count before and after
  - Example: scripts/ui-suite-in-chrome.sh's Playwright run produces the same pass count before and after
- The developer stage MUST report this work in the exact field shape loom's state validator accepts, so a correct diff is not rejected a second time.
  - Example: List each renamed file once, by its new path only, in filesCreated: apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts, home.page.ts, cart.page.ts, login.page.ts, account.page.ts, product-details.page.ts, reagent-map.page.ts
  - Example: Never list an old path and never list an 'old -> new' pair in any field — loom checks that each listed path exists and has no field for a removed file
  - Example: List files whose import lines changed but which were not themselves renamed as plain paths in filesModified: flows/draw-line.flow.ts, features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, tests/reagent-map-ui.spec.ts
  - Example: The prior attempt's rejection (developer.rejected.json) was caused exactly by violating this: it put 'AccountPage.ts -> account.page.ts'-style pairs into filesModified

### Non-Functional Requirements

None

## Proposed Fitness Functions

- **Every file in apps/ye-olde-magic-shop/lib/site/ matches kebab-case name.type.ts (page/site/element/flow suffix), with no PascalCase .ts file remaining in that directory.**
  - **Verification**: judgment-only — code-reviewer confirms the rename map and naming shape at review time; a machine-enforced lint for this convention is the separate, downstream intake/golden-lint-and-exemplars.md chore and is out of scope for this run.
  - **Owner**: code-reviewer

## Out of Scope

- Any behavior change.
- Changing any .feature file or test assertion.
- Wiring a lint rule or CI check that enforces the name.type.ts convention going forward — tracked separately in intake/golden-lint-and-exemplars.md, which depends on this chore landing first.
- Splitting CartPage.ts's (now cart.page.ts's) inline TextElement class into a standalone *.element.ts file — no standalone element file exists today, so there is nothing to rename under that rule.

## Technical Breakdown

None

### Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: apps/ye-olde-magic-shop/lib/site/** (the renamed site model) is imported by apps/ye-olde-magic-shop/features/support/world.ts and features/step_definitions/webapp.steps.ts (Cucumber World/step bindings) and by apps/ye-olde-magic-shop/tests/cart.spec.ts and tests/reagent-map-ui.spec.ts (Playwright specs) by file path, so every rename requires a matching import update in all four.

### Domain Events (Event Storming Lite)

- **Commands**: None
- **Events Produced**: None
- **Owning Aggregates**: None
- **Read Models / Projections**: None

### Affected Components

- `apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts` — Already renamed from MagicShop.ts on disk; its six relative page imports already updated to the new kebab-case paths; class MagicShop unchanged. Developer stage must report this path in filesCreated (new path only).
- `apps/ye-olde-magic-shop/lib/site/home.page.ts` — Already renamed from HomePage.ts on disk; class HomePage unchanged; no internal imports to update. Report in filesCreated.
- `apps/ye-olde-magic-shop/lib/site/cart.page.ts` — Already renamed from CartPage.ts on disk; class CartPage (and its inline, non-exported TextElement) unchanged; no internal imports to update. Report in filesCreated.
- `apps/ye-olde-magic-shop/lib/site/login.page.ts` — Already renamed from LoginPage.ts on disk; class LoginPage unchanged; no internal imports to update. Report in filesCreated.
- `apps/ye-olde-magic-shop/lib/site/account.page.ts` — Already renamed from AccountPage.ts on disk; class AccountPage unchanged; no internal imports to update. Report in filesCreated.
- `apps/ye-olde-magic-shop/lib/site/product-details.page.ts` — Already renamed from ProductDetailsPage.ts on disk; class ProductDetailsPage unchanged; no internal imports to update. Report in filesCreated.
- `apps/ye-olde-magic-shop/lib/site/reagent-map.page.ts` — Already renamed from ReagentMapPage.ts on disk; class ReagentMapPage unchanged; imported by both magic-shop.site.ts and flows/draw-line.flow.ts, both already updated. Report in filesCreated.
- `apps/ye-olde-magic-shop/lib/site/flows/draw-line.flow.ts` — Already name.type.ts-shaped; not renamed. Its import of ReagentMapPage already updated to '../reagent-map.page'. Report as a plain path in filesModified, never as an old -> new pair.
- `apps/ye-olde-magic-shop/features/support/world.ts` — Its MagicShop import already updated to '../../lib/site/magic-shop.site'. Report as a plain path in filesModified.
- `apps/ye-olde-magic-shop/features/step_definitions/webapp.steps.ts` — Its MagicShop import already updated to '../../lib/site/magic-shop.site.ts', preserving the explicit .ts extension this file already used. Report as a plain path in filesModified.
- `apps/ye-olde-magic-shop/tests/cart.spec.ts` — Its MagicShop import already updated to '../lib/site/magic-shop.site'. Report as a plain path in filesModified.
- `apps/ye-olde-magic-shop/tests/reagent-map-ui.spec.ts` — Its MagicShop import already updated to '../lib/site/magic-shop.site'. Report as a plain path in filesModified.
- `docs/features/golden-site-file-names/` — Delivery record: the full old-path -> new-path rename map and the before/after suite-pass results; the one location allowed to change outside scope.paths besides the renames themselves.

### Data Model Changes

None

### API Changes

None

### New Dependencies

None

## Task List

None

### Developer Tasks

1. Verify, do not redo: confirm apps/ye-olde-magic-shop/lib/site/ already contains only magic-shop.site.ts, home.page.ts, cart.page.ts, login.page.ts, account.page.ts, product-details.page.ts, reagent-map.page.ts, and flows/draw-line.flow.ts — no PascalCase file remains (this was done by a prior attempt).
2. Verify class names are unchanged: MagicShop, HomePage, CartPage, LoginPage, AccountPage, ProductDetailsPage, ReagentMapPage each still `export class <Name>` in their new file.
3. Verify magic-shop.site.ts's six relative page imports already read ./home.page, ./cart.page, ./product-details.page, ./login.page, ./account.page, ./reagent-map.page.
4. Verify flows/draw-line.flow.ts's import already reads `import type { ReagentMapPage } from '../reagent-map.page'`.
5. Verify features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, and tests/reagent-map-ui.spec.ts already import MagicShop from the new magic-shop.site path, each preserving its own prior '.ts'-extension style.
6. Grep the whole apps/ye-olde-magic-shop tree for any remaining reference to the seven old PascalCase paths (lib/site/MagicShop, lib/site/HomePage, lib/site/CartPage, lib/site/LoginPage, lib/site/AccountPage, lib/site/ProductDetailsPage, lib/site/ReagentMapPage) and fix any found; confirm none remain.
7. Run `npx tsc --noEmit` (or the project's existing type-check path) against the tree to confirm no import resolves incorrectly.
8. Report state correctly this time: list the seven renamed files once each, by new path only, in filesCreated (magic-shop.site.ts, home.page.ts, cart.page.ts, login.page.ts, account.page.ts, product-details.page.ts, reagent-map.page.ts). List flows/draw-line.flow.ts, features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, and tests/reagent-map-ui.spec.ts as plain paths in filesModified. Never list an old path, and never write an 'old -> new' pair in any field — that exact mistake caused the prior attempt's rejection.

### QA Tasks

1. Run apps/ye-olde-magic-shop's test:bdd script; confirm every scenario that passed before this change still passes, with zero .feature or assertion diffs.
2. Run the Playwright suite through scripts/ui-suite-in-chrome.sh, one runner image at a time, never two at once, waiting for each run (about two minutes) to finish inside the stage's 30-minute budget; confirm the same specs pass as before.
3. Confirm no .feature file and no assertion in any .spec.ts or .steps.ts file changed — only import paths changed.
4. Confirm the diff touches no file outside scope.paths (apps/ye-olde-magic-shop/lib/**, features/**, tests/**, docs/**) except the delivery record under docs/features/golden-site-file-names/.

### Tech Writer Tasks

1. Write the delivery record under docs/features/golden-site-file-names/ capturing the full old-path -> new-path rename map and the before/after test-suite pass results.
2. State explicitly that no README, API doc, or ADR changes are needed beyond the delivery record, since this is a pure rename with no behavior change.

### DevOps Tasks

None

## Edge Cases and Risks

- apps/ye-olde-magic-shop/reports/cucumber/*.html and *.json (generated run artifacts) and apps/ye-olde-magic-shop/docs/USING_SATURDAY_MCP.md mention 'MagicShop'.: Leave untouched. The reports/ files are generated outputs outside scope.paths, not source. The docs mention is an unrelated sample class name (MagicShopSite) inside a code snippet, not an import of the renamed file.
- features/step_definitions/webapp.steps.ts imports MagicShop with an explicit '.ts' extension, while world.ts, cart.spec.ts, and reagent-map-ui.spec.ts omit it.: Each file's existing extension style is already preserved in the rewritten path — only the filename segment changed (MagicShop -> magic-shop.site), not whether '.ts' is present.
- cart.page.ts defines a second, non-exported class TextElement inline, which matches the *.element.ts rule's subject matter but is not its own file.: Leave it inside cart.page.ts. The spec renames files, not class structure, and no standalone element file exists today for this rule to apply to.
- flows/draw-line.flow.ts already matches name.type.ts and must not be renamed, but imports the renamed reagent-map.page.ts.: Only the import statement inside draw-line.flow.ts changes; the file itself is not renamed or otherwise modified. Already verified correct on disk.
- A prior developer stage performed every rename correctly with `git mv` (git status shows clean R/RM renames, not delete+add) but was still rejected.: The rejection was purely a state-reporting format error (old -> new pairs inside filesModified), not a code or git-history problem. Do not redo the renames; fix only how the stage reports filesCreated/filesModified.

## Definition of Done

- Every file in apps/ye-olde-magic-shop/lib/site/ matches kebab-case name.type.ts (pages as *.page.ts, the site as magic-shop.site.ts, the flow unchanged), with no PascalCase .ts file remaining in that directory.
- Every class name is unchanged, and a grep for the seven old PascalCase import paths across apps/ye-olde-magic-shop returns no matches.
- test:bdd and the Playwright suite (via scripts/ui-suite-in-chrome.sh) pass with the same scenario results as before the rename, with zero .feature or assertion diffs — recorded in docs/features/golden-site-file-names/.
- No file outside scope.paths changed, except the delivery record under docs/features/golden-site-file-names/.
- `npx tsc --noEmit` (or equivalent) resolves cleanly against the renamed tree.
- The developer stage's filesCreated/filesModified report matches the actual git rename state exactly — new paths only in filesCreated, plain paths in filesModified, no old path or 'old -> new' pair anywhere — so loom accepts the stage on this attempt.
