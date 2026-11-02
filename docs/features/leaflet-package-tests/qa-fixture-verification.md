# QA Verification: Fixture-Break-and-Revert Sweep (AC3) and AC4 Docker Run

This is the delivery record required by AC3 and AC4 of `leaflet-package-tests`: per-helper proof
that each test can fail, and the captured output of the literal AC4 docker command. Every edit
below was made, run against a single targeted test, confirmed failing, then reverted to the exact
original text before moving to the next helper. `git status` / `git diff` after the full sweep
shows no residual change to any tracked or untracked file.

## AC4 — `docker run ... saturday-runner:debian ... pnpm test` (real Chrome, `--network none`)

Command:
```
docker run --rm --platform linux/amd64 --network none -v "$PWD:/work" saturday-runner:debian bash -lc 'cd /work && pnpm --filter @orieken/saturday-playwright-leaflet test'
```

Output:
```
Openbox-Message: Unable to find a valid menu file "/var/lib/openbox/debian-menu.xml"

> @orieken/saturday-playwright-leaflet@0.1.0 test /work/packages/saturday-playwright-leaflet
> playwright test --project=chrome


Running 11 tests using 7 workers

  ✓   3 [chrome] › tests/fixture.spec.ts:13:7 › LeafletFixture against the local fixture page › waitForMapLoad resolves once the fixture map is exposed and its tiles are ready (3.5s)
  ✓   6 [chrome] › tests/fixture.spec.ts:38:7 › LeafletFixture against the local fixture page › getMarkers returns the fixture's three known markers with their saturdayMetadata (3.5s)
  ✓   2 [chrome] › tests/fixture.spec.ts:22:7 › LeafletFixture against the local fixture page › getCenter returns the fixture map's known center (3.6s)
  ✓   5 [chrome] › tests/fixture.spec.ts:90:7 › LeafletFixture against the local fixture page › distanceBetween returns the known geodesic distance between two fixture markers (3.7s)
  ✓   4 [chrome] › tests/fixture.spec.ts:63:7 › LeafletFixture against the local fixture page › drawLine adds one more interactive path to the map (3.6s)
  ✓   7 [chrome] › tests/fixture.spec.ts:29:7 › LeafletFixture against the local fixture page › getBounds returns a viewport that contains Marker A (3.7s)
  ✓   1 [chrome] › tests/fixture.spec.ts:76:7 › LeafletFixture against the local fixture page › getMarkersAlongLine returns only the markers within the buffer of a known line (3.8s)
  ✓   8 [chrome] › tests/fixture.spec.ts:100:7 › LeafletFixture against the local fixture page › isPointInPolygon matches the fixture's known polygon bounds (1.9s)
  ✓   9 [chrome] › tests/fixture.spec.ts:107:7 › LeafletFixture against the local fixture page › assertPopupContent opens a known marker's popup and matches its fixture text (1.9s)
  ✓  10 [chrome] › tests/fixture.spec.ts:114:7 › LeafletFixture against the local fixture page › isLayerVisible reflects the fixture's known overlay toggle state (1.9s)
  ✓  11 [chrome] › tests/fixture.spec.ts:128:7 › LeafletFixture against the local fixture page › checkTileFailures returns false when the fixture's local tile image loads successfully (1.9s)

  11 passed (10.6s)
```

Run against: dependencies already installed, package already built (`pnpm --filter @orieken/saturday-playwright-leaflet build`) before the container run, per AC4's precondition. No network was available inside the container (`--network none`); the fixture page is served entirely from `fixtures/serve.mjs` against local files, so this is the proof the package needs no network access at all.

## AC3 — per-helper fixture-break-and-revert sweep

All 11 README helpers below. For 9 of them the break is a value baked into
`fixtures/assets/fixture-map.js` that the helper reads back from the page. For 2
(`drawLine`, `distanceBetween`) the helper operates purely on caller-supplied coordinates and
reads nothing from the fixture file — see the note at the end of this section for why those two
are proven via a targeted implementation mutation instead, and why that is the honest answer
rather than a fabricated fixture dependency.

### 1. `waitForMapLoad`
**Change**: `fixtures/assets/fixture-map.js` — `window.__SATURDAY_LEAFLET_MAP__ = buildFixtureMap();` → `buildFixtureMap();` (map built but never exposed on `window`).
**Result**: FAILED.
```
Error: page.evaluate: Map not found
    at ../src/fixture.ts:96
```
**Reverted**: restored the original assignment line exactly.

### 2. `getCenter`
**Change**: `fixtures/assets/fixture-map.js` — `setView([45.0, -120.0], 5)` → `setView([10.0, -10.0], 5)`.
**Result**: FAILED.
```
Error: expect(received).toBeCloseTo(expected, precision)
Expected: 45
Received: 10
```
**Reverted**: restored `setView([45.0, -120.0], 5)`.

### 3. `getBounds`
**Change**: same edit as #2 (`setView([45.0, -120.0], 5)` → `setView([10.0, -10.0], 5)`), run against the `getBounds` test alone.
**Result**: FAILED.
```
Error: expect(received).toBeGreaterThan(expected)
Expected: > 45
Received:   25.045792240303445
```
**Reverted**: restored `setView([45.0, -120.0], 5)`.

### 4. `getMarkers`
**Change**: `fixtures/assets/fixture-map.js` — Marker A's `lat: 45.0` → `lat: 44.0`.
**Result**: FAILED.
```
Error: expect(received).toMatchObject(expected)
- "lat": 45,
+ "lat": 44,
```
**Reverted**: restored `lat: 45.0`.

### 5. `drawLine` — see note below (implementation mutation, not a fixture value)
**Change**: `src/fixture.ts` — `L.polyline(points, options).addTo(map);` → `L.polyline(points, options);` (polyline created but never attached to the map).
**Result**: FAILED.
```
Error: expect(received).toBe(expected)
Expected: 3
Received: 2
```
**Reverted**: restored `.addTo(map)`.

### 6. `getMarkersAlongLine`
**Change**: `fixtures/assets/fixture-map.js` — Marker C moved from `{ lat: 47.5, lng: -121.0 }` to `{ lat: 10.0, lng: -60.0 }` (outside the test's buffer).
**Result**: FAILED.
```
Error: expect(received).toEqual(expected) // deep equality
- Expected  - 1
+ Received  + 0
  Array [
    "Marker A",
-   "Marker C",
  ]
```
**Reverted**: restored Marker C's original coordinates.

### 7. `distanceBetween` — see note below (implementation mutation, not a fixture value)
**Change**: `src/fixture.ts` — `turf.distance(p1, p2, { units: 'kilometers' })` → `{ units: 'miles' }`.
**Result**: FAILED.
```
Error: expect(received).toBeCloseTo(expected, precision)
Expected: 336.012
Received: 208.78826840899717
```
**Reverted**: restored `units: 'kilometers'`.

### 8. `isPointInPolygon`
**Change**: `fixtures/assets/fixture-map.js` — `FIXTURE_POLYGON` moved from spanning lat `[41,48]`/lng `[-122,-118]` to spanning lat `[1,5]`/lng `[-2,2]` (far from the test's probe point).
**Result**: FAILED.
```
Error: expect(received).toBe(expected)
Expected: true
Received: false
```
**Reverted**: restored the original `FIXTURE_POLYGON` coordinates.

### 9. `assertPopupContent`
**Change**: `fixtures/assets/fixture-map.js` — Marker A's popup text `'Fixture content A'` → `'DIFFERENT TEXT'`.
**Result**: FAILED.
```
Error: Popup content did not contain "Fixture content A". Actual content: "DIFFERENT TEXT"
```
**Reverted**: restored `'Fixture content A'`.

### 10. `isLayerVisible`
**Change**: `fixtures/assets/fixture-map.js` — `FIXTURE_OVERLAY_NAME` renamed from `'Fixture Overlay'` to `'Renamed Overlay'` (test still looks up the old literal name).
**Result**: FAILED.
```
Error: expect(received).toBe(expected)
Expected: true
Received: false
```
**Reverted**: restored `'Fixture Overlay'`.

### 11. `checkTileFailures`
**Change**: `fixtures/assets/fixture-map.js` — tile layer URL `'/assets/tile.png'` → `'/assets/tile-missing.png'` (404s against `fixtures/serve.mjs`).
**Result**: FAILED.
```
Error: expect(received).toBe(expected)
Expected: false
Received: true
```
**Reverted**: restored `'/assets/tile.png'`.

### Note on `drawLine` and `distanceBetween`

Both helpers take coordinates directly from the caller and never read
`window.__SATURDAY_LEAFLET_MAP__`'s rendered state or any fixture-sourced value — `distanceBetween`
calls `turf.distance` on the two points its own caller passed in; `drawLine` draws a polyline from
the points its caller passed in. There is no line in `fixtures/assets/fixture-map.js` whose change
would alter either helper's output, so a literal "change the fixture's corresponding value" break
does not exist for these two. Rather than invent a fictitious fixture dependency, each was proven
able to fail by mutating the helper's own implementation in `src/fixture.ts` — `drawLine`'s
`.addTo(map)` and `distanceBetween`'s `units` — which demonstrates the same thing the fixture-break
sweep demonstrates for the other 9: the assertion is live and will catch a real regression. This is
flagged here, rather than silently presented as an equivalent fixture-value change, per
`shared/rules/test-repair-contract.md`'s requirement to state what evidence backs a "proves it can
fail" claim.

## Full-suite confirmation after the sweep

After every edit above was reverted, the full local suite was re-run to confirm no residual state:

```
Running 11 tests using 7 workers
...
  11 passed (4.4s)
```

`git status` / `git diff` against `src/fixture.ts` show only the developer's original Turf-resolution
and tile-failure-detection changes (see `implementation-notes.md`) — no leftover `BROKEN-FOR-QA`
markers from this sweep. `fixtures/` is untracked, so its content was verified by direct read-back
against the original text captured at the start of this sweep; it matches exactly.
