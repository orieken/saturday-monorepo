---
feature: "leaflet-local-tiles"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: leaflet-local-tiles

## Locations view loads map imagery from its own origin only

Verifies: The locations view's map imagery MUST come only from the app's own origin; no tile or image request MUST leave the page.

```gherkin
Scenario: Locations view loads map imagery from its own origin only
  Given the locations view is open
  When the map finishes loading
  Then no network request is made to any origin other than the app's own
  And the map's imagery is served from the app's own build output, for example a same-origin request to /sword-coast-map.png
```

## Markers, territory polygon and popups keep their relative positions after the tile layer is replaced

Verifies: The existing markers, territory polygon and popups MUST keep the same position relative to each other after the tile layer is replaced.

```gherkin
Scenario: Markers, territory polygon and popups keep their relative positions after the tile layer is replaced
  Given the locations view has loaded with the local map imagery in place
  When the shop markers, the Sword Coast territory polygon, and their popups are inspected
  Then the Waterdeep Branch marker stays at [45.0, -120.0]
  And the Sword Coast polygon still spans [-122.0, 48.0] to [-118.0, 41.0]
  And each marker, polygon and popup sits in the same position relative to the others as it did before the imagery changed
```

## Map centre stays at the documented coordinates under the default CRS

Verifies: Every test in tests/map.spec.ts MUST pass unchanged under the default Leaflet CRS.

```gherkin
Scenario: Map centre stays at the documented coordinates under the default CRS
  When the locations view loads
  Then the map centre stays close to latitude 45.0 and longitude -120.0
```

## Zoom level 10 bounds include Waterdeep and exclude the southern boundary

Verifies: Every test in tests/map.spec.ts MUST pass unchanged under the default Leaflet CRS.

```gherkin
Scenario: Zoom level 10 bounds include Waterdeep and exclude the southern boundary
  When the map is zoomed to level 10
  Then the visible bounds include latitude 45.0, longitude -120.0
  And the visible bounds exclude latitude 42.0
```

## Geolocation pan moves the map centre to the user's location

Verifies: Every test in tests/map.spec.ts MUST pass unchanged under the default Leaflet CRS.

```gherkin
Scenario: Geolocation pan moves the map centre to the user's location
  Given the browser reports a geolocation position
  When the locations view requests the user's location
  Then the map centre pans to latitude 47.6, longitude -121.1
```

## Distance between Waterdeep and Baldur's Gate stays within the expected range

Verifies: Every test in tests/map.spec.ts MUST pass unchanged under the default Leaflet CRS.

```gherkin
Scenario: Distance between Waterdeep and Baldur's Gate stays within the expected range
  When the distance between the Waterdeep and Baldur's Gate markers is calculated
  Then the calculated distance stays between 300 and 400 km
```

## Locations view loads successfully with every off-origin request blocked

Verifies: A Playwright test that aborts every request not addressed to the app's own origin MUST show the locations view loading, leaflet.waitForMapLoad() resolving, and leaflet.checkTileFailures() returning false.

```gherkin
Scenario: Locations view loads successfully with every off-origin request blocked
  Given every network request whose origin is not the app's own origin is blocked
  When the locations view loads
  Then the page's map helper waitForMapLoad() resolves
  And the page's map helper checkTileFailures() returns false
```

## The map image's licence is recorded next to the asset

Verifies: The new map image's licence MUST be recorded next to it, either as an original in-repo creation or under a licence that permits redistribution.

```gherkin
Scenario: The map image's licence is recorded next to the asset
  Given the local map image exists in the repository
  When the image's licence information is checked
  Then a licence record sits next to the image
  And the record either marks the image as an original in-repo creation or names a licence that permits redistribution
```

## Changed code stays within complexity, coverage and scope limits

Verifies: Every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except the run's own delivery record under docs/features/leaflet-local-tiles/.

```gherkin
Scenario: Changed code stays within complexity, coverage and scope limits
  When the implementation for this feature is complete
  Then every new or changed function has cyclomatic complexity below 7
  And every changed line is covered by a test that can fail
  And no file outside scope.paths changes, except the delivery record under docs/features/leaflet-local-tiles/
```
