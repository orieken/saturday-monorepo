<template>
  <div class="reagent-map-view">
    <div class="mb-8 border-b border-dnd-slate pb-4">
      <h2 class="text-3xl font-magical text-dnd-parchment">Aurelis Reagent Map</h2>
      <p class="text-gray-400 mt-2">Every known reagent site on Aurelis, shown on a 3D globe.</p>
    </div>
    <div
      v-if="mapError"
      data-testid="reagent-map-error"
      role="alert"
      class="mb-4 rounded-lg border-2 border-red-700 bg-red-950 px-4 py-3 text-red-200"
    >
      {{ mapError }}
    </div>
    <div class="flex gap-6">
      <div
        id="cesiumContainer"
        role="img"
        aria-label="3D globe of Aurelis with every reagent site marked and every known expedition drawn as a line between its two sites; the full list of sites is in the adjacent list"
        class="h-[600px] w-full rounded-lg border-2 border-dnd-slate shadow-xl"
      ></div>
      <ul
        data-testid="reagent-site-list"
        aria-label="Reagent sites"
        class="w-64 shrink-0 text-dnd-parchment"
      >
        <li v-for="site in sites" :key="site.id" class="py-1 border-b border-dnd-slate">
          <button
            type="button"
            class="block w-full cursor-pointer text-left"
            :data-testid="`fly-to-site-${site.id}`"
            :aria-label="`Fly to ${site.name}`"
            @click="flyToSite(site)"
          >
            {{ site.name }}
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue';
import 'cesium/Build/Cesium/Widgets/widgets.css';
// Type-only: erased at compile time, so this never becomes a runtime import
// of the 'cesium' module (see loadCesiumGlobal's comment for why that
// matters) — it only gives us real types for the global Cesium exposes.
import type * as CesiumModule from 'cesium';
import { fetchReagentSites, type ReagentSite } from '../api/reagentSites';
import { fetchExpeditions, type Expedition } from '../api/expeditions';

declare global {
  interface Window {
    Cesium?: typeof CesiumModule;
    __SATURDAY_CESIUM__?: { viewer: CesiumModule.Viewer; Cesium: typeof CesiumModule };
  }
}

const AURELIS_RADIUS_METERS = 6_000_000;
const CESIUM_SCRIPT_URL = '/cesium/Cesium.js';
const SITE_MARKER_PIXEL_SIZE = 10;
// The spec's default view is centred over Moonwell Spring. Without this, Cesium
// falls back to its Earth-sized home view, far from the planet's seeded sites.
const DEFAULT_VIEW = { lat: 10, lng: 20, heightMeters: 12_000_000 };
const FLY_TO_HEIGHT_METERS = 3_000_000;

const sites = ref<ReagentSite[]>([]);
const expeditions = ref<Expedition[]>([]);
const mapError = ref<string | null>(null);
let viewer: CesiumModule.Viewer | null = null;
let cesium: typeof CesiumModule | null = null;
let ellipsoid: CesiumModule.Ellipsoid | null = null;

// Cesium's npm package re-exports its engine package's internals in a way
// Rollup cannot resolve for a production bundle (an upstream packaging
// mismatch between cesium@1.146.0 and its installed @cesium/engine).
// Loading the prebuilt, self-contained Build/Cesium/Cesium.js as a classic
// global script — the same file the Vite plugin already serves at
// /cesium/ — sidesteps bundling it and matches how fixture-globe.js's
// already-reviewed page loads Cesium.
function loadCesiumGlobal(): Promise<typeof CesiumModule> {
  if (window.Cesium) return Promise.resolve(window.Cesium);

  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = CESIUM_SCRIPT_URL;
    script.onload = () => resolve(window.Cesium as typeof CesiumModule);
    script.onerror = () => reject(new Error('Failed to load Cesium script'));
    document.head.appendChild(script);
  });
}

async function buildViewer() {
  const Cesium = await loadCesiumGlobal();
  const ellipsoid = new Cesium.Ellipsoid(AURELIS_RADIUS_METERS, AURELIS_RADIUS_METERS, AURELIS_RADIUS_METERS);
  const globe = new Cesium.Globe(ellipsoid);
  const imageryProvider = await Cesium.SingleTileImageryProvider.fromUrl('/aurelis.jpg', {
    rectangle: Cesium.Rectangle.fromDegrees(-180, -90, 180, 90),
    ellipsoid
  });

  const newViewer = new Cesium.Viewer('cesiumContainer', {
    // The scene's ellipsoid, not only the globe's. Without it the camera
    // controller grabs and zooms against Earth's ellipsoid, so a dragged
    // point drifts away from the cursor.
    ellipsoid,
    globe,
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
    shouldAnimate: false
  });
  newViewer.scene.mode = Cesium.SceneMode.SCENE3D;
  // No ion in this view, so its default credit logo would only mislead.
  (newViewer.cesiumWidget.creditContainer as HTMLElement).style.display = 'none';
  newViewer.camera.setView({
    destination: Cesium.Cartesian3.fromDegrees(DEFAULT_VIEW.lng, DEFAULT_VIEW.lat, DEFAULT_VIEW.heightMeters, ellipsoid)
  });

  return { Cesium, viewer: newViewer, ellipsoid };
}

function addSiteEntities(
  Cesium: typeof CesiumModule,
  cesiumViewer: CesiumModule.Viewer,
  ellipsoid: CesiumModule.Ellipsoid,
  reagentSites: ReagentSite[]
) {
  reagentSites.forEach((site) => {
    const entity = cesiumViewer.entities.add({
      position: Cesium.Cartesian3.fromDegrees(site.lng, site.lat, 0, ellipsoid),
      point: { pixelSize: SITE_MARKER_PIXEL_SIZE }
    });
    (entity as unknown as { saturdayMetadata: unknown }).saturdayMetadata = {
      id: site.id,
      name: site.name,
      kind: 'site'
    };
  });
}

function addExpeditionEntities(
  Cesium: typeof CesiumModule,
  cesiumViewer: CesiumModule.Viewer,
  cesiumEllipsoid: CesiumModule.Ellipsoid,
  items: Expedition[]
) {
  items.forEach((expedition) => {
    const entity = cesiumViewer.entities.add({
      polyline: {
        positions: Cesium.Cartesian3.fromDegreesArray(
          [expedition.from.lng, expedition.from.lat, expedition.to.lng, expedition.to.lat],
          cesiumEllipsoid
        ),
        arcType: Cesium.ArcType.GEODESIC
      }
    });
    (entity as unknown as { saturdayMetadata: unknown }).saturdayMetadata = {
      id: expedition.id,
      name: expedition.name ?? '',
      kind: 'expedition'
    };
  });
}

function flyToSite(site: ReagentSite) {
  if (!viewer || !cesium || !ellipsoid) return;
  viewer.camera.flyTo({
    destination: cesium.Cartesian3.fromDegrees(site.lng, site.lat, FLY_TO_HEIGHT_METERS, ellipsoid)
  });
}

onMounted(async () => {
  const [viewerResult, sitesResult, expeditionsResult] = await Promise.allSettled([
    buildViewer(),
    fetchReagentSites(),
    fetchExpeditions()
  ]);

  if (viewerResult.status === 'rejected') {
    throw viewerResult.reason;
  }

  const { Cesium, viewer: newViewer, ellipsoid: newEllipsoid } = viewerResult.value;

  if (sitesResult.status === 'rejected' || expeditionsResult.status === 'rejected') {
    newViewer.destroy();
    mapError.value = 'Failed to load reagent sites. Try refreshing the page.';
    throw sitesResult.status === 'rejected' ? sitesResult.reason : expeditionsResult.reason;
  }

  viewer = newViewer;
  cesium = Cesium;
  ellipsoid = newEllipsoid;
  sites.value = sitesResult.value;
  expeditions.value = expeditionsResult.value;
  addSiteEntities(Cesium, viewer, ellipsoid, sites.value);
  addExpeditionEntities(Cesium, viewer, ellipsoid, expeditions.value);

  window.__SATURDAY_CESIUM__ = { viewer, Cesium };
});

onBeforeUnmount(() => {
  if (viewer) {
    viewer.destroy();
    window.__SATURDAY_CESIUM__ = undefined;
  }
});
</script>
