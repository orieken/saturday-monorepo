// Pure, radius-parameterized geometry math. Zero imports from 'cesium' or
// '@playwright/test' by design — see architecture-notes.md's structural
// decision to keep this package-local rather than a new packages/saturday-geo,
// written so it stays a cheap extraction later if a second consumer appears.
//
// The fixture globe (and Aurelis) is declared as a spherical ellipsoid with no
// flattening, so these are spherical (haversine-style) formulas taking a single
// radiusMeters value, not a full oblate-ellipsoid geodesic algorithm.

export interface GeoPoint {
  lat: number;
  lng: number;
}

function toRadians(degrees: number): number {
  return (degrees * Math.PI) / 180;
}

/**
 * Great-circle distance between two points, in metres, on a sphere of the
 * given radius. Callers must pass the globe's own radius — never a
 * hardcoded Earth radius (see architecture-notes.md).
 */
export function distance(a: GeoPoint, b: GeoPoint, radiusMeters: number): number {
  const lat1 = toRadians(a.lat);
  const lat2 = toRadians(b.lat);
  const halfDeltaLat = toRadians(b.lat - a.lat) / 2;
  const halfDeltaLng = toRadians(b.lng - a.lng) / 2;
  const h =
    Math.sin(halfDeltaLat) ** 2 +
    Math.cos(lat1) * Math.cos(lat2) * Math.sin(halfDeltaLng) ** 2;
  return 2 * radiusMeters * Math.asin(Math.sqrt(h));
}

/**
 * Total length, in metres, of the polyline formed by consecutive points.
 */
export function lineLength(points: GeoPoint[], radiusMeters: number): number {
  let total = 0;
  for (let index = 1; index < points.length; index++) {
    total += distance(points[index - 1], points[index], radiusMeters);
  }
  return total;
}

/**
 * Whether `point` falls within `circleRadiusMeters` of `center`, measured as
 * great-circle distance on a globe of `globeRadiusMeters`.
 */
export function pointInCircle(
  point: GeoPoint,
  center: GeoPoint,
  circleRadiusMeters: number,
  globeRadiusMeters: number,
): boolean {
  return distance(point, center, globeRadiusMeters) <= circleRadiusMeters;
}
