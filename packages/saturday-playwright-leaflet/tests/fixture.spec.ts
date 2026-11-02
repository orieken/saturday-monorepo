import { expect } from '@playwright/test';
import { test } from '../src/fixture';

// Every assertion below targets a literal, known value baked into
// fixtures/assets/fixture-map.js — change that value and the matching test
// here must fail (shared/rules/test-repair-contract.md's "prove it can fail").

test.describe('LeafletFixture against the local fixture page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('waitForMapLoad resolves once the fixture map is exposed and its tiles are ready', async ({
    leaflet,
    page,
  }) => {
    await leaflet.waitForMapLoad();
    const mapExists = await page.evaluate(() => Boolean((window as any).__SATURDAY_LEAFLET_MAP__));
    expect(mapExists).toBe(true);
  });

  test('getCenter returns the fixture map\'s known center', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    const center = await leaflet.getCenter();
    expect(center.lat).toBeCloseTo(45.0, 5);
    expect(center.lng).toBeCloseTo(-120.0, 5);
  });

  test('getBounds returns a viewport that contains Marker A', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    const bounds = await leaflet.getBounds();
    expect(bounds._southWest.lat).toBeLessThan(45.0);
    expect(bounds._northEast.lat).toBeGreaterThan(45.0);
    expect(bounds._southWest.lng).toBeLessThan(-120.0);
    expect(bounds._northEast.lng).toBeGreaterThan(-120.0);
  });

  test('getMarkers returns the fixture\'s three known markers with their saturdayMetadata', async ({
    leaflet,
  }) => {
    await leaflet.waitForMapLoad();
    const markers = await leaflet.getMarkers();
    expect(markers).toHaveLength(3);

    const byId = Object.fromEntries(markers.map((marker) => [marker.metadata?.id, marker]));
    expect(byId['marker-a']).toMatchObject({
      lat: 45.0,
      lng: -120.0,
      metadata: { id: 'marker-a', name: 'Marker A' },
    });
    expect(byId['marker-b']).toMatchObject({
      lat: 42.0,
      lng: -119.5,
      metadata: { id: 'marker-b', name: 'Marker B' },
    });
    expect(byId['marker-c']).toMatchObject({
      lat: 47.5,
      lng: -121.0,
      metadata: { id: 'marker-c', name: 'Marker C' },
    });
  });

  test('drawLine adds one more interactive path to the map', async ({ leaflet, page }) => {
    await leaflet.waitForMapLoad();
    const pathCountBefore = await page.locator('path.leaflet-interactive').count();

    await leaflet.drawLine([
      [45.0, -120.0],
      [47.5, -121.0],
    ]);

    const pathCountAfter = await page.locator('path.leaflet-interactive').count();
    expect(pathCountAfter).toBe(pathCountBefore + 1);
  });

  test('getMarkersAlongLine returns only the markers within the buffer of a known line', async ({
    leaflet,
  }) => {
    await leaflet.waitForMapLoad();
    const route: Array<[number, number]> = [
      [45.0, -120.0], // Marker A
      [47.5, -121.0], // Marker C
    ];

    const markers = await leaflet.getMarkersAlongLine(route, 100000);
    const names = markers.map((marker) => marker.metadata?.name).sort();
    expect(names).toEqual(['Marker A', 'Marker C']);
  });

  test('distanceBetween returns the known geodesic distance between two fixture markers', async ({
    leaflet,
  }) => {
    await leaflet.waitForMapLoad();
    // Marker A (45.0, -120.0) to Marker B (42.0, -119.5), computed with @turf/turf
    // (the same library the helper uses) ahead of time: ~336.012 km.
    const distance = await leaflet.distanceBetween([45.0, -120.0], [42.0, -119.5]);
    expect(distance).toBeCloseTo(336.012, 1);
  });

  test('isPointInPolygon matches the fixture\'s known polygon bounds', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    // Polygon spans lat [41, 48], lng [-122, -118].
    expect(await leaflet.isPointInPolygon(45.0, -120.0)).toBe(true);
    expect(await leaflet.isPointInPolygon(50.0, -115.0)).toBe(false);
  });

  test('assertPopupContent opens a known marker\'s popup and matches its fixture text', async ({
    leaflet,
  }) => {
    await leaflet.waitForMapLoad();
    await leaflet.assertPopupContent('marker-a', 'Fixture content A');
  });

  test('isLayerVisible reflects the fixture\'s known overlay toggle state', async ({ leaflet, page }) => {
    await leaflet.waitForMapLoad();
    expect(await leaflet.isLayerVisible('Fixture Overlay')).toBe(true);

    await page.evaluate(() => {
      const labels = Array.from(document.querySelectorAll('.leaflet-control-layers-overlays label'));
      const target = labels.find((label) => label.textContent?.includes('Fixture Overlay'));
      const checkbox = target?.querySelector('input');
      if (checkbox?.checked) checkbox.click();
    });

    expect(await leaflet.isLayerVisible('Fixture Overlay')).toBe(false);
  });

  test('checkTileFailures returns false when the fixture\'s local tile image loads successfully', async ({
    leaflet,
  }) => {
    await leaflet.waitForMapLoad();
    expect(await leaflet.checkTileFailures()).toBe(false);
  });
});
