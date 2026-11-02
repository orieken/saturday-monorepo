// Builds the known Leaflet map every fixture.spec.ts test asserts against.
// Markers, polygon and layer-control name here are the "known values" the
// README-helper tests check — change one and its test must fail (see
// shared/rules/test-repair-contract.md's "prove it can fail").

var FIXTURE_MARKERS = [
  { id: 'marker-a', name: 'Marker A', lat: 45.0, lng: -120.0, popup: 'Fixture content A' },
  { id: 'marker-b', name: 'Marker B', lat: 42.0, lng: -119.5, popup: 'Fixture content B' },
  { id: 'marker-c', name: 'Marker C', lat: 47.5, lng: -121.0, popup: 'Fixture content C' }
];

var FIXTURE_POLYGON = [
  [41, -122],
  [41, -118],
  [48, -118],
  [48, -122]
];

var FIXTURE_OVERLAY_NAME = 'Fixture Overlay';

function addFixtureMarkers(map) {
  FIXTURE_MARKERS.forEach(function (definition) {
    var marker = L.marker([definition.lat, definition.lng]);
    marker.saturdayMetadata = { id: definition.id, name: definition.name };
    marker.bindPopup(definition.popup);
    marker.addTo(map);
  });
}

function addFixtureOverlay(map) {
  var overlay = L.layerGroup([L.circle([45.0, -120.0], { radius: 1000 })]);
  overlay.addTo(map);
  var overlays = {};
  overlays[FIXTURE_OVERLAY_NAME] = overlay;
  L.control.layers(null, overlays).addTo(map);
}

function buildFixtureMap() {
  var map = L.map('map', { zoomControl: false }).setView([45.0, -120.0], 5);
  L.tileLayer('/assets/tile.png', { attribution: '', minZoom: 0, maxZoom: 18 }).addTo(map);
  addFixtureMarkers(map);
  L.polygon(FIXTURE_POLYGON).addTo(map);
  addFixtureOverlay(map);
  return map;
}

window.__SATURDAY_LEAFLET_MAP__ = buildFixtureMap();
