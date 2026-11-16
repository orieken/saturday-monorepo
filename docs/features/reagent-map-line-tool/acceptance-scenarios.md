---
feature: "reagent-map-line-tool"
bounded_context: ""
domain_terms: ["Adventurer", "Reagent Site", "Expedition", "Seeded Expedition", "Aurelis"]
files_touched: []
issue_refs: ["reagent-map-line-tool", "AC6", "AC7", "AC8", "AC9", "AC10", "AC14"]
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: reagent-map-line-tool

## Placing two points on the land completes a line ready to save

Verifies: The line tool MUST join two points placed on the land with a completed geodesic line that is ready to save.

```gherkin
Scenario: Placing two points on the land completes a line ready to save
  Given The adventurer has started the line tool on the reagent map
  When The adventurer places a point at Moonwell Spring
  And The adventurer places a point at Ember Crag
  Then A completed line joins Moonwell Spring and Ember Crag
  And The completed line is ready to be saved
```

## Placing only one point shows no completed line

Verifies: After only one point is placed, the line tool MUST NOT show a completed line and that line MUST NOT be saveable.

```gherkin
Scenario: Placing only one point shows no completed line
  Given The adventurer has started the line tool
  When The adventurer places a point at Moonwell Spring
  Then No completed line is shown
  And There is nothing to save
```

## The completed line follows the shortest path over Aurelis's surface

Verifies: A completed line MUST follow the shortest path over Aurelis's surface between its two points, not a straight line on a flat projection.

```gherkin
Scenario: The completed line follows the shortest path over Aurelis's surface
  Given The adventurer has started the line tool
  When The adventurer places a point at Thornveil Glade
  And The adventurer places a point at Glimmerdeep Caves
  Then The line between Thornveil Glade and Glimmerdeep Caves follows the shortest path over Aurelis's curved surface, not a straight line on a flat map
```

## A point can be placed anywhere on the visible surface, including the sea

Verifies: The line tool MUST place a point anywhere on the near side of the globe, sea included, since 'the land' means the planet's surface for now.

```gherkin
Scenario: A point can be placed anywhere on the visible surface, including the sea
  Given The adventurer has started the line tool
  When The adventurer places a point over open sea on the visible side of Aurelis
  Then The point is placed there
```

## Clicking beside the globe places nothing

Verifies: A click on space beside the globe MUST place no point, show nothing, and leave the line tool waiting for its first point.

```gherkin
Scenario: Clicking beside the globe places nothing
  Given The adventurer has started the line tool
  When The adventurer clicks on space beside the globe
  Then No point is placed
  And The line tool still waits for its first point
```

## Clicking the globe's far edge places nothing

Verifies: A click on the globe's edge, where the surface faces away from the camera, MUST place no point.

```gherkin
Scenario: Clicking the globe's far edge places nothing
  Given The adventurer has started the line tool
  When The adventurer clicks on the globe's edge, where the surface faces away from them
  Then No point is placed
  And The line tool still waits for its first point
```

## Starting the line tool again discards the point already placed

Verifies: Starting the line tool again MUST discard any points from the line already in progress, so the new line has no points.

```gherkin
Scenario: Starting the line tool again discards the point already placed
  Given The adventurer has started the line tool and placed one point
  When The adventurer starts the line tool again
  Then The new line has no points
  And The previously placed point is discarded
```

## Placed points and completed lines are marked as drafts, not yet saved

Verifies: Every point and line the tool draws MUST carry saturdayMetadata with kind 'draft-point' or 'draft-line' so tests can read them through getEntities().

```gherkin
Scenario: Placed points and completed lines are marked as drafts, not yet saved
  Given The adventurer has started the line tool
  When The adventurer places a point at Moonwell Spring
  And The adventurer places a point at Ember Crag
  Then Each placed point and the completed line are identifiable as drafts, distinct from any saved expedition
```
