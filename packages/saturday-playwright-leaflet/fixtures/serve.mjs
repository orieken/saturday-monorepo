// Minimal static file server for the Playwright webServer block — no new
// dependency, plain node:http/node:fs. Serves the fixture page plus the
// local Leaflet/Turf/image assets it references, entirely offline.
import { createServer } from 'node:http';
import { readFile, stat } from 'node:fs/promises';
import { extname, join, resolve, sep } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = fileURLToPath(new URL('.', import.meta.url));
const packageRoot = resolve(here, '..');

const ROUTES = [
  { prefix: '/leaflet/', root: join(packageRoot, 'node_modules/leaflet/dist') },
  { prefix: '/assets/', root: join(here, 'assets') }
];

const CONTENT_TYPES = {
  '.html': 'text/html',
  '.js': 'text/javascript',
  '.css': 'text/css',
  '.png': 'image/png',
  '.map': 'application/json'
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
  if (pathname === '/turf.min.js') {
    return join(packageRoot, 'node_modules/@turf/turf/turf.min.js');
  }
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
  const pathname = decodeURIComponent(new URL(req.url, 'http://localhost').pathname);
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

const port = Number(process.env.FIXTURE_SERVER_PORT ?? 4173);
createServer((req, res) => {
  handleRequest(req, res);
}).listen(port, () => {
  console.log(`saturday-playwright-leaflet fixture server listening on http://localhost:${port}`);
});
