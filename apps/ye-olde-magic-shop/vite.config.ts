import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';

// Vite config that proxies /api to the local mock API on port 8001
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
  ],
  server: {
    port: 8000,
    proxy: {
      '/api': {
        target: 'http://localhost:8001',
        changeOrigin: true,
        secure: false
      },
      '/images': {
        target: 'http://localhost:8001',
        changeOrigin: true,
        secure: false
      }
    },
    allowedHosts: true
  },
});
