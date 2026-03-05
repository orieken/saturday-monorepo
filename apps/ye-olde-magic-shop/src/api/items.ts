export type Item = { id: string; name: string; price: number; rarity?: string; image?: string; description?: string };

// Determine API base at runtime:
// 1) build-time env VITE_API_BASE (import.meta.env)
// 2) if running in browser on localhost -> use http://localhost:8001
// 3) otherwise use in-cluster service name http://mock-api:8001
const runtimeHost = typeof window !== 'undefined' ? window.location.hostname : '';
const isLocal = runtimeHost === 'localhost' || runtimeHost === '127.0.0.1' || runtimeHost === 'host.docker.internal';
const API_BASE = import.meta.env.VITE_API_BASE || (isLocal ? 'http://localhost:8001' : 'http://mock-api:8001');

export async function fetchItems(): Promise<Item[]> {
  const res = await fetch(`${API_BASE}/api/items`);
  if (!res.ok) throw new Error('Failed to fetch items');
  return (await res.json()) as Item[];
}

export async function fetchItemById(id: string): Promise<Item | null> {
  const res = await fetch(`${API_BASE}/api/items/${encodeURIComponent(id)}`);
  if (res.status === 404) return null;
  if (!res.ok) throw new Error('Failed to fetch item');
  return (await res.json()) as Item;
}
