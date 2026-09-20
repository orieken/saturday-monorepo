import { test } from '@orieken/saturday-playwright-leaflet';
import { expect } from '@playwright/test';

test.describe('Ye Olde Magic Shop - Map Locations', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/locations');
  });

  test('map should load and center on the correct region', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    const center = await leaflet.getCenter();
    expect(center.lat).toBeCloseTo(45.0, 1);
    expect(center.lng).toBeCloseTo(-120.0, 1);
  });

  test('markers should exist and popups should display correct information', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    const markers = await leaflet.getMarkers();
    expect(markers.length).toBeGreaterThanOrEqual(3);
    
    // Check Waterdeep Branch
    await leaflet.assertPopupContent('Waterdeep Branch', 'Potions of Healing');
  });

  test('should be able to draw a line and find markers along it', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    const route: Array<[number, number]> = [
      [42.0, -119.5], // Baldur's Gate
      [45.0, -120.0], // Waterdeep
      [47.5, -121.0]  // Neverwinter
    ];
    
    await leaflet.drawLine(route, { color: 'red', weight: 4 });
    
    // The buffer is in meters. 100km buffer.
    const markers = await leaflet.getMarkersAlongLine(route, 100000);
    expect(markers.length).toBe(3);
    
    const names = markers.map(m => m.metadata?.name);
    expect(names).toContain('Waterdeep Branch');
    expect(names).toContain('Baldur\'s Gate Branch');
    expect(names).toContain('Neverwinter Branch');
  });
  
  test('should check point in polygon', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    // Check a point inside the Sword Coast territory
    // Polygon bounds: [-122, 41] to [-118, 48]
    const inside = await leaflet.isPointInPolygon(45.0, -120.0);
    expect(inside).toBe(true);

    // Check a point outside
    const outside = await leaflet.isPointInPolygon(50.0, -115.0);
    expect(outside).toBe(false);
  });

  test('should verify bounds and visibility', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    // Zoom in heavily to Waterdeep
    await leaflet.zoomMap(10);
    
    const bounds = await leaflet.getBounds();
    
    // Waterdeep is at 45.0, -120.0
    expect(45.0).toBeGreaterThan(bounds._southWest.lat);
    expect(45.0).toBeLessThan(bounds._northEast.lat);
    expect(-120.0).toBeGreaterThan(bounds._southWest.lng);
    expect(-120.0).toBeLessThan(bounds._northEast.lng);
    
    // Baldur's Gate (42.0, -119.5) should be outside the bounds at zoom level 10
    expect(42.0).toBeLessThan(bounds._southWest.lat); 
  });

  test('should use geolocation to find nearest shop', async ({ page, context, leaflet }) => {
    await context.grantPermissions(['geolocation']);
    // Mock user location near Neverwinter (47.5, -121.0)
    await context.setGeolocation({ latitude: 47.6, longitude: -121.1 });
    
    await leaflet.waitForMapLoad();
    
    await page.getByRole('button', { name: 'Find Nearest Shop' }).click();
    
    // Give it a moment to pan
    await page.waitForTimeout(1000);
    
    const center = await leaflet.getCenter();
    // It should have panned to the spoofed location
    expect(center.lat).toBeCloseTo(47.6, 1);
    expect(center.lng).toBeCloseTo(-121.1, 1);
  });

  test('should not have broken tiles', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    const hasFailures = await leaflet.checkTileFailures();
    expect(hasFailures).toBe(false);
  });

  test('should calculate distance between locations', async ({ leaflet }) => {
    await leaflet.waitForMapLoad();
    
    // Waterdeep: [45.0, -120.0]
    // Baldur's Gate: [42.0, -119.5]
    const distance = await leaflet.distanceBetween([45.0, -120.0], [42.0, -119.5]);
    
    // Should be roughly 336 kilometers (very rough due to projection, but Turf computes it mathematically)
    // We just verify it's a positive number greater than 0
    expect(distance).toBeGreaterThan(300);
    expect(distance).toBeLessThan(400);
  });

  test('should be able to toggle layer visibility', async ({ page, leaflet }) => {
    await leaflet.waitForMapLoad();
    
    // Initially should be visible
    const isVisibleBefore = await leaflet.isLayerVisible('The Sword Coast');
    expect(isVisibleBefore).toBe(true);
    
    // Toggle layer control via evaluate to avoid Playwright hover/interception issues
    await page.evaluate(() => {
      const labels = Array.from(document.querySelectorAll('.leaflet-control-layers-overlays label'));
      const targetLabel = labels.find(label => label.textContent && label.textContent.includes('The Sword Coast'));
      const checkbox = targetLabel?.querySelector('input');
      if (checkbox && checkbox.checked) {
        checkbox.click();
      }
    });
    
    // Verify it is no longer visible
    const isVisibleAfter = await leaflet.isLayerVisible('The Sword Coast');
    expect(isVisibleAfter).toBe(false);
  });
});
