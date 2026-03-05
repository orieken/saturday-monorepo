import { createApp } from 'vue';
import { createPinia } from 'pinia';
import './style.css';
import App from './App.vue';
import router from './router';
import useCucumberStore from './stores/cucumber';

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);

// Load cucumber index on startup
const cucumber = useCucumberStore(pinia);
cucumber.load();

app.mount('#app');
