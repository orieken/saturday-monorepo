import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const apiBase = env.VITE_API_BASE || 'http://localhost:9001';

  return {
    plugins: [
      vue(),
      tailwindcss(),
    ],
    server: {
      port: 9000,
      proxy: {
        // Proxy API calls to the backend in dev. Frontend should call 
        // relative paths like "/api/..." so this kicks in automatically.
        '/api': {
          target: apiBase,
          changeOrigin: true
        },
        '/reports': {
          target: apiBase,
          changeOrigin: true
        }
      }
    }
  };
});
