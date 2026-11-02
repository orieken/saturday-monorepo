import { describe, expect, it } from 'vitest';
import { distance, lineLength, pointInCircle } from '../src/geometry';
import type { GeoPoint } from '../src/geometry';

// Aurelis, the fixture globe's declared planet, per docs/features/saturday-playwright-cesium.
const AURELIS_RADIUS_METERS = 6_000_000;
// Earth's mean radius — must never be the value these checks actually use.
const EARTH_RADIUS_METERS = 6_371_000;

describe('distance', () => {
  const equator: GeoPoint = { lat: 0, lng: 0 };
  const quarterWay: GeoPoint = { lat: 0, lng: 90 };

  it('returns a quarter of the great circle for two points 90 degrees apart on the equator', () => {
    const expected = (Math.PI / 2) * AURELIS_RADIUS_METERS;
    expect(distance(equator, quarterWay, AURELIS_RADIUS_METERS)).toBeCloseTo(expected, 3);
  });

  it('fails if Earth’s radius is substituted for the globe’s own radius', () => {
    const aurelisResult = distance(equator, quarterWay, AURELIS_RADIUS_METERS);
    const earthResult = distance(equator, quarterWay, EARTH_RADIUS_METERS);
    expect(aurelisResult).not.toBe(earthResult);
    expect(earthResult).toBeCloseTo((Math.PI / 2) * EARTH_RADIUS_METERS, 3);
  });

  it('returns zero for two identical points', () => {
    expect(distance(equator, { ...equator }, AURELIS_RADIUS_METERS)).toBeCloseTo(0, 9);
  });
});

describe('lineLength', () => {
  it('sums the great-circle distance of each consecutive segment', () => {
    const points: GeoPoint[] = [
      { lat: 0, lng: 0 },
      { lat: 0, lng: 90 },
      { lat: 90, lng: 90 },
    ];
    const expectedFirstLeg = (Math.PI / 2) * AURELIS_RADIUS_METERS;
    const expectedSecondLeg = (Math.PI / 2) * AURELIS_RADIUS_METERS;
    expect(lineLength(points, AURELIS_RADIUS_METERS)).toBeCloseTo(
      expectedFirstLeg + expectedSecondLeg,
      3,
    );
  });

  it('returns zero for a single point', () => {
    expect(lineLength([{ lat: 10, lng: 20 }], AURELIS_RADIUS_METERS)).toBe(0);
  });
});

describe('pointInCircle', () => {
  const center: GeoPoint = { lat: 0, lng: 0 };
  // Central angle of 0.08 radians along the equator (central angle equals the
  // longitude delta exactly when both points sit on lat = 0).
  const point: GeoPoint = { lat: 0, lng: (0.08 * 180) / Math.PI };
  const circleRadiusMeters = 500_000;

  it('is true when the globe’s own (smaller) radius keeps the point inside the circle', () => {
    expect(pointInCircle(point, center, circleRadiusMeters, AURELIS_RADIUS_METERS)).toBe(true);
  });

  it('is false when Earth’s (larger) radius pushes the same point outside the circle', () => {
    expect(pointInCircle(point, center, circleRadiusMeters, EARTH_RADIUS_METERS)).toBe(false);
  });
});
