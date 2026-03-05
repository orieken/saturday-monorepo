export type Item = {
  id: string;
  name: string;
  price: number;
  rarity?: string;
  image?: string;
  description?: string;
};

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:9001';

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
