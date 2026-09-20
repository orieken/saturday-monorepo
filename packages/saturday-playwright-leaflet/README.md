# @orieken/saturday-playwright-leaflet

A Playwright testing fixture designed to interrogate Leaflet maps. This plugin bypasses the typical flakiness of testing maps via the DOM (like trying to click SVG paths or canvas elements) and instead exposes a robust API to directly query Leaflet's internal state using `Turf.js` for spatial calculations.

## 1. Application Setup (Wiring the Map)

For this plugin to work, the application under test must expose the instantiated Leaflet map to the `window` object and optionally attach metadata to markers.

### Expose the Map Instance
When initializing your map, attach it to `window.__SATURDAY_LEAFLET_MAP__`:

```typescript
import L from 'leaflet';

onMounted(() => {
  const map = L.map('map-container').setView([45.0, -120.0], 5);
  
  // Expose map for Playwright testing
  if (window) {
    (window as any).__SATURDAY_LEAFLET_MAP__ = map;
  }
});

onBeforeUnmount(() => {
  if (window) {
    (window as any).__SATURDAY_LEAFLET_MAP__ = undefined;
  }
});
```

### Attach Metadata to Markers
If you want to assert against specific markers or target them by name, attach a `saturdayMetadata` object to the marker instance before adding it to the map.

```typescript
const marker = L.marker([45.0, -120.0]);

// Attach test metadata
(marker as any).saturdayMetadata = {
  id: 'waterdeep-branch',
  name: 'Waterdeep Branch'
};

marker.addTo(map);
```

## 2. Test Setup (Wiring the Plugin)

Import the extended `test` object from this package instead of directly from `@playwright/test`. This injects the `leaflet` fixture into your tests.

```typescript
// import { test } from '@playwright/test'; // <-- Replace this
import { test } from '@orieken/saturday-playwright-leaflet';
import { expect } from '@playwright/test';

test('my map test', async ({ page, leaflet }) => {
  await page.goto('/locations');
  
  // Wait for all tiles and layers to finish loading
  await leaflet.waitForMapLoad();
  
  // Now use the leaflet helpers!
});
```

## Available Helpers

- `leaflet.waitForMapLoad()`: Waits for the map to exist and all tiles to finish downloading.
- `leaflet.getCenter()`: Returns `{ lat, lng }` of the current map center.
- `leaflet.getBounds()`: Returns the NorthEast and SouthWest coordinates of the current viewport.
- `leaflet.getMarkers()`: Returns an array of all markers on the map, including their coordinates and attached `saturdayMetadata`.
- `leaflet.drawLine(points)`: Draws a mock polyline on the map.
- `leaflet.getMarkersAlongLine(points, bufferMeters)`: Uses Turf.js to create a buffer polygon around a line and returns all markers found inside that zone.
- `leaflet.distanceBetween(pt1, pt2)`: Returns the distance in kilometers between two coordinates.
- `leaflet.isPointInPolygon(lat, lng)`: Verifies if a given coordinate falls within any rendered GeoJSON polygon on the map.
- `leaflet.assertPopupContent(markerId, text)`: Programmatically opens a marker's popup and asserts the DOM content contains the text.
- `leaflet.isLayerVisible(layerName)`: Interrogates the Leaflet Layer Control to see if a specific layer is toggled on.
- `leaflet.checkTileFailures()`: Returns true if any map tiles failed to load (HTTP 404/500).
