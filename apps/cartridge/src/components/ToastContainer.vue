<template>
  <div class="fixed bottom-4 right-4 flex flex-col gap-2 z-50">
    <div v-for="t in toasts" :key="t.id" :class="['px-3 py-2 rounded shadow', typeClass(t.type)]">
      <div class="text-sm">{{ t.message }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import useToastStore from '../stores/toast';

const toast = useToastStore();
const toasts = computed(() => toast.toasts);

function typeClass(type: string | undefined) {
  if (type === 'success') return 'bg-emerald-600 text-white';
  if (type === 'error') return 'bg-red-600 text-white';
  return 'bg-slate-700 text-white';
}
</script>
