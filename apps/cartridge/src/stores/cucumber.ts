import { defineStore } from 'pinia';
import { ref } from 'vue';
import fetchCucumberIndex, { CucumberIndex } from '../api/cucumber';

export default defineStore('cucumber', () => {
  const indexes = ref<CucumberIndex[] | null>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);

  async function load() {
    loading.value = true;
    error.value = null;
    try {
      indexes.value = await fetchCucumberIndex();
    } catch (e: any) {
      error.value = e?.message || String(e);
    } finally {
      loading.value = false;
    }
  }

  return {
    indexes,
    loading,
    error,
    load,
  };
});
