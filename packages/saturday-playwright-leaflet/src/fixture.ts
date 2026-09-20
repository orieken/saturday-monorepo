import { test as base } from '@playwright/test';
import type { Page } from '@playwright/test';
import type * as L from 'leaflet';

export interface LeafletFixture {
  /**
   * Wait for all tiles to finish loading and the map to be ready.
   */
  waitForMapLoad(): Promise<void>;
  
  /**
   * Gets the current center of the map.
   */
  getCenter(): Promise<{ lat: number; lng: number }>;
  
  /**
   * Gets all markers on the map, including their coordinates and any attached metadata.
   */
  getMarkers(): Promise<Array<{ lat: number; lng: number; metadata?: any }>>;
  
  /**
   * Draws a polyline on the map between the given points.
   */
  drawLine(points: Array<[number, number]>, options?: any): Promise<void>;
  
  /**
   * Gets all markers that are within a certain buffer distance of a polyline.
   */
  getMarkersAlongLine(points: Array<[number, number]>, bufferMeters?: number): Promise<Array<{ lat: number; lng: number; metadata?: any }>>;

  /**
   * Click a marker (by metadata id or name) and assert its popup text contains the given string.
   */
  assertPopupContent(markerIdOrName: string, expectedText: string): Promise<void>;

  /**
   * Simulate dragging the map by a certain pixel offset.
   */
  dragMap(pixelOffsetX: number, pixelOffsetY: number): Promise<void>;

  /**
   * Zoom the map to a specific level.
   */
  zoomMap(level: number): Promise<void>;

  /**
   * Verify if a given lat/lng point falls within any rendered GeoJSON polygon on the map.
   */
  isPointInPolygon(lat: number, lng: number): Promise<boolean>;

  /**
   * Check if any tile requests have failed with 404s.
   */
  checkTileFailures(): Promise<boolean>;

  /**
   * Gets the current viewport bounds of the map.
   */
  getBounds(): Promise<{ _southWest: { lat: number; lng: number }; _northEast: { lat: number; lng: number } }>;

  /**
   * Calculates the distance in kilometers between two lat/lng points.
   */
  distanceBetween(point1: [number, number], point2: [number, number]): Promise<number>;

  /**
   * Checks if a specific layer (e.g., from a layer control) is currently visible on the map.
   */
  isLayerVisible(layerName: string): Promise<boolean>;
}

export const test = base.extend<{ leaflet: LeafletFixture }>({
  leaflet: async ({ page }, use) => {
    
    // Inject Turf.js for spatial operations like point-in-polygon and line buffers
    await page.addScriptTag({ url: 'https://cdn.jsdelivr.net/npm/@turf/turf@6/turf.min.js' });

    const getMap = () => {
      const map = (window as any).__SATURDAY_LEAFLET_MAP__;
      if (!map) throw new Error('Leaflet map not found on window.__SATURDAY_LEAFLET_MAP__');
      return map as L.Map;
    };

    const fixture: LeafletFixture = {
      async waitForMapLoad() {
        await page.evaluate(() => {
          return new Promise<void>((resolve, reject) => {
            const map = (window as any).__SATURDAY_LEAFLET_MAP__;
            if (!map) return reject('Map not found');
            
            // Wait for tiles to load if they are currently loading
            let isReady = true;
            map.eachLayer((layer: any) => {
              if (layer.isLoading && layer.isLoading()) {
                isReady = false;
                layer.once('load', resolve);
              }
            });
            if (isReady) resolve();
          });
        });
      },

      async getCenter() {
        return page.evaluate(() => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const center = map.getCenter();
          return { lat: center.lat, lng: center.lng };
        });
      },

      async getMarkers() {
        return page.evaluate(() => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const markers: any[] = [];
          map.eachLayer((layer: any) => {
            if (layer instanceof (window as any).L.Marker) {
              const pos = layer.getLatLng();
              markers.push({
                lat: pos.lat,
                lng: pos.lng,
                metadata: layer.saturdayMetadata || null
              });
            }
          });
          return markers;
        });
      },

      async drawLine(points, options = { color: 'blue', weight: 3 }) {
        await page.evaluate(({ points, options }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const L = (window as any).L;
          L.polyline(points, options).addTo(map);
        }, { points, options });
      },

      async getMarkersAlongLine(points, bufferMeters = 50000) {
        return page.evaluate(({ points, bufferMeters }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const L = (window as any).L;
          const turf = (window as any).turf;
          
          // Create Turf line string
          const line = turf.lineString(points.map((p: number[]) => [p[1], p[0]])); // turf uses [lng, lat]
          const buffered = turf.buffer(line, bufferMeters / 1000, { units: 'kilometers' });

          const markersInZone: any[] = [];
          map.eachLayer((layer: any) => {
            if (layer instanceof L.Marker) {
              const pos = layer.getLatLng();
              const pt = turf.point([pos.lng, pos.lat]);
              if (turf.booleanPointInPolygon(pt, buffered)) {
                markersInZone.push({
                  lat: pos.lat,
                  lng: pos.lng,
                  metadata: layer.saturdayMetadata || null
                });
              }
            }
          });
          return markersInZone;
        }, { points, bufferMeters });
      },

      async assertPopupContent(markerIdOrName, expectedText) {
        await page.evaluate(({ markerIdOrName }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const L = (window as any).L;
          let targetMarker: any = null;
          
          map.eachLayer((layer: any) => {
            if (layer instanceof L.Marker && layer.saturdayMetadata) {
              if (layer.saturdayMetadata.id === markerIdOrName || layer.saturdayMetadata.name === markerIdOrName) {
                targetMarker = layer;
              }
            }
          });
          
          if (!targetMarker) throw new Error(`Marker ${markerIdOrName} not found`);
          targetMarker.openPopup();
        }, { markerIdOrName });

        // Wait for the popup to appear in the DOM and check its text
        const popupLocator = page.locator('.leaflet-popup-content');
        await popupLocator.waitFor({ state: 'visible' });
        const content = await popupLocator.textContent();
        
        if (!content || !content.includes(expectedText)) {
          throw new Error(`Popup content did not contain "${expectedText}". Actual content: "${content}"`);
        }
      },

      async dragMap(pixelOffsetX, pixelOffsetY) {
        await page.evaluate(({ pixelOffsetX, pixelOffsetY }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          map.panBy([pixelOffsetX, pixelOffsetY]);
        }, { pixelOffsetX, pixelOffsetY });
        
        // Wait for pan animation to finish
        await page.waitForTimeout(500); 
      },

      async zoomMap(level) {
        await page.evaluate(({ level }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          map.setZoom(level);
        }, { level });
        
        await page.waitForTimeout(500);
      },

      async isPointInPolygon(lat, lng) {
        return page.evaluate(({ lat, lng }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const turf = (window as any).turf;
          const L = (window as any).L;
          
          let isInPolygon = false;
          const pt = turf.point([lng, lat]);

          map.eachLayer((layer: any) => {
            // Leaflet polygons
            if (layer instanceof L.Polygon) {
              const geojson = layer.toGeoJSON();
              if (turf.booleanPointInPolygon(pt, geojson)) {
                isInPolygon = true;
              }
            }
          });
          return isInPolygon;
        }, { lat, lng });
      },

      async checkTileFailures() {
        // Playwright can intercept network requests, but doing it post-hoc here 
        // requires we had setup an interceptor earlier.
        // Instead, we can check if Leaflet has any error-class tiles in the DOM
        return page.evaluate(() => {
          const errorTiles = document.querySelectorAll('.leaflet-tile-loaded.leaflet-tile-error');
          return errorTiles.length > 0;
        });
      },

      async getBounds() {
        return page.evaluate(() => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          const bounds = map.getBounds();
          return {
            _southWest: { lat: bounds._southWest.lat, lng: bounds._southWest.lng },
            _northEast: { lat: bounds._northEast.lat, lng: bounds._northEast.lng }
          };
        });
      },

      async distanceBetween(point1, point2) {
        return page.evaluate(({ point1, point2 }) => {
          const turf = (window as any).turf;
          // turf expects [lng, lat]
          const p1 = turf.point([point1[1], point1[0]]);
          const p2 = turf.point([point2[1], point2[0]]);
          return turf.distance(p1, p2, { units: 'kilometers' });
        }, { point1, point2 });
      },

      async isLayerVisible(layerName) {
        return page.evaluate(({ layerName }) => {
          const map = (window as any).__SATURDAY_LEAFLET_MAP__;
          let isVisible = false;
          
          // One way to check is to look at the layer control in the DOM if we are using the default L.control.layers
          // But a more robust way is to check the map's layers, though Leaflet doesn't natively tag layers with their control names.
          // In our implementation, we bind popups to the territory layer, and we can check the popup content or layer options.
          // Since it's tricky to map arbitrary layer names, we can check the DOM for the layer control inputs.
          const labels = document.querySelectorAll('.leaflet-control-layers-overlays label');
          labels.forEach((label) => {
            if (label.textContent && label.textContent.trim() === layerName) {
              const input = label.querySelector('input');
              if (input && input.checked) {
                isVisible = true;
              }
            }
          });
          return isVisible;
        }, { layerName });
      }
    };

    await use(fixture);
  }
});
