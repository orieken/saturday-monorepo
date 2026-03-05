<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-xl font-semibold text-slate-100">Run History</h2>
      <button @click="refresh" class="text-sm text-cyan-400 hover:text-cyan-300">Refresh</button>
    </div>

    <div class="bg-slate-900 rounded-lg border border-slate-800 overflow-hidden">
      <table class="w-full text-left text-sm text-slate-300">
        <thead class="bg-slate-950 text-slate-400 uppercase text-xs font-semibold">
          <tr>
            <th class="px-4 py-3">Run ID</th>
            <th class="px-4 py-3">Scenario</th>
            <th class="px-4 py-3">Executor</th>
            <th class="px-4 py-3">Started At</th>
            <th class="px-4 py-3">Status</th>
            <th class="px-4 py-3 text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-800">
          <tr v-for="run in runs" :key="run.id" class="hover:bg-slate-800/50 transition-colors">
            <td class="px-4 py-3 font-mono text-xs text-slate-500">{{ run.id.slice(0, 8) }}...</td>
            <td class="px-4 py-3">
              <div class="font-medium text-slate-200">{{ run.scenarioId }}</div>
              <div class="text-xs text-slate-500">{{ run.suiteId }}</div>
            </td>
            <td class="px-4 py-3">
              <span class="px-2 py-0.5 rounded bg-slate-800 text-xs">{{ run.executor }}</span>
            </td>
            <td class="px-4 py-3 text-slate-400">{{ formatDate(run.startedAt) }}</td>
            <td class="px-4 py-3">
              <span :class="statusClass(run.status)" class="font-medium">{{ run.status }}</span>
            </td>
            <td class="px-4 py-3 text-right">
              <a 
                v-if="run.reportUrl" 
                :href="formatReportUrl(run.reportUrl)" 
                target="_blank" 
                class="inline-flex items-center gap-1 text-cyan-400 hover:text-cyan-300 hover:underline"
              >
                View Report ↗
              </a>
              <span v-else class="text-slate-600">-</span>
            </td>
          </tr>
          <tr v-if="runs.length === 0">
            <td colspan="6" class="px-4 py-8 text-center text-slate-500">No runs recorded yet.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRunsStore } from '../stores/runs';

const runsStore = useRunsStore();
const runs = computed(() => runsStore.listRuns());

function refresh() {
  // The store is reactive, but if we needed to re-fetch from server we could do it here.
  // Currently runs are persisted in local storage and updated via polling.
}

function formatDate(iso?: string) {
  if (!iso) return '-';
  return new Date(iso).toLocaleString();
}

function formatReportUrl(url: string) {
  if (url.startsWith('/')) return `http://localhost:9001${url}`;
  return url;
}

function statusClass(s: string) {
  if (s === 'running') return 'text-amber-400';
  if (s === 'passed') return 'text-emerald-400';
  if (s === 'failed') return 'text-rose-400';
  return 'text-slate-400';
}
</script>
