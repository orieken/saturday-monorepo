# Acceptance criteria: reagent-map-spec

Generated from `acceptance-criteria.json`. Settled by the analyst at `answer-questions`:

- **tolerance**: 1 km on the planet's surface (Q3)
- **line**: geodesic: the shortest path over the surface (Q4)
- **offLand**: nothing happens; the line tool keeps waiting (Q5)
- **visibility**: an adventurer sees the shop's seeded expeditions and their own; signed out, seeded only (Q6)
- **samePlace**: two points closer than the tolerance are the same place (Q11)

## Seed

Planet **Aurelis**, a sphere of radius 6,000 km.

| Reagent site | Latitude | Longitude |
|---|---|---|
| Moonwell Spring | 10 | 20 |
| Ember Crag | 25 | 60 |
| Frostcap Hollow | 85 | 0 |
| Ashen Reach | -84 | 90 |
| Tidewrack Shoals | 0 | 180 |
| Thornveil Glade | -30 | -45 |
| Glimmerdeep Caves | 40 | -100 |
| Saltmarsh Cairn | -15 | 130 |

| Seeded expedition | From | To | Owner |
|---|---|---|---|
| Spring to Crag | Moonwell Spring | Ember Crag | shop |
| Glade to Caves | Thornveil Glade | Glimmerdeep Caves | shop |
| Cairn to Shoals | Saltmarsh Cairn | Tidewrack Shoals | shop |

Adventurers for the tests: Elara (saved: Moonwell Spring to Frostcap Hollow); Borin (saved: Ember Crag to Saltmarsh Cairn).

## Criteria

### AC1 (rule R1, observable via both)

- **Given** the shop publishes the eight reagent sites of Aurelis
- **When** an adventurer opens the reagent map
- **Then** all eight sites are listed, and each is shown within 1 km of its published position

### AC2 (rule R1, observable via ui)

- **Given** Frostcap Hollow lies near the north pole and Tidewrack Shoals on the 180th meridian
- **When** the adventurer turns the globe to each of them
- **Then** each is shown at its position, neither clipped nor repeated

### AC3 (rule R2, observable via both)

- **Given** no adventurer is signed in
- **When** the reagent map is opened
- **Then** exactly the three seeded expeditions are shown, each with both ends within 1 km of its seeded sites

### AC4 (rule R2, observable via both)

- **Given** Elara is signed in and has saved one expedition, and Borin has saved another
- **When** Elara opens the reagent map
- **Then** the three seeded expeditions and Elara's expedition are shown, and Borin's is not

### AC5 (rule R3, observable via ui)

- **Given** the default view, in which Tidewrack Shoals is on the far side of the globe
- **When** the adventurer flies to Tidewrack Shoals
- **Then** Tidewrack Shoals is on screen, on the near side of the globe

### AC6 (rule R4, observable via ui)

- **Given** the adventurer has started a line
- **When** they place a point at Moonwell Spring and a second point at Ember Crag
- **Then** a completed line joins the two points, and it can be saved

### AC7 (rule R4, observable via ui)

- **Given** the adventurer has started a line and placed one point
- **When** no second point has been placed
- **Then** no line is complete and it cannot be saved

### AC8 (rule R4, observable via ui)

- **Given** a completed line from Thornveil Glade to Glimmerdeep Caves
- **When** the line is shown
- **Then** it follows the shortest path over the surface of Aurelis between them, not a straight line on a flat map

### AC9 (rule R5, observable via ui)

- **Given** the adventurer has started a line
- **When** they try to place a point on space beside the globe
- **Then** no point is placed and the line tool is still waiting for its first point

### AC10 (rule R5, observable via ui)

- **Given** the adventurer has started a line
- **When** they try to place a point on the edge of the globe, where the land faces away
- **Then** no point is placed

### AC11 (rule R6, observable via both)

- **Given** Elara is signed in
- **When** she draws a line from Moonwell Spring to Ember Crag, saves it, and opens the reagent map again
- **Then** her new expedition is shown, and the shop holds it for Elara with both ends within 1 km of where she placed them

### AC12 (rule R7, observable via both)

- **Given** no adventurer is signed in and a completed line is on the map
- **When** they try to save it
- **Then** they are asked to sign in, and the shop holds no new expedition

### AC13 (rule R7, observable via both)

- **Given** Elara is signed in and has placed both points less than 1 km apart
- **When** she tries to save the line
- **Then** the save is refused as a zero-length line, and the shop holds no new expedition

### AC14 (rule R8, observable via both)

- **Given** Elara has placed one point of a line
- **When** she starts a new line
- **Then** the new line has no points, and the shop holds no new expedition

### AC15 (rule R8, observable via both)

- **Given** Elara has drawn a line without saving it
- **When** she leaves the reagent map and opens it again
- **Then** the unsaved line is not shown, and the shop holds no new expedition

## Traceability

| Criterion | Rule | Events | Scenarios |
|---|---|---|---|
| AC1 | R1 | ReagentMapOpened | reagent-sites.feature: The shop publishes every reagent site (api)<br>reagent-sites.feature: Every reagent site is shown at its position (ui) |
| AC2 | R1 | ReagentMapOpened | reagent-sites.feature: Sites near a pole and on the 180th meridian are shown once, where they are (ui) |
| AC3 | R2 | ReagentMapOpened | expeditions-visibility.feature: Signed out, only the seeded expeditions are visible (api)<br>expeditions-visibility.feature: Seeded expeditions are drawn where they were seeded (ui) |
| AC4 | R2 | ReagentMapOpened | expeditions-visibility.feature: An adventurer sees their own expeditions and not another adventurer's (api)<br>expeditions-visibility.feature: A signed-in adventurer's own expedition is drawn with the seeded ones (ui) |
| AC5 | R3 | SiteFocused | fly-to-site.feature: Flying to a site on the far side brings it into view (ui) |
| AC6 | R4 | LineStarted, PointPlaced | draw-line.feature: Two points on the land make a line that can be saved (ui) |
| AC7 | R4 | LineStarted, PointPlaced | draw-line.feature: One point is not a line (ui) |
| AC8 | R4 | PointPlaced | draw-line.feature: A line follows the shortest path over the surface (ui) |
| AC9 | R5 | PointRejected | draw-line.feature: A point on space beside the globe is not placed (ui) |
| AC10 | R5 | PointRejected | draw-line.feature: A point on the edge of the globe, where the land faces away, is not placed (ui) |
| AC11 | R6 | ExpeditionSaved | save-expedition.feature: A saved expedition is shown again and held by the shop where it was drawn (ui)<br>save-expedition.feature: The shop keeps an expedition for the adventurer who saved it (api) |
| AC12 | R7 | ExpeditionRejected | save-expedition.feature: Saving while signed out is refused (api)<br>save-expedition.feature: A signed-out adventurer is asked to sign in when saving (ui) |
| AC13 | R7 | ExpeditionRejected | save-expedition.feature: A zero-length expedition is refused (api)<br>save-expedition.feature: Placing both points at the same place cannot be saved (ui) |
| AC14 | R8 | LineStarted | save-expedition.feature: Starting a new line abandons the unfinished one (ui) |
| AC15 | R8 | LineStarted | save-expedition.feature: An unsaved line is not kept after leaving the map (ui) |

## spec-critic

Deterministic checks (`docs/specs/spec-critic.py`): 0 findings. Shown able to fail: a rule with its negative example removed and a `Then` reading "the system processes the sites" both fail it.
All 21 scenarios parse with `@cucumber/gherkin` 28. `gherkin-lint` was not run (not installed in the repo).

Model critique, round 1:

- AC8 was marked observable via both; the shop holds only the two endpoints, so whether the line is geodesic is visible only in how it is drawn. → AC8 observableVia: ui
- AC1 'listed' depends on H5 (sites listed beside the globe); kept, since H5 was accepted. → none
- No criterion says what a signed-out adventurer sees while drawing (H3: they can draw). AC12 covers the save; drawing itself is covered by AC6 regardless of sign-in. → none
