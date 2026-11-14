export type GeoPoint = { lat: number; lng: number };
export type Expedition = { id: string; name?: string; owner: string; from: GeoPoint; to: GeoPoint };

// Determine API base at runtime:
// 1) build-time env VITE_API_BASE (import.meta.env)
// 2) if running in browser on localhost -> use http://localhost:8001
// 3) otherwise use in-cluster service name http://mock-api:8001
const runtimeHost = typeof window !== 'undefined' ? window.location.hostname : '';
const isLocal = runtimeHost === 'localhost' || runtimeHost === '127.0.0.1' || runtimeHost === 'host.docker.internal';
const API_BASE = import.meta.env.VITE_API_BASE || (isLocal ? 'http://localhost:8001' : 'http://mock-api:8001');

const REQUEST_TIMEOUT_MS = 5_000;

// GET /api/expeditions tolerates an absent or invalid token (it only affects
// which expeditions come back), so the header is attached when present and
// omitted otherwise — never required. No Content-Type: this is a bodyless GET.
function getAuthHeaders(): HeadersInit {
  const token = localStorage.getItem('dnd_auth_token');
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export async function fetchExpeditions(): Promise<Expedition[]> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  try {
    const res = await fetch(`${API_BASE}/api/expeditions`, {
      headers: getAuthHeaders(),
      signal: controller.signal
    });
    if (!res.ok) throw new Error('Failed to fetch expeditions');
    return (await res.json()) as Expedition[];
  } finally {
    clearTimeout(timer);
  }
}
