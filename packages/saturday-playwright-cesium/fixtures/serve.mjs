// Minimal static file server for the Playwright webServer block — no new
// dependency, plain node:http/node:fs. Serves the fixture page plus Cesium's
// own Build/Cesium assets and the local planet texture, entirely offline.
//
// Binds to 127.0.0.1 explicitly (see architecture-notes.md's structural
// decision): saturday-playwright-leaflet's serve.mjs calls .listen(port) with
// no host, which binds all interfaces — that call must not be copied here.
import { createServer } from 'node:http';
import { readFile, stat } from 'node:fs/promises';
import { extname, join, resolve, sep } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = fileURLToPath(new URL('.', import.meta.url));
const packageRoot = resolve(here, '..');

const ROUTES = [
  { prefix: '/cesium/', root: join(packageRoot, 'node_modules/cesium/Build/Cesium') },
  { prefix: '/assets/', root: join(here, 'assets') },
];

const CONTENT_TYPES = {
  '.html': 'text/html',
  '.js': 'text/javascript',
  '.css': 'text/css',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.json': 'application/json',
  '.map': 'application/json',
  '.wasm': 'application/wasm',
  '.ktx2': 'application/octet-stream',
};

function contentTypeFor(filePath) {
  return CONTENT_TYPES[extname(filePath)] ?? 'application/octet-stream';
}

function resolveWithinRoot(root, requestPath) {
  const resolved = resolve(root, `.${requestPath}`);
  const boundary = root.endsWith(sep) ? root : root + sep;
  if (resolved !== root && !resolved.startsWith(boundary)) return null;
  return resolved;
}

function resolveRequestedFile(pathname) {
  if (pathname === '/') return join(here, 'index.html');
  const route = ROUTES.find((candidate) => pathname.startsWith(candidate.prefix));
  if (!route) return null;
  return resolveWithinRoot(route.root, pathname.slice(route.prefix.length - 1));
}

async function readIfFile(filePath) {
  const info = await stat(filePath);
  if (!info.isFile()) return null;
  return readFile(filePath);
}

async function handleRequest(req, res) {
  const pathname = decodeURIComponent(new URL(req.url, 'http://127.0.0.1').pathname);
  const filePath = resolveRequestedFile(pathname);

  try {
    const body = filePath ? await readIfFile(filePath) : null;
    if (!body) throw new Error('not found');
    res.writeHead(200, { 'Content-Type': contentTypeFor(filePath) });
    res.end(body);
  } catch {
    res.writeHead(404, { 'Content-Type': 'text/plain' });
    res.end('Not Found');
  }
}

const port = Number(process.env.FIXTURE_SERVER_PORT ?? 4174);
createServer((req, res) => {
  handleRequest(req, res);
}).listen(port, '127.0.0.1', () => {
  console.log(`saturday-playwright-cesium fixture server listening on http://127.0.0.1:${port}`);
});
