import type { CartEntry } from '../stores/cartStore';
import type { Address } from '../stores/authStore';

export interface Order {
  id: string;
  userId: string;
  items: CartEntry[];
  total: number;
  status: 'pending' | 'processing' | 'completed' | 'cancelled';
  createdAt: string;
  shippingAddress?: Address | null;
}

// Determine API base at runtime
const runtimeHost = typeof window !== 'undefined' ? window.location.hostname : '';
const isLocal = runtimeHost === 'localhost' || runtimeHost === '127.0.0.1' || runtimeHost === 'host.docker.internal';
const API_BASE = import.meta.env.VITE_API_BASE || (isLocal ? 'http://localhost:8001' : 'http://mock-api:8001');

function getAuthHeaders(): HeadersInit {
  const token = localStorage.getItem('dnd_auth_token');
  return {
    'Content-Type': 'application/json',
    ...(token ? { 'Authorization': `Bearer ${token}` } : {})
  };
}

export async function fetchOrders(): Promise<Order[]> {
  const res = await fetch(`${API_BASE}/api/orders`, {
    headers: getAuthHeaders()
  });
  
  if (!res.ok) throw new Error('Failed to fetch orders');
  return (await res.json()) as Order[];
}

export async function fetchOrderById(id: string): Promise<Order | null> {
  const res = await fetch(`${API_BASE}/api/orders/${encodeURIComponent(id)}`, {
    headers: getAuthHeaders()
  });
  
  if (res.status === 404) return null;
  if (!res.ok) throw new Error('Failed to fetch order');
  return (await res.json()) as Order;
}

export async function createOrder(
  items: CartEntry[],
  total: number,
  shippingAddress?: Address
): Promise<Order> {
  const res = await fetch(`${API_BASE}/api/orders`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify({ items, total, shippingAddress })
  });
  
  if (!res.ok) throw new Error('Failed to create order');
  return (await res.json()) as Order;
}
