<template>
  <div class="locations-view">
    <div class="mb-8 border-b border-dnd-slate pb-4 flex justify-between items-end">
      <div>
        <h2 class="text-3xl font-magical text-dnd-parchment">Shop Locations</h2>
        <p class="text-gray-400 mt-2">Find a Ye Olde Magic Shop branch near your adventuring party.</p>
      </div>
      <button 
        @click="findNearestShop"
        class="bg-dnd-gold text-dnd-dark font-bold py-2 px-4 rounded hover:bg-yellow-600 transition"
      >
        Find Nearest Shop
      </button>
    </div>
    <div id="map" class="h-[600px] w-full rounded-lg border-2 border-dnd-slate shadow-xl z-0"></div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue';
import L from 'leaflet';

const map = ref<L.Map | null>(null);

// Custom icons
const magicShopIcon = L.divIcon({
  html: `<div style="font-size: 24px; text-shadow: 0 0 5px #ff00ff;">🧙‍♂️</div>`,
  className: 'custom-magic-icon',
  iconSize: [30, 30],
  iconAnchor: [15, 15],
  popupAnchor: [0, -15]
});

const shops = [
  {
    name: 'Waterdeep Branch',
    coords: [45.0, -120.0] as L.LatLngExpression,
    inventory: 'Potions of Healing, Scrolls of Fireball',
    rumor: 'The shopkeeper is actually a polymorphed dragon.'
  },
  {
    name: 'Baldur\'s Gate Branch',
    coords: [42.0, -119.5] as L.LatLngExpression,
    inventory: 'Poisons, Daggers of Venom, Cloaks of Elvenkind',
    rumor: 'Ask for the "special reserve" under the counter.'
  },
  {
    name: 'Neverwinter Branch',
    coords: [47.5, -121.0] as L.LatLngExpression,
    inventory: 'Wands of Magic Missiles, Rings of Protection',
    rumor: 'Prices are high due to the recent goblin raids.'
  }
];

const findNearestShop = () => {
  if (navigator.geolocation && map.value) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        // Just panning to the user's location for simplicity in this demo,
        // though "Find Nearest Shop" implies we might calculate distance.
        const lat = position.coords.latitude;
        const lng = position.coords.longitude;
        map.value!.flyTo([lat, lng], 6);
      },
      (error) => {
        console.error("Geolocation error:", error);
      }
    );
  }
};

onMounted(() => {
  // Initialize the map
  map.value = L.map('map').setView([45.0, -120.0], 5);

  // Expose for playwright testing
  (window as any).__SATURDAY_LEAFLET_MAP__ = map.value;

  // Use OSM with a CSS class we can style to look like parchment
  const tileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    className: 'map-tiles-parchment'
  }).addTo(map.value);

  // Add a GeoJSON layer for a magical territory
  const swordCoastGeoJSON = {
    type: 'FeatureCollection',
    features: [
      {
        type: 'Feature',
        properties: { name: 'The Sword Coast' },
        geometry: {
          type: 'Polygon',
          coordinates: [
            [
              [-122.0, 48.0],
              [-118.0, 48.0],
              [-118.0, 41.0],
              [-122.0, 41.0],
              [-122.0, 48.0]
            ]
          ]
        }
      }
    ]
  };

  const territoryLayer = L.geoJSON(swordCoastGeoJSON as any, {
    style: {
      color: '#8b0000',
      weight: 2,
      opacity: 0.6,
      fillColor: '#d2b48c',
      fillOpacity: 0.3
    }
  }).bindPopup('<b>The Sword Coast</b><br>Danger Level: High');

  territoryLayer.addTo(map.value);

  const shopLayer = L.layerGroup();

  // Add shop markers
  shops.forEach(shop => {
    const marker = L.marker(shop.coords, { icon: magicShopIcon });
    
    // Attach metadata for Playwright
    (marker as any).saturdayMetadata = {
      id: shop.name.toLowerCase().replace(/[\s']/g, '-'),
      name: shop.name
    };

    marker.bindPopup(`
      <div class="font-sans text-black">
        <h3 class="font-bold text-lg mb-1">${shop.name}</h3>
        <p class="text-sm mb-2"><b>Inventory:</b> ${shop.inventory}</p>
        <p class="text-sm italic text-gray-600">Rumor: ${shop.rumor}</p>
      </div>
    `);

    marker.addTo(shopLayer);
  });

  shopLayer.addTo(map.value);

  // Add Layer Controls
  const overlayMaps = {
    "Shops": shopLayer,
    "The Sword Coast": territoryLayer
  };

  L.control.layers(undefined, overlayMaps).addTo(map.value);
});

onBeforeUnmount(() => {
  if (map.value) {
    map.value.remove();
    (window as any).__SATURDAY_LEAFLET_MAP__ = undefined;
  }
});
</script>

<style>
/* Fantasy parchment style map tiles */
.map-tiles-parchment {
  filter: sepia(0.8) hue-rotate(-15deg) saturate(1.2) contrast(0.9) brightness(0.9);
}

.custom-magic-icon {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
