---
feature: "local-fonts"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: local-fonts

## Fonts load from the application's own origin with exact-pinned packages

Verifies: Cinzel (400, 700) and Inter (400, 600) MUST load from the app's own origin, through @fontsource/cinzel and @fontsource/inter pinned to exact versions, and src/styles.css MUST have no @import from another origin.

```gherkin
Scenario: Fonts load from the application's own origin with exact-pinned packages
  Given The application serves its own stylesheet and font packages
  When The browser loads the page
  Then Cinzel and Inter font files are requested only from the app's own origin
  And src/styles.css contains no @import pointing at another origin
  And package.json pins @fontsource/cinzel and @fontsource/inter to exact versions with no ^ or ~
```

## No-request test passes without any off-origin host exception

Verifies: tests/reagent-map-globe.spec.ts MUST drop PRE_EXISTING_OFF_ORIGIN_HOSTS and the code that skips those hosts, and its no-request test MUST pass with no exception.

```gherkin
Scenario: No-request test passes without any off-origin host exception
  Given tests/reagent-map-globe.spec.ts no longer defines PRE_EXISTING_OFF_ORIGIN_HOSTS or any host-skipping logic
  When The no-request test is run
  Then The test passes with no exception raised
  And No request to any off-origin host is observed or excused
```

## No-request test fails if the cross-origin font import returns

Verifies: The no-request test MUST fail when the old cross-origin @import is restored. Record the failing output in the delivery record.

```gherkin
Scenario: No-request test fails if the cross-origin font import returns
  Given The old cross-origin Google Fonts @import is temporarily restored in styles.css
  When The no-request test is run
  Then The test fails
  And The failing output is captured and recorded in the delivery record under docs/features/local-fonts/
  And The @import is reverted afterward so the suite returns to green
```

## Headings render in Cinzel and body text renders in Inter at the correct weights

Verifies: A page opened in the runner image MUST still render headings in Cinzel and body text in Inter, checked in tests/local-fonts.spec.ts with document.fonts.load at the exact weight, never document.fonts.check.

```gherkin
Scenario: Headings render in Cinzel and body text renders in Inter at the correct weights
  Given A page is opened in the runner image
  When tests/local-fonts.spec.ts checks the fonts in use via document.fonts.load
  Then document.fonts.load confirms a Cinzel face loaded at weight 400 and at weight 700
  And document.fonts.load confirms an Inter face loaded at weight 400 and at weight 600
  And The test never calls document.fonts.check
```

## Removing all self-hosted font imports is caught by the font test

Verifies: A page opened in the runner image MUST still render headings in Cinzel and body text in Inter, checked in tests/local-fonts.spec.ts with document.fonts.load at the exact weight, never document.fonts.check.

```gherkin
Scenario: Removing all self-hosted font imports is caught by the font test
  Given All four @fontsource imports are temporarily removed from main.ts
  When tests/local-fonts.spec.ts is run
  Then The test fails because document.fonts.load resolves no matching face for each family/weight pair
  And The failing output is recorded in the delivery record
  And The imports are reverted afterward
```

## Removing only the Inter 600 import is caught by weight, not by document.fonts.check

Verifies: A page opened in the runner image MUST still render headings in Cinzel and body text in Inter, checked in tests/local-fonts.spec.ts with document.fonts.load at the exact weight, never document.fonts.check.

```gherkin
Scenario: Removing only the Inter 600 import is caught by weight, not by document.fonts.check
  Given Only the @fontsource/inter/600.css import is temporarily removed from main.ts
  When tests/local-fonts.spec.ts is run
  Then The test fails because the resolved Inter FontFace's weight does not match 600
  And The failure is detected even though document.fonts.check would have returned true for this case
  And The failing output is recorded in the delivery record
  And The import is reverted afterward
```

## Documentation no longer lists Google Fonts as a known gap

Verifies: docs/testing-cesium.md's Google Fonts entry under "Known gaps" MUST be gone, and the file MUST still follow .claude/rules/writing-style.md.

```gherkin
Scenario: Documentation no longer lists Google Fonts as a known gap
  Given docs/testing-cesium.md previously listed Google Fonts under its "Known gaps" section
  When The local-fonts feature is delivered
  Then The Google Fonts entry is removed from "Known gaps"
  And The rest of docs/testing-cesium.md still follows the Simplified Technical English rules in .claude/rules/writing-style.md
```

## The UI suite passes in both the debian and ubi9 runner images

Verifies: scripts/ui-suite-in-chrome.sh and scripts/ui-suite-in-chrome.sh --ubi9 MUST both pass.

```gherkin
Scenario: The UI suite passes in both the debian and ubi9 runner images
  Given The local-fonts change is complete
  When scripts/ui-suite-in-chrome.sh is run
  And scripts/ui-suite-in-chrome.sh --ubi9 is run
  Then Both runs pass
  And tests/reagent-map-globe.spec.ts's no-request test and tests/local-fonts.spec.ts pass in both runs
```

## The change stays within complexity, coverage, and file-scope budgets

Verifies: Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change (the run's own delivery record under docs/features/local-fonts/ excepted).

```gherkin
Scenario: The change stays within complexity, coverage, and file-scope budgets
  Given The delivered change set for local-fonts
  When Complexity, test coverage, and file scope are checked
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test that can fail
  And No file outside scope.paths changed, except the delivery record under docs/features/local-fonts/
```
