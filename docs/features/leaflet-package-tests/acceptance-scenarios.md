---
feature: "leaflet-package-tests"
bounded_context: ""
domain_terms: ["saturdayMetadata", "Site-Centric pattern"]
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: leaflet-package-tests

## Fixture page exposes a fully loaded map with known markers, polygon, and layer control

Verifies: Given packages/saturday-playwright-leaflet/fixtures/, a static page MUST load Leaflet from the package's own dependencies (no CDN), a local image, and a known set of markers carrying saturdayMetadata, a polygon, and a layer control, and MUST expose window.__SATURDAY_LEAFLET_MAP__.

```gherkin
Scenario: Fixture page exposes a fully loaded map with known markers, polygon, and layer control
  Given the fixture page at packages/saturday-playwright-leaflet/fixtures/index.html is served locally with no network access
  And Leaflet is referenced only from the package's own node_modules/leaflet/dist, never a CDN URL
  When the fixture page finishes loading the map, its tile image, and all layers
  Then window.__SATURDAY_LEAFLET_MAP__ is set
  And at least 3 markers exist, each carrying saturdayMetadata with an id and a name
  And one L.polygon exists with fixed, known coordinates
  And one L.control.layers overlay exists with a fixed, known name
```

## Every README-documented helper has at least one test asserting a known fixture value

Verifies: Given the package's test script, it MUST serve the fixture page and run a Playwright suite with at least one test per helper documented in README.md, and each test MUST assert a known value taken from the fixture.

```gherkin
Scenario: Every README-documented helper has at least one test asserting a known fixture value
  Given the fixture page is served by the package's test script
  When the package's Playwright test suite runs
  Then each of the 11 README helpers — waitForMapLoad, getCenter, getBounds, getMarkers, drawLine, getMarkersAlongLine, distanceBetween, isPointInPolygon, assertPopupContent, isLayerVisible, checkTileFailures — has at least one test
  And each test's assertion matches a known value taken directly from the fixture
```

## Each helper's test fails when its corresponding fixture value is deliberately changed

Verifies: Given each helper's test, it MUST fail when the fixture's corresponding value is changed; the delivery record MUST document, per helper, the fixture change made and the resulting failing output.

```gherkin
Scenario: Each helper's test fails when its corresponding fixture value is deliberately changed
  Given a helper's test passes against the unmodified fixture
  When the fixture value that helper's test depends on is deliberately changed, then that single test is re-run
  Then the test fails
  And the delivery record documents, for each of the 11 helpers, the exact fixture change made and the resulting failing output
```

## Suite passes in real Chrome with no network access inside the runner image

Verifies: Given the suite, it MUST run in real Chrome (channel: 'chrome') inside the runner image with no network, and MUST pass with the literal docker command below after dependencies are installed and the package is built; the delivery record MUST capture that command's output.

```gherkin
Scenario: Suite passes in real Chrome with no network access inside the runner image
  Given package dependencies are installed and the package is built
  When the literal command is run: docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-leaflet test'
  Then the suite runs using real Chrome (channel: 'chrome') with no network access available
  And the suite passes
  And the delivery record captures that command's full output
```

## Change stays within scope and meets complexity and coverage bars

Verifies: Given the change is complete, every new or changed function MUST have cyclomatic complexity below 7, every changed line MUST be covered by a test that can fail, and no file outside scope.paths MUST change except the run's own delivery record under docs/features/leaflet-package-tests/.

```gherkin
Scenario: Change stays within scope and meets complexity and coverage bars
  Given the feature implementation is complete
  When complexity, test coverage, and the set of changed files are checked
  Then every new or changed function has cyclomatic complexity below 7
  And every changed line is covered by a test proven able to fail
  And no file outside packages/saturday-playwright-leaflet/**, pnpm-lock.yaml, and docs/features/leaflet-package-tests/** changed
```
