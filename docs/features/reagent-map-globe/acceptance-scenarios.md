---
feature: "reagent-map-globe"
bounded_context: ""
domain_terms: []
files_touched: []
issue_refs: []
linked_adrs: []
linked_kis: []
---

# Acceptance Scenarios: reagent-map-globe

## Reagent map route and nav link appear when the feature flag is explicitly true

Verifies: The /reagent-map route and its nav link exist only when VITE_FEATURE_REAGENT_MAP is exactly "true"; with it unset, empty, or any other value, there is no route and no link, and the app is otherwise unchanged. .env sets the flag to false.

```gherkin
Scenario: Reagent map route and nav link appear when the feature flag is explicitly true
  Given VITE_FEATURE_REAGENT_MAP is set to the literal string "true"
  When a visitor loads the application
  Then the /reagent-map route exists
  And the nav link to the reagent map is visible
```

## Reagent map route and nav link are absent when the feature flag is off, unset, or any other value

Verifies: The /reagent-map route and its nav link exist only when VITE_FEATURE_REAGENT_MAP is exactly "true"; with it unset, empty, or any other value, there is no route and no link, and the app is otherwise unchanged. .env sets the flag to false.

```gherkin
Scenario: Reagent map route and nav link are absent when the feature flag is off, unset, or any other value
  Given VITE_FEATURE_REAGENT_MAP is unset, empty, "false", or any value other than the literal string "true"
  When a visitor loads the application
  Then there is no /reagent-map route
  And there is no nav link to the reagent map
  And the rest of the application is unchanged
```

## Reagent map globe renders entirely from the app's own files with no external network dependency

Verifies: The view runs Cesium with no network dependency: no ion token or ion imagery, baseLayerPicker and geocoder off, one equirectangular image of Aurelis served from the app's own files, Cesium's own static assets served same-origin via CESIUM_BASE_URL, scene mode SCENE3D, and an ellipsoid that is a sphere of radius 6,000 km.

```gherkin
Scenario: Reagent map globe renders entirely from the app's own files with no external network dependency
  Given the feature flag is enabled
  And a visitor opens the reagent map view
  When the globe finishes loading
  Then the globe uses no Cesium Ion token and no Cesium Ion imagery
  And the base layer picker and geocoder controls are off
  And the globe surface shows one equirectangular image of Aurelis served from the app's own files
  And Cesium's own static assets load from the app's own origin
  And the globe is shown in 3D scene mode as a sphere with a 6,000 km radius
```

## Reagent map exposes the Saturday test seam and tags each site with its metadata

Verifies: The view sets window.__SATURDAY_CESIUM__ = { viewer, Cesium } and attaches saturdayMetadata { id, name, kind: 'site' } to every reagent site entity, per packages/saturday-playwright-cesium/README.md's seam contract.

```gherkin
Scenario: Reagent map exposes the Saturday test seam and tags each site with its metadata
  Given the feature flag is enabled
  And a visitor opens the reagent map view
  When the globe finishes loading
  Then window.__SATURDAY_CESIUM__ exposes the viewer and the Cesium library
  And every reagent site entity carries saturdayMetadata with its id, name, and kind "site"
```

## All eight reagent sites are shown exactly once at their published positions

Verifies: Every one of the eight reagent sites is shown within 1 km of its published position and listed by name beside the globe, including Frostcap Hollow near the north pole and Tidewrack Shoals on the 180th meridian, each shown exactly once with no clipping or duplication.

```gherkin
Scenario: All eight reagent sites are shown exactly once at their published positions
  Given the feature flag is enabled
  And a visitor opens the reagent map view
  When the globe finishes loading
  Then all eight reagent sites are shown, each within 1 km of its published position
  And each site appears exactly once, with no clipping or duplication
  And each site's name is listed beside the globe
```

## Sites at the pole and the antimeridian are shown once, not duplicated

Verifies: Every one of the eight reagent sites is shown within 1 km of its published position and listed by name beside the globe, including Frostcap Hollow near the north pole and Tidewrack Shoals on the 180th meridian, each shown exactly once with no clipping or duplication.

```gherkin
Scenario: Sites at the pole and the antimeridian are shown once, not duplicated
  Given the feature flag is enabled
  And a visitor opens the reagent map view
  And Frostcap Hollow is published near the north pole at latitude 85, longitude 0
  And Tidewrack Shoals is published on the 180th meridian at latitude 0, longitude 180
  When the globe finishes loading
  Then Frostcap Hollow appears exactly once, within 1 km of its published position
  And Tidewrack Shoals appears exactly once, within 1 km of its published position, with no duplicate entity at the antimeridian wrap
```

## No request leaves the app's own origin or the mock API's origin while the reagent map view is open

Verifies: No request made while the reagent map view is open leaves the app's own origin or the mock API's own origin.

```gherkin
Scenario: No request leaves the app's own origin or the mock API's origin while the reagent map view is open
  Given the feature flag is enabled
  And a visitor opens the reagent map view
  When the view is open and in use
  Then every network request made stays on the app's own origin or the mock API's own origin
  And no request is made to any other origin
```

## The feature flag decision is testable without a browser

Verifies: A small pure function decides whether the /reagent-map route exists, tested for "true", "false", and unset without a browser. scripts/ui-suite-in-chrome.sh starts the app with VITE_FEATURE_REAGENT_MAP=true so the pre-ship check exercises the view, while .env keeps the flag false. New specs named tests/reagent-map-*.spec.ts run only in the chrome Playwright project; npm run test:e2e (bundled Chromium) still passes without them.

```gherkin
Scenario: The feature flag decision is testable without a browser
  Given a small pure function decides whether the /reagent-map route exists
  When that function is called with "true", "false", and an unset value
  Then it returns true only for the literal string "true", and false for "false" and unset, without needing a browser
```

## New reagent-map specs run only under the chrome Playwright project, leaving the default suite unaffected

Verifies: A small pure function decides whether the /reagent-map route exists, tested for "true", "false", and unset without a browser. scripts/ui-suite-in-chrome.sh starts the app with VITE_FEATURE_REAGENT_MAP=true so the pre-ship check exercises the view, while .env keeps the flag false. New specs named tests/reagent-map-*.spec.ts run only in the chrome Playwright project; npm run test:e2e (bundled Chromium) still passes without them.

```gherkin
Scenario: New reagent-map specs run only under the chrome Playwright project, leaving the default suite unaffected
  Given new spec files named tests/reagent-map-*.spec.ts exist
  And apps/ye-olde-magic-shop/.env keeps the feature flag false
  When npm run test:e2e runs the bundled-Chromium suite
  And scripts/ui-suite-in-chrome.sh starts the app with the feature flag set to true and runs the chrome project
  Then npm run test:e2e still passes without running the new reagent-map specs
  And the chrome project run starts the app with the flag on and exercises the reagent map view
```

## New and changed code stays within complexity, coverage, and scope limits

Verifies: Every new or changed function has cyclomatic complexity below 7, every changed line is covered by a test that can fail, and no file outside scope.paths changes except this run's own delivery record under docs/features/reagent-map-globe/.

```gherkin
Scenario: New and changed code stays within complexity, coverage, and scope limits
  Given the reagent map feature has been implemented
  When the change is reviewed against the delivery's definition of done
  Then every new or changed function has cyclomatic complexity below 7
  And every changed line is covered by a test that can fail
  And no file outside scope.paths changes, except this run's own delivery record under docs/features/reagent-map-globe/
```
