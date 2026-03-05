import { defineStore } from 'pinia';
import type { Item } from '../api/items';

const STORAGE_KEY = 'dnd_cart_v1';

export type CartEntry = { id: string; name: string; priceSnapshot: number; quantity: number };

export const useCartStore = defineStore('cart', {
  state: () => ({ items: {} as Record<string, CartEntry> }),
  getters: {
    itemsArray: (state): CartEntry[] => Object.values(state.items),
    totalItems: (state): number => Object.values(state.items).reduce((s, e) => s + e.quantity, 0),
    totalPrice: (state): number => Object.values(state.items).reduce((s, e) => s + e.quantity * e.priceSnapshot, 0)
  },
  actions: {
    hydrate() {
      try {
        const raw = localStorage.getItem(STORAGE_KEY);
        if (!raw) return;
        const arr = JSON.parse(raw) as CartEntry[];
        this.items = arr.reduce((acc, e) => { acc[e.id] = e; return acc; }, {} as Record<string, CartEntry>);
      } catch (err) {}
    },
    persist() {
      try { localStorage.setItem(STORAGE_KEY, JSON.stringify(Object.values(this.items))); } catch (err) {}
    },
    addItem(item: Item, qty = 1) {
      const existing = this.items[item.id];
      if (existing) existing.quantity += qty;
      else this.items[item.id] = { id: item.id, name: item.name, priceSnapshot: item.price, quantity: qty };
      this.persist();
    },
    removeItem(id: string) { delete this.items[id]; this.persist(); },
    setQuantity(id: string, q: number) { if (!this.items[id]) return; const n = Math.max(0, Math.floor(q)); if (n<=0) delete this.items[id]; else this.items[id].quantity = n; this.persist(); },
    clearCart() { this.items = {}; this.persist(); }
  }
});

