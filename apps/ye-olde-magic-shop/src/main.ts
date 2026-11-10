import { createApp } from 'vue';
import { createPinia } from 'pinia';
import router from './router';
import App from './App.vue';
import './styles.css';
import 'leaflet/dist/leaflet.css';
import '@fontsource/cinzel/400.css';
import '@fontsource/cinzel/700.css';
import '@fontsource/inter/400.css';
import '@fontsource/inter/600.css';

declare global {
  interface Window {
    CESIUM_BASE_URL?: string;
  }
}

// Read by the globally-loaded Cesium script (loaded via
// <script src="/cesium/Cesium.js"> in ReagentMapView.vue's loadCesiumGlobal)
// to resolve its own Workers/Assets/ThirdParty paths.
window.CESIUM_BASE_URL = '/cesium/';

const app = createApp(App);
app.use(createPinia());
app.use(router);
app.mount('#app');

