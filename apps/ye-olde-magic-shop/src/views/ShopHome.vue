<template>
  <div>
    <div class="mb-8 border-b border-dnd-slate pb-4">
      <h2 class="text-3xl font-magical text-dnd-parchment">Wares for Sale</h2>
      <p class="text-gray-400 mt-2">Browse our collection of legendary artifacts and magical trinkets.</p>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <ItemCard v-for="it in items" :key="it.id" :item="it" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { fetchItems } from '../api/items';
import type { Item } from '../api/items';
import ItemCard from '../components/ItemCard.vue';

const items = ref<Item[]>([]);

onMounted(async () => {
  try { items.value = await fetchItems(); } catch (e) { console.error(e); }
});
</script>

