---
feature: "golden-site-file-names"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: golden-site-file-names

## Site module files follow kebab-case name.type.ts naming

Verifies: Every file under apps/ye-olde-magic-shop/lib/site/ MUST use kebab-case name.type.ts naming: pages as *.page.ts, the site as magic-shop.site.ts, elements as *.element.ts.

```gherkin
Scenario: Site module files follow kebab-case name.type.ts naming
  Given the ye-olde-magic-shop site module lives under apps/ye-olde-magic-shop/lib/site/
  When the file names in apps/ye-olde-magic-shop/lib/site/ (including flows/) are inspected
  Then every page file is named *.page.ts, the site file is named magic-shop.site.ts, any element file is named *.element.ts, and the existing flows/draw-line.flow.ts is unchanged
  And no PascalCase .ts file remains anywhere in apps/ye-olde-magic-shop/lib/site/
```

## Renamed site files keep their class names and every importer points to the new path

Verifies: Class names MUST stay unchanged and every import of a renamed file MUST be updated to its new path.

```gherkin
Scenario: Renamed site files keep their class names and every importer points to the new path
  Given the site, page, and flow files under apps/ye-olde-magic-shop/lib/site/ have been renamed to kebab-case
  When the exported class names in the renamed files and every file that imports them are inspected
  Then each exported class name (for example MagicShop) is exactly the same as before the rename
  And every importer, including magic-shop.site.ts, flows/draw-line.flow.ts, features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, and tests/reagent-map-ui.spec.ts, imports from the new kebab-case path
  And a search for any of the old PascalCase import paths (lib/site/MagicShop, HomePage, CartPage, LoginPage, AccountPage, ProductDetailsPage, ReagentMapPage) across apps/ye-olde-magic-shop returns zero matches
```

## Existing BDD and Playwright suites pass unchanged after the rename

Verifies: Every scenario that passed before the change MUST still pass after it, with no .feature file or assertion changed.

```gherkin
Scenario: Existing BDD and Playwright suites pass unchanged after the rename
  Given the test:bdd Cucumber suite and the Playwright suite run via scripts/ui-suite-in-chrome.sh were passing before the file rename
  When the same suites are run again against the renamed site files
  Then test:bdd produces the same scenario pass count as before the rename
  And the Playwright suite run through scripts/ui-suite-in-chrome.sh produces the same spec pass count as before the rename
  And no .feature file and no assertion in any .spec.ts or .steps.ts file differs from before the rename
```

## Developer stage reports renamed files in the shape the state validator accepts

Verifies: The developer stage MUST report this work in the exact field shape loom's state validator accepts, so a correct diff is not rejected a second time.

```gherkin
Scenario: Developer stage reports renamed files in the shape the state validator accepts
  Given the developer stage has renamed the seven PascalCase site files to kebab-case using git mv
  When the developer stage records its filesCreated and filesModified fields for this change
  Then filesCreated lists only the new kebab-case paths (magic-shop.site.ts, home.page.ts, cart.page.ts, login.page.ts, account.page.ts, product-details.page.ts, reagent-map.page.ts) and never an old path or an 'old -> new' pair
  And filesModified lists only plain paths for files whose imports changed but which were not themselves renamed (flows/draw-line.flow.ts, features/support/world.ts, features/step_definitions/webapp.steps.ts, tests/cart.spec.ts, tests/reagent-map-ui.spec.ts)
  And loom's state validator accepts the stage on this attempt
```
