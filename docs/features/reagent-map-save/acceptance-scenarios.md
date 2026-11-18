---
feature: "reagent-map-save"
bounded_context: ""
domain_terms: ["adventurer", "expedition", "reagent map", "draft line", "draft point"]
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: reagent-map-save

## Signed-in adventurer sees a saved expedition again after reopening the reagent map

Verifies: Elara, signed in, MUST see her saved expedition drawn again after reopening the reagent map, and the shop MUST hold it with both ends within 1 km of where she placed them (spec AC11, rule R6).

```gherkin
Scenario: Signed-in adventurer sees a saved expedition again after reopening the reagent map
  Given Elara is signed in
  And Elara has drawn a line on the reagent map from Moonwell Spring to Ember Crag and saved it
  When Elara reopens the reagent map
  Then The saved expedition is drawn on the map
  And The shop holds the expedition with both ends within 1 km of where Elara placed them
```

## Saving without being signed in asks the adventurer to sign in and stores nothing

Verifies: When no adventurer is signed in and a completed line is on the map, saving it MUST ask the adventurer to sign in and MUST NOT send a POST or store anything (spec AC12, rule R7).

```gherkin
Scenario: Saving without being signed in asks the adventurer to sign in and stores nothing
  Given No adventurer is signed in
  And A completed line is drawn on the reagent map
  When The adventurer tries to save the line
  Then The adventurer is asked to sign in
  And Nothing is sent to the shop
  And Nothing new is stored
```

## Saving a zero-length line shows the shop's own message and stores nothing

Verifies: When both points of a completed line are less than 1 km apart, saving it MUST show the server's own zero-length message and MUST NOT store anything (spec AC13, rule R7).

```gherkin
Scenario: Saving a zero-length line shows the shop's own message and stores nothing
  Given An adventurer is signed in
  And Both points of a completed line are less than 1 km apart
  When The adventurer saves the line
  Then The shop's own zero-length message is shown to the adventurer
  And Nothing new is stored
```

## An unsaved line is gone after leaving and reopening the reagent map

Verifies: An unsaved line MUST NOT be shown, and nothing MUST be stored, after the adventurer leaves the reagent map and opens it again (spec AC15, rule R8).

```gherkin
Scenario: An unsaved line is gone after leaving and reopening the reagent map
  Given An adventurer has placed one or two points on the reagent map
  And The adventurer has not saved the line
  When The adventurer leaves the reagent map and opens it again
  Then No draft point or draft line remains on the map
  And Nothing new is stored
```

## The reagent map is visible by default

Verifies: Once this order ships, apps/ye-olde-magic-shop/.env MUST set VITE_FEATURE_REAGENT_MAP=true so the nav link and route show by default, and setting the variable to false MUST still hide both.

```gherkin
Scenario: The reagent map is visible by default
  Given The shop is configured with its shipped default settings
  When An adventurer opens the shop
  Then The reagent map nav link is visible
  And The reagent map route is reachable
```

## Turning the reagent map feature flag off hides the map

Verifies: Once this order ships, apps/ye-olde-magic-shop/.env MUST set VITE_FEATURE_REAGENT_MAP=true so the nav link and route show by default, and setting the variable to false MUST still hide both.

```gherkin
Scenario: Turning the reagent map feature flag off hides the map
  Given The shop's reagent map feature flag is set to false
  When An adventurer opens the shop
  Then The reagent map nav link is hidden
  And The reagent map route is not reachable
```

## Concurrently saving expeditions keeps each adventurer's own data intact

Verifies: Every test that saves an expedition MUST restore apps/mock-api/data/expeditions.json afterward, and no test MUST see an expedition another test saved, proven by running the reagent map specs at --workers=3 --repeat-each=3 with reagent-map-expeditions.spec.ts still counting exactly Elara's expeditions. No existing count assertion MUST be loosened to make this pass.

```gherkin
Scenario: Concurrently saving expeditions keeps each adventurer's own data intact
  Given Two adventurers are signed in on separate sessions
  When Both adventurers save a completed line at the same time
  Then Each adventurer sees only their own newly saved expedition
  And Elara's expedition count is exactly what it was before, plus her own save, unaffected by the other adventurer's save
```

## The reagent map's appearance is unchanged by the save feature unless a reviewed change is recorded

Verifies: Only a @visual baseline this order changes on purpose MUST be regenerated, in both runner images, via scripts/ui-suite-in-chrome.sh --update-baselines and --ubi9 --update-baselines; every changed PNG MUST be viewed by eye and the delivery record MUST state what changed and why; no tolerance MUST be raised.

```gherkin
Scenario: The reagent map's appearance is unchanged by the save feature unless a reviewed change is recorded
  Given The reagent map's existing visual baselines
  When The save feature ships
  Then Each baseline is either unchanged or its change was viewed by eye and recorded with what changed and why
  And No visual comparison tolerance was raised to make a baseline pass
```

## The delivery stays within its agreed scope and quality bars

Verifies: Every new or changed function MUST keep cyclomatic complexity below 7, every changed line MUST be covered by a test proven able to fail, and no file outside scope.paths MUST change except this order's own delivery record under docs/features/reagent-map-save/.

```gherkin
Scenario: The delivery stays within its agreed scope and quality bars
  Given The save feature implementation is complete
  When The change is reviewed
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test proven able to fail
  And No file outside the agreed scope paths changed except this order's own delivery record
```
