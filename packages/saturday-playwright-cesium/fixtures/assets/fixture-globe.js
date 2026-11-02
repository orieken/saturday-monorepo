// Builds the known Cesium globe every fixture.spec.ts test asserts against:
// an invented planet, Aurelis, radius 6,000 km (not WGS84 Earth). The
// entities and camera values below are the "known values" the README-helper
// tests check — change one and its test must fail (see
// shared/rules/test-repair-contract.md's "prove it can fail").
(function () {
  var AURELIS_RADIUS_METERS = 6000000;
  var CAMERA_HEIGHT_METERS = AURELIS_RADIUS_METERS / 2;

  var FIXTURE_ENTITIES = [
    {
      id: 'reagent-site-1',
      name: 'Reagent Site 1',
      kind: 'point',
      lat: 5,
      lng: 10,
    },
    {
      id: 'supply-route-1',
      name: 'Supply Route 1',
      kind: 'line',
      points: [
        { lat: 5, lng: 10 },
        { lat: -10, lng: 30 },
      ],
    },
    {
      id: 'exclusion-zone-1',
      name: 'Exclusion Zone 1',
      kind: 'circle',
      center: { lat: 0, lng: 0 },
      radiusMeters: 400000,
    },
    {
      // Deliberately not 'point'/'line'/'circle' — proves getEntities()
      // dispatches on the entity's own Cesium shape, never on this string
      // (see architecture-notes.md's structural decision on saturdayMetadata.kind).
      id: 'mystery-site-1',
      name: 'Mystery Site',
      kind: 'unexpected-kind-value',
      lat: -5,
      lng: -5,
    },
  ];

  var UNTAGGED_POINT = { lat: 2, lng: 2 };

  function addPointEntity(viewer, ellipsoid, definition) {
    var entity = viewer.entities.add({
      position: Cesium.Cartesian3.fromDegrees(definition.lng, definition.lat, 0, ellipsoid),
      point: { pixelSize: 10 },
    });
    entity.saturdayMetadata = { id: definition.id, name: definition.name, kind: definition.kind };
  }

  function addLineEntity(viewer, ellipsoid, definition) {
    var positions = definition.points.map(function (point) {
      return Cesium.Cartesian3.fromDegrees(point.lng, point.lat, 0, ellipsoid);
    });
    var entity = viewer.entities.add({
      polyline: { positions: positions, width: 2 },
    });
    entity.saturdayMetadata = { id: definition.id, name: definition.name, kind: definition.kind };
  }

  function addCircleEntity(viewer, ellipsoid, definition) {
    var entity = viewer.entities.add({
      position: Cesium.Cartesian3.fromDegrees(definition.center.lng, definition.center.lat, 0, ellipsoid),
      ellipse: {
        semiMajorAxis: definition.radiusMeters,
        semiMinorAxis: definition.radiusMeters,
      },
    });
    entity.saturdayMetadata = { id: definition.id, name: definition.name, kind: definition.kind };
  }

  function addFixtureEntity(viewer, ellipsoid, definition) {
    if (definition.kind === 'line') return addLineEntity(viewer, ellipsoid, definition);
    if (definition.kind === 'circle') return addCircleEntity(viewer, ellipsoid, definition);
    return addPointEntity(viewer, ellipsoid, definition);
  }

  function addUntaggedEntity(viewer, ellipsoid) {
    // No saturdayMetadata attached — getEntities() must exclude this one.
    viewer.entities.add({
      position: Cesium.Cartesian3.fromDegrees(UNTAGGED_POINT.lng, UNTAGGED_POINT.lat, 0, ellipsoid),
      point: { pixelSize: 10 },
    });
  }

  async function buildFixtureViewer() {
    var ellipsoid = new Cesium.Ellipsoid(AURELIS_RADIUS_METERS, AURELIS_RADIUS_METERS, AURELIS_RADIUS_METERS);
    var globe = new Cesium.Globe(ellipsoid);

    var imageryProvider = await Cesium.SingleTileImageryProvider.fromUrl('/assets/aurelis.jpg', {
      rectangle: Cesium.Rectangle.fromDegrees(-180, -90, 180, 90),
      ellipsoid: ellipsoid,
    });

    var viewer = new Cesium.Viewer('cesiumContainer', {
      globe: globe,
      baseLayerPicker: false,
      geocoder: false,
      baseLayer: new Cesium.ImageryLayer(imageryProvider),
      animation: false,
      timeline: false,
      fullscreenButton: false,
      homeButton: false,
      sceneModePicker: false,
      navigationHelpButton: false,
      infoBox: false,
      selectionIndicator: false,
      skyBox: false,
      skyAtmosphere: false,
      shouldAnimate: false,
    });

    FIXTURE_ENTITIES.forEach(function (definition) {
      addFixtureEntity(viewer, ellipsoid, definition);
    });
    addUntaggedEntity(viewer, ellipsoid);

    viewer.camera.setView({
      destination: Cesium.Cartesian3.fromDegrees(0, 0, CAMERA_HEIGHT_METERS, ellipsoid),
      orientation: {
        heading: Cesium.Math.toRadians(0),
        pitch: Cesium.Math.toRadians(-90),
        roll: 0,
      },
    });

    window.__SATURDAY_CESIUM__ = { viewer: viewer, Cesium: Cesium };
  }

  buildFixtureViewer();
})();
