<template>
  <div class="bg-slate-800 rounded p-4">
    <h2 class="text-lg font-semibold mb-2">Runs</h2>
    <div v-if="runs.length === 0" class="text-sm text-slate-400">No runs yet</div>
    <ul v-else class="text-sm">
      <li v-for="r in runs" :key="r.id" class="flex items-center justify-between gap-4 py-2 border-b border-slate-700">
        <div>
          <div class="font-medium">{{ r.scenarioId }} <span class="text-xs text-slate-400">({{ r.suiteId }})</span></div>
          <div class="text-xs text-slate-400">{{ r.startedAt }}</div>
        </div>
        <div class="flex items-center gap-2">
          <span :class="statusClass(r.status)">{{ r.status }}</span>
          <a v-if="r.reportUrl" :href="formatReportUrl(r.reportUrl)" target="_blank" class="text-xs text-emerald-300 hover:underline">report</a>
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRunsStore } from '../stores/runs';

const runsStore = useRunsStore();
const runs = computed(() => runsStore.listRuns());

function formatReportUrl(url: string) {
  if (url.startsWith('/')) return `http://localhost:9001${url}`;
  return url;
}

function statusClass(s: string) {
  if (s === 'running') return 'text-yellow-300 font-semibold';
  if (s === 'passed') return 'text-emerald-300 font-semibold';
  if (s === 'failed') return 'text-red-400 font-semibold';
  return 'text-slate-300';
}
</script>

