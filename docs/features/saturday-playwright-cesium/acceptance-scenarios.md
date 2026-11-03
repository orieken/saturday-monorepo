---
feature: "saturday-playwright-cesium"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: saturday-playwright-cesium

## Fixture page exposes the Cesium seam for test discovery

Verifies: An app under test exposes window.__SATURDAY_CESIUM__ = { viewer, Cesium } and attaches saturdayMetadata (id, name, kind) to each entity it wants tests to find, as documented in the package's own README seam contract.

```gherkin
Scenario: Fixture page exposes the Cesium seam for test discovery
  Given an app under test has loaded a page with a Cesium Viewer
  When the app follows the package's seam contract to attach the viewer and tag entities
  Then window.__SATURDAY_CESIUM__ exposes the viewer and the Cesium namespace
  And every entity the app wants tests to find carries saturdayMetadata with an id, a name, and a kind
```

## Fixture page renders a self-contained globe with known tagged entities

Verifies: packages/saturday-playwright-cesium/fixtures/ serves a static page with a Cesium Viewer on a spherical ellipsoid of a declared radius, one equirectangular image covering the whole globe from the package's own files, baseLayerPicker/geocoder/ion imagery off, known point/line/circle entities carrying metadata, and makes no request off its own origin.

```gherkin
Scenario: Fixture page renders a self-contained globe with known tagged entities
  Given the fixture's static page is requested
  When the page finishes loading
  Then the globe uses a spherical ellipsoid of the fixture's declared radius
  And the globe's whole surface is covered by one equirectangular image served from the package's own files
  And baseLayerPicker, geocoder, and ion imagery are all off
  And a point, a line, and a circle entity each carry saturdayMetadata, plus one entity whose kind is not a recognized value
  And one additional point entity carries no saturdayMetadata
  And no request made while loading the page leaves the page's own origin
```

## waitForSettled resolves once tiles, camera, and render are stable

Verifies: cesium.waitForSettled() resolves once the globe's tiles are loaded, the camera has stopped, and no render is pending, and rejects with a message naming the stuck condition (tiles/camera/render/seam) after a timeout.

```gherkin
Scenario: waitForSettled resolves once tiles, camera, and render are stable
  Given the fixture page has loaded and the Cesium seam is present
  When the globe finishes loading tiles, the camera stops moving, and no render is pending, and the test calls cesium.waitForSettled()
  Then the call resolves
```

## waitForSettled rejects naming the camera as the stuck condition

Verifies: cesium.waitForSettled() resolves once the globe's tiles are loaded, the camera has stopped, and no render is pending, and rejects with a message naming the stuck condition (tiles/camera/render/seam) after a timeout.

```gherkin
Scenario: waitForSettled rejects naming the camera as the stuck condition
  Given the fixture page has loaded and the Cesium seam is present
  And the camera is forced to appear still moving
  When the test calls cesium.waitForSettled() and waits past the timeout
  Then the call rejects with a message that names the camera as the condition preventing settlement
```

## getEntities returns metadata and geometry for every tagged entity

Verifies: cesium.getEntities() returns each saturdayMetadata-tagged entity's metadata and geometry in degrees and metres: point position, line endpoints, circle centre and radius.

```gherkin
Scenario: getEntities returns metadata and geometry for every tagged entity
  Given the fixture page has loaded with its tagged point, line, circle, and unrecognized-kind entities, plus one untagged point
  When the test calls cesium.getEntities()
  Then the result includes each tagged entity's saturdayMetadata
  And the point entity's position is reported in degrees
  And the line entity's endpoints are reported in degrees
  And the circle entity's centre is reported in degrees and its radius in metres
  And the entity whose kind is not a recognized value is still included with its geometry
  And the untagged entity is not included in the result
```

## toScreen returns page coordinates for a point facing the camera

Verifies: cesium.toScreen(lat, lng) returns page coordinates, or null when the point is on the far side of the globe or off-screen.

```gherkin
Scenario: toScreen returns page coordinates for a point facing the camera
  Given the fixture page has loaded with the camera looking at a known location
  When the test calls cesium.toScreen(lat, lng) for a point facing the camera
  Then it returns the point's page coordinates
```

## toScreen returns null for a point on the globe's far side

Verifies: cesium.toScreen(lat, lng) returns page coordinates, or null when the point is on the far side of the globe or off-screen.

```gherkin
Scenario: toScreen returns null for a point on the globe's far side
  Given the fixture page has loaded
  When the test calls cesium.toScreen(lat, lng) for a point on the far side of the globe
  Then it returns null
```

## toScreen returns null for a point that is visible but off-canvas

Verifies: cesium.toScreen(lat, lng) returns page coordinates, or null when the point is on the far side of the globe or off-screen.

```gherkin
Scenario: toScreen returns null for a point that is visible but off-canvas
  Given the fixture page has loaded
  When the test calls cesium.toScreen(lat, lng) for a point that faces the camera but falls outside the canvas bounds
  Then it returns null
```

## clickAt clicks through to the correct location near the globe's centre

Verifies: cesium.clickAt(lat, lng) clicks through page.mouse, and a fixture test shows the globe's pick returned that location within the tolerance metersPerPixelAt gives.

```gherkin
Scenario: clickAt clicks through to the correct location near the globe's centre
  Given the fixture page has loaded and a pick handler is installed on the globe
  When the test calls cesium.clickAt(lat, lng) for a visible point near the centre of the view
  Then the pick handler reports a location within metersPerPixelAt(lat, lng) times two metres of the requested point
```

## clickAt clicks through to the correct location near the horizon

Verifies: cesium.clickAt(lat, lng) clicks through page.mouse, and a fixture test shows the globe's pick returned that location within the tolerance metersPerPixelAt gives.

```gherkin
Scenario: clickAt clicks through to the correct location near the horizon
  Given the fixture page has loaded and a pick handler is installed on the globe
  When the test calls cesium.clickAt(lat, lng) for a visible point near the horizon
  Then the pick handler reports a location within metersPerPixelAt(lat, lng) times two metres of the requested point
```

## flyTo moves the camera to the requested destination and waits until settled

Verifies: cesium.flyTo(lat, lng, heightMeters) moves the camera and waits until settled.

```gherkin
Scenario: flyTo moves the camera to the requested destination and waits until settled
  Given the fixture page has loaded
  When the test calls cesium.flyTo(lat, lng, heightMeters)
  Then the camera's resulting latitude, longitude, and height match the requested destination within tolerance
  And the call does not resolve until the globe reports settled
```

## Geometry helpers use the globe's own radius instead of assuming Earth's

Verifies: distance, lineLength and point-in-circle checks take the globe's radius from the seam rather than assuming Earth's, and are unit-tested in Node without a browser; whether they live in this package or a shared saturday-geo is decided at the design gate.

```gherkin
Scenario: Geometry helpers use the globe's own radius instead of assuming Earth's
  Given a declared globe radius of 6,000,000 metres, differing from Earth's 6,371,000 metres
  When distance, lineLength, and point-in-circle are computed for the same coordinates once with the declared globe radius and once with Earth's radius, in a Node environment with no browser
  Then the two radii produce different results, showing each helper used the radius it was given rather than an assumed Earth radius
```

## The fixture suite passes in real Chrome via the published docker command

Verifies: The package's test script runs the fixture suite in real Chrome (channel: 'chrome', headed, --use-gl=angle --use-angle=swiftshader) and passes with the literal AC9 docker command after dependencies are installed and the package is built; record its output in the delivery record.

```gherkin
Scenario: The fixture suite passes in real Chrome via the published docker command
  Given dependencies are installed and the package is built
  When the docker command `docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-cesium test'` is run
  Then both the Vitest suite and the Playwright chrome-project suite pass with no network access
```

## Each helper's test fails when its proof value changes

Verifies: Each helper's test fails when the fixture's corresponding value changes; for a helper reading nothing from the fixture, a change to the helper's own code that the test catches is the proof instead. Record, per helper, the change made and the failing output in the delivery record.

```gherkin
Scenario: Each helper's test fails when its proof value changes
  Given a helper under test and, where one exists, the fixture value its test depends on
  When the fixture's corresponding value is deliberately changed, or, for a helper reading no fixed fixture value, the helper's own code is deliberately changed, and the test is re-run
  Then the test fails, showing it can detect a real regression, and passes again once the change is reverted
```

## The fixture server binds the loopback interface only

Verifies: The fixture's static server listens on 127.0.0.1 only.

```gherkin
Scenario: The fixture server binds the loopback interface only
  Given the fixture's static server is started
  When the server begins listening
  Then it binds to 127.0.0.1 and does not accept connections on other interfaces
```

## Changed code stays within complexity, coverage, and scope limits

Verifies: Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail, and no file outside scope.paths changes except the run's own delivery record under docs/features/saturday-playwright-cesium/.

```gherkin
Scenario: Changed code stays within complexity, coverage, and scope limits
  Given the implementation for this feature is complete
  When complexity analysis and coverage tooling run against the changed files, and the changed file list is compared against scope.paths
  Then every new or changed function has cyclomatic complexity below 7
  And every changed line is covered by a test proven able to fail
  And no file outside scope.paths changed, except the feature's own delivery record under docs/features/saturday-playwright-cesium/
```
