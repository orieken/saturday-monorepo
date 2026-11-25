---
feature: "golden-site-file-names"
bounded_context: "Test Automation"
domain_terms: ["Site", "Page", "BaseSite", "BasePage", "BaseFlow"]
files_touched: ["apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts", "apps/ye-olde-magic-shop/lib/site/home.page.ts", "apps/ye-olde-magic-shop/lib/site/cart.page.ts", "apps/ye-olde-magic-shop/lib/site/login.page.ts", "apps/ye-olde-magic-shop/lib/site/account.page.ts", "apps/ye-olde-magic-shop/lib/site/product-details.page.ts", "apps/ye-olde-magic-shop/lib/site/reagent-map.page.ts"]
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Architecture Notes: golden-site-file-names

## Structural Decisions

### No new structural decision is required: this feature renames existing files to match the already-established name.type.ts kebab-case convention (demonstrated by flows/draw-line.flow.ts), with no class rename, no new abstraction, no layer boundary change, and no new package.
- **Rationale**: CLAUDE.md's architect-invocation criteria (new package, new base class, layer boundary change, cross-cutting concern, or constraining decision) do not apply. Verified on disk: apps/ye-olde-magic-shop/lib/site/ already contains only account.page.ts, cart.page.ts, home.page.ts, login.page.ts, magic-shop.site.ts, product-details.page.ts, reagent-map.page.ts, and flows/draw-line.flow.ts. A grep for the seven old PascalCase paths across apps/ye-olde-magic-shop returned zero matches, and magic-shop.site.ts's imports already read from the six new kebab-case page paths. This is additive-within-an-established-pattern work, which CLAUDE.md explicitly instructs the architect to skip.
- **Fitness Function**: judgment-only — CLAUDE.md's architect-invocation criteria (new package, new base class, layer boundary change, cross-cutting concern, or constraining decision) do not apply. Verified on disk: apps/ye-olde-magic-shop/lib/site/ already contains only account.page.ts, cart.page.ts, home.page.ts, login.page.ts, magic-shop.site.ts, product-details.page.ts, reagent-map.page.ts, and flows/draw-line.flow.ts. A grep for the seven old PascalCase paths across apps/ye-olde-magic-shop returned zero matches, and magic-shop.site.ts's imports already read from the six new kebab-case page paths. This is additive-within-an-established-pattern work, which CLAUDE.md explicitly instructs the architect to skip.

## Component Placement

- `MagicShop` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/magic-shop.site.ts` — App-specific Site composition root; unchanged placement, filename now matches name.type.ts convention.
- `HomePage` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/home.page.ts` — App-specific Page object extending BasePage; unchanged placement.
- `CartPage` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/cart.page.ts` — App-specific Page object extending BasePage; unchanged placement.
- `LoginPage` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/login.page.ts` — App-specific Page object extending BasePage; unchanged placement.
- `AccountPage` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/account.page.ts` — App-specific Page object extending BasePage; unchanged placement.
- `ProductDetailsPage` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/product-details.page.ts` — App-specific Page object extending BasePage; unchanged placement.
- `ReagentMapPage` → Interface Adapters layer, `apps/ye-olde-magic-shop/lib/site/reagent-map.page.ts` — App-specific Page object extending BasePage; unchanged placement.

## Bounded Context

- **Owning Context**: Test Automation
- **Context Crossings**: apps/ye-olde-magic-shop/lib/site/** is imported by features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, and tests/reagent-map-ui.spec.ts by file path — all four already updated to the new paths, verified on disk.

## Stability Design

None

## Observability Design

None

## Layer Boundary Checks

- [PASS] Site/Page objects (Interface Adapters layer) must not introduce new dependencies on domain or framework internals as part of this rename. — Confirmed by reading magic-shop.site.ts: only import changes are the six relative page paths and the existing BaseSite/Page/DrawLineFlow imports, all unchanged in target, only path changed. mcp__loom__* graph tools were not used (no need — a plain import read sufficed for a pure rename with no new edges).

## Anti-Pattern Check

None

## Fitness Functions

- **Every file in apps/ye-olde-magic-shop/lib/site/ matches kebab-case name.type.ts (page/site/element/flow suffix), with no PascalCase .ts file remaining in that directory.**
  - **Verification**: judgment-only — code-reviewer confirms the rename map and naming shape at review time; a machine-enforced lint for this convention is the separate, downstream intake/golden-lint-and-exemplars.md chore and is out of scope for this run.
  - **Owner**: code-reviewer

## Refactoring Opportunities

None

## Developer Handoff Notes

- No architect-driven changes are needed. Verified on disk that the rename and import updates described in analysis.md are already complete and correct: lib/site/ contains exactly the seven target files plus flows/draw-line.flow.ts, and a repo-wide grep for the seven old PascalCase paths returns zero matches.
- Developer's only remaining obligation is the reporting-format fix analysis.md already specifies: list the seven renamed files once each, by new path only, in filesCreated, and list flows/draw-line.flow.ts, features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, and tests/reagent-map-ui.spec.ts as plain paths in filesModified — never an 'old -> new' pair in any field.

## Open Architectural Questions

None
