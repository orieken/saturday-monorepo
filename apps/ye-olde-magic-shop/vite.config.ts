import { defineConfig, type Plugin } from 'vite';
import { cpSync, createReadStream, existsSync, statSync } from 'node:fs';
import { dirname, extname, join, normalize } from 'node:path';
import { fileURLToPath } from 'node:url';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';

const __dirname = dirname(fileURLToPath(import.meta.url));
const CESIUM_SOURCE_DIR = join(__dirname, 'node_modules', 'cesium', 'Build', 'Cesium');
const CESIUM_URL_PREFIX = '/cesium/';

const CONTENT_TYPES_BY_EXTENSION: Record<string, string> = {
  '.js': 'text/javascript',
  '.css': 'text/css',
  '.json': 'application/json',
  '.wasm': 'application/wasm',
  '.html': 'text/html',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.gif': 'image/gif',
  '.svg': 'image/svg+xml',
  '.ktx2': 'image/ktx2',
  '.glb': 'model/gltf-binary',
};

function resolveRequestedFile(urlPath: string): string | null {
  const relativePath = normalize(decodeURIComponent(urlPath.slice(CESIUM_URL_PREFIX.length)));
  if (relativePath.startsWith('..')) return null;
  const resolved = join(CESIUM_SOURCE_DIR, relativePath);
  return existsSync(resolved) && statSync(resolved).isFile() ? resolved : null;
}

// Serves Cesium's own prebuilt static assets (workers, Widgets CSS, Assets
// images) same-origin, without committing a second copy of them into
// public/ or into source control, and with no new dependency. See
// architecture-notes.md's structural decision 3 — scope.paths for this
// order excludes index.html and .gitignore, which rules out the two more
// conventional alternatives.
function cesiumStaticAssetsPlugin(): Plugin {
  return {
    name: 'cesium-static-assets',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        if (!req.url?.startsWith(CESIUM_URL_PREFIX)) return next();
        const filePath = resolveRequestedFile(req.url);
        if (!filePath) return next();
        res.setHeader('Content-Type', CONTENT_TYPES_BY_EXTENSION[extname(filePath)] ?? 'application/octet-stream');
        createReadStream(filePath).pipe(res);
      });
    },
    closeBundle() {
      if (!existsSync(CESIUM_SOURCE_DIR)) return;
      cpSync(CESIUM_SOURCE_DIR, join(__dirname, 'dist', 'cesium'), { recursive: true });
    },
  };
}

// Vite config that proxies /api to the local mock API on port 8001
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    cesiumStaticAssetsPlugin(),
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
