export type ReagentSite = { id: string; name: string; lat: number; lng: number };

// Determine API base at runtime:
// 1) build-time env VITE_API_BASE (import.meta.env)
// 2) if running in browser on localhost -> use http://localhost:8001
// 3) otherwise use in-cluster service name http://mock-api:8001
const runtimeHost = typeof window !== 'undefined' ? window.location.hostname : '';
const isLocal = runtimeHost === 'localhost' || runtimeHost === '127.0.0.1' || runtimeHost === 'host.docker.internal';
const API_BASE = import.meta.env.VITE_API_BASE || (isLocal ? 'http://localhost:8001' : 'http://mock-api:8001');

const REQUEST_TIMEOUT_MS = 5_000;

export async function fetchReagentSites(): Promise<ReagentSite[]> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  try {
    const res = await fetch(`${API_BASE}/api/reagent-sites`, { signal: controller.signal });
    if (!res.ok) throw new Error('Failed to fetch reagent sites');
    return (await res.json()) as ReagentSite[];
  } finally {
    clearTimeout(timer);
  }
}
