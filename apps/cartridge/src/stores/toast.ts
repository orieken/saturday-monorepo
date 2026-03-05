import { defineStore } from 'pinia';
import { ref } from 'vue';

export type Toast = {
  id: number;
  message: string;
  type?: 'info' | 'success' | 'error';
  ttl?: number;
};

export default defineStore('toast', () => {
  const toasts = ref<Toast[]>([]);

  function add(message: string, type: Toast['type'] = 'info', ttl = 5000) {
    const t = { id: Date.now() + Math.floor(Math.random() * 1000), message, type, ttl };
    toasts.value.push(t);
    // auto remove
    setTimeout(() => remove(t.id), ttl);
    return t.id;
  }

  function remove(id: number) {
    const idx = toasts.value.findIndex(t => t.id === id);
    if (idx !== -1) toasts.value.splice(idx, 1);
  }

  return { toasts, add, remove };
});
