---
feature: "reagent-map-expeditions"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: ["AC3:ui", "AC4:ui", "AC5"]
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: reagent-map-expeditions

## Signed-out adventurer sees only the shop's seeded expeditions drawn on the globe

Verifies: The view MUST draw every expedition GET /api/expeditions returns for the current adventurer as a geodesic polyline, both ends within 1 km of its stored ends, tagged with saturdayMetadata { id, name, kind: 'expedition' }.

```gherkin
Scenario: Signed-out adventurer sees only the shop's seeded expeditions drawn on the globe
  Given The adventurer is signed out
  And The reagent map view has loaded
  When The globe finishes loading expeditions from GET /api/expeditions
  Then Exactly three expeditions are drawn: 'Spring to Crag', 'Glade to Caves', and 'Cairn to Shoals'
  And Each drawn expedition is a geodesic line with both ends within 1 km of its stored reagent site coordinates
  And Each drawn expedition carries saturdayMetadata with id, name, and kind: 'expedition'
```

## Signed-in adventurer sees the seeded expeditions plus their own, never another adventurer's

Verifies: The view MUST draw every expedition GET /api/expeditions returns for the current adventurer as a geodesic polyline, both ends within 1 km of its stored ends, tagged with saturdayMetadata { id, name, kind: 'expedition' }.

```gherkin
Scenario: Signed-in adventurer sees the seeded expeditions plus their own, never another adventurer's
  Given Elara is signed in
  And The reagent map view has loaded
  When The globe finishes loading expeditions from GET /api/expeditions
  Then The three seeded expeditions ('Spring to Crag', 'Glade to Caves', 'Cairn to Shoals') are drawn
  And 'Moonwell Spring to Frostcap Hollow' is also drawn
  And 'Ember Crag to Saltmarsh Cairn' is not drawn
```

## Choosing a far-side site from the list flies the camera so it lands on screen on the near side

Verifies: Choosing a site from the list MUST fly the camera so that site lands on screen on the near side of the globe, including a site that starts on the far side.

```gherkin
Scenario: Choosing a far-side site from the list flies the camera so it lands on screen on the near side
  Given The default view is showing
  And Tidewrack Shoals is on the far side of the globe and not on screen
  When The adventurer chooses Tidewrack Shoals from the site list
  Then The camera flies so Tidewrack Shoals lands on screen on the near side of the globe
```

## Drawing expeditions on the globe changes only the globe baseline image

Verifies: The @visual globe baseline MUST change only because the newly drawn expeditions now appear in it. Regenerate it in both runner images, view each changed PNG, and record what changed and why. No other baseline changes, and no tolerance is raised.

```gherkin
Scenario: Drawing expeditions on the globe changes only the globe baseline image
  Given The expeditions are now drawn on the globe by this change
  When The visual baseline suite is regenerated with --update-baselines in both runner images
  Then Only globe-default-view.png differs from its prior baseline, in each image
  And reagent-map-header.png, reagent-site-list.png, and reagent-map-error.png remain byte-identical to their prior baselines
  And No existing visual-diff tolerance is raised to make a difference pass
```

## The change stays within its complexity, coverage, and scope limits

Verifies: Every new or changed function MUST keep cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except this order's own delivery record under docs/features/reagent-map-expeditions/.

```gherkin
Scenario: The change stays within its complexity, coverage, and scope limits
  Given The implementation of this order is complete
  When Complexity analysis and test coverage are checked against the changed lines, and the changed file list is checked against scope.paths
  Then Every new or changed function has cyclomatic complexity below 7
  And Every changed line is covered by a test that can fail
  And No file outside scope.paths changed, except this order's own delivery record under docs/features/reagent-map-expeditions/
```
