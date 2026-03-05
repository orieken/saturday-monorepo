<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
    <div class="w-full max-w-6xl rounded-lg bg-slate-900 border border-slate-700 shadow-xl flex flex-col h-[85vh]">
      <!-- Header -->
      <div class="flex items-center justify-between px-4 py-3 border-b border-slate-700">
        <div>
          <div class="text-xs text-slate-400">Running Scenario</div>
          <div class="text-base font-semibold text-slate-50">{{ scenario?.name }}</div>
          <div class="text-xs text-slate-500">{{ scenario?.file }}:{{ scenario?.line }}</div>
        </div>
        <div class="flex items-center gap-3">
          <!-- Link removed -->
          <button class="px-3 py-1.5 text-sm rounded bg-slate-800 hover:bg-slate-700 text-slate-200" @click="$emit('close')">
            ✕ Close
          </button>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex border-b border-slate-700 px-4">
        <button 
          @click="activeTab = 'logs'"
          class="px-4 py-2 text-sm font-medium border-b-2 transition-colors"
          :class="activeTab === 'logs' ? 'border-sky-500 text-sky-400' : 'border-transparent text-slate-400 hover:text-slate-300'">
          <span class="flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Logs
          </span>
        </button>
        <button 
          @click="activeTab = 'report'"
          :disabled="!localRun?.reportUrl"
          class="px-4 py-2 text-sm font-medium border-b-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          :class="activeTab === 'report' ? 'border-sky-500 text-sky-400' : 'border-transparent text-slate-400 hover:text-slate-300'">
          <span class="flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Report
            <span v-if="!localRun?.reportUrl" class="text-xs">(not ready)</span>
          </span>
        </button>
      </div>

      <!-- Tab Content -->
      <div class="flex-1 overflow-hidden">
        <!-- Logs Tab -->
        <div v-show="activeTab === 'logs'" class="h-full overflow-auto bg-black/80 p-4">
          <div class="bg-black/90 p-4 rounded text-xs text-slate-200 font-mono">
            <pre class="whitespace-pre-wrap"><span v-for="(line, idx) in visibleLines" :key="idx">{{ line }}
</span></pre>
            <div v-if="visibleLines.length === 0" class="text-slate-500 text-center py-8">
              <div class="inline-block animate-pulse">Waiting for logs...</div>
            </div>
          </div>
        </div>

        <!-- Report Tab -->
        <div v-show="activeTab === 'report'" class="h-full overflow-auto bg-slate-950 p-4">
          <CucumberReportViewer v-if="localRun?.id" :runId="localRun.id" />
        </div>
      </div>

      <!-- Footer -->
      <div class="px-4 py-2 border-t border-slate-700 flex items-center justify-between text-xs text-slate-400">
        <span>
          Status:
          <strong :class="statusClass">{{ localRun?.status?.toUpperCase() || 'UNKNOWN' }}</strong>
        </span>
        <span v-if="localRun && localRun.status !== 'running'">Test completed. You can close this window.</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, computed, watch } from 'vue';
import type { Run } from '../api/runs';
import type { CucumberScenarioRef } from '../api/cucumber';
import CucumberReportViewer from './CucumberReportViewer.vue';

const props = defineProps<{
  run: Run | null;
  scenario?: CucumberScenarioRef | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

// Create a local reactive copy of the run that we can update
const localRun = ref<Run | null>(props.run);
const scenario = props.scenario;

const visibleLines = ref<string[]>([]);
const activeTab = ref<'logs' | 'report'>('logs');
let es: EventSource | null = null;
let pollingInterval: number | null = null;

// API base URL (same as used in api/runs.ts)
const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:9001';

const reportUrlFull = computed(() => {
  if (!localRun.value || !localRun.value.reportUrl) return '#';
  if (localRun.value.reportUrl.startsWith('/')) return `${API_BASE}${localRun.value.reportUrl}`;
  return localRun.value.reportUrl;
});


const statusClass = computed(() => {
  if (!localRun.value) return 'text-slate-300';
  if (localRun.value.status === 'running') return 'text-amber-400';
  if (localRun.value.status === 'passed') return 'text-emerald-400';
  if (localRun.value.status === 'failed') return 'text-rose-400';
  return 'text-slate-300';
});

function startSSE() {
  if (!localRun.value) return;
  // Clear existing lines as we will stream everything from scratch
  visibleLines.value = [];
  
  try {
    // Use full API URL for SSE connection
    const streamUrl = `${API_BASE}/api/runs/${encodeURIComponent(localRun.value.id)}/stream`;
    console.log('[SSE] Connecting to:', streamUrl);
    es = new EventSource(streamUrl);
    
    es.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data);
        if (msg.type === 'log' && msg.text) {
          visibleLines.value.push(msg.text); // Use push for reactivity
        } else if (msg.type === 'status' && msg.status) {
          visibleLines.value.push(`[status] ${msg.status}`);
          
          // Update local run status
          if (localRun.value) {
            localRun.value.status = msg.status;
          }
          
          // Close SSE connection when run completes
          if (msg.status !== 'running') {
            if (es) {
              es.close();
              es = null;
            }
          }
        }
      } catch (e) {
        // ignore
      }
    };
    es.onerror = (err) => {
      console.error('[SSE] Connection error:', err);
      if (es) {
        es.close();
        es = null;
      }
    };
  } catch (e) {
    // ignore
  }
}

onMounted(() => {
  startSSE();
});

onBeforeUnmount(() => {
  if (es) {
    es.close();
    es = null;
  }
});

watch(
  () => props.run,
  (newRun) => {
    // Update local run when prop changes
    localRun.value = newRun;
    
    // restart SSE when run changes
    if (es) {
        es.close();
        es = null;
    }
    startSSE();
  },
  { immediate: true }
);
</script>
